// streaming_bench benchmarks canvas.Image vs canvas.StreamingImage.
//
// Measures actual time inside GL draw functions (instrumented).
// Test 1: Single 1080p stream — baseline.
// Test 2: Single 4K stream — stresses GPU memory bandwidth.
// Test 3: 4x 1080p streams — simulates multi-camera view.
// Test 4: Shifting resolutions — 10 seconds, resolution changes every second.
package main

import (
	"fmt"
	"image"
	"math/rand"
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

type testResult struct {
	label     string
	drawCount int64
	drawNs    int64
	wallTime  time.Duration
}

func (r testResult) avgDrawUs() float64 {
	if r.drawCount == 0 {
		return 0
	}
	return float64(r.drawNs) / float64(r.drawCount) / 1000.0
}

func (r testResult) actualFps() float64 {
	if r.wallTime == 0 {
		return 0
	}
	return float64(r.drawCount) / r.wallTime.Seconds()
}

func resetCounters() {
	glpainter.BenchImageDrawNs.Store(0)
	glpainter.BenchImageDrawCount.Store(0)
	glpainter.BenchStreamDrawNs.Store(0)
	glpainter.BenchStreamDrawCount.Store(0)
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

type benchConfig struct {
	name           string
	duration       time.Duration
	numStreams     int
	resolutions    []image.Point
	changeInterval time.Duration
}

func runBench(w fyne.Window, cfg benchConfig) (imgResult, stmResult testResult) {
	fmt.Printf("\n--- %s ---\n", cfg.name)
	if len(cfg.resolutions) > 1 {
		fmt.Printf("  Resolutions: ")
		for i, r := range cfg.resolutions {
			if i > 0 {
				fmt.Printf(" -> ")
			}
			fmt.Printf("%dx%d", r.X, r.Y)
		}
		fmt.Println()
	} else {
		fmt.Printf("  Resolution: %dx%d  Streams: %d\n", cfg.resolutions[0].X, cfg.resolutions[0].Y, cfg.numStreams)
	}

	pools := pregenFrames(cfg.resolutions, 16)
	runtime.GC()

	// --- canvas.Image ---
	fmt.Println("  Running canvas.Image...")
	{
		statusLabel := widget.NewLabel("canvas.Image...")
		images := make([]*canvas.Image, cfg.numStreams)
		containers := make([]fyne.CanvasObject, cfg.numStreams)
		for i := range images {
			images[i] = canvas.NewImageFromImage(pools[cfg.resolutions[0]][0])
			images[i].FillMode = canvas.ImageFillStretch
			images[i].ScaleMode = canvas.ImageScaleFastest
			containers[i] = images[i]
		}

		grid := container.NewGridWithColumns(max(1, cfg.numStreams/max(1, (cfg.numStreams+1)/2)), containers...)
		content := container.NewBorder(statusLabel, nil, nil, nil, grid)
		fyne.DoAndWait(func() {
			w.SetContent(content)
			w.SetTitle(fmt.Sprintf("Bench: canvas.Image - %s", cfg.name))
		})
		resetCounters()

		done := make(chan struct{})
		go func() {
			time.Sleep(300 * time.Millisecond)
			resIdx := 0
			currentRes := cfg.resolutions[0]
			lastResChange := time.Now()
			start := time.Now()
			var pushCount int64

			for time.Since(start) < cfg.duration {
				if len(cfg.resolutions) > 1 && cfg.changeInterval > 0 && time.Since(lastResChange) >= cfg.changeInterval {
					resIdx = (resIdx + 1) % len(cfg.resolutions)
					currentRes = cfg.resolutions[resIdx]
					lastResChange = time.Now()
				}

				pool := pools[currentRes]
				frame := pool[pushCount%int64(len(pool))]

				fyne.DoAndWait(func() {
					for _, img := range images {
						img.Image = frame
						img.Refresh()
					}
				})
				pushCount++
			}

			close(done)
		}()
		<-done
		time.Sleep(100 * time.Millisecond) // let final paint complete

		imgResult = testResult{
			label:     "canvas.Image",
			drawCount: glpainter.BenchImageDrawCount.Load(),
			drawNs:    glpainter.BenchImageDrawNs.Load(),
			wallTime:  cfg.duration,
		}
	}

	// --- StreamingImage ---
	fmt.Println("  Running StreamingImage...")
	{
		statusLabel := widget.NewLabel("StreamingImage...")
		streams := make([]*canvas.StreamingImage, cfg.numStreams)
		containers := make([]fyne.CanvasObject, cfg.numStreams)
		for i := range streams {
			streams[i] = canvas.NewStreamingImage()
			streams[i].FillMode = canvas.ImageFillStretch
			streams[i].ScaleMode = canvas.ImageScaleFastest
			containers[i] = streams[i]
		}

		grid := container.NewGridWithColumns(max(1, cfg.numStreams/max(1, (cfg.numStreams+1)/2)), containers...)
		content := container.NewBorder(statusLabel, nil, nil, nil, grid)
		fyne.DoAndWait(func() {
			w.SetContent(content)
			w.SetTitle(fmt.Sprintf("Bench: StreamingImage - %s", cfg.name))
		})
		resetCounters()

		done := make(chan struct{})
		go func() {
			time.Sleep(300 * time.Millisecond)
			resIdx := 0
			currentRes := cfg.resolutions[0]
			lastResChange := time.Now()
			start := time.Now()
			var pushCount int64

			for time.Since(start) < cfg.duration {
				if len(cfg.resolutions) > 1 && cfg.changeInterval > 0 && time.Since(lastResChange) >= cfg.changeInterval {
					resIdx = (resIdx + 1) % len(cfg.resolutions)
					currentRes = cfg.resolutions[resIdx]
					lastResChange = time.Now()
				}

				pool := pools[currentRes]
				frame := pool[pushCount%int64(len(pool))]

				fyne.DoAndWait(func() {
					for _, s := range streams {
						s.UpdateFrame(frame)
					}
				})
				pushCount++
			}

			close(done)
		}()
		<-done
		time.Sleep(100 * time.Millisecond)

		stmResult = testResult{
			label:     "StreamingImage",
			drawCount: glpainter.BenchStreamDrawCount.Load(),
			drawNs:    glpainter.BenchStreamDrawNs.Load(),
			wallTime:  cfg.duration,
		}
	}

	return imgResult, stmResult
}

func printResult(img, stm testResult) {
	speedup := float64(0)
	if stm.avgDrawUs() > 0 {
		speedup = img.avgDrawUs() / stm.avgDrawUs()
	}
	fmt.Printf("  canvas.Image:   %5d draws | avg %9.1f us/draw | total %8.1f ms in GL\n",
		img.drawCount, img.avgDrawUs(), float64(img.drawNs)/1e6)
	fmt.Printf("  StreamingImage: %5d draws | avg %9.1f us/draw | total %8.1f ms in GL\n",
		stm.drawCount, stm.avgDrawUs(), float64(stm.drawNs)/1e6)
	fmt.Printf("  Draw time ratio: %.2fx  (%s)\n",
		speedup, func() string {
			if speedup > 1.0 {
				return "StreamingImage faster"
			}
			return "canvas.Image faster"
		}())
}

func main() {
	fmt.Println("================================================")
	fmt.Println("  canvas.Image vs StreamingImage GL Benchmark")
	fmt.Println("  (instrumented draw function timing)")
	fmt.Println("================================================")

	a := app.New()
	w := a.NewWindow("Streaming Benchmark")
	w.Resize(fyne.NewSize(960, 720))

	type resultPair struct {
		cfg benchConfig
		img testResult
		stm testResult
	}
	var results []resultPair

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
				name:       "4x 1080p streams (multi-camera)",
				duration:   10 * time.Second,
				numStreams: 4,
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
			img, stm := runBench(w, cfg)
			printResult(img, stm)
			results = append(results, resultPair{cfg: cfg, img: img, stm: stm})
		}

		// Summary
		fmt.Println("\n================================================")
		fmt.Println("  SUMMARY")
		fmt.Println("================================================")
		for _, r := range results {
			speedup := float64(0)
			if r.stm.avgDrawUs() > 0 {
				speedup = r.img.avgDrawUs() / r.stm.avgDrawUs()
			}
			fmt.Printf("%-35s  Image=%7.1fus  Stream=%7.1fus  Ratio=%.2fx\n",
				r.cfg.name, r.img.avgDrawUs(), r.stm.avgDrawUs(), speedup)
		}

		fyne.DoAndWait(func() {
			items := []fyne.CanvasObject{widget.NewLabel("Benchmark Complete!")}
			for _, r := range results {
				speedup := r.img.avgDrawUs() / r.stm.avgDrawUs()
				items = append(items, widget.NewLabel(fmt.Sprintf("%-30s  %.1fus vs %.1fus  %.2fx",
					r.cfg.name, r.img.avgDrawUs(), r.stm.avgDrawUs(), speedup)))
			}
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
