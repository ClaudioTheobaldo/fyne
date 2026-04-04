package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	sW float32 = 200
	sH float32 = 110
)

// ── colors ────────────────────────────────────────────────────────────────────

var (
	clRed    = color.NRGBA{220, 60, 60, 255}
	clBlue   = color.NRGBA{37, 99, 235, 255}
	clGreen  = color.NRGBA{40, 180, 80, 255}
	clPurple = color.NRGBA{124, 58, 237, 255}
	clOrange = color.NRGBA{240, 140, 40, 255}
	clCyan   = color.NRGBA{40, 200, 220, 255}
	clPink   = color.NRGBA{240, 100, 180, 255}
	clDark   = color.NRGBA{10, 10, 20, 255}
	clWhite  = color.NRGBA{240, 240, 240, 255}
)

// ── sample factories ──────────────────────────────────────────────────────────

type sampleFn func() fyne.CanvasObject

func hGrad(a, b color.Color) sampleFn {
	return func() fyne.CanvasObject {
		g := canvas.NewHorizontalGradient(a, b)
		g.SetMinSize(fyne.NewSize(sW, sH))
		return g
	}
}
func vGrad(a, b color.Color) sampleFn {
	return func() fyne.CanvasObject {
		g := canvas.NewVerticalGradient(a, b)
		g.SetMinSize(fyne.NewSize(sW, sH))
		return g
	}
}
func solid(c color.Color) sampleFn {
	return func() fyne.CanvasObject {
		r := canvas.NewRectangle(c)
		r.SetMinSize(fyne.NewSize(sW, sH))
		return r
	}
}

// ── addEffecter interface ─────────────────────────────────────────────────────

type addEffecter interface {
	fyne.CanvasObject
	AddEffect(effect.EffectType, ...float32) *effect.Effect
}

// ── row builder ───────────────────────────────────────────────────────────────

type sliderDef struct {
	label   string
	uniform string
	min     float64
	max     float64
	def     float64
}

func buildRow(
	name string,
	code string,
	mk sampleFn,
	kind effect.EffectType,
	initParams []float32,
	sliders []sliderDef,
) fyne.CanvasObject {

	beforeObj := mk()

	afterRaw := mk()
	eff := afterRaw.(addEffecter).AddEffect(kind, initParams...)

	nameLbl := canvas.NewText(name, clWhite)
	nameLbl.TextSize = 13
	nameLbl.TextStyle = fyne.TextStyle{Bold: true}

	codeLbl := canvas.NewText(code, color.NRGBA{100, 220, 100, 255})
	codeLbl.TextStyle = fyne.TextStyle{Monospace: true}
	codeLbl.TextSize = 11

	ctrlBox := container.NewVBox(codeLbl)
	for _, s := range sliders {
		s := s
		sl := widget.NewSlider(s.min, s.max)
		sl.Value = s.def
		sl.Step = (s.max - s.min) / 200.0
		valLbl := widget.NewLabel(fmt.Sprintf("%.2f", s.def))
		valLbl.TextStyle = fyne.TextStyle{Monospace: true}
		sl.OnChanged = func(v float64) {
			valLbl.SetText(fmt.Sprintf("%.2f", v))
			eff.SetFloat(s.uniform, float32(v))
			canvas.Refresh(afterRaw)
		}
		ctrlBox.Add(container.NewBorder(nil, nil,
			widget.NewLabel(s.label), valLbl, sl))
	}

	beforeLbl := widget.NewLabel("Before")
	beforeLbl.Alignment = fyne.TextAlignCenter
	afterLbl := widget.NewLabel("After")
	afterLbl.Alignment = fyne.TextAlignCenter
	arrow := canvas.NewText("▶", color.NRGBA{80, 180, 255, 220})
	arrow.TextSize = 28

	samples := container.NewHBox(
		container.NewVBox(beforeLbl, beforeObj),
		container.NewCenter(arrow),
		container.NewVBox(afterLbl, afterRaw),
	)

	body := container.NewBorder(nil, nil, samples, nil, ctrlBox)

	// Effects without sliders never get a Refresh call. Schedule one after the
	// window has had a chance to lay out and assign a canvas to the object.
	go func() {
		time.Sleep(50 * time.Millisecond)
		canvas.Refresh(afterRaw)
	}()

	return container.NewVBox(
		container.NewPadded(nameLbl),
		container.NewPadded(body),
		widget.NewSeparator(),
	)
}

// ── animated row builder ──────────────────────────────────────────────────────

func buildAnimRow(
	name string,
	code string,
	mk sampleFn,
	kind effect.EffectType,
	initParams []float32,
	sliders []sliderDef,
	animate func(eff *effect.Effect, obj fyne.CanvasObject),
) fyne.CanvasObject {

	staticObj := mk() // before: frozen at time=0

	liveObj := mk()
	eff := liveObj.(addEffecter).AddEffect(kind, initParams...)
	animate(eff, liveObj) // auto-start

	nameLbl := canvas.NewText(name, clWhite)
	nameLbl.TextSize = 13
	nameLbl.TextStyle = fyne.TextStyle{Bold: true}

	codeLbl := canvas.NewText(code, color.NRGBA{100, 220, 100, 255})
	codeLbl.TextStyle = fyne.TextStyle{Monospace: true}
	codeLbl.TextSize = 11

	ctrlBox := container.NewVBox(codeLbl)
	for _, s := range sliders {
		s := s
		sl := widget.NewSlider(s.min, s.max)
		sl.Value = s.def
		sl.Step = (s.max - s.min) / 200.0
		valLbl := widget.NewLabel(fmt.Sprintf("%.2f", s.def))
		valLbl.TextStyle = fyne.TextStyle{Monospace: true}
		sl.OnChanged = func(v float64) {
			valLbl.SetText(fmt.Sprintf("%.2f", v))
			eff.SetFloat(s.uniform, float32(v))
			canvas.Refresh(liveObj)
		}
		ctrlBox.Add(container.NewBorder(nil, nil,
			widget.NewLabel(s.label), valLbl, sl))
	}

	staticLbl := widget.NewLabel("Static (t=0)")
	staticLbl.Alignment = fyne.TextAlignCenter
	liveLbl := widget.NewLabel("Animated ▶")
	liveLbl.Alignment = fyne.TextAlignCenter
	arrow := canvas.NewText("▶", color.NRGBA{80, 180, 255, 220})
	arrow.TextSize = 28

	samples := container.NewHBox(
		container.NewVBox(staticLbl, staticObj),
		container.NewCenter(arrow),
		container.NewVBox(liveLbl, liveObj),
	)

	body := container.NewBorder(nil, nil, samples, nil, ctrlBox)
	return container.NewVBox(
		container.NewPadded(nameLbl),
		container.NewPadded(body),
		widget.NewSeparator(),
	)
}

// ── tabs ──────────────────────────────────────────────────────────────────────

func colorReplaceTab() fyne.CanvasObject {
	mk := solid(clRed)
	rows := container.NewVBox(
		buildRow("Color Replace — Red → Blue",
			`rect.AddEffect(effect.ColorReplace, 0.86,0.24,0.24, 0.15,0.39,0.92, 0.3)`,
			mk, effect.ColorReplace,
			[]float32{0.86, 0.24, 0.24, 0.15, 0.39, 0.92, 0.3},
			[]sliderDef{{"Tolerance", "tolerance", 0, 1, 0.3}},
		),
		buildRow("Color Replace — Red → Green",
			`rect.AddEffect(effect.ColorReplace, 0.86,0.24,0.24, 0.16,0.71,0.31, 0.4)`,
			mk, effect.ColorReplace,
			[]float32{0.86, 0.24, 0.24, 0.16, 0.71, 0.31, 0.4},
			[]sliderDef{{"Tolerance", "tolerance", 0, 1, 0.4}},
		),
		buildRow("Color Replace — Red → Gold",
			`rect.AddEffect(effect.ColorReplace, 0.86,0.24,0.24, 1.0,0.84,0.0, 0.3)`,
			solid(clRed), effect.ColorReplace,
			[]float32{0.86, 0.24, 0.24, 1.0, 0.84, 0.0, 0.3},
			[]sliderDef{{"Tolerance", "tolerance", 0, 1, 0.3}},
		),
	)
	return container.NewScroll(rows)
}

func duotoneTab() fyne.CanvasObject {
	mk := vGrad(clDark, clWhite)
	rows := container.NewVBox(
		buildRow("Duotone — Cyan / Magenta",
			`rect.AddEffect(effect.Duotone, 0,0.8,0.8,  0.8,0,0.8)`,
			mk, effect.Duotone, []float32{0, 0.8, 0.8, 0.8, 0, 0.8}, nil),
		buildRow("Duotone — Navy / Gold",
			`rect.AddEffect(effect.Duotone, 0.05,0.05,0.2,  1,0.84,0)`,
			mk, effect.Duotone, []float32{0.05, 0.05, 0.2, 1, 0.84, 0}, nil),
		buildRow("Duotone — Black / White",
			`rect.AddEffect(effect.Duotone, 0,0,0,  1,1,1)`,
			mk, effect.Duotone, []float32{0, 0, 0, 1, 1, 1}, nil),
		buildRow("Duotone — Purple / Pink",
			`rect.AddEffect(effect.Duotone, 0.3,0,0.5,  1,0.4,0.7)`,
			mk, effect.Duotone, []float32{0.3, 0, 0.5, 1, 0.4, 0.7}, nil),
		buildRow("Duotone — Teal / Orange",
			`rect.AddEffect(effect.Duotone, 0,0.5,0.4,  1,0.5,0)`,
			mk, effect.Duotone, []float32{0, 0.5, 0.4, 1, 0.5, 0}, nil),
	)
	return container.NewScroll(rows)
}

func splitToneTab() fyne.CanvasObject {
	mk := vGrad(clDark, clWhite)
	rows := container.NewVBox(
		buildRow("Split Tone — Cool shadows / Warm highlights",
			`rect.AddEffect(effect.SplitTone, 0,0.2,0.8, 0.8,0.6,0, 0.5)`,
			mk, effect.SplitTone,
			[]float32{0, 0.2, 0.8, 0.8, 0.6, 0, 0.5},
			[]sliderDef{{"Balance", "balance", 0, 1, 0.5}},
		),
		buildRow("Split Tone — Teal shadows / Orange highlights",
			`rect.AddEffect(effect.SplitTone, 0,0.6,0.5, 1,0.5,0, 0.4)`,
			mk, effect.SplitTone,
			[]float32{0, 0.6, 0.5, 1, 0.5, 0, 0.4},
			[]sliderDef{{"Balance", "balance", 0, 1, 0.4}},
		),
		buildRow("Split Tone — Blue shadows / Gold highlights",
			`rect.AddEffect(effect.SplitTone, 0,0,0.8, 0.8,0.7,0, 0.5)`,
			mk, effect.SplitTone,
			[]float32{0, 0, 0.8, 0.8, 0.7, 0, 0.5},
			[]sliderDef{{"Balance", "balance", 0, 1, 0.5}},
		),
	)
	return container.NewScroll(rows)
}

func channelMixerTab() fyne.CanvasObject {
	mk := hGrad(clRed, clBlue)
	rows := container.NewVBox(
		buildRow("Channel Mixer — R↔B Swap",
			`rect.AddEffect(effect.ChannelMixer, 0,0,1, 0,1,0, 1,0,0)`,
			mk, effect.ChannelMixer, []float32{0, 0, 1, 0, 1, 0, 1, 0, 0}, nil),
		buildRow("Channel Mixer — Green Only",
			`rect.AddEffect(effect.ChannelMixer, 0,1,0, 0,1,0, 0,1,0)`,
			hGrad(clRed, clGreen), effect.ChannelMixer, []float32{0, 1, 0, 0, 1, 0, 0, 1, 0}, nil),
		buildRow("Channel Mixer — Invert Red",
			`rect.AddEffect(effect.ChannelMixer, -1,0,0, 0,1,0, 0,0,1)`,
			hGrad(clRed, clOrange), effect.ChannelMixer, []float32{-1, 0, 0, 0, 1, 0, 0, 0, 1}, nil),
		buildRow("Channel Mixer — Luminance (grayscale)",
			`rect.AddEffect(effect.ChannelMixer, 0.3,0.6,0.1, 0.3,0.6,0.1, 0.3,0.6,0.1)`,
			hGrad(clRed, clBlue), effect.ChannelMixer, []float32{0.3, 0.6, 0.1, 0.3, 0.6, 0.1, 0.3, 0.6, 0.1}, nil),
		buildRow("Channel Mixer — Warm Shift",
			`rect.AddEffect(effect.ChannelMixer, 1.2,0.1,0, 0,1,0, 0,0,0.8)`,
			hGrad(clCyan, clBlue), effect.ChannelMixer, []float32{1.2, 0.1, 0, 0, 1, 0, 0, 0, 0.8}, nil),
	)
	return container.NewScroll(rows)
}

func gradientMapTab() fyne.CanvasObject {
	mk := vGrad(clDark, clWhite)
	rows := container.NewVBox(
		buildRow("Gradient Map — Thermal (black→blue→green→yellow→red)",
			`rect.AddEffect(effect.GradientMap, 0,0,0.2, 0,0,0.8, 0,0.8,0, 0.8,0.8,0, 0.8,0,0)`,
			mk, effect.GradientMap,
			[]float32{0, 0, 0.2, 0, 0, 0.8, 0, 0.8, 0, 0.8, 0.8, 0, 0.8, 0, 0}, nil),
		buildRow("Gradient Map — Sunset (dark purple→pink→orange→gold→white)",
			`rect.AddEffect(effect.GradientMap, 0.1,0,0.2, 0.5,0,0.3, 0.9,0.3,0, 1,0.8,0.2, 1,1,0.5)`,
			mk, effect.GradientMap,
			[]float32{0.1, 0, 0.2, 0.5, 0, 0.3, 0.9, 0.3, 0, 1, 0.8, 0.2, 1, 1, 0.5}, nil),
		buildRow("Gradient Map — Neon (black→blue→green→yellow→magenta)",
			`rect.AddEffect(effect.GradientMap, 0,0,0, 0,0,1, 0,1,0, 1,1,0, 1,0,1)`,
			mk, effect.GradientMap,
			[]float32{0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1}, nil),
		buildRow("Gradient Map — Desert (dark brown→sand→gold→cream)",
			`rect.AddEffect(effect.GradientMap, 0.1,0.05,0, 0.4,0.2,0, 0.8,0.5,0.1, 1,0.8,0.4, 1,1,0.8)`,
			mk, effect.GradientMap,
			[]float32{0.1, 0.05, 0, 0.4, 0.2, 0, 0.8, 0.5, 0.1, 1, 0.8, 0.4, 1, 1, 0.8}, nil),
		buildRow("Gradient Map — Ice (dark navy→teal→sky→pale→white)",
			`rect.AddEffect(effect.GradientMap, 0,0,0.1, 0,0.2,0.5, 0.2,0.6,0.9, 0.7,0.9,1, 1,1,1)`,
			mk, effect.GradientMap,
			[]float32{0, 0, 0.1, 0, 0.2, 0.5, 0.2, 0.6, 0.9, 0.7, 0.9, 1, 1, 1, 1}, nil),
	)
	return container.NewScroll(rows)
}

func patternTab() fyne.CanvasObject {
	mk := solid(clBlue)
	rows := container.NewVBox(
		buildRow("Pattern Overlay — Checkerboard",
			`rect.AddEffect(effect.PatternOverlay, 10, 0, 0.3)`,
			mk, effect.PatternOverlay,
			[]float32{10, 0, 0.3},
			[]sliderDef{
				{"Size", "patternSize", 4, 40, 10},
				{"Opacity", "opacity", 0, 1, 0.3},
			},
		),
		buildRow("Pattern Overlay — Horizontal Stripes",
			`rect.AddEffect(effect.PatternOverlay, 8, 1, 0.3)`,
			solid(clGreen), effect.PatternOverlay,
			[]float32{8, 1, 0.3},
			[]sliderDef{
				{"Size", "patternSize", 2, 40, 8},
				{"Opacity", "opacity", 0, 1, 0.3},
			},
		),
		buildRow("Pattern Overlay — Vertical Stripes",
			`rect.AddEffect(effect.PatternOverlay, 8, 2, 0.3)`,
			solid(clPurple), effect.PatternOverlay,
			[]float32{8, 2, 0.3},
			[]sliderDef{
				{"Size", "patternSize", 2, 40, 8},
				{"Opacity", "opacity", 0, 1, 0.3},
			},
		),
		buildRow("Pattern Overlay — Dots",
			`rect.AddEffect(effect.PatternOverlay, 12, 3, 0.3)`,
			solid(clOrange), effect.PatternOverlay,
			[]float32{12, 3, 0.3},
			[]sliderDef{
				{"Size", "patternSize", 4, 40, 12},
				{"Opacity", "opacity", 0, 1, 0.3},
			},
		),
	)
	return container.NewScroll(rows)
}

func noiseTab() fyne.CanvasObject {
	mk := hGrad(clCyan, clPink)
	rows := container.NewVBox(
		buildRow("Noise Displacement — Perlin UV warp",
			`rect.AddEffect(effect.NoiseDisplacement, 5, 8, 0)`,
			mk, effect.NoiseDisplacement,
			[]float32{5, 8, 0},
			[]sliderDef{
				{"Amount", "amount", 0, 20, 5},
				{"Scale", "scale", 1, 30, 8},
			},
		),
		buildRow("Noise Displacement — Fine grain",
			`rect.AddEffect(effect.NoiseDisplacement, 4, 20, 42)`,
			mk, effect.NoiseDisplacement,
			[]float32{4, 20, 42},
			[]sliderDef{
				{"Amount", "amount", 0, 20, 4},
				{"Scale", "scale", 1, 40, 20},
			},
		),
		buildRow("Noise Displacement — Heavy warp",
			`rect.AddEffect(effect.NoiseDisplacement, 12, 5, 7)`,
			hGrad(clPurple, clOrange), effect.NoiseDisplacement,
			[]float32{12, 5, 7},
			[]sliderDef{
				{"Amount", "amount", 0, 30, 12},
				{"Scale", "scale", 1, 20, 5},
			},
		),
	)
	return container.NewScroll(rows)
}

func animatedTab() fyne.CanvasObject {
	startTime := time.Now()
	elapsed := func() float32 { return float32(time.Since(startTime).Seconds()) }

	driveTime := func(eff *effect.Effect, obj fyne.CanvasObject) {
		a := fyne.NewAnimation(200*time.Second, func(_ float32) {
			eff.SetFloat("time", elapsed())
			canvas.Refresh(obj)
		})
		a.RepeatCount = fyne.AnimationRepeatForever
		a.Start()
	}

	rows := container.NewVBox(
		buildAnimRow(
			"Shimmer — sweeping highlight band",
			`rect.AddEffect(effect.Shimmer, 0, 0.08, 30, 0.6)`,
			solid(clBlue),
			effect.Shimmer,
			[]float32{0, 0.08, 30, 0.6},
			[]sliderDef{
				{"Width", "width", 0.02, 0.3, 0.08},
				{"Intensity", "intensity", 0, 1, 0.6},
				{"Angle °", "angle", 0, 90, 30},
			},
			driveTime,
		),
		buildAnimRow(
			"Pulse — breathing brightness/scale",
			`rect.AddEffect(effect.Pulse, 0, 1.0, 0.7, 1.3, 0.9, 1.1)`,
			solid(clPurple),
			effect.Pulse,
			[]float32{0, 1.0, 0.7, 1.3, 0.9, 1.1},
			[]sliderDef{
				{"Speed", "speed", 0.1, 4, 1.0},
				{"Bright Min", "brightMin", 0.3, 1, 0.7},
				{"Bright Max", "brightMax", 1, 2, 1.3},
			},
			driveTime,
		),
		buildAnimRow(
			"Glitch — block corruption",
			`rect.AddEffect(effect.Glitch, 0, 3, 20)`,
			hGrad(clRed, clCyan),
			effect.Glitch,
			[]float32{0, 3, 20},
			[]sliderDef{
				{"Amount", "amount", 0, 10, 3},
				{"Block Size", "blockSize", 4, 60, 20},
			},
			driveTime,
		),
		buildAnimRow(
			"Matrix Rain — falling characters",
			`rect.AddEffect(effect.MatrixRain, 0, 8, 0.5, 0.8)`,
			solid(clDark),
			effect.MatrixRain,
			[]float32{0, 8, 0.5, 0.8},
			[]sliderDef{
				{"Density", "density", 1, 20, 8},
				{"Speed", "speed", 0.1, 2, 0.5},
				{"Opacity", "opacity", 0, 1, 0.8},
			},
			driveTime,
		),
	)
	return container.NewScroll(rows)
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	a := app.New()
	w := a.NewWindow("Compositing & Animated Effects — Interactive Showcase")
	w.Resize(fyne.NewSize(1100, 820))

	title := canvas.NewText("Compositing · Procedural · Animated Effects", clWhite)
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText(
		"Left = original · Right = with effect · Sliders update live · Animated tab auto-plays",
		color.NRGBA{160, 160, 160, 255},
	)
	subtitle.TextSize = 11

	header := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		widget.NewSeparator(),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Color Replace",   colorReplaceTab()),
		container.NewTabItem("Duotone",         duotoneTab()),
		container.NewTabItem("Split Tone",      splitToneTab()),
		container.NewTabItem("Channel Mixer",   channelMixerTab()),
		container.NewTabItem("Gradient Map",    gradientMapTab()),
		container.NewTabItem("Pattern Overlay", patternTab()),
		container.NewTabItem("Noise Displace",  noiseTab()),
		container.NewTabItem("Animated ▶",      animatedTab()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(container.NewBorder(header, nil, nil, nil, tabs))
	w.ShowAndRun()
}
