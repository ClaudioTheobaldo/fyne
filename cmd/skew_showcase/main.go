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

func skewBtn(label string, bgColor, textColor color.Color, skewX, skewY float32) *fyne.Container {
	rect := canvas.NewRectangle(bgColor)
	rect.SetMinSize(fyne.NewSize(140, 50))
	rect.CornerRadius = 8
	rect.AddEffect(effect.Skew, skewX, skewY)

	text := canvas.NewText(label, textColor)
	text.TextSize = 13
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Alignment = fyne.TextAlignCenter

	desc := widget.NewLabel(fmt.Sprintf("skewX=%.1f  skewY=%.1f", skewX, skewY))
	desc.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		container.NewStack(rect, container.NewCenter(text)),
		desc,
	)
}

func main() {
	a := app.New()
	w := a.NewWindow("Skew Effect Showcase")
	w.Resize(fyne.NewSize(900, 750))

	title := canvas.NewText("Skew Effect — UV Distortion Shader", color.White)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Visual skew via fragment shader UV remapping (not geometry transform)", color.NRGBA{160, 160, 160, 255})
	subtitle.TextSize = 12

	blue := color.NRGBA{37, 99, 235, 255}
	purple := color.NRGBA{124, 58, 237, 255}
	red := color.NRGBA{220, 60, 60, 255}
	green := color.NRGBA{40, 180, 80, 255}
	orange := color.NRGBA{240, 140, 40, 255}
	cyan := color.NRGBA{40, 200, 220, 255}
	pink := color.NRGBA{240, 100, 180, 255}
	dark := color.NRGBA{30, 30, 40, 255}
	gold := color.NRGBA{255, 215, 0, 255}
	white := color.White

	// === Section 1: SkewX only ===
	sec1 := canvas.NewText("Horizontal Skew (skewX)", color.NRGBA{100, 200, 255, 255})
	sec1.TextSize = 15
	sec1.TextStyle = fyne.TextStyle{Bold: true}

	skewXRow := container.NewGridWrap(fyne.NewSize(160, 100),
		skewBtn("Skew X -1.0", blue, white, -1.0, 0),
		skewBtn("Skew X -0.5", blue, white, -0.5, 0),
		skewBtn("Skew X -0.2", blue, white, -0.2, 0),
		skewBtn("No Skew", blue, white, 0, 0),
		skewBtn("Skew X +0.2", blue, white, 0.2, 0),
		skewBtn("Skew X +0.5", blue, white, 0.5, 0),
		skewBtn("Skew X +1.0", blue, white, 1.0, 0),
	)

	// === Section 2: SkewY only ===
	sec2 := canvas.NewText("Vertical Skew (skewY)", color.NRGBA{100, 200, 255, 255})
	sec2.TextSize = 15
	sec2.TextStyle = fyne.TextStyle{Bold: true}

	skewYRow := container.NewGridWrap(fyne.NewSize(160, 100),
		skewBtn("Skew Y -1.0", purple, white, 0, -1.0),
		skewBtn("Skew Y -0.5", purple, white, 0, -0.5),
		skewBtn("Skew Y -0.2", purple, white, 0, -0.2),
		skewBtn("No Skew", purple, white, 0, 0),
		skewBtn("Skew Y +0.2", purple, white, 0, 0.2),
		skewBtn("Skew Y +0.5", purple, white, 0, 0.5),
		skewBtn("Skew Y +1.0", purple, white, 0, 1.0),
	)

	// === Section 3: Combined X+Y ===
	sec3 := canvas.NewText("Combined Skew (X + Y)", color.NRGBA{255, 200, 100, 255})
	sec3.TextSize = 15
	sec3.TextStyle = fyne.TextStyle{Bold: true}

	combinedRow := container.NewGridWrap(fyne.NewSize(160, 100),
		skewBtn("Italic", green, white, 0.3, 0),
		skewBtn("Lean Back", red, white, -0.3, 0),
		skewBtn("Shear ↗", orange, white, 0.4, 0.4),
		skewBtn("Shear ↙", cyan, white, -0.4, -0.4),
		skewBtn("Diamond", pink, white, 0.5, -0.5),
		skewBtn("Rhombus", gold, dark, -0.5, 0.5),
		skewBtn("Extreme", dark, gold, 0.8, 0.3),
	)

	// === Section 4: Skew + other effects stacked ===
	sec4 := canvas.NewText("Skew + Stacked Effects", color.NRGBA{100, 255, 200, 255})
	sec4.TextSize = 15
	sec4.TextStyle = fyne.TextStyle{Bold: true}

	// Skew + blur
	skewBlur := canvas.NewRectangle(red)
	skewBlur.SetMinSize(fyne.NewSize(140, 50))
	skewBlur.CornerRadius = 8
	skewBlur.AddEffect(effect.Skew, 0.4, 0)
	skewBlur.AddEffect(effect.GaussianBlur, 2.0)
	skewBlurText := canvas.NewText("Skew + Blur", white)
	skewBlurText.TextSize = 13
	skewBlurText.TextStyle = fyne.TextStyle{Bold: true}
	skewBlurText.Alignment = fyne.TextAlignCenter

	// Skew + sepia
	skewSepia := canvas.NewRectangle(green)
	skewSepia.SetMinSize(fyne.NewSize(140, 50))
	skewSepia.CornerRadius = 8
	skewSepia.AddEffect(effect.Skew, -0.3, 0)
	skewSepia.AddEffect(effect.Sepia, 0.8)
	skewSepiaText := canvas.NewText("Skew + Sepia", white)
	skewSepiaText.TextSize = 13
	skewSepiaText.TextStyle = fyne.TextStyle{Bold: true}
	skewSepiaText.Alignment = fyne.TextAlignCenter

	// Skew + vignette
	skewVig := canvas.NewRectangle(orange)
	skewVig.SetMinSize(fyne.NewSize(140, 50))
	skewVig.CornerRadius = 8
	skewVig.AddEffect(effect.Skew, 0.3, 0.1)
	skewVig.AddEffect(effect.Vignette, 1.2, 0.4)
	skewVigText := canvas.NewText("Skew + Vignette", white)
	skewVigText.TextSize = 13
	skewVigText.TextStyle = fyne.TextStyle{Bold: true}
	skewVigText.Alignment = fyne.TextAlignCenter

	// Skew + chromatic aberration
	skewChroma := canvas.NewRectangle(dark)
	skewChroma.SetMinSize(fyne.NewSize(140, 50))
	skewChroma.CornerRadius = 8
	skewChroma.AddEffect(effect.Skew, 0.2, -0.2)
	skewChroma.AddEffect(effect.ChromaticAberration, 4.0)
	skewChroma.AddEffect(effect.Brightness, 1.3)
	skewChromaText := canvas.NewText("Skew + Glitch", color.NRGBA{0, 255, 200, 255})
	skewChromaText.TextSize = 13
	skewChromaText.TextStyle = fyne.TextStyle{Bold: true}
	skewChromaText.Alignment = fyne.TextAlignCenter

	// Skew + scanlines
	skewScan := canvas.NewRectangle(color.NRGBA{10, 30, 10, 255})
	skewScan.SetMinSize(fyne.NewSize(140, 50))
	skewScan.CornerRadius = 4
	skewScan.AddEffect(effect.Skew, -0.2, 0)
	skewScan.AddEffect(effect.Scanlines, 1.0, 0.3)
	skewScan.AddEffect(effect.Brightness, 1.3)
	skewScanText := canvas.NewText("Skew + Scanlines", color.NRGBA{0, 255, 0, 255})
	skewScanText.TextSize = 13
	skewScanText.TextStyle = fyne.TextStyle{Bold: true}
	skewScanText.Alignment = fyne.TextAlignCenter

	stackedRow := container.NewGridWrap(fyne.NewSize(160, 100),
		container.NewVBox(container.NewStack(skewBlur, container.NewCenter(skewBlurText)), widget.NewLabel("Skew 0.4 + Blur 2")),
		container.NewVBox(container.NewStack(skewSepia, container.NewCenter(skewSepiaText)), widget.NewLabel("Skew -0.3 + Sepia")),
		container.NewVBox(container.NewStack(skewVig, container.NewCenter(skewVigText)), widget.NewLabel("Skew 0.3,0.1 + Vig")),
		container.NewVBox(container.NewStack(skewChroma, container.NewCenter(skewChromaText)), widget.NewLabel("Skew + Chroma + Bright")),
		container.NewVBox(container.NewStack(skewScan, container.NewCenter(skewScanText)), widget.NewLabel("Skew + Scanlines")),
	)

	content := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		widget.NewSeparator(),
		sec1, skewXRow,
		widget.NewSeparator(),
		sec2, skewYRow,
		widget.NewSeparator(),
		sec3, combinedRow,
		widget.NewSeparator(),
		sec4, stackedRow,
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
