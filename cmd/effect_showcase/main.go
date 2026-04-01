package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// colors used throughout
var (
	red     = color.NRGBA{220, 60, 60, 255}
	orange  = color.NRGBA{240, 140, 40, 255}
	yellow  = color.NRGBA{240, 220, 60, 255}
	green   = color.NRGBA{60, 200, 80, 255}
	cyan    = color.NRGBA{60, 200, 220, 255}
	blue    = color.NRGBA{60, 100, 240, 255}
	purple  = color.NRGBA{160, 60, 220, 255}
	pink    = color.NRGBA{240, 100, 180, 255}
	white   = color.NRGBA{240, 240, 240, 255}
	gray    = color.NRGBA{140, 140, 140, 255}
	darkGray = color.NRGBA{60, 60, 60, 255}
)

func effectCard(label string, c color.Color, effects ...effectDef) *fyne.Container {
	rect := canvas.NewRectangle(c)
	rect.SetMinSize(fyne.NewSize(90, 70))
	for _, e := range effects {
		rect.AddEffect(e.kind, e.params...)
	}
	lbl := widget.NewLabel(label)
	lbl.Wrapping = fyne.TextWrapWord
	lbl.Alignment = fyne.TextAlignCenter
	return container.NewVBox(rect, lbl)
}

type effectDef struct {
	kind   effect.EffectType
	params []float32
}

func ef(kind effect.EffectType, params ...float32) effectDef {
	return effectDef{kind, params}
}

func sectionTitle(text string) *canvas.Text {
	t := canvas.NewText(text, color.NRGBA{80, 180, 255, 255})
	t.TextSize = 16
	t.TextStyle = fyne.TextStyle{Bold: true}
	return t
}

func main() {
	a := app.New()
	w := a.NewWindow("Fyne Effect Showcase — Complete")
	w.Resize(fyne.NewSize(1000, 800))

	// ==================== HEADER ====================
	title := canvas.NewText("Fyne Effect System — Complete Showcase", color.White)
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("56 built-in effects · FBO ping-pong pipeline · Stacking · Animation · Custom shaders", color.NRGBA{170, 170, 170, 255})
	subtitle.TextSize = 12

	apiLine := canvas.NewText(`blur := rect.AddEffect(effect.GaussianBlur, 4.0)`, color.NRGBA{140, 255, 140, 255})
	apiLine.TextSize = 12
	apiLine.TextStyle = fyne.TextStyle{Monospace: true}

	header := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		container.NewCenter(apiLine),
		widget.NewSeparator(),
	)

	// ==================== COLOR EFFECTS ====================
	colorSection := container.NewVBox(
		sectionTitle("Color Effects (14)"),
		container.NewGridWrap(fyne.NewSize(110, 110),
			effectCard("Brightness\n1.5", orange, ef(effect.Brightness, 1.5)),
			effectCard("Brightness\n0.5", orange, ef(effect.Brightness, 0.5)),
			effectCard("Contrast\n2.0", green, ef(effect.Contrast, 2.0)),
			effectCard("Contrast\n0.3", green, ef(effect.Contrast, 0.3)),
			effectCard("Grayscale", red, ef(effect.Grayscale, 1.0)),
			effectCard("Grayscale\n50%", red, ef(effect.Grayscale, 0.5)),
			effectCard("Sepia", blue, ef(effect.Sepia, 1.0)),
			effectCard("Invert", yellow, ef(effect.Invert, 1.0)),
			effectCard("Invert\n50%", yellow, ef(effect.Invert, 0.5)),
			effectCard("HueRotate\n90°", orange, ef(effect.HueRotate, 90)),
			effectCard("HueRotate\n180°", orange, ef(effect.HueRotate, 180)),
			effectCard("Saturate\n0.0", pink, ef(effect.Saturate, 0.0)),
			effectCard("Saturate\n3.0", pink, ef(effect.Saturate, 3.0)),
			effectCard("Opacity\n0.3", red, ef(effect.Opacity, 0.3)),
			effectCard("Posterize\n4", cyan, ef(effect.Posterize, 4)),
			effectCard("Threshold\n0.5", purple, ef(effect.Threshold, 0.5)),
			effectCard("Gamma\n0.5", green, ef(effect.Gamma, 0.5)),
			effectCard("Gamma\n2.0", green, ef(effect.Gamma, 2.0)),
			effectCard("Vibrance\n1.0", gray, ef(effect.Vibrance, 1.0)),
			effectCard("Temperature\n+3", white, ef(effect.Temperature, 3.0)),
			effectCard("Temperature\n-3", white, ef(effect.Temperature, -3.0)),
		),
		widget.NewSeparator(),
	)

	// ==================== BLUR EFFECTS ====================
	blurSection := container.NewVBox(
		sectionTitle("Blur Effects (5)"),
		container.NewGridWrap(fyne.NewSize(110, 110),
			effectCard("Gaussian\n2px", red, ef(effect.GaussianBlur, 2.0)),
			effectCard("Gaussian\n6px", red, ef(effect.GaussianBlur, 6.0)),
			effectCard("Gaussian\n12px", red, ef(effect.GaussianBlur, 12.0)),
			effectCard("Box Blur\n4", green, ef(effect.BoxBlur, 4.0)),
			effectCard("Directional\nH", blue, ef(effect.DirectionalBlur, 8.0, 1.0, 0.0)),
			effectCard("Directional\nV", blue, ef(effect.DirectionalBlur, 8.0, 0.0, 1.0)),
			effectCard("Zoom Blur", yellow, ef(effect.ZoomBlur, 5.0, 0.5, 0.5)),
			effectCard("Bilateral", orange, ef(effect.BilateralBlur, 4.0, 3.0, 0.2)),
		),
		widget.NewSeparator(),
	)

	// ==================== DISTORTION EFFECTS ====================
	distortionSection := container.NewVBox(
		sectionTitle("Distortion Effects (10)"),
		container.NewGridWrap(fyne.NewSize(110, 110),
			effectCard("Pixelate\n6", green, ef(effect.Pixelate, 6.0)),
			effectCard("Pixelate\n12", green, ef(effect.Pixelate, 12.0)),
			effectCard("Ripple", cyan, ef(effect.Ripple, 3.0, 20.0, 0.0)),
			effectCard("Swirl", yellow, ef(effect.Swirl, 0.4, 4.0, 0.5, 0.5)),
			effectCard("Barrel\n+0.5", blue, ef(effect.Barrel, 0.5)),
			effectCard("Barrel\n-0.3", blue, ef(effect.Barrel, -0.3)),
			effectCard("Spherize", purple, ef(effect.Spherize, 0.4, 0.5, 0.5)),
			effectCard("Fisheye", orange, ef(effect.Fisheye, 0.5)),
			effectCard("Chromatic\nAberration", red, ef(effect.ChromaticAberration, 5.0)),
			effectCard("RGB Shift", pink, ef(effect.RGBShift, 4.0, 45.0)),
			effectCard("Frosted", white, ef(effect.Frosted, 3.0, 2.0)),
			effectCard("Tilt Shift", green, ef(effect.TiltShift, 8.0, 0.5, 0.15)),
		),
		widget.NewSeparator(),
	)

	// ==================== DETAIL EFFECTS ====================
	detailSection := container.NewVBox(
		sectionTitle("Edge & Detail Effects (3)"),
		container.NewGridWrap(fyne.NewSize(110, 110),
			effectCard("Sharpen\n1.0", purple, ef(effect.Sharpen, 1.0)),
			effectCard("Sharpen\n3.0", purple, ef(effect.Sharpen, 3.0)),
			effectCard("Emboss\n1.0", cyan, ef(effect.Emboss, 1.0)),
			effectCard("Emboss\n3.0", cyan, ef(effect.Emboss, 3.0)),
			effectCard("Edge Detect\n0.1", orange, ef(effect.EdgeDetect, 0.1)),
			effectCard("Edge Detect\n0.5", orange, ef(effect.EdgeDetect, 0.5)),
		),
		widget.NewSeparator(),
	)

	// ==================== SHADOW EFFECTS ====================
	shadowSection := container.NewVBox(
		sectionTitle("Shadow Effects (4)"),
		container.NewGridWrap(fyne.NewSize(130, 120),
			effectCard("Drop Shadow", blue, ef(effect.DropShadow, 4, 4, 6, 0, 0, 0, 0.7)),
			effectCard("Box Shadow", green, ef(effect.BoxShadow, 3, 3, 5, 2, 0, 0, 0, 0.6)),
			effectCard("Inner Shadow", white, ef(effect.InnerShadow, 2, 2, 4, 0, 0, 0, 0.5)),
			effectCard("Outer Glow", darkGray, ef(effect.OuterGlow, 8, 1, 0.5, 0, 1)),
		),
		widget.NewSeparator(),
	)

	// ==================== STYLIZATION EFFECTS ====================
	styleSection := container.NewVBox(
		sectionTitle("Stylization Effects (7)"),
		container.NewGridWrap(fyne.NewSize(110, 110),
			effectCard("Vignette\nStrong", white, ef(effect.Vignette, 1.5, 0.4)),
			effectCard("Vignette\nSubtle", white, ef(effect.Vignette, 0.5, 0.6)),
			effectCard("Film Grain", blue, ef(effect.FilmGrain, 0.3, 0.0)),
			effectCard("Scanlines", green, ef(effect.Scanlines, 1.0, 0.4)),
			effectCard("CRT", orange, ef(effect.CRT, 0.02, 0.3, 0.3)),
			effectCard("Halftone", red, ef(effect.Halftone, 6.0, 0.0)),
			effectCard("Dot Matrix", cyan, ef(effect.DotMatrix, 5.0, 2.0)),
			effectCard("Oil Painting", purple, ef(effect.OilPainting, 3.0, 8.0)),
		),
		widget.NewSeparator(),
	)

	// ==================== LIGHTING EFFECTS ====================
	lightSection := container.NewVBox(
		sectionTitle("Lighting Effects (3)"),
		container.NewGridWrap(fyne.NewSize(130, 120),
			effectCard("Diffuse Light", gray, ef(effect.DiffuseLight, 0.5, 0.5, 1.0, 1.0, 0.5)),
			effectCard("Specular Light", gray, ef(effect.SpecularLight, 0.5, 0.5, 1.0, 0.8, 16.0, 0.5)),
			effectCard("Ambient Light", darkGray, ef(effect.AmbientLight, 1.0, 0.8, 0.5, 0.4)),
		),
		widget.NewSeparator(),
	)

	// ==================== GRADIENT & MASK EFFECTS ====================
	gradientSection := container.NewVBox(
		sectionTitle("Gradient Overlay & Mask (4)"),
		container.NewGridWrap(fyne.NewSize(130, 120),
			effectCard("Linear Grad", white, ef(effect.LinearGradientOverlay, 1, 0, 0, 0.7, 0, 0, 1, 0.7, 0)),
			effectCard("Radial Grad", white, ef(effect.RadialGradientOverlay, 1, 0.8, 0, 0.6, 0, 0.2, 1, 0.6, 0.5, 0.5)),
			effectCard("Conic Grad", white, ef(effect.ConicGradientOverlay, 1, 0, 0, 0.5, 0, 1, 0, 0.5, 0.5, 0.5, 0)),
			effectCard("Gradient Mask", red, ef(effect.GradientMask, 1.0, 0.0, 0)),
		),
		widget.NewSeparator(),
	)

	// ==================== PROCEDURAL & ADVANCED ====================
	advancedSection := container.NewVBox(
		sectionTitle("Procedural & Advanced (4)"),
		container.NewGridWrap(fyne.NewSize(130, 120),
			effectCard("Turbulence", blue, ef(effect.Turbulence, 5.0, 5.0, 4.0, 42.0)),
			effectCard("Bloom", white, ef(effect.Bloom, 0.3, 1.5, 6.0)),
			effectCard("Tone Map", orange, ef(effect.ToneMapping, 2.0, 2.2)),
			effectCard("FXAA", red, ef(effect.FXAA)),
		),
		widget.NewSeparator(),
	)

	// ==================== BLEND MODES ====================
	blendSection := container.NewVBox(
		sectionTitle("Blend Modes (2 shown)"),
		container.NewGridWrap(fyne.NewSize(130, 120),
			effectCard("Blend Multiply", orange, ef(effect.BlendMultiply, 1.0)),
			effectCard("Blend Screen", blue, ef(effect.BlendScreen, 1.0)),
		),
		widget.NewSeparator(),
	)

	// ==================== STACKED EFFECTS — CLEVER COMBOS ====================
	stackTitle := sectionTitle("Stacked Effects — Creative Combinations")

	// Frosted glass panel
	frostedGlass := canvas.NewRectangle(white)
	frostedGlass.SetMinSize(fyne.NewSize(140, 100))
	frostedGlass.AddEffect(effect.GaussianBlur, 4.0)
	frostedGlass.AddEffect(effect.Brightness, 1.1)
	frostedGlass.AddEffect(effect.Opacity, 0.85)

	// Vintage photo
	vintagePhoto := canvas.NewRectangle(color.NRGBA{180, 120, 80, 255})
	vintagePhoto.SetMinSize(fyne.NewSize(140, 100))
	vintagePhoto.AddEffect(effect.Sepia, 0.7)
	vintagePhoto.AddEffect(effect.Vignette, 1.2, 0.4)
	vintagePhoto.AddEffect(effect.FilmGrain, 0.15, 0.0)
	vintagePhoto.AddEffect(effect.Contrast, 1.2)

	// Disabled/inactive state
	disabledState := canvas.NewRectangle(green)
	disabledState.SetMinSize(fyne.NewSize(140, 100))
	disabledState.AddEffect(effect.Grayscale, 1.0)
	disabledState.AddEffect(effect.Opacity, 0.4)

	// Glitch/cyberpunk
	glitchRect := canvas.NewRectangle(red)
	glitchRect.SetMinSize(fyne.NewSize(140, 100))
	glitchRect.AddEffect(effect.ChromaticAberration, 6.0)
	glitchRect.AddEffect(effect.Scanlines, 1.0, 0.25)
	glitchRect.AddEffect(effect.Contrast, 1.3)

	// Dream/soft focus
	dreamRect := canvas.NewRectangle(purple)
	dreamRect.SetMinSize(fyne.NewSize(140, 100))
	dreamRect.AddEffect(effect.GaussianBlur, 2.0)
	dreamRect.AddEffect(effect.Brightness, 1.15)
	dreamRect.AddEffect(effect.Saturate, 1.4)

	// Embossed metal
	metalRect := canvas.NewRectangle(gray)
	metalRect.SetMinSize(fyne.NewSize(140, 100))
	metalRect.AddEffect(effect.Emboss, 2.0)
	metalRect.AddEffect(effect.Contrast, 1.5)
	metalRect.AddEffect(effect.Brightness, 1.2)

	// Neon glow
	neonRect := canvas.NewRectangle(darkGray)
	neonRect.SetMinSize(fyne.NewSize(140, 100))
	neonRect.AddEffect(effect.OuterGlow, 10, 0, 1, 0.5, 1)
	neonRect.AddEffect(effect.Brightness, 1.3)

	// Thermal vision
	thermalRect := canvas.NewRectangle(orange)
	thermalRect.SetMinSize(fyne.NewSize(140, 100))
	thermalRect.AddEffect(effect.Posterize, 5)
	thermalRect.AddEffect(effect.HueRotate, 120)
	thermalRect.AddEffect(effect.Saturate, 2.5)

	// Pencil sketch
	sketchRect := canvas.NewRectangle(color.NRGBA{180, 160, 140, 255})
	sketchRect.SetMinSize(fyne.NewSize(140, 100))
	sketchRect.AddEffect(effect.EdgeDetect, 0.15)
	sketchRect.AddEffect(effect.Invert, 1.0)
	sketchRect.AddEffect(effect.Grayscale, 1.0)

	// Old TV
	oldTVRect := canvas.NewRectangle(green)
	oldTVRect.SetMinSize(fyne.NewSize(140, 100))
	oldTVRect.AddEffect(effect.CRT, 0.03, 0.4, 0.4)
	oldTVRect.AddEffect(effect.Grayscale, 0.8)
	oldTVRect.AddEffect(effect.FilmGrain, 0.2, 0.0)

	// Watercolor
	watercolorRect := canvas.NewRectangle(cyan)
	watercolorRect.SetMinSize(fyne.NewSize(140, 100))
	watercolorRect.AddEffect(effect.GaussianBlur, 1.5)
	watercolorRect.AddEffect(effect.Posterize, 6)
	watercolorRect.AddEffect(effect.Saturate, 1.8)

	// Miniature/tilt-shift scene
	miniatureRect := canvas.NewRectangle(color.NRGBA{100, 180, 60, 255})
	miniatureRect.SetMinSize(fyne.NewSize(140, 100))
	miniatureRect.AddEffect(effect.TiltShift, 6.0, 0.5, 0.12)
	miniatureRect.AddEffect(effect.Saturate, 1.6)
	miniatureRect.AddEffect(effect.Contrast, 1.2)

	stackSection := container.NewVBox(
		stackTitle,
		container.NewGridWrap(fyne.NewSize(160, 130),
			container.NewVBox(frostedGlass, widget.NewLabel("Frosted Glass\nBlur+Bright+Opacity")),
			container.NewVBox(vintagePhoto, widget.NewLabel("Vintage Photo\nSepia+Vig+Grain+Contrast")),
			container.NewVBox(disabledState, widget.NewLabel("Disabled State\nGrayscale+Opacity")),
			container.NewVBox(glitchRect, widget.NewLabel("Glitch/Cyberpunk\nChroma+Scan+Contrast")),
			container.NewVBox(dreamRect, widget.NewLabel("Dream Soft Focus\nBlur+Bright+Saturate")),
			container.NewVBox(metalRect, widget.NewLabel("Embossed Metal\nEmboss+Contrast+Bright")),
			container.NewVBox(neonRect, widget.NewLabel("Neon Glow\nOuterGlow+Bright")),
			container.NewVBox(thermalRect, widget.NewLabel("Thermal Vision\nPosterize+Hue+Saturate")),
			container.NewVBox(sketchRect, widget.NewLabel("Pencil Sketch\nEdge+Invert+Gray")),
			container.NewVBox(oldTVRect, widget.NewLabel("Old TV\nCRT+Gray+Grain")),
			container.NewVBox(watercolorRect, widget.NewLabel("Watercolor\nBlur+Posterize+Saturate")),
			container.NewVBox(miniatureRect, widget.NewLabel("Miniature Scene\nTiltShift+Sat+Contrast")),
		),
		widget.NewSeparator(),
	)

	// ==================== ANIMATED EFFECTS ====================
	animTitle := sectionTitle("Animated Effects")

	// Pulsing blur
	pulseRect := canvas.NewRectangle(red)
	pulseRect.SetMinSize(fyne.NewSize(140, 80))
	pulseBlur := pulseRect.AddEffect(effect.GaussianBlur, 0.0)
	effect.SetOwner(pulseBlur, pulseRect)

	// Breathing opacity
	breatheRect := canvas.NewRectangle(blue)
	breatheRect.SetMinSize(fyne.NewSize(140, 80))
	breatheOp := breatheRect.AddEffect(effect.Opacity, 1.0)
	effect.SetOwner(breatheOp, breatheRect)

	// Color cycling via hue rotation
	hueRect := canvas.NewRectangle(red)
	hueRect.SetMinSize(fyne.NewSize(140, 80))
	hueEff := hueRect.AddEffect(effect.HueRotate, 0)
	effect.SetOwner(hueEff, hueRect)

	// Glitch pulse
	glitchAnimRect := canvas.NewRectangle(cyan)
	glitchAnimRect.SetMinSize(fyne.NewSize(140, 80))
	glitchChroma := glitchAnimRect.AddEffect(effect.ChromaticAberration, 0.0)
	effect.SetOwner(glitchChroma, glitchAnimRect)
	glitchAnimRect.AddEffect(effect.Scanlines, 1.0, 0.2)

	startBtn := widget.NewButton("Start All Animations", func() {
		// Pulsing blur: 0 → 8 → 0
		a1 := fyne.NewAnimation(2*time.Second, func(t float32) {
			v := t * 2
			if v > 1.0 {
				v = 2.0 - v
			}
			pulseBlur.SetFloat("radius", v*8.0)
		})
		a1.AutoReverse = true
		a1.RepeatCount = fyne.AnimationRepeatForever
		a1.Start()

		// Breathing opacity: 1.0 → 0.2 → 1.0
		a2 := fyne.NewAnimation(3*time.Second, func(t float32) {
			v := t * 2
			if v > 1.0 {
				v = 2.0 - v
			}
			breatheOp.SetFloat("opacity", 0.2+v*0.8)
		})
		a2.AutoReverse = true
		a2.RepeatCount = fyne.AnimationRepeatForever
		a2.Start()

		// Hue cycling: 0 → 360
		a3 := fyne.NewAnimation(4*time.Second, func(t float32) {
			hueEff.SetFloat("angle", t*360.0)
		})
		a3.RepeatCount = fyne.AnimationRepeatForever
		a3.Start()

		// Glitch pulse: 0 → 10 → 0
		a4 := fyne.NewAnimation(time.Second, func(t float32) {
			v := t * 2
			if v > 1.0 {
				v = 2.0 - v
			}
			glitchChroma.SetFloat("offset", v*10.0)
		})
		a4.AutoReverse = true
		a4.RepeatCount = fyne.AnimationRepeatForever
		a4.Start()
	})

	animSection := container.NewVBox(
		animTitle,
		container.NewHBox(
			container.NewVBox(pulseRect, widget.NewLabel("Pulsing Blur")),
			container.NewVBox(breatheRect, widget.NewLabel("Breathing Opacity")),
			container.NewVBox(hueRect, widget.NewLabel("Hue Cycling")),
			container.NewVBox(glitchAnimRect, widget.NewLabel("Glitch Pulse")),
			layout.NewSpacer(),
			startBtn,
		),
		widget.NewSeparator(),
	)

	// ==================== TEXT WITH EFFECTS ====================
	textTitle := sectionTitle("Text with Effects")

	blurredText := canvas.NewText("Blurred Text", color.White)
	blurredText.TextSize = 20
	blurredText.SetMinSize(fyne.NewSize(160, 30))
	blurredText.AddEffect(effect.GaussianBlur, 2.0)

	sepiaText := canvas.NewText("Vintage Text", color.NRGBA{200, 180, 140, 255})
	sepiaText.TextSize = 20
	sepiaText.SetMinSize(fyne.NewSize(160, 30))
	sepiaText.AddEffect(effect.Sepia, 0.8)

	invertText := canvas.NewText("Inverted Text", color.White)
	invertText.TextSize = 20
	invertText.SetMinSize(fyne.NewSize(160, 30))
	invertText.AddEffect(effect.Invert, 1.0)

	textSection := container.NewVBox(
		textTitle,
		container.NewHBox(blurredText, sepiaText, invertText),
		widget.NewSeparator(),
	)

	// ==================== STATS ====================
	statsText := canvas.NewText("56 built-in effects · 112 shader files (GL+ES) · 77 effect types · 4 platforms · 36 unit tests", color.NRGBA{120, 120, 120, 255})
	statsText.TextSize = 11
	statsText.Alignment = fyne.TextAlignCenter

	// ==================== ASSEMBLE ====================
	content := container.NewVBox(
		header,
		colorSection,
		blurSection,
		distortionSection,
		detailSection,
		shadowSection,
		styleSection,
		lightSection,
		gradientSection,
		advancedSection,
		blendSection,
		stackSection,
		animSection,
		textSection,
		container.NewCenter(statsText),
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
