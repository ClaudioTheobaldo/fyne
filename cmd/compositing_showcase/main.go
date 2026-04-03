package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func section(text string) *canvas.Text {
	t := canvas.NewText(text, color.NRGBA{80, 200, 255, 255})
	t.TextSize = 16
	t.TextStyle = fyne.TextStyle{Bold: true}
	return t
}

func card(label string, c color.Color, kind effect.EffectType, params ...float32) *fyne.Container {
	rect := canvas.NewRectangle(c)
	rect.SetMinSize(fyne.NewSize(120, 80))
	rect.CornerRadius = 6
	rect.AddEffect(kind, params...)
	lbl := widget.NewLabel(label)
	lbl.Alignment = fyne.TextAlignCenter
	return container.NewVBox(rect, lbl)
}

func gradCard(label string, c1, c2 color.Color, angle float64, kind effect.EffectType, params ...float32) *fyne.Container {
	g := canvas.NewLinearGradient(c1, c2, angle)
	g.SetMinSize(fyne.NewSize(120, 80))
	g.AddEffect(kind, params...)
	lbl := widget.NewLabel(label)
	lbl.Alignment = fyne.TextAlignCenter
	return container.NewVBox(g, lbl)
}

func main() {
	a := app.New()
	w := a.NewWindow("Compositing, Procedural & Animated Effects")
	w.Resize(fyne.NewSize(1000, 800))

	heading := canvas.NewText("Compositing · Procedural · Animated Effects", color.White)
	heading.TextSize = 20
	heading.TextStyle = fyne.TextStyle{Bold: true}

	red := color.NRGBA{220, 60, 60, 255}
	blue := color.NRGBA{37, 99, 235, 255}
	green := color.NRGBA{40, 180, 80, 255}
	purple := color.NRGBA{124, 58, 237, 255}
	orange := color.NRGBA{240, 140, 40, 255}
	cyan := color.NRGBA{40, 200, 220, 255}
	pink := color.NRGBA{240, 100, 180, 255}

	// ==================== COLOR REPLACE ====================
	colorReplaceRow := container.NewGridWrap(fyne.NewSize(140, 110),
		// Replace red with blue (tolerance 0.3)
		card("Red→Blue\ntol=0.3", red, effect.ColorReplace,
			0.86, 0.24, 0.24, 0.15, 0.39, 0.92, 0.3),
		card("Red→Green\ntol=0.5", red, effect.ColorReplace,
			0.86, 0.24, 0.24, 0.16, 0.71, 0.31, 0.5),
		card("Blue→Pink\ntol=0.3", blue, effect.ColorReplace,
			0.15, 0.39, 0.92, 0.94, 0.39, 0.71, 0.3),
	)

	// ==================== DUOTONE ====================
	duotoneRow := container.NewGridWrap(fyne.NewSize(140, 110),
		gradCard("Cyan/Magenta", red, blue, 45, effect.Duotone,
			0, 0.8, 0.8, 0.8, 0, 0.8),
		gradCard("Navy/Gold", green, purple, 30, effect.Duotone,
			0.05, 0.05, 0.2, 1, 0.84, 0),
		gradCard("Black/White", orange, cyan, 60, effect.Duotone,
			0, 0, 0, 1, 1, 1),
		gradCard("Purple/Pink", red, green, 90, effect.Duotone,
			0.3, 0, 0.5, 1, 0.4, 0.7),
	)

	// ==================== SPLIT TONE ====================
	splitRow := container.NewGridWrap(fyne.NewSize(140, 110),
		gradCard("Cool/Warm", blue, orange, 45, effect.SplitTone,
			0, 0.2, 0.8, 0.8, 0.6, 0, 0.5),
		gradCard("Teal/Orange", purple, green, 30, effect.SplitTone,
			0, 0.6, 0.5, 1, 0.5, 0, 0.4),
		gradCard("Blue/Gold", red, cyan, 60, effect.SplitTone,
			0, 0, 0.8, 0.8, 0.7, 0, 0.5),
	)

	// ==================== CHANNEL MIXER ====================
	channelRow := container.NewGridWrap(fyne.NewSize(140, 110),
		gradCard("R↔B swap", red, blue, 45, effect.ChannelMixer,
			0, 0, 1, 0, 1, 0, 1, 0, 0),
		gradCard("Green only", red, green, 30, effect.ChannelMixer,
			0, 1, 0, 0, 1, 0, 0, 1, 0),
		gradCard("Inverted R", orange, purple, 60, effect.ChannelMixer,
			-1, 0, 0, 0, 1, 0, 0, 0, 1),
	)

	// ==================== GRADIENT MAP ====================
	gradMapRow := container.NewGridWrap(fyne.NewSize(140, 110),
		gradCard("Thermal", blue, red, 45, effect.GradientMap,
			0, 0, 0.2, 0, 0, 0.8, 0, 0.8, 0, 0.8, 0.8, 0, 0.8, 0, 0),
		gradCard("Sunset", green, orange, 30, effect.GradientMap,
			0.1, 0, 0.2, 0.5, 0, 0.3, 0.9, 0.3, 0, 1, 0.8, 0.2, 1, 1, 0.5),
		gradCard("Neon", purple, cyan, 60, effect.GradientMap,
			0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1),
	)

	// ==================== PATTERN OVERLAY ====================
	patternRow := container.NewGridWrap(fyne.NewSize(140, 110),
		card("Checkerboard", blue, effect.PatternOverlay, 10, 0, 0.3),
		card("H Stripes", green, effect.PatternOverlay, 8, 1, 0.3),
		card("V Stripes", purple, effect.PatternOverlay, 8, 2, 0.3),
		card("Dots", orange, effect.PatternOverlay, 12, 3, 0.3),
	)

	// ==================== NOISE DISPLACEMENT ====================
	noiseRow := container.NewGridWrap(fyne.NewSize(140, 110),
		gradCard("Subtle", cyan, pink, 45, effect.NoiseDisplacement, 2, 5, 0),
		gradCard("Medium", cyan, pink, 45, effect.NoiseDisplacement, 5, 8, 0),
		gradCard("Strong", cyan, pink, 45, effect.NoiseDisplacement, 10, 10, 0),
		gradCard("Fine", cyan, pink, 45, effect.NoiseDisplacement, 4, 20, 42),
	)

	// ==================== ANIMATED EFFECTS ====================
	// These need time-driven uniforms — we use fyne.Animation to drive them

	// Shimmer
	shimmerRect := canvas.NewRectangle(blue)
	shimmerRect.SetMinSize(fyne.NewSize(180, 70))
	shimmerRect.CornerRadius = 10
	shimmerEff := shimmerRect.AddEffect(effect.Shimmer, 0, 0.08, 30, 0.6)
	effect.SetOwner(shimmerEff, shimmerRect)

	// Pulse
	pulseRect := canvas.NewRectangle(purple)
	pulseRect.SetMinSize(fyne.NewSize(180, 70))
	pulseRect.CornerRadius = 10
	pulseEff := pulseRect.AddEffect(effect.Pulse, 0, 1.0, 0.7, 1.3, 0.9, 1.1)
	effect.SetOwner(pulseEff, pulseRect)

	// Glitch
	glitchRect := canvas.NewRectangle(red)
	glitchRect.SetMinSize(fyne.NewSize(180, 70))
	glitchRect.CornerRadius = 10
	glitchEff := glitchRect.AddEffect(effect.Glitch, 0, 3, 20)
	effect.SetOwner(glitchEff, glitchRect)

	// Matrix Rain
	matrixRect := canvas.NewRectangle(color.NRGBA{5, 15, 5, 255})
	matrixRect.SetMinSize(fyne.NewSize(180, 70))
	matrixRect.CornerRadius = 4
	matrixEff := matrixRect.AddEffect(effect.MatrixRain, 0, 8, 0.5, 0.8)
	effect.SetOwner(matrixEff, matrixRect)

	startBtn := widget.NewButton("Start Animations", func() {
		startTime := time.Now()
		anim := fyne.NewAnimation(100*time.Second, func(t float32) {
			elapsed := float32(time.Since(startTime).Seconds())
			shimmerEff.SetFloat("time", elapsed*0.5)
			pulseEff.SetFloat("time", elapsed)
			glitchEff.SetFloat("time", elapsed)
			matrixEff.SetFloat("time", elapsed)
		})
		anim.RepeatCount = fyne.AnimationRepeatForever
		anim.Start()
	})

	animRow := container.NewGridWrap(fyne.NewSize(200, 100),
		container.NewVBox(shimmerRect, widget.NewLabel("Shimmer")),
		container.NewVBox(pulseRect, widget.NewLabel("Pulse")),
		container.NewVBox(glitchRect, widget.NewLabel("Glitch")),
		container.NewVBox(matrixRect, widget.NewLabel("Matrix Rain")),
	)

	content := container.NewVBox(
		container.NewCenter(heading),
		widget.NewSeparator(),
		section("Color Replace (chroma key)"), colorReplaceRow,
		widget.NewSeparator(),
		section("Duotone (luminance → 2 colors)"), duotoneRow,
		widget.NewSeparator(),
		section("Split Tone (shadow/highlight tinting)"), splitRow,
		widget.NewSeparator(),
		section("Channel Mixer (RGB remap)"), channelRow,
		widget.NewSeparator(),
		section("Gradient Map (luminance → 5-stop ramp)"), gradMapRow,
		widget.NewSeparator(),
		section("Pattern Overlay (checker/stripes/dots)"), patternRow,
		widget.NewSeparator(),
		section("Noise Displacement (Perlin UV warp)"), noiseRow,
		widget.NewSeparator(),
		section("Time-Based Animated Effects (click Start)"), animRow,
		container.NewCenter(startBtn),
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
