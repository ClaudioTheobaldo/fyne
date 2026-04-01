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

var (
	gold      = color.NRGBA{255, 215, 0, 255}
	darkGold  = color.NRGBA{191, 160, 0, 255}
	deepBlue  = color.NRGBA{37, 99, 235, 255}
	lightBlue = color.NRGBA{96, 165, 250, 255}
	deepPurple = color.NRGBA{124, 58, 237, 255}
	richBlack  = color.NRGBA{17, 17, 17, 255}
	softWhite  = color.NRGBA{246, 247, 248, 255}
	darkText   = color.NRGBA{30, 30, 30, 255}
)

func main() {
	a := app.New()
	w := a.NewWindow("CSS Button Effects → Fyne Replication")
	w.Resize(fyne.NewSize(1000, 750))

	title := canvas.NewText("CSS → Fyne Button Effect Replication", color.White)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Source: veebilehed24.ee — 20 Modern CSS Buttons (2026)", color.NRGBA{160, 160, 160, 255})
	subtitle.TextSize = 12

	// ============================================================
	// #10 — FLOATING GLASS BUTTON
	// CSS: backdrop-filter:blur(8px), translateY(-6px), box-shadow on hover
	// Fyne: GaussianBlur + Opacity + Brightness + DropShadow animated
	// ============================================================
	glass := NewEffectButton("Floating Glass", color.NRGBA{255, 255, 255, 90}, darkText, nil)
	glass.Radius = 24
	glass.MinWidth = 180

	var glassBlur, glassShadow *effect.Effect
	glass.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if glassBlur == nil {
			glassBlur = bg.AddEffect(effect.GaussianBlur, 0.0)
			effect.SetOwner(glassBlur, bg)
			glassShadow = bg.AddEffect(effect.DropShadow, 0, 4, 0, 0, 0, 0, 0)
			effect.SetOwner(glassShadow, bg)
		}
		animateEffect(glassBlur, "radius", 0, 3, 250*time.Millisecond)
		animateEffect(glassShadow, "blur", 0, 8, 250*time.Millisecond)
		animateEffect(glassShadow, "offsetY", 0, 6, 250*time.Millisecond)
	}
	glass.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		if glassBlur != nil {
			animateEffect(glassBlur, "radius", 3, 0, 200*time.Millisecond)
			animateEffect(glassShadow, "blur", 8, 0, 200*time.Millisecond)
			animateEffect(glassShadow, "offsetY", 6, 0, 200*time.Millisecond)
		}
	}

	glassDesc := widget.NewLabel("CSS: backdrop-filter:blur + box-shadow\nFyne: GaussianBlur + DropShadow animated")
	glassDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// #13 — NEUMORPHIC BUTTON
	// CSS: Soft outer shadows that invert to inset on press
	// Fyne: DropShadow + InnerShadow toggled
	// ============================================================
	neu := NewEffectButton("Neumorphic", softWhite, color.NRGBA{202, 162, 74, 255}, nil)
	neu.Radius = 30
	neu.MinWidth = 180

	var neuShadow *effect.Effect
	neu.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if neuShadow == nil {
			neuShadow = bg.AddEffect(effect.DropShadow, 0, 0, 6, 0.2, 0.2, 0.2, 0.3)
			effect.SetOwner(neuShadow, bg)
		}
		animateEffect(neuShadow, "blur", 6, 12, 250*time.Millisecond)
	}
	neu.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		if neuShadow != nil {
			animateEffect(neuShadow, "blur", 12, 6, 200*time.Millisecond)
		}
	}
	neu.OnPress = func(bg *canvas.Rectangle, label *canvas.Text) {
		if neuShadow != nil {
			neuShadow.SetFloat("blur", 2)
		}
	}
	neu.OnRelease = func(bg *canvas.Rectangle, label *canvas.Text) {
		if neuShadow != nil {
			animateEffect(neuShadow, "blur", 2, 6, 150*time.Millisecond)
		}
	}

	neuDesc := widget.NewLabel("CSS: box-shadow outer/inset toggle\nFyne: DropShadow blur animated on hover/press")
	neuDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// #14 — 3D EXTRUDE BUTTON
	// CSS: box-shadow:0 6px 0 color, translateY(-6px) hover, translateY(3px) active
	// Fyne: DropShadow + Brightness shift
	// ============================================================
	extrude := NewEffectButton("3D Extrude", richBlack, gold, nil)
	extrude.Radius = 6
	extrude.MinWidth = 180

	var extrudeShadow *effect.Effect
	extrude.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if extrudeShadow == nil {
			extrudeShadow = bg.AddEffect(effect.DropShadow, 0, 6, 1, 0.75, 0.6, 0, 0.8)
			effect.SetOwner(extrudeShadow, bg)
		}
		animateEffect(extrudeShadow, "offsetY", 6, 10, 150*time.Millisecond)
		animateEffect(extrudeShadow, "blur", 1, 3, 150*time.Millisecond)
	}
	extrude.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		if extrudeShadow != nil {
			animateEffect(extrudeShadow, "offsetY", 10, 6, 150*time.Millisecond)
			animateEffect(extrudeShadow, "blur", 3, 1, 150*time.Millisecond)
		}
	}
	extrude.OnPress = func(bg *canvas.Rectangle, label *canvas.Text) {
		if extrudeShadow != nil {
			extrudeShadow.SetFloat("offsetY", 2)
			extrudeShadow.SetFloat("blur", 0)
		}
	}
	extrude.OnRelease = func(bg *canvas.Rectangle, label *canvas.Text) {
		if extrudeShadow != nil {
			animateEffect(extrudeShadow, "offsetY", 2, 6, 150*time.Millisecond)
			animateEffect(extrudeShadow, "blur", 0, 1, 150*time.Millisecond)
		}
	}

	extrudeDesc := widget.NewLabel("CSS: box-shadow:0 6px 0 color\nFyne: DropShadow animated on hover/press")
	extrudeDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// #19 — GRADIENT MOTION BUTTON
	// CSS: background-size:200% 200%, background-position shifts on hover
	// Fyne: LinearGradientOverlay + Brightness animated
	// ============================================================
	gradient := NewEffectButton("Gradient Motion", deepPurple, color.White, nil)
	gradient.Radius = 24
	gradient.MinWidth = 180

	var gradBright *effect.Effect
	gradient.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if gradBright == nil {
			bg.AddEffect(effect.LinearGradientOverlay, 0.66, 0.36, 0.93, 0.4, 0.48, 0.26, 0.89, 0.4, 0)
			gradBright = bg.AddEffect(effect.Brightness, 1.0)
			effect.SetOwner(gradBright, bg)
		}
		animateEffect(gradBright, "brightness", 1.0, 1.25, 300*time.Millisecond)
	}
	gradient.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		if gradBright != nil {
			animateEffect(gradBright, "brightness", 1.25, 1.0, 200*time.Millisecond)
		}
	}

	gradientDesc := widget.NewLabel("CSS: background-size:200% + position shift\nFyne: LinearGradientOverlay + Brightness")
	gradientDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// #18 — FLOATING ISLAND BUTTON
	// CSS: shadow deepens, translateY(-6px), padding expands
	// Fyne: DropShadow animated
	// ============================================================
	island := NewEffectButton("Floating Island", deepBlue, color.White, nil)
	island.Radius = 30
	island.MinWidth = 180

	var islandShadow *effect.Effect
	island.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if islandShadow == nil {
			islandShadow = bg.AddEffect(effect.DropShadow, 0, 6, 8, 0.15, 0.39, 0.92, 0.35)
			effect.SetOwner(islandShadow, bg)
		}
		animateEffect(islandShadow, "blur", 8, 16, 250*time.Millisecond)
		animateEffect(islandShadow, "offsetY", 6, 12, 250*time.Millisecond)
	}
	island.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		if islandShadow != nil {
			animateEffect(islandShadow, "blur", 16, 8, 200*time.Millisecond)
			animateEffect(islandShadow, "offsetY", 12, 6, 200*time.Millisecond)
		}
	}

	islandDesc := widget.NewLabel("CSS: box-shadow deepens + translateY\nFyne: DropShadow blur+offset animated")
	islandDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// #20 — OUTLINE EXPAND BUTTON
	// CSS: transparent border fills with color on hover
	// Fyne: Opacity animated from 0 to 1 on a colored rect
	// ============================================================
	outline := NewEffectButton("Outline Expand", color.NRGBA{124, 58, 237, 40}, deepPurple, nil)
	outline.Radius = 4
	outline.MinWidth = 180

	var outlineBright *effect.Effect
	outline.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if outlineBright == nil {
			outlineBright = bg.AddEffect(effect.Brightness, 1.0)
			effect.SetOwner(outlineBright, bg)
		}
		animateEffect(outlineBright, "brightness", 1.0, 2.5, 250*time.Millisecond)
		label.Color = color.White
		canvas.Refresh(label)
	}
	outline.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		if outlineBright != nil {
			animateEffect(outlineBright, "brightness", 2.5, 1.0, 200*time.Millisecond)
		}
		label.Color = deepPurple
		canvas.Refresh(label)
	}

	outlineDesc := widget.NewLabel("CSS: border fills with color on hover\nFyne: Brightness animated + text color swap")
	outlineDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// BONUS — GLOW PULSE BUTTON (no CSS equivalent — our unique effect)
	// Uses OuterGlow + ChromaticAberration animated continuously
	// ============================================================
	glow := NewEffectButton("Glow Pulse", color.NRGBA{30, 30, 50, 255}, color.NRGBA{100, 255, 200, 255}, nil)
	glow.Radius = 14
	glow.MinWidth = 180

	var glowEff *effect.Effect
	glow.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		if glowEff == nil {
			glowEff = bg.AddEffect(effect.OuterGlow, 0, 0.4, 1.0, 0.8, 1.0)
			effect.SetOwner(glowEff, bg)
			bg.AddEffect(effect.ChromaticAberration, 0)
		}
		// Pulse the glow
		anim := fyne.NewAnimation(time.Second, func(t float32) {
			v := t * 2
			if v > 1 {
				v = 2 - v
			}
			glowEff.SetFloat("radius", v * 12)
		})
		anim.AutoReverse = true
		anim.RepeatCount = fyne.AnimationRepeatForever
		anim.Start()
	}

	glowDesc := widget.NewLabel("BONUS — No CSS equivalent!\nFyne: OuterGlow + ChromaticAberration pulsing")
	glowDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// BONUS 2 — CRT RETRO BUTTON
	// ============================================================
	crt := NewEffectButton("RETRO CRT", color.NRGBA{20, 40, 20, 255}, color.NRGBA{0, 255, 0, 255}, nil)
	crt.Radius = 4
	crt.MinWidth = 180

	crt.OnHoverIn = func(bg *canvas.Rectangle, label *canvas.Text) {
		bg.ClearEffects()
		bg.AddEffect(effect.CRT, 0.02, 0.3, 0.3)
		bg.AddEffect(effect.Scanlines, 1.0, 0.3)
		bg.AddEffect(effect.FilmGrain, 0.15, 0)
		bg.AddEffect(effect.Brightness, 1.3)
		bg.Refresh()
	}
	crt.OnHoverOut = func(bg *canvas.Rectangle, label *canvas.Text) {
		bg.ClearEffects()
		bg.Refresh()
	}

	crtDesc := widget.NewLabel("BONUS — Retro terminal look\nFyne: CRT + Scanlines + FilmGrain + Brightness")
	crtDesc.Wrapping = fyne.TextWrapWord

	// ============================================================
	// LAYOUT — COMPARISON EXPLANATION
	// ============================================================
	explanationTitle := canvas.NewText("What Can and Can't Be Replicated", color.NRGBA{255, 200, 100, 255})
	explanationTitle.TextSize = 16
	explanationTitle.TextStyle = fyne.TextStyle{Bold: true}

	canDo := widget.NewLabel(
		"CAN replicate:\n" +
			"· Shadow effects (drop, box, inner, outer glow) — all via FBO shaders\n" +
			"· Color transitions (brightness, opacity, hue) — animated uniforms\n" +
			"· Glassmorphism (blur + opacity) — GaussianBlur effect\n" +
			"· Gradient overlays — LinearGradientOverlay effect\n" +
			"· Stylization (CRT, scanlines, grain) — beyond what CSS can do\n" +
			"· Stacked multi-effect compositions — FBO ping-pong pipeline")
	canDo.Wrapping = fyne.TextWrapWord

	cantDo := widget.NewLabel(
		"CANNOT replicate (yet):\n" +
			"· CSS ::before/::after pseudo-elements (sliding fills) — need geometry layer\n" +
			"· border-radius morphing animation — needs animated corner radius\n" +
			"· translateY for physical movement — effects are post-render, not layout\n" +
			"· Text flip/slide — need per-glyph or layout animation\n" +
			"· SVG stroke animation — no SVG path support in effect system")
	cantDo.Wrapping = fyne.TextWrapWord

	// Build grid
	buttonGrid := container.NewGridWrap(fyne.NewSize(450, 120),
		container.NewHBox(glass, glassDesc),
		container.NewHBox(neu, neuDesc),
		container.NewHBox(extrude, extrudeDesc),
		container.NewHBox(gradient, gradientDesc),
		container.NewHBox(island, islandDesc),
		container.NewHBox(outline, outlineDesc),
		container.NewHBox(glow, glowDesc),
		container.NewHBox(crt, crtDesc),
	)

	content := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		widget.NewSeparator(),
		buttonGrid,
		widget.NewSeparator(),
		explanationTitle,
		canDo,
		cantDo,
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
