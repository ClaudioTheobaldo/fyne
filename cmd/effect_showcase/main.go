package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	sW float32 = 190
	sH float32 = 110
)

// ctrl is one interactive slider for an effect parameter.
type ctrl struct {
	label     string
	uniform   string
	uniformSz int // 1 = SetFloat, 2 = SetVec2 (symmetric, both components equal)
	min, max  float64
	def       float64
}

// effectEntry describes one effect row in a tab.
type effectEntry struct {
	name       string
	kind       effect.EffectType
	initParams []float32
	ctrls      []ctrl
	code       string
}

// sampleFn creates a fresh canvas object used as the before/after sample.
type sampleFn func() fyne.CanvasObject

// canvasAdder is implemented by any canvas object that supports AddEffect.
type canvasAdder interface {
	fyne.CanvasObject
	AddEffect(effect.EffectType, ...float32) *effect.Effect
}

// ── sample factory helpers ────────────────────────────────────────────────────

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
func dGrad(a, b color.Color, angle float64) sampleFn {
	return func() fyne.CanvasObject {
		g := canvas.NewLinearGradient(a, b, angle)
		g.SetMinSize(fyne.NewSize(sW, sH))
		return g
	}
}
func rGrad(a, b color.Color) sampleFn {
	return func() fyne.CanvasObject {
		g := canvas.NewRadialGradient(a, b)
		g.SetMinSize(fyne.NewSize(sW, sH))
		return g
	}
}
// ── colors ───────────────────────────────────────────────────────────────────

var (
	clRed    = color.NRGBA{220, 60, 60, 255}
	clOrange = color.NRGBA{240, 140, 40, 255}
	clYellow = color.NRGBA{240, 220, 60, 255}
	clCyan   = color.NRGBA{60, 200, 220, 255}
	clBlue   = color.NRGBA{60, 100, 240, 255}
	clPurple = color.NRGBA{160, 60, 220, 255}
	clPink   = color.NRGBA{240, 100, 180, 255}
	clWhite  = color.NRGBA{240, 240, 240, 255}
	clGray   = color.NRGBA{140, 140, 140, 255}
	clDark   = color.NRGBA{30, 30, 30, 255}
	clTeal   = color.NRGBA{40, 200, 180, 255}
)

// ── row builder ──────────────────────────────────────────────────────────────

func buildRow(e effectEntry, mk sampleFn) fyne.CanvasObject {
	// Before: clean sample, no effect
	beforeObj := mk()

	// After: same sample + effect applied
	afterRaw := mk()
	ae := afterRaw.(canvasAdder)
	eff := ae.AddEffect(e.kind, e.initParams...)

	// Effect name
	nameLbl := canvas.NewText(e.name, clWhite)
	nameLbl.TextSize = 13
	nameLbl.TextStyle = fyne.TextStyle{Bold: true}

	// Code snippet
	codeLbl := canvas.NewText(e.code, color.NRGBA{120, 230, 120, 255})
	codeLbl.TextStyle = fyne.TextStyle{Monospace: true}
	codeLbl.TextSize = 11

	// Sliders
	ctrlsBox := container.NewVBox(codeLbl)
	for _, c := range e.ctrls {
		c := c // capture loop var
		s := widget.NewSlider(c.min, c.max)
		s.Value = c.def
		s.Step = (c.max - c.min) / 200.0

		valLbl := widget.NewLabel(fmt.Sprintf("%.2f", c.def))
		valLbl.TextStyle = fyne.TextStyle{Monospace: true}

		s.OnChanged = func(v float64) {
			valLbl.SetText(fmt.Sprintf("%.2f", v))
			f := float32(v)
			if c.uniformSz == 2 {
				eff.SetVec2(c.uniform, f, f)
			} else {
				eff.SetFloat(c.uniform, f)
			}
			canvas.Refresh(afterRaw)
		}

		ctrlsBox.Add(container.NewBorder(nil, nil,
			widget.NewLabel(c.label),
			valLbl,
			s,
		))
	}

	// Before label
	beforeLbl := widget.NewLabel("Before")
	beforeLbl.Alignment = fyne.TextAlignCenter

	// After label
	afterLbl := widget.NewLabel("After")
	afterLbl.Alignment = fyne.TextAlignCenter

	arrow := canvas.NewText("▶", color.NRGBA{80, 180, 255, 220})
	arrow.TextSize = 28

	samplesCol := func(obj fyne.CanvasObject, lbl *widget.Label) fyne.CanvasObject {
		return container.NewVBox(lbl, obj)
	}

	samplesRow := container.NewHBox(
		samplesCol(beforeObj, beforeLbl),
		container.NewCenter(arrow),
		samplesCol(afterRaw, afterLbl),
	)

	body := container.NewBorder(nil, nil, samplesRow, nil, ctrlsBox)

	return container.NewVBox(
		container.NewPadded(nameLbl),
		container.NewPadded(body),
		widget.NewSeparator(),
	)
}

func buildTab(entries []effectEntry, mk sampleFn) fyne.CanvasObject {
	vbox := container.NewVBox()
	for _, e := range entries {
		vbox.Add(buildRow(e, mk))
	}
	return container.NewScroll(vbox)
}

// ── effect definitions per tab ───────────────────────────────────────────────

func colorTab() ([]effectEntry, sampleFn) {
	mk := hGrad(clRed, clBlue)
	entries := []effectEntry{
		{
			name: "Brightness", kind: effect.Brightness, initParams: []float32{1.5},
			ctrls: []ctrl{{"Brightness", "brightness", 1, 0.1, 3.0, 1.5}},
			code:  "rect.AddEffect(effect.Brightness, 1.5)",
		},
		{
			name: "Contrast", kind: effect.Contrast, initParams: []float32{2.0},
			ctrls: []ctrl{{"Contrast", "contrast", 1, 0.0, 3.0, 2.0}},
			code:  "rect.AddEffect(effect.Contrast, 2.0)",
		},
		{
			name: "Saturate", kind: effect.Saturate, initParams: []float32{2.0},
			ctrls: []ctrl{{"Saturation", "saturation", 1, 0.0, 3.0, 2.0}},
			code:  "rect.AddEffect(effect.Saturate, 2.0)",
		},
		{
			name: "Grayscale", kind: effect.Grayscale, initParams: []float32{1.0},
			ctrls: []ctrl{{"Amount", "amount", 1, 0.0, 1.0, 1.0}},
			code:  "rect.AddEffect(effect.Grayscale, 1.0)",
		},
		{
			name: "Hue Rotate", kind: effect.HueRotate, initParams: []float32{90},
			ctrls: []ctrl{{"Angle °", "angle", 1, 0, 360, 90}},
			code:  "rect.AddEffect(effect.HueRotate, 90.0)",
		},
		{
			name: "Sepia", kind: effect.Sepia, initParams: []float32{1.0},
			ctrls: []ctrl{{"Amount", "amount", 1, 0.0, 1.0, 1.0}},
			code:  "rect.AddEffect(effect.Sepia, 1.0)",
		},
		{
			name: "Invert", kind: effect.Invert, initParams: []float32{1.0},
			ctrls: []ctrl{{"Amount", "amount", 1, 0.0, 1.0, 1.0}},
			code:  "rect.AddEffect(effect.Invert, 1.0)",
		},
		{
			name: "Opacity", kind: effect.Opacity, initParams: []float32{0.5},
			ctrls: []ctrl{{"Opacity", "opacity", 1, 0.0, 1.0, 0.5}},
			code:  "rect.AddEffect(effect.Opacity, 0.5)",
		},
		{
			name: "Posterize", kind: effect.Posterize, initParams: []float32{4},
			ctrls: []ctrl{{"Levels", "levels", 1, 2, 16, 4}},
			code:  "rect.AddEffect(effect.Posterize, 4.0)",
		},
		{
			name: "Threshold", kind: effect.Threshold, initParams: []float32{0.5},
			ctrls: []ctrl{{"Threshold", "threshold", 1, 0.0, 1.0, 0.5}},
			code:  "rect.AddEffect(effect.Threshold, 0.5)",
		},
		{
			name: "Gamma", kind: effect.Gamma, initParams: []float32{0.5},
			ctrls: []ctrl{{"Gamma", "gamma", 1, 0.1, 3.0, 0.5}},
			code:  "rect.AddEffect(effect.Gamma, 0.5)",
		},
		{
			name: "Vibrance", kind: effect.Vibrance, initParams: []float32{1.5},
			ctrls: []ctrl{{"Vibrance", "vibrance", 1, 0.0, 3.0, 1.5}},
			code:  "rect.AddEffect(effect.Vibrance, 1.5)",
		},
		{
			name: "Temperature", kind: effect.Temperature, initParams: []float32{3.0},
			ctrls: []ctrl{{"Temperature", "temperature", 1, -5.0, 5.0, 3.0}},
			code:  "rect.AddEffect(effect.Temperature, 3.0)",
		},
	}
	return entries, mk
}

func blurTab() ([]effectEntry, sampleFn) {
	mk := dGrad(clOrange, clTeal, 45)
	entries := []effectEntry{
		{
			name: "Gaussian Blur", kind: effect.GaussianBlur, initParams: []float32{5.0},
			ctrls: []ctrl{{"Radius", "radius", 1, 0, 20, 5}},
			code:  "rect.AddEffect(effect.GaussianBlur, 5.0)",
		},
		{
			name: "Box Blur", kind: effect.BoxBlur, initParams: []float32{6.0},
			ctrls: []ctrl{{"Radius", "radius", 1, 0, 20, 6}},
			code:  "rect.AddEffect(effect.BoxBlur, 6.0)",
		},
		{
			name: "Directional Blur (horizontal)", kind: effect.DirectionalBlur, initParams: []float32{8.0, 1.0, 0.0},
			ctrls: []ctrl{{"Radius", "radius", 1, 0, 20, 8}},
			code:  "rect.AddEffect(effect.DirectionalBlur, 8.0, 1.0, 0.0)",
		},
		{
			name: "Directional Blur (vertical)", kind: effect.DirectionalBlur, initParams: []float32{8.0, 0.0, 1.0},
			ctrls: []ctrl{{"Radius", "radius", 1, 0, 20, 8}},
			code:  "rect.AddEffect(effect.DirectionalBlur, 8.0, 0.0, 1.0)",
		},
		{
			name: "Zoom Blur", kind: effect.ZoomBlur, initParams: []float32{5.0, 0.5, 0.5},
			ctrls: []ctrl{{"Strength", "strength", 1, 0, 20, 5}},
			code:  "rect.AddEffect(effect.ZoomBlur, 5.0, 0.5, 0.5)",
		},
		{
			name: "Bilateral Blur", kind: effect.BilateralBlur, initParams: []float32{4.0, 3.0, 0.2},
			ctrls: []ctrl{
				{"Radius", "radius", 1, 0, 15, 4},
				{"Sigma Space", "sigmaSpace", 1, 0, 10, 3},
				{"Sigma Color", "sigmaColor", 1, 0, 1, 0.2},
			},
			code: "rect.AddEffect(effect.BilateralBlur, 4.0, 3.0, 0.2)",
		},
	}
	return entries, mk
}

func distortionTab() ([]effectEntry, sampleFn) {
	mk := hGrad(clPurple, clYellow)
	entries := []effectEntry{
		{
			name: "Pixelate", kind: effect.Pixelate, initParams: []float32{8.0},
			ctrls: []ctrl{{"Pixel Size", "pixelSize", 1, 2, 40, 8}},
			code:  "rect.AddEffect(effect.Pixelate, 8.0)",
		},
		{
			name: "Ripple", kind: effect.Ripple, initParams: []float32{4.0, 20.0, 0.0},
			ctrls: []ctrl{
				{"Amplitude", "amplitude", 1, 0, 10, 4},
				{"Frequency", "frequency", 1, 1, 50, 20},
			},
			code: "rect.AddEffect(effect.Ripple, 4.0, 20.0, 0.0)",
		},
		{
			name: "Swirl", kind: effect.Swirl, initParams: []float32{0.6, 5.0, 0.5, 0.5},
			ctrls: []ctrl{
				{"Radius", "radius", 1, 0, 1.5, 0.6},
				{"Angle", "angle", 1, 0, 12, 5},
			},
			code: "rect.AddEffect(effect.Swirl, 0.6, 5.0, 0.5, 0.5)",
		},
		{
			name: "Barrel Distortion", kind: effect.Barrel, initParams: []float32{0.5},
			ctrls: []ctrl{{"Distortion", "distortion", 1, -1, 1, 0.5}},
			code:  "rect.AddEffect(effect.Barrel, 0.5)",
		},
		{
			name: "Spherize", kind: effect.Spherize, initParams: []float32{0.6, 0.5, 0.5},
			ctrls: []ctrl{{"Radius", "radius", 1, 0, 1.5, 0.6}},
			code:  "rect.AddEffect(effect.Spherize, 0.6, 0.5, 0.5)",
		},
		{
			name: "Fisheye", kind: effect.Fisheye, initParams: []float32{0.7},
			ctrls: []ctrl{{"Strength", "strength", 1, 0, 2, 0.7}},
			code:  "rect.AddEffect(effect.Fisheye, 0.7)",
		},
		{
			name: "Chromatic Aberration", kind: effect.ChromaticAberration, initParams: []float32{7.0},
			ctrls: []ctrl{{"Offset px", "offset", 1, 0, 20, 7}},
			code:  "rect.AddEffect(effect.ChromaticAberration, 7.0)",
		},
		{
			name: "RGB Shift", kind: effect.RGBShift, initParams: []float32{6.0, 45.0},
			ctrls: []ctrl{
				{"Amount px", "amount", 1, 0, 20, 6},
				{"Angle °", "angle", 1, 0, 360, 45},
			},
			code: "rect.AddEffect(effect.RGBShift, 6.0, 45.0)",
		},
		{
			name: "Frosted Glass", kind: effect.Frosted, initParams: []float32{4.0, 2.5},
			ctrls: []ctrl{
				{"Blur Radius", "radius", 1, 0, 10, 4},
				{"Noise Scale", "noiseScale", 1, 0, 5, 2.5},
			},
			code: "rect.AddEffect(effect.Frosted, 4.0, 2.5)",
		},
		{
			name: "Tilt Shift", kind: effect.TiltShift, initParams: []float32{8.0, 0.5, 0.15},
			ctrls: []ctrl{
				{"Blur Radius", "radius", 1, 0, 20, 8},
				{"Focus Y", "focusY", 1, 0, 1, 0.5},
				{"Focus Width", "focusWidth", 1, 0.02, 0.5, 0.15},
			},
			code: "rect.AddEffect(effect.TiltShift, 8.0, 0.5, 0.15)",
		},
	}
	return entries, mk
}

func edgeTab() ([]effectEntry, sampleFn) {
	mk := dGrad(clDark, clWhite, 135)
	entries := []effectEntry{
		{
			name: "Sharpen", kind: effect.Sharpen, initParams: []float32{2.0},
			ctrls: []ctrl{{"Strength", "strength", 1, 0, 5, 2}},
			code:  "rect.AddEffect(effect.Sharpen, 2.0)",
		},
		{
			name: "Emboss", kind: effect.Emboss, initParams: []float32{2.5},
			ctrls: []ctrl{{"Strength", "strength", 1, 0, 5, 2.5}},
			code:  "rect.AddEffect(effect.Emboss, 2.5)",
		},
		{
			name: "Edge Detect", kind: effect.EdgeDetect, initParams: []float32{0.2},
			ctrls: []ctrl{{"Threshold", "threshold", 1, 0, 1, 0.2}},
			code:  "rect.AddEffect(effect.EdgeDetect, 0.2)",
		},
	}
	return entries, mk
}

func shadowTab() ([]effectEntry, sampleFn) {
	mk := vGrad(clBlue, clDark)
	entries := []effectEntry{
		{
			name: "Drop Shadow", kind: effect.DropShadow, initParams: []float32{4, 4, 6, 0, 0, 0, 0.8},
			ctrls: []ctrl{
				{"Offset X", "offsetX", 1, 0, 20, 4},
				{"Offset Y", "offsetY", 1, 0, 20, 4},
				{"Blur", "blur", 1, 0, 20, 6},
			},
			code: "rect.AddEffect(effect.DropShadow, 4, 4, 6,  0,0,0, 0.8)",
		},
		{
			name: "Box Shadow", kind: effect.BoxShadow, initParams: []float32{3, 3, 5, 2, 0, 0, 0, 0.7},
			ctrls: []ctrl{
				{"Offset X", "offsetX", 1, 0, 15, 3},
				{"Offset Y", "offsetY", 1, 0, 15, 3},
				{"Blur", "blur", 1, 0, 15, 5},
				{"Spread", "spread", 1, 0, 10, 2},
			},
			code: "rect.AddEffect(effect.BoxShadow, 3, 3, 5, 2,  0,0,0, 0.7)",
		},
		{
			name: "Inner Shadow", kind: effect.InnerShadow, initParams: []float32{2, 2, 4, 0, 0, 0, 0.6},
			ctrls: []ctrl{
				{"Offset X", "offsetX", 1, 0, 10, 2},
				{"Offset Y", "offsetY", 1, 0, 10, 2},
				{"Blur", "blur", 1, 0, 15, 4},
			},
			code: "rect.AddEffect(effect.InnerShadow, 2, 2, 4,  0,0,0, 0.6)",
		},
		{
			name: "Outer Glow", kind: effect.OuterGlow, initParams: []float32{12, 0, 1, 0.5, 1},
			ctrls: []ctrl{{"Radius", "radius", 1, 0, 30, 12}},
			code:  "rect.AddEffect(effect.OuterGlow, 12, 0,1,0.5, 1.0)",
		},
	}
	return entries, mk
}

func styleTab() ([]effectEntry, sampleFn) {
	mk := rGrad(clWhite, clDark)
	entries := []effectEntry{
		{
			name: "Vignette", kind: effect.Vignette, initParams: []float32{1.5, 0.5},
			ctrls: []ctrl{
				{"Intensity", "intensity", 1, 0, 3, 1.5},
				{"Smoothness", "smoothness", 1, 0, 1, 0.5},
			},
			code: "rect.AddEffect(effect.Vignette, 1.5, 0.5)",
		},
		{
			name: "Film Grain", kind: effect.FilmGrain, initParams: []float32{0.4, 0.0},
			ctrls: []ctrl{{"Intensity", "intensity", 1, 0, 1, 0.4}},
			code:  "rect.AddEffect(effect.FilmGrain, 0.4, 0.0)",
		},
		{
			name: "Scanlines", kind: effect.Scanlines, initParams: []float32{1.0, 0.5},
			ctrls: []ctrl{
				{"Density", "density", 1, 0, 5, 1},
				{"Opacity", "opacity", 1, 0, 1, 0.5},
			},
			code: "rect.AddEffect(effect.Scanlines, 1.0, 0.5)",
		},
		{
			name: "CRT", kind: effect.CRT, initParams: []float32{0.02, 0.4, 0.3},
			ctrls: []ctrl{
				{"Curvature", "curvature", 1, 0, 0.1, 0.02},
				{"Scanlines", "scanlineIntensity", 1, 0, 1, 0.4},
				{"Vignette", "vignetteStrength", 1, 0, 1, 0.3},
			},
			code: "rect.AddEffect(effect.CRT, 0.02, 0.4, 0.3)",
		},
		{
			name: "Halftone", kind: effect.Halftone, initParams: []float32{6.0, 0.0},
			ctrls: []ctrl{
				{"Dot Size", "dotSize", 1, 2, 20, 6},
				{"Angle °", "angle", 1, 0, 90, 0},
			},
			code: "rect.AddEffect(effect.Halftone, 6.0, 0.0)",
		},
		{
			name: "Dot Matrix", kind: effect.DotMatrix, initParams: []float32{5.0, 2.0},
			ctrls: []ctrl{
				{"Dot Size", "dotSize", 1, 2, 20, 5},
				{"Spacing", "spacing", 1, 1, 10, 2},
			},
			code: "rect.AddEffect(effect.DotMatrix, 5.0, 2.0)",
		},
		{
			name: "Oil Painting", kind: effect.OilPainting, initParams: []float32{3.0, 8.0},
			ctrls: []ctrl{
				{"Radius", "radius", 1, 1, 8, 3},
				{"Levels", "levels", 1, 2, 32, 8},
			},
			code: "rect.AddEffect(effect.OilPainting, 3.0, 8.0)",
		},
	}
	return entries, mk
}

func lightTab() ([]effectEntry, sampleFn) {
	mk := rGrad(clGray, clDark)
	entries := []effectEntry{
		{
			name: "Diffuse Light", kind: effect.DiffuseLight, initParams: []float32{0.5, 0.5, 1.0, 1.0, 1.0},
			ctrls: []ctrl{
				{"Diffuse", "diffuseConstant", 1, 0, 2, 1.0},
				{"Surface Scale", "surfaceScale", 1, 0, 5, 1.0},
			},
			code: "rect.AddEffect(effect.DiffuseLight, 0.5,0.5,1.0, 1.0, 1.0)",
		},
		{
			name: "Specular Light", kind: effect.SpecularLight, initParams: []float32{0.5, 0.5, 1.0, 0.8, 16.0, 1.0},
			ctrls: []ctrl{
				{"Specular", "specularConstant", 1, 0, 2, 0.8},
				{"Shininess", "specularExponent", 1, 1, 64, 16},
				{"Surface Scale", "surfaceScale", 1, 0, 5, 1.0},
			},
			code: "rect.AddEffect(effect.SpecularLight, 0.5,0.5,1.0, 0.8, 16.0, 1.0)",
		},
		{
			name: "Ambient Light", kind: effect.AmbientLight, initParams: []float32{1.0, 0.8, 0.5, 0.5},
			ctrls: []ctrl{{"Intensity", "intensity", 1, 0, 2, 0.5}},
			code:  "rect.AddEffect(effect.AmbientLight, 1.0, 0.8, 0.5, 0.5)",
		},
	}
	return entries, mk
}

func advancedTab() ([]effectEntry, sampleFn) {
	mk := dGrad(clPurple, clOrange, 60)
	entries := []effectEntry{
		{
			name: "Turbulence", kind: effect.Turbulence, initParams: []float32{5.0, 5.0, 4.0, 42.0},
			ctrls: []ctrl{
				{"Frequency", "baseFrequency", 2, 0, 20, 5}, // vec2, symmetric
				{"Octaves", "numOctaves", 1, 1, 8, 4},
			},
			code: "rect.AddEffect(effect.Turbulence, 5.0, 5.0, 4.0, 42.0)",
		},
		{
			name: "Bloom", kind: effect.Bloom, initParams: []float32{0.3, 1.5, 6.0},
			ctrls: []ctrl{
				{"Threshold", "threshold", 1, 0, 1, 0.3},
				{"Intensity", "intensity", 1, 0, 3, 1.5},
				{"Radius", "radius", 1, 0, 20, 6},
			},
			code: "rect.AddEffect(effect.Bloom, 0.3, 1.5, 6.0)",
		},
		{
			name: "Tone Mapping", kind: effect.ToneMapping, initParams: []float32{2.0, 2.2},
			ctrls: []ctrl{
				{"Exposure", "exposure", 1, 0, 5, 2.0},
				{"Gamma", "gamma", 1, 0.1, 5, 2.2},
			},
			code: "rect.AddEffect(effect.ToneMapping, 2.0, 2.2)",
		},
		{
			name:  "FXAA (anti-alias)",
			kind:  effect.FXAA,
			code:  "rect.AddEffect(effect.FXAA)",
		},
	}
	return entries, mk
}

func uvTab() ([]effectEntry, sampleFn) {
	mk := hGrad(clCyan, clPink)
	entries := []effectEntry{
		{
			name: "Skew", kind: effect.Skew, initParams: []float32{0.4, 0.0},
			ctrls: []ctrl{
				{"Skew X", "skewX", 1, -1, 1, 0.4},
				{"Skew Y", "skewY", 1, -1, 1, 0},
			},
			code: "rect.AddEffect(effect.Skew, 0.4, 0.0)",
		},
		{
			name: "Rotate", kind: effect.Rotate, initParams: []float32{0.3},
			ctrls: []ctrl{{"Angle rad", "angle", 1, -3.14, 3.14, 0.3}},
			code:  "rect.AddEffect(effect.Rotate, 0.3)",
		},
		{
			name: "Scale / Zoom", kind: effect.Scale, initParams: []float32{1.5, 1.5},
			ctrls: []ctrl{
				{"Scale X", "scaleX", 1, 0.1, 3, 1.5},
				{"Scale Y", "scaleY", 1, 0.1, 3, 1.5},
			},
			code: "rect.AddEffect(effect.Scale, 1.5, 1.5)",
		},
		{
			name: "Perspective", kind: effect.PerspectiveTransform, initParams: []float32{0.3, 0.0},
			ctrls: []ctrl{
				{"Amount", "amount", 1, -1, 1, 0.3},
				{"Direction", "direction", 1, 0, 1, 0},
			},
			code: "rect.AddEffect(effect.PerspectiveTransform, 0.3, 0.0)",
		},
		{
			name: "Wave / Flag", kind: effect.Wave, initParams: []float32{0.05, 8.0, 0.0, 0.0},
			ctrls: []ctrl{
				{"Amplitude", "amplitude", 1, 0, 0.2, 0.05},
				{"Frequency", "frequency", 1, 1, 20, 8},
			},
			code: "rect.AddEffect(effect.Wave, 0.05, 8.0, 0.0, 0.0)",
		},
		{
			name: "Page Curl", kind: effect.PageCurl, initParams: []float32{0.7, 0.1},
			ctrls: []ctrl{
				{"Curl", "curl", 1, 0, 1, 0.7},
				{"Radius", "radius", 1, 0.01, 0.3, 0.1},
			},
			code: "rect.AddEffect(effect.PageCurl, 0.7, 0.1)",
		},
		{
			name:       "Mirror (flip X)",
			kind:       effect.Mirror,
			initParams: []float32{1, 0},
			code:       "rect.AddEffect(effect.Mirror, 1.0, 0.0)",
		},
		{
			name:       "Mirror (flip Y)",
			kind:       effect.Mirror,
			initParams: []float32{0, 1},
			code:       "rect.AddEffect(effect.Mirror, 0.0, 1.0)",
		},
		{
			name: "Kaleidoscope", kind: effect.Kaleidoscope, initParams: []float32{6, 0},
			ctrls: []ctrl{
				{"Segments", "segments", 1, 2, 16, 6},
				{"Rotation", "rotation", 1, 0, 6.28, 0},
			},
			code: "rect.AddEffect(effect.Kaleidoscope, 6.0, 0.0)",
		},
	}
	return entries, mk
}

// ── main ─────────────────────────────────────────────────────────────────────

func main() {
	a := app.New()
	w := a.NewWindow("Fyne Effect System — Interactive Showcase")
	w.Resize(fyne.NewSize(1100, 820))

	colorE, colorMk := colorTab()
	blurE, blurMk := blurTab()
	distE, distMk := distortionTab()
	edgeE, edgeMk := edgeTab()
	shadowE, shadowMk := shadowTab()
	styleE, styleMk := styleTab()
	lightE, lightMk := lightTab()
	advE, advMk := advancedTab()
	uvE, uvMk := uvTab()

	tabs := container.NewAppTabs(
		container.NewTabItem("Color",       buildTab(colorE,  colorMk)),
		container.NewTabItem("Blur",        buildTab(blurE,   blurMk)),
		container.NewTabItem("Distortion",  buildTab(distE,   distMk)),
		container.NewTabItem("Edge",        buildTab(edgeE,   edgeMk)),
		container.NewTabItem("Shadow",      buildTab(shadowE, shadowMk)),
		container.NewTabItem("Stylization", buildTab(styleE,  styleMk)),
		container.NewTabItem("Lighting",    buildTab(lightE,  lightMk)),
		container.NewTabItem("Advanced",    buildTab(advE,    advMk)),
		container.NewTabItem("UV / Warp",   buildTab(uvE,     uvMk)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	title := canvas.NewText("Fyne Effect System — Interactive Showcase", clWhite)
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText(
		"Left = original · Right = with effect · Sliders update live · Copy the green code snippet",
		color.NRGBA{160, 160, 160, 255},
	)
	subtitle.TextSize = 11

	header := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		widget.NewSeparator(),
	)

	w.SetContent(container.NewBorder(header, nil, nil, nil, tabs))
	w.ShowAndRun()
}
