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

var (
	blue   = color.NRGBA{37, 99, 235, 255}
	purple = color.NRGBA{124, 58, 237, 255}
	red    = color.NRGBA{220, 60, 60, 255}
	green  = color.NRGBA{40, 180, 80, 255}
	orange = color.NRGBA{240, 140, 40, 255}
	cyan   = color.NRGBA{40, 200, 220, 255}
	pink   = color.NRGBA{240, 100, 180, 255}
	teal   = color.NRGBA{0, 160, 140, 255}
	white  = color.White
)

func card(label string, c color.Color, kind effect.EffectType, params ...float32) *fyne.Container {
	rect := canvas.NewRectangle(c)
	rect.SetMinSize(fyne.NewSize(100, 70))
	rect.CornerRadius = 6
	rect.AddEffect(kind, params...)
	lbl := widget.NewLabel(label)
	lbl.Alignment = fyne.TextAlignCenter
	return container.NewVBox(rect, lbl)
}

func section(title string) *canvas.Text {
	t := canvas.NewText(title, color.NRGBA{80, 200, 255, 255})
	t.TextSize = 16
	t.TextStyle = fyne.TextStyle{Bold: true}
	return t
}

func main() {
	a := app.New()
	w := a.NewWindow("UV Transform Effects — Complete Showcase")
	w.Resize(fyne.NewSize(1000, 800))

	heading := canvas.NewText("UV Transform Effects — 7 New Transforms", color.White)
	heading.TextSize = 20
	heading.TextStyle = fyne.TextStyle{Bold: true}

	// ==================== ROTATE ====================
	rotateRow := container.NewGridWrap(fyne.NewSize(130, 110),
		card("0° (none)", blue, effect.Rotate, 0),
		card("15°", blue, effect.Rotate, 15),
		card("30°", blue, effect.Rotate, 30),
		card("45°", blue, effect.Rotate, 45),
		card("90°", blue, effect.Rotate, 90),
		card("180°", blue, effect.Rotate, 180),
		card("-30°", blue, effect.Rotate, -30),
	)

	// ==================== SCALE ====================
	scaleRow := container.NewGridWrap(fyne.NewSize(130, 110),
		card("0.5x, 0.5x", purple, effect.Scale, 0.5, 0.5),
		card("0.75x, 0.75x", purple, effect.Scale, 0.75, 0.75),
		card("1x, 1x (none)", purple, effect.Scale, 1.0, 1.0),
		card("1.5x, 1.5x", purple, effect.Scale, 1.5, 1.5),
		card("2.0x, 1.0x", purple, effect.Scale, 2.0, 1.0),
		card("1.0x, 2.0x", purple, effect.Scale, 1.0, 2.0),
		card("0.5x, 1.5x", purple, effect.Scale, 0.5, 1.5),
	)

	// ==================== PERSPECTIVE ====================
	_ = fmt.Sprintf // suppress unused
	perspRow := container.NewGridWrap(fyne.NewSize(130, 110),
		card("Bot narrow\n0.3", red, effect.PerspectiveTransform, 0.3, 0),
		card("Bot narrow\n0.6", red, effect.PerspectiveTransform, 0.6, 0),
		card("Top narrow\n0.3", red, effect.PerspectiveTransform, 0.3, 1),
		card("Top narrow\n0.6", red, effect.PerspectiveTransform, 0.6, 1),
		card("Left narrow\n0.3", red, effect.PerspectiveTransform, 0.3, 2),
		card("Right narrow\n0.3", red, effect.PerspectiveTransform, 0.3, 3),
	)

	// ==================== WAVE ====================
	waveRow := container.NewGridWrap(fyne.NewSize(130, 110),
		card("Gentle H", green, effect.Wave, 2, 3, 0, 0),
		card("Strong H", green, effect.Wave, 5, 3, 0, 0),
		card("Fast H", green, effect.Wave, 3, 8, 0, 0),
		card("Gentle V", green, effect.Wave, 2, 3, 0, 1),
		card("Strong V", green, effect.Wave, 5, 3, 0, 1),
		card("Fast V", green, effect.Wave, 3, 8, 0, 1),
	)

	// ==================== PAGE CURL ====================
	curlRow := container.NewGridWrap(fyne.NewSize(130, 110),
		card("Curl 0.0\n(flat)", orange, effect.PageCurl, 0.0, 0.1),
		card("Curl 0.2", orange, effect.PageCurl, 0.2, 0.1),
		card("Curl 0.4", orange, effect.PageCurl, 0.4, 0.1),
		card("Curl 0.6", orange, effect.PageCurl, 0.6, 0.1),
		card("Curl 0.8", orange, effect.PageCurl, 0.8, 0.1),
		card("Curl 1.0\n(full)", orange, effect.PageCurl, 1.0, 0.1),
	)

	// ==================== MIRROR ====================
	mirrorRow := container.NewGridWrap(fyne.NewSize(130, 110),
		card("Original", cyan, effect.Mirror, 0, 0),
		card("Flip X", cyan, effect.Mirror, 1, 0),
		card("Flip Y", cyan, effect.Mirror, 0, 1),
		card("Flip X+Y", cyan, effect.Mirror, 1, 1),
	)

	// Use a gradient to make mirror more visible
	mirrorGrad := canvas.NewLinearGradient(cyan, pink, 30)
	mirrorGrad.SetMinSize(fyne.NewSize(100, 70))

	mirrorGradFlipX := canvas.NewLinearGradient(cyan, pink, 30)
	mirrorGradFlipX.SetMinSize(fyne.NewSize(100, 70))
	mirrorGradFlipX.AddEffect(effect.Mirror, 1, 0)

	mirrorGradFlipY := canvas.NewLinearGradient(cyan, pink, 30)
	mirrorGradFlipY.SetMinSize(fyne.NewSize(100, 70))
	mirrorGradFlipY.AddEffect(effect.Mirror, 0, 1)

	mirrorGradFlipXY := canvas.NewLinearGradient(cyan, pink, 30)
	mirrorGradFlipXY.SetMinSize(fyne.NewSize(100, 70))
	mirrorGradFlipXY.AddEffect(effect.Mirror, 1, 1)

	mirrorGradRow := container.NewGridWrap(fyne.NewSize(130, 110),
		container.NewVBox(mirrorGrad, widget.NewLabel("Original")),
		container.NewVBox(mirrorGradFlipX, widget.NewLabel("Flip X")),
		container.NewVBox(mirrorGradFlipY, widget.NewLabel("Flip Y")),
		container.NewVBox(mirrorGradFlipXY, widget.NewLabel("Flip X+Y")),
	)

	// ==================== KALEIDOSCOPE ====================
	kaleidRow := container.NewGridWrap(fyne.NewSize(130, 110),
		card("3 segments", teal, effect.Kaleidoscope, 3, 0),
		card("4 segments", teal, effect.Kaleidoscope, 4, 0),
		card("6 segments", teal, effect.Kaleidoscope, 6, 0),
		card("8 segments", teal, effect.Kaleidoscope, 8, 0),
		card("12 segments", teal, effect.Kaleidoscope, 12, 0),
		card("6 + rot 30°", teal, effect.Kaleidoscope, 6, 30),
	)

	// Use a gradient for more visible kaleidoscope
	kGrad := func(segs, rot float32, label string) *fyne.Container {
		g := canvas.NewLinearGradient(teal, pink, 45)
		g.SetMinSize(fyne.NewSize(100, 100))
		g.AddEffect(effect.Kaleidoscope, segs, rot)
		return container.NewVBox(g, widget.NewLabel(label))
	}

	kaleidGradRow := container.NewGridWrap(fyne.NewSize(130, 130),
		kGrad(3, 0, "3 segs"),
		kGrad(4, 0, "4 segs"),
		kGrad(6, 0, "6 segs"),
		kGrad(8, 0, "8 segs"),
		kGrad(12, 0, "12 segs"),
		kGrad(6, 45, "6 + rot 45°"),
	)

	content := container.NewVBox(
		container.NewCenter(heading),
		widget.NewSeparator(),
		section("1. Rotate (angle in degrees)"), rotateRow,
		widget.NewSeparator(),
		section("2. Scale / Zoom (scaleX, scaleY)"), scaleRow,
		widget.NewSeparator(),
		section("3. Perspective Transform (amount, direction)"), perspRow,
		widget.NewSeparator(),
		section("4. Wave / Flag (amplitude, frequency, phase, direction)"), waveRow,
		widget.NewSeparator(),
		section("5. Page Curl (curl progress, radius)"), curlRow,
		widget.NewSeparator(),
		section("6. Mirror / Flip (flipX, flipY) — solid + gradient"), mirrorRow, mirrorGradRow,
		widget.NewSeparator(),
		section("7. Kaleidoscope (segments, rotation) — solid + gradient"), kaleidRow, kaleidGradRow,
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
