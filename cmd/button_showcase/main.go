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

func main() {
	a := app.New()
	w := a.NewWindow("CSS Button Effects — Fyne Replication")
	w.Resize(fyne.NewSize(850, 700))

	title := canvas.NewText("CSS Button Effects Replicated in Fyne", color.White)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}

	// ===== 1. FLOATING GLASS =====
	glass := NewEffectButton("Floating Glass", color.NRGBA{255, 255, 255, 80}, color.NRGBA{40, 40, 40, 255}, nil)
	glass.Radius = 24

	var glassBlur, glassShadow *effect.Effect
	glass.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		glassShadow = bg.AddEffect(effect.DropShadow, 0, 3, 4, 0, 0, 0, 0.15)
		effect.SetOwner(glassShadow, bg)
	}
	glass.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if glassBlur == nil {
			glassBlur = bg.AddEffect(effect.GaussianBlur, 0)
			effect.SetOwner(glassBlur, bg)
		}
		animateEffect(glassBlur, "radius", 0, 3, 250*time.Millisecond)
		animateEffect(glassShadow, "blur", 4, 12, 250*time.Millisecond)
		animateEffect(glassShadow, "offsetY", 3, 8, 250*time.Millisecond)
	}
	glass.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		if glassBlur != nil {
			animateEffect(glassBlur, "radius", 3, 0, 200*time.Millisecond)
		}
		animateEffect(glassShadow, "blur", 12, 4, 200*time.Millisecond)
		animateEffect(glassShadow, "offsetY", 8, 3, 200*time.Millisecond)
	}

	// ===== 2. NEUMORPHIC =====
	neu := NewEffectButton("Neumorphic", color.NRGBA{235, 235, 240, 255}, color.NRGBA{180, 140, 60, 255}, nil)
	neu.Radius = 20

	var neuShadow *effect.Effect
	neu.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		neuShadow = bg.AddEffect(effect.DropShadow, 3, 3, 6, 0.3, 0.3, 0.3, 0.3)
		effect.SetOwner(neuShadow, bg)
	}
	neu.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(neuShadow, "blur", 6, 14, 250*time.Millisecond)
	}
	neu.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(neuShadow, "blur", 14, 6, 200*time.Millisecond)
	}
	neu.OnPress = func(bg *canvas.Rectangle, label *canvas.Text) {
		neuShadow.SetFloat("blur", 2)
		neuShadow.SetFloat("offsetX", 1)
		neuShadow.SetFloat("offsetY", 1)
	}
	neu.OnRelease = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(neuShadow, "blur", 2, 6, 150*time.Millisecond)
		animateEffect(neuShadow, "offsetX", 1, 3, 150*time.Millisecond)
		animateEffect(neuShadow, "offsetY", 1, 3, 150*time.Millisecond)
	}

	// ===== 3. 3D EXTRUDE =====
	extrude := NewEffectButton("3D Extrude", color.NRGBA{20, 20, 20, 255}, color.NRGBA{255, 215, 0, 255}, nil)
	extrude.Radius = 8

	var extShadow *effect.Effect
	extrude.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		extShadow = bg.AddEffect(effect.DropShadow, 0, 5, 0.5, 0.75, 0.6, 0, 0.9)
		effect.SetOwner(extShadow, bg)
	}
	extrude.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(extShadow, "offsetY", 5, 8, 150*time.Millisecond)
		animateEffect(extShadow, "blur", 0.5, 2, 150*time.Millisecond)
	}
	extrude.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(extShadow, "offsetY", 8, 5, 150*time.Millisecond)
		animateEffect(extShadow, "blur", 2, 0.5, 150*time.Millisecond)
	}
	extrude.OnPress = func(bg *canvas.Rectangle, label *canvas.Text) {
		extShadow.SetFloat("offsetY", 1)
		extShadow.SetFloat("blur", 0)
	}
	extrude.OnRelease = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(extShadow, "offsetY", 1, 5, 100*time.Millisecond)
		animateEffect(extShadow, "blur", 0, 0.5, 100*time.Millisecond)
	}

	// ===== 4. GRADIENT GLOW =====
	gradGlow := NewEffectButton("Gradient Glow", color.NRGBA{124, 58, 237, 255}, color.White, nil)
	gradGlow.Radius = 16

	var gradBright *effect.Effect
	gradGlow.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		bg.AddEffect(effect.LinearGradientOverlay, 0.66, 0.36, 0.93, 0.3, 0.48, 0.26, 0.89, 0.3, 45)
		gradBright = bg.AddEffect(effect.Brightness, 1.0)
		effect.SetOwner(gradBright, bg)
	}
	gradGlow.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(gradBright, "brightness", 1.0, 1.3, 300*time.Millisecond)
	}
	gradGlow.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(gradBright, "brightness", 1.3, 1.0, 200*time.Millisecond)
	}

	// ===== 5. FLOATING ISLAND =====
	island := NewEffectButton("Floating Island", color.NRGBA{37, 99, 235, 255}, color.White, nil)
	island.Radius = 24

	var islShadow *effect.Effect
	island.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		islShadow = bg.AddEffect(effect.DropShadow, 0, 5, 8, 0.15, 0.39, 0.92, 0.35)
		effect.SetOwner(islShadow, bg)
	}
	island.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(islShadow, "blur", 8, 18, 250*time.Millisecond)
		animateEffect(islShadow, "offsetY", 5, 12, 250*time.Millisecond)
	}
	island.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(islShadow, "blur", 18, 8, 200*time.Millisecond)
		animateEffect(islShadow, "offsetY", 12, 5, 200*time.Millisecond)
	}

	// ===== 6. OUTLINE FILL =====
	outline := NewEffectButton("Outline Fill", color.NRGBA{124, 58, 237, 30}, color.NRGBA{124, 58, 237, 255}, nil)
	outline.Radius = 6

	var outBright *effect.Effect
	outline.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		outBright = bg.AddEffect(effect.Brightness, 1.0)
		effect.SetOwner(outBright, bg)
	}
	outline.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(outBright, "brightness", 1.0, 3.0, 250*time.Millisecond)
		label.Color = color.White
		canvas.Refresh(label)
	}
	outline.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(outBright, "brightness", 3.0, 1.0, 200*time.Millisecond)
		label.Color = color.NRGBA{124, 58, 237, 255}
		canvas.Refresh(label)
	}

	// ===== 7. NEON GLOW (beyond CSS) =====
	neon := NewEffectButton("Neon Glow", color.NRGBA{15, 15, 30, 255}, color.NRGBA{0, 255, 200, 255}, nil)
	neon.Radius = 10

	var neonGlow *effect.Effect
	neon.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		neonGlow = bg.AddEffect(effect.OuterGlow, 4, 0.0, 1.0, 0.8, 0.8)
		effect.SetOwner(neonGlow, bg)
	}
	neon.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(neonGlow, "radius", 4, 12, 300*time.Millisecond)
	}
	neon.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		animateEffect(neonGlow, "radius", 12, 4, 200*time.Millisecond)
	}

	// ===== 8. RETRO CRT (beyond CSS) =====
	crt := NewEffectButton("RETRO CRT", color.NRGBA{10, 30, 10, 255}, color.NRGBA{0, 255, 0, 255}, nil)
	crt.Radius = 4

	crt.OnInit = func(bg *canvas.Rectangle, label *canvas.Text) {
		bg.AddEffect(effect.Scanlines, 1.0, 0.15)
	}
	crt.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		bg.ClearEffects()
		bg.AddEffect(effect.CRT, 0.02, 0.35, 0.3)
		bg.AddEffect(effect.Scanlines, 1.0, 0.3)
		bg.AddEffect(effect.FilmGrain, 0.12, 0)
		bg.AddEffect(effect.Brightness, 1.4)
		bg.Refresh()
	}
	crt.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		bg.ClearEffects()
		bg.AddEffect(effect.Scanlines, 1.0, 0.15)
		bg.Refresh()
	}

	// ===== LAYOUT =====
	makeRow := func(btn *EffectButton, cssText, fyneText string) *fyne.Container {
		lbl := widget.NewLabel("CSS: " + cssText + "\nFyne: " + fyneText)
		lbl.Wrapping = fyne.TextWrapWord
		return container.NewGridWithColumns(2, container.NewCenter(btn), lbl)
	}

	content := container.NewVBox(
		container.NewCenter(title),
		widget.NewSeparator(),
		makeRow(glass,
			"backdrop-filter:blur + box-shadow + translateY",
			"GaussianBlur + DropShadow animated on hover"),
		widget.NewSeparator(),
		makeRow(neu,
			"box-shadow outer → inset on press",
			"DropShadow blur animated, shrinks on press"),
		widget.NewSeparator(),
		makeRow(extrude,
			"box-shadow:0 6px 0 color + translateY",
			"DropShadow hard offset, animated on hover/press"),
		widget.NewSeparator(),
		makeRow(gradGlow,
			"linear-gradient + brightness on hover",
			"LinearGradientOverlay + Brightness animated"),
		widget.NewSeparator(),
		makeRow(island,
			"box-shadow deepens + translateY(-6px)",
			"DropShadow blur+offset deepens on hover"),
		widget.NewSeparator(),
		makeRow(outline,
			"transparent border → fills with color",
			"Brightness 1→3 animated + text color swap"),
		widget.NewSeparator(),
		makeRow(neon,
			"N/A — beyond CSS without WebGL",
			"OuterGlow radius animated on hover"),
		widget.NewSeparator(),
		makeRow(crt,
			"N/A — beyond CSS without WebGL",
			"CRT + Scanlines + FilmGrain + Brightness stacked"),
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
