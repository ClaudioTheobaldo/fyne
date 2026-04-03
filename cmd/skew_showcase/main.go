package main

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Skew on Every Primitive")
	w.Resize(fyne.NewSize(950, 700))

	title := canvas.NewText("effect.Skew on Every Canvas Primitive", color.White)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}

	skewX := float32(0.5)
	skewY := float32(0.0)

	// ===== 1. Rectangle =====
	rect := canvas.NewRectangle(color.NRGBA{37, 99, 235, 255})
	rect.SetMinSize(fyne.NewSize(150, 70))
	rect.CornerRadius = 10
	rect.AddEffect(effect.Skew, skewX, skewY)

	// ===== 2. Text =====
	txt := canvas.NewText("Hello Fyne!", color.NRGBA{255, 215, 0, 255})
	txt.TextSize = 28
	txt.TextStyle = fyne.TextStyle{Bold: true}
	txt.SetMinSize(fyne.NewSize(200, 40))
	txt.AddEffect(effect.Skew, skewX, skewY)

	// ===== 3. Circle =====
	circle := canvas.NewCircle(color.NRGBA{220, 60, 60, 255})
	circle.StrokeColor = color.NRGBA{255, 200, 200, 255}
	circle.StrokeWidth = 3
	circle.Resize(fyne.NewSize(100, 100))
	circle.AddEffect(effect.Skew, skewX, skewY)

	// ===== 4. Line =====
	line := canvas.NewLine(color.NRGBA{0, 255, 200, 255})
	line.StrokeWidth = 4
	line.Resize(fyne.NewSize(150, 60))
	line.AddEffect(effect.Skew, skewX, skewY)

	// ===== 5. Arc (Pie) =====
	arc := canvas.NewPieArc(0, 270, color.NRGBA{240, 140, 40, 255})
	arc.StrokeColor = color.NRGBA{255, 200, 100, 255}
	arc.StrokeWidth = 2
	arc.Resize(fyne.NewSize(100, 100))
	arc.SetMinSize(fyne.NewSize(100, 100))
	arc.AddEffect(effect.Skew, skewX, skewY)

	// ===== 6. Polygon (Hexagon) =====
	poly := canvas.NewPolygon(6, color.NRGBA{124, 58, 237, 255})
	poly.StrokeColor = color.NRGBA{200, 160, 255, 255}
	poly.StrokeWidth = 2
	poly.Resize(fyne.NewSize(100, 100))
	poly.SetMinSize(fyne.NewSize(100, 100))
	poly.AddEffect(effect.Skew, skewX, skewY)

	// ===== 7. LinearGradient =====
	linGrad := canvas.NewLinearGradient(
		color.NRGBA{255, 0, 100, 255},
		color.NRGBA{0, 100, 255, 255},
		45,
	)
	linGrad.SetMinSize(fyne.NewSize(150, 70))
	linGrad.AddEffect(effect.Skew, skewX, skewY)

	// ===== 8. RadialGradient =====
	radGrad := canvas.NewRadialGradient(
		color.NRGBA{255, 255, 0, 255},
		color.NRGBA{50, 0, 100, 255},
	)
	radGrad.SetMinSize(fyne.NewSize(100, 100))
	radGrad.AddEffect(effect.Skew, skewX, skewY)

	// ===== 9. Raster (procedural) =====
	raster := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		fx := float64(x) / float64(w)
		fy := float64(y) / float64(h)
		r := uint8(128 + 127*math.Sin(fx*10))
		g := uint8(128 + 127*math.Cos(fy*10))
		b := uint8(128 + 127*math.Sin((fx+fy)*8))
		return color.NRGBA{r, g, b, 255}
	})
	raster.SetMinSize(fyne.NewSize(120, 80))
	raster.AddEffect(effect.Skew, skewX, skewY)

	// ===== 10. Image (generated) =====
	img := image.NewNRGBA(image.Rect(0, 0, 120, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 120; x++ {
			r := uint8((x * 255) / 120)
			g := uint8((y * 255) / 80)
			b := uint8(128)
			img.SetNRGBA(x, y, color.NRGBA{r, g, b, 255})
		}
	}
	fyneImg := canvas.NewImageFromImage(img)
	fyneImg.SetMinSize(fyne.NewSize(120, 80))
	fyneImg.FillMode = canvas.ImageFillStretch
	fyneImg.AddEffect(effect.Skew, skewX, skewY)

	// ===== 11. ShaderRect =====
	shaderSrc := `#version 110
uniform vec4 fill_color;
void main() {
    gl_FragColor = fill_color;
}`
	shaderRect := canvas.NewShaderRect(shaderSrc, color.NRGBA{0, 200, 150, 255})
	shaderRect.SetMinSize(fyne.NewSize(150, 70))
	shaderRect.CornerRadius = 8
	shaderRect.AddEffect(effect.Skew, skewX, skewY)

	// ===== LAYOUT =====
	card := func(label string, obj fyne.CanvasObject) *fyne.Container {
		lbl := widget.NewLabel(label)
		lbl.Alignment = fyne.TextAlignCenter
		lbl.TextStyle = fyne.TextStyle{Bold: true}
		return container.NewVBox(container.NewCenter(obj), lbl)
	}

	grid := container.NewGridWrap(fyne.NewSize(180, 140),
		card("Rectangle", rect),
		card("Text", txt),
		card("Circle", circle),
		card("Line", line),
		card("Arc (Pie)", arc),
		card("Polygon (Hex)", poly),
		card("LinearGradient", linGrad),
		card("RadialGradient", radGrad),
		card("Raster (Procedural)", raster),
		card("Image", fyneImg),
		card("ShaderRect", shaderRect),
	)

	note := canvas.NewText("All primitives: AddEffect(effect.Skew, 0.5, 0.0) — same API, same pipeline", color.NRGBA{140, 255, 140, 255})
	note.TextSize = 12
	note.TextStyle = fyne.TextStyle{Monospace: true}

	content := container.NewVBox(
		container.NewCenter(title),
		widget.NewSeparator(),
		grid,
		widget.NewSeparator(),
		container.NewCenter(note),
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
