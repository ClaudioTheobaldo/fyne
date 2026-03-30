// streaming_bench benchmarks 4 frame upload strategies:
//
//  1. canvas.Image          — delete + create + TexImage2D from CPU (baseline)
//  2. TexSubImage2D only    — reuse texture, sync CPU→GPU copy
//  3. PBO only              — new texture each frame, async DMA from PBO
//  4. PBO + TexSubImage2D   — reuse texture + async DMA from PBO
//
// Tests: single 1080p, single 4K, 4x 1080p multi-cam, shifting resolutions.
package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	glpainter "fyne.io/fyne/v2/internal/painter/gl"
	"fyne.io/fyne/v2/widget"
)

func generateFrame(w, h int, seed byte) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i] = seed
		img.Pix[i+1] = byte(i >> 8)
		img.Pix[i+2] = byte(i >> 16)
		img.Pix[i+3] = 255
	}
	return img
}

type result struct {
	draws int64
	ns    int64
}

func (r result) avgUs() float64 {
	if r.draws == 0 {
		return 0
	}
	return float64(r.ns) / float64(r.draws) / 1000.0
}

func (r result) totalMs() float64 {
	return float64(r.ns) / 1e6
}

func pregenFrames(resolutions []image.Point, count int) map[image.Point][]*image.RGBA {
	pools := make(map[image.Point][]*image.RGBA)
	for _, r := range resolutions {
		if _, ok := pools[r]; ok {
			continue
		}
		frames := make([]*image.RGBA, count)
		for i := range frames {
			frames[i] = generateFrame(r.X, r.Y, byte(i))
		}
		pools[r] = frames
	}
	return pools
}

func resetCounters() {
	glpainter.BenchImageDrawNs.Store(0)
	glpainter.BenchImageDrawCount.Store(0)
	glpainter.BenchStreamDrawNs.Store(0)
	glpainter.BenchStreamDrawCount.Store(0)
}

type benchConfig struct {
	name           string
	duration       time.Duration
	numStreams     int
	resolutions    []image.Point
	changeInterval time.Duration
}

// pushFrames runs the main benchmark loop, pushing frames for the given duration.
func pushFrames(duration time.Duration, resolutions []image.Point, changeInterval time.Duration, pools map[image.Point][]*image.RGBA, pushFn func(*image.RGBA)) {
	resIdx := 0
	currentRes := resolutions[0]
	lastResChange := time.Now()
	start := time.Now()
	var count int64

	for time.Since(start) < duration {
		if len(resolutions) > 1 && changeInterval > 0 && time.Since(lastResChange) >= changeInterval {
			resIdx = (resIdx + 1) % len(resolutions)
			currentRes = resolutions[resIdx]
			lastResChange = time.Now()
		}
		pool := pools[currentRes]
		frame := pool[count%int64(len(pool))]

		fyne.DoAndWait(func() {
			pushFn(frame)
		})
		count++
	}
}

func runVariant(w fyne.Window, label string, cfg benchConfig, pools map[image.Point][]*image.RGBA, setup func() (fyne.CanvasObject, func(*image.RGBA)), useImageCounter bool) result {
	fmt.Printf("    %-24s ", label+"...")

	obj, pushFn := setup()

	statusLabel := widget.NewLabel(label)
	content := container.NewBorder(statusLabel, nil, nil, nil, obj)
	fyne.DoAndWait(func() {
		w.SetContent(content)
		w.SetTitle("Bench: " + label)
	})
	runtime.GC()
	resetCounters()

	done := make(chan result)
	go func() {
		time.Sleep(200 * time.Millisecond)
		pushFrames(cfg.duration, cfg.resolutions, cfg.changeInterval, pools, pushFn)
		time.Sleep(100 * time.Millisecond) // let final paint complete

		var r result
		if useImageCounter {
			r = result{draws: glpainter.BenchImageDrawCount.Load(), ns: glpainter.BenchImageDrawNs.Load()}
		} else {
			r = result{draws: glpainter.BenchStreamDrawCount.Load(), ns: glpainter.BenchStreamDrawNs.Load()}
		}
		done <- r
	}()

	r := <-done
	fmt.Printf("%5d draws | avg %8.1f us | total %7.1f ms\n", r.draws, r.avgUs(), r.totalMs())
	return r
}

func runTest(w fyne.Window, cfg benchConfig) [4]result {
	fmt.Printf("\n=== %s ===\n", cfg.name)
	if len(cfg.resolutions) > 1 {
		fmt.Printf("    Resolutions: ")
		for i, r := range cfg.resolutions {
			if i > 0 {
				fmt.Printf(" -> ")
			}
			fmt.Printf("%dx%d", r.X, r.Y)
		}
		fmt.Println()
	} else {
		fmt.Printf("    Resolution: %dx%d  Streams: %d\n", cfg.resolutions[0].X, cfg.resolutions[0].Y, cfg.numStreams)
	}

	pools := pregenFrames(cfg.resolutions, 16)
	runtime.GC()

	var results [4]result

	// 1. canvas.Image (baseline)
	results[0] = runVariant(w, "canvas.Image", cfg, pools, func() (fyne.CanvasObject, func(*image.RGBA)) {
		images := make([]*canvas.Image, cfg.numStreams)
		objs := make([]fyne.CanvasObject, cfg.numStreams)
		for i := range images {
			images[i] = canvas.NewImageFromImage(pools[cfg.resolutions[0]][0])
			images[i].FillMode = canvas.ImageFillStretch
			images[i].ScaleMode = canvas.ImageScaleFastest
			objs[i] = images[i]
		}
		grid := container.NewGridWithColumns(max(1, cfg.numStreams), objs...)
		return grid, func(frame *image.RGBA) {
			for _, img := range images {
				img.Image = frame
				img.Refresh()
			}
		}
	}, true)

	// 2. TexSubImage2D only (no PBO, reuse texture)
	results[1] = runVariant(w, "TexSubImage2D only", cfg, pools, func() (fyne.CanvasObject, func(*image.RGBA)) {
		streams := make([]*canvas.StreamingImage, cfg.numStreams)
		objs := make([]fyne.CanvasObject, cfg.numStreams)
		for i := range streams {
			streams[i] = canvas.NewStreamingImage()
			streams[i].FillMode = canvas.ImageFillStretch
			streams[i].ScaleMode = canvas.ImageScaleFastest
			streams[i].DisablePBO = true
			streams[i].DisableTexReuse = false
			objs[i] = streams[i]
		}
		grid := container.NewGridWithColumns(max(1, cfg.numStreams), objs...)
		return grid, func(frame *image.RGBA) {
			for _, s := range streams {
				s.UpdateFrame(frame)
			}
		}
	}, false)

	// 3. PBO only (PBO + new texture each frame, no reuse)
	results[2] = runVariant(w, "PBO only", cfg, pools, func() (fyne.CanvasObject, func(*image.RGBA)) {
		streams := make([]*canvas.StreamingImage, cfg.numStreams)
		objs := make([]fyne.CanvasObject, cfg.numStreams)
		for i := range streams {
			streams[i] = canvas.NewStreamingImage()
			streams[i].FillMode = canvas.ImageFillStretch
			streams[i].ScaleMode = canvas.ImageScaleFastest
			streams[i].DisablePBO = false
			streams[i].DisableTexReuse = true
			objs[i] = streams[i]
		}
		grid := container.NewGridWithColumns(max(1, cfg.numStreams), objs...)
		return grid, func(frame *image.RGBA) {
			for _, s := range streams {
				s.UpdateFrame(frame)
			}
		}
	}, false)

	// 4. PBO + TexSubImage2D (full optimization)
	results[3] = runVariant(w, "PBO + TexSubImage2D", cfg, pools, func() (fyne.CanvasObject, func(*image.RGBA)) {
		streams := make([]*canvas.StreamingImage, cfg.numStreams)
		objs := make([]fyne.CanvasObject, cfg.numStreams)
		for i := range streams {
			streams[i] = canvas.NewStreamingImage()
			streams[i].FillMode = canvas.ImageFillStretch
			streams[i].ScaleMode = canvas.ImageScaleFastest
			streams[i].DisablePBO = false
			streams[i].DisableTexReuse = false
			objs[i] = streams[i]
		}
		grid := container.NewGridWithColumns(max(1, cfg.numStreams), objs...)
		return grid, func(frame *image.RGBA) {
			for _, s := range streams {
				s.UpdateFrame(frame)
			}
		}
	}, false)

	return results
}

var labels = [4]string{"canvas.Image", "TexSubImage2D", "PBO only", "PBO+TexSubImage2D"}

func main() {
	// Uncap both the software ticker and GPU VSync
	os.Setenv("FYNE_TICK_RATE", "10000")
	os.Setenv("FYNE_VSYNC", "0")

	fmt.Println("============================================================")
	fmt.Println("  Frame Upload Strategy Benchmark (4-way comparison)")
	fmt.Println("  VSync OFF, render loop uncapped")
	fmt.Println("============================================================")

	a := app.New()
	w := a.NewWindow("Streaming Benchmark")
	w.Resize(fyne.NewSize(960, 720))

	type testRun struct {
		cfg     benchConfig
		results [4]result
	}
	var runs []testRun

	go func() {
		time.Sleep(500 * time.Millisecond)

		configs := []benchConfig{
			{
				name:       "Single 1080p stream",
				duration:   10 * time.Second,
				numStreams: 1,
				resolutions: []image.Point{{X: 1920, Y: 1080}},
			},
			{
				name:       "Single 4K stream",
				duration:   10 * time.Second,
				numStreams: 1,
				resolutions: []image.Point{{X: 3840, Y: 2160}},
			},
			{
				name:       "4x 1080p streams",
				duration:   10 * time.Second,
				numStreams: 4,
				resolutions: []image.Point{{X: 1920, Y: 1080}},
			},
			{
				name:       "16x 1080p streams",
				duration:   10 * time.Second,
				numStreams: 16,
				resolutions: []image.Point{{X: 1920, Y: 1080}},
			},
			{
				name:           "Shifting resolutions (every 1s)",
				duration:       10 * time.Second,
				numStreams:     1,
				changeInterval: 1 * time.Second,
				resolutions: func() []image.Point {
					res := []image.Point{
						{X: 1920, Y: 1080}, {X: 1280, Y: 720}, {X: 640, Y: 480},
						{X: 3840, Y: 2160}, {X: 800, Y: 600}, {X: 1024, Y: 768},
						{X: 1600, Y: 900}, {X: 2560, Y: 1440}, {X: 720, Y: 576},
						{X: 1920, Y: 1080},
					}
					rand.Shuffle(len(res), func(i, j int) { res[i], res[j] = res[j], res[i] })
					return res
				}(),
			},
		}

		for _, cfg := range configs {
			results := runTest(w, cfg)
			runs = append(runs, testRun{cfg: cfg, results: results})
		}

		// Summary table
		fmt.Println("\n============================================================")
		fmt.Println("  SUMMARY  (avg microseconds per draw call)")
		fmt.Println("============================================================")
		fmt.Printf("%-30s  %12s  %12s  %12s  %12s\n", "Test", labels[0], labels[1], labels[2], labels[3])
		fmt.Printf("%-30s  %12s  %12s  %12s  %12s\n", "----", "----------", "----------", "----------", "----------")
		for _, r := range runs {
			fmt.Printf("%-30s", r.cfg.name)
			for i := 0; i < 4; i++ {
				fmt.Printf("  %9.1f us", r.results[i].avgUs())
			}
			fmt.Println()
		}

		fmt.Println("\nSpeedups vs canvas.Image:")
		fmt.Printf("%-30s  %12s  %12s  %12s  %12s\n", "Test", labels[0], labels[1], labels[2], labels[3])
		fmt.Printf("%-30s  %12s  %12s  %12s  %12s\n", "----", "----------", "----------", "----------", "----------")
		for _, r := range runs {
			fmt.Printf("%-30s", r.cfg.name)
			baseline := r.results[0].avgUs()
			for i := 0; i < 4; i++ {
				if r.results[i].avgUs() > 0 && baseline > 0 {
					fmt.Printf("  %9.2fx   ", baseline/r.results[i].avgUs())
				} else {
					fmt.Printf("  %12s", "N/A")
				}
			}
			fmt.Println()
		}

		fyne.DoAndWait(func() {
			items := []fyne.CanvasObject{widget.NewLabel("Benchmark Complete! Check terminal for full results.")}
			items = append(items, widget.NewButton("Close", func() { w.Close() }))
			w.SetContent(container.NewVBox(items...))
			w.SetTitle("Results")
		})
	}()

	w.ShowAndRun()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
