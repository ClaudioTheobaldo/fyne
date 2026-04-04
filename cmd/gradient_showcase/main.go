package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Gradient Showcase — Elliptical + Multi-Stop + Scanlines")
	w.Resize(fyne.NewSize(960, 640))

	// ─────────────────────────────────────────────────────
	// Tab 1: Config-screen background (CSS replica)
	// ─────────────────────────────────────────────────────

	configBg := buildConfigBackground()

	// Overlay some content to show transparency
	title := canvas.NewText("SETTINGS", color.NRGBA{232, 228, 223, 255})
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Gradient background with scanlines — replicated from CSS", color.NRGBA{232, 228, 223, 128})
	subtitle.TextSize = 12

	configTab := container.NewStack(
		configBg,
		container.NewVBox(
			widget.NewSeparator(),
			container.NewCenter(title),
			container.NewCenter(subtitle),
			widget.NewSeparator(),
		),
	)

	// ─────────────────────────────────────────────────────
	// Tab 2: Elliptical radial gradient comparison
	// ─────────────────────────────────────────────────────

	circularGrad := canvas.NewRadialGradient(
		color.NRGBA{R: 37, G: 99, B: 235, A: 255},
		color.Transparent,
	)
	circularGrad.SetMinSize(fyne.NewSize(200, 200))

	ellipticalGrad := canvas.NewRadialGradient(
		color.NRGBA{R: 37, G: 99, B: 235, A: 255},
		color.Transparent,
	)
	ellipticalGrad.ScaleX = 2.0
	ellipticalGrad.ScaleY = 0.7
	ellipticalGrad.SetMinSize(fyne.NewSize(200, 200))

	circLabel := canvas.NewText("Circular (default)", color.White)
	circLabel.Alignment = fyne.TextAlignCenter
	ellipLabel := canvas.NewText("Elliptical (ScaleX=2.0, ScaleY=0.7)", color.White)
	ellipLabel.Alignment = fyne.TextAlignCenter

	bg2 := canvas.NewRectangle(color.NRGBA{20, 20, 25, 255})
	ellipTab := container.NewStack(
		bg2,
		container.NewCenter(container.NewHBox(
			container.NewVBox(circLabel, circularGrad),
			widget.NewSeparator(),
			container.NewVBox(ellipLabel, ellipticalGrad),
		)),
	)

	// ─────────────────────────────────────────────────────
	// Tab 3: Multi-stop linear gradients
	// ─────────────────────────────────────────────────────

	rainbow := &canvas.LinearGradient{
		Angle: 270, // horizontal
		Stops: []canvas.GradientStop{
			{Color: color.NRGBA{R: 255, A: 255}, Position: 0.0},
			{Color: color.NRGBA{R: 255, G: 165, A: 255}, Position: 0.2},
			{Color: color.NRGBA{R: 255, G: 255, A: 255}, Position: 0.4},
			{Color: color.NRGBA{G: 255, A: 255}, Position: 0.6},
			{Color: color.NRGBA{B: 255, A: 255}, Position: 0.8},
		},
	}
	rainbow.SetMinSize(fyne.NewSize(600, 60))

	sunset := &canvas.LinearGradient{
		Angle: 0, // vertical
		Stops: []canvas.GradientStop{
			{Color: color.NRGBA{R: 25, G: 25, B: 50, A: 255}, Position: 0.0},
			{Color: color.NRGBA{R: 180, G: 60, B: 100, A: 255}, Position: 0.4},
			{Color: color.NRGBA{R: 255, G: 150, B: 50, A: 255}, Position: 0.7},
			{Color: color.NRGBA{R: 255, G: 220, B: 100, A: 255}, Position: 1.0},
		},
	}
	sunset.SetMinSize(fyne.NewSize(600, 120))

	twoStopWithEarlyEnd := &canvas.LinearGradient{
		Angle: 270,
		Stops: []canvas.GradientStop{
			{Color: color.NRGBA{R: 232, G: 67, B: 42, A: 255}, Position: 0.0},
			{Color: color.Transparent, Position: 0.5},
		},
	}
	twoStopWithEarlyEnd.SetMinSize(fyne.NewSize(600, 60))

	bg3 := canvas.NewRectangle(color.NRGBA{20, 20, 25, 255})
	multiStopTab := container.NewStack(
		bg3,
		container.NewVBox(
			canvas.NewText("Rainbow (5 stops, horizontal)", color.White),
			rainbow,
			widget.NewSeparator(),
			canvas.NewText("Sunset (4 stops, vertical)", color.White),
			sunset,
			widget.NewSeparator(),
			canvas.NewText("Fade to transparent at 50% (2 stops with early end)", color.White),
			twoStopWithEarlyEnd,
		),
	)

	// ─────────────────────────────────────────────────────
	// Tab 4: Effect-based gradients (GPU) + Scanlines
	// ─────────────────────────────────────────────────────

	effectRect := canvas.NewRectangle(color.NRGBA{10, 10, 14, 255})
	effectRect.SetMinSize(fyne.NewSize(600, 400))

	// Layer 1: linear gradient base
	effectRect.AddEffect(effect.LinearGradientOverlay,
		// startColor RGBA
		10.0/255.0, 10.0/255.0, 14.0/255.0, 1.0,
		// endColor RGBA
		18.0/255.0, 14.0/255.0, 12.0/255.0, 1.0,
		// angle
		180.0,
	)

	// Layer 2: elliptical radial glow (warm)
	e2 := effectRect.AddEffect(effect.RadialGradientOverlay,
		// startColor
		40.0/255.0, 20.0/255.0, 10.0/255.0, 0.4,
		// endColor
		0.0, 0.0, 0.0, 0.0,
		// center
		0.8, 0.3,
		// radius (elliptical)
		0.8, 1.0,
	)
	_ = e2

	// Layer 3: red accent glow
	effectRect.AddEffect(effect.RadialGradientOverlay,
		// startColor
		232.0/255.0, 67.0/255.0, 42.0/255.0, 0.06,
		// endColor
		0.0, 0.0, 0.0, 0.0,
		// center
		0.2, 0.5,
		// radius (elliptical)
		1.2, 0.8,
	)

	// Layer 4: scanlines
	effectRect.AddEffect(effect.Scanlines, 0.5, 0.03)

	effectLabel := canvas.NewText("GPU effect pipeline: LinearGradient + 2x RadialGradient (elliptical) + Scanlines", color.NRGBA{232, 228, 223, 200})
	effectLabel.TextSize = 11

	effectTab := container.NewVBox(
		effectLabel,
		effectRect,
	)

	// ─────────────────────────────────────────────────────
	// Assemble tabs
	// ─────────────────────────────────────────────────────

	tabs := container.NewAppTabs(
		container.NewTabItem("Config Background", configTab),
		container.NewTabItem("Elliptical Gradients", ellipTab),
		container.NewTabItem("Multi-Stop", multiStopTab),
		container.NewTabItem("GPU Effects", effectTab),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}

// buildConfigBackground replicates the CSS configscreen.html background:
//
//	radial-gradient(ellipse 120% 80% at 20% 50%, rgba(232,67,42,0.06) 0%, transparent 60%),
//	radial-gradient(ellipse 80% 100% at 80% 30%, rgba(40,20,10,0.4) 0%, transparent 50%),
//	linear-gradient(180deg, rgba(10,10,14,1) 0%, rgba(18,14,12,1) 100%)
//	+ repeating scan-line overlay at 3% opacity
func buildConfigBackground() fyne.CanvasObject {
	// Base: vertical linear gradient
	bg := canvas.NewVerticalGradient(
		color.NRGBA{10, 10, 14, 255},
		color.NRGBA{18, 14, 12, 255},
	)

	// Warm tint (elliptical radial, center at 80% 30%)
	warmGlow := canvas.NewRadialGradient(
		color.NRGBA{40, 20, 10, 102}, // 0.4 alpha
		color.Transparent,
	)
	warmGlow.CenterOffsetX = 0.3  // push center to 80%
	warmGlow.CenterOffsetY = -0.2 // push center to 30%
	warmGlow.ScaleX = 0.8
	warmGlow.ScaleY = 1.0
	warmGlow.Stops = []canvas.GradientStop{
		{Color: color.NRGBA{40, 20, 10, 102}, Position: 0.0},
		{Color: color.Transparent, Position: 0.5},
	}

	// Red accent glow (elliptical, center at 20% 50%)
	redGlow := canvas.NewRadialGradient(
		color.NRGBA{232, 67, 42, 15}, // 0.06 alpha
		color.Transparent,
	)
	redGlow.CenterOffsetX = -0.3
	redGlow.ScaleX = 1.2
	redGlow.ScaleY = 0.8
	redGlow.Stops = []canvas.GradientStop{
		{Color: color.NRGBA{232, 67, 42, 15}, Position: 0.0},
		{Color: color.Transparent, Position: 0.6},
	}

	// Scan-lines on a transparent rectangle
	scanRect := canvas.NewRectangle(color.Transparent)
	scanRect.AddEffect(effect.Scanlines, 0.5, 0.03)

	return container.NewStack(bg, warmGlow, redGlow, scanRect)
}
