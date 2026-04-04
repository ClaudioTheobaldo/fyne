package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ── hover tap widget ─────────────────────────────────────────────────────────
// A zero-overhead transparent overlay that fires hover/press callbacks.
// Place it on top of a canvas rect via container.NewStack.

type hotspot struct {
	widget.BaseWidget
	onIn    func()
	onOut   func()
	onPress func()
	onUp    func()
}

var _ desktop.Hoverable = (*hotspot)(nil)
var _ desktop.Mouseable = (*hotspot)(nil)
var _ fyne.Tappable = (*hotspot)(nil)

func newHotspot(in, out, press, up func()) *hotspot {
	h := &hotspot{onIn: in, onOut: out, onPress: press, onUp: up}
	h.ExtendBaseWidget(h)
	return h
}
func (h *hotspot) CreateRenderer() fyne.WidgetRenderer { return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent)) }
func (h *hotspot) MouseIn(_ *desktop.MouseEvent)       { if h.onIn != nil { h.onIn() } }
func (h *hotspot) MouseOut()                           { if h.onOut != nil { h.onOut() } }
func (h *hotspot) MouseMoved(_ *desktop.MouseEvent)    {}
func (h *hotspot) MouseDown(_ *desktop.MouseEvent)     { if h.onPress != nil { h.onPress() } }
func (h *hotspot) MouseUp(_ *desktop.MouseEvent)       { if h.onUp != nil { h.onUp() } }
func (h *hotspot) Tapped(_ *fyne.PointEvent)           {}

// ── animation helper ─────────────────────────────────────────────────────────

func anim(eff *effect.Effect, name string, from, to float32, dur time.Duration) {
	a := fyne.NewAnimation(dur, func(t float32) {
		eff.SetFloat(name, from+(to-from)*t)
	})
	a.Curve = fyne.AnimationEaseOut
	a.Start()
}

// ── button builder ────────────────────────────────────────────────────────────
// Returns a padded stack: [bg rect] + [label] + [transparent hotspot]
// Padding gives drop shadows room to bleed outside the rect bounds.

func makeButton(
	label string,
	bgColor color.Color,
	textColor color.Color,
	radius float32,
	pad float32, // extra space for shadow bleed
	setup func(bg *canvas.Rectangle, lbl *canvas.Text) (onIn, onOut, onPress, onUp func()),
) fyne.CanvasObject {

	bg := canvas.NewRectangle(bgColor)
	bg.CornerRadius = radius
	bg.SetMinSize(fyne.NewSize(200, 56))

	lbl := canvas.NewText(label, textColor)
	lbl.TextSize = 15
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	lbl.Alignment = fyne.TextAlignCenter

	var onIn, onOut, onPress, onUp func()
	if setup != nil {
		onIn, onOut, onPress, onUp = setup(bg, lbl)
	}

	hs := newHotspot(onIn, onOut, onPress, onUp)

	// Stack: bg fills naturally, label centered, hotspot covers all
	inner := container.NewStack(bg, container.NewCenter(lbl), hs)

	if pad > 0 {
		return container.NewPadded(container.NewPadded(inner))
	}
	return container.NewPadded(inner)
}

// ── code label ───────────────────────────────────────────────────────────────

func codeLabel(s string) *canvas.Text {
	t := canvas.NewText(s, color.NRGBA{100, 220, 100, 255})
	t.TextStyle = fyne.TextStyle{Monospace: true}
	t.TextSize = 11
	return t
}

func descLabel(s string) *widget.Label {
	l := widget.NewLabel(s)
	l.Wrapping = fyne.TextWrapWord
	return l
}

func row(btn fyne.CanvasObject, code, desc string) fyne.CanvasObject {
	info := container.NewVBox(
		codeLabel(code),
		descLabel(desc),
	)
	return container.NewGridWithColumns(2,
		container.NewCenter(btn),
		container.NewPadded(info),
	)
}

// ── buttons ───────────────────────────────────────────────────────────────────

func floatingGlass() fyne.CanvasObject {
	return makeButton("Floating Glass",
		color.NRGBA{255, 255, 255, 60},
		color.NRGBA{30, 30, 30, 255},
		24, 16,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			blur := bg.AddEffect(effect.GaussianBlur, 0)
			shadow := bg.AddEffect(effect.DropShadow, 0, 4, 8, 0, 0, 0, 0.18)
			effect.SetOwner(blur, bg)
			effect.SetOwner(shadow, bg)
			return func() { // in
					anim(blur, "radius", 0, 4, 250*time.Millisecond)
					anim(shadow, "blur", 8, 18, 250*time.Millisecond)
					anim(shadow, "offsetY", 4, 12, 250*time.Millisecond)
				}, func() { // out
					anim(blur, "radius", 4, 0, 200*time.Millisecond)
					anim(shadow, "blur", 18, 8, 200*time.Millisecond)
					anim(shadow, "offsetY", 12, 4, 200*time.Millisecond)
				}, nil, nil
		},
	)
}

func neumorphic() fyne.CanvasObject {
	return makeButton("Neumorphic",
		color.NRGBA{225, 225, 235, 255},
		color.NRGBA{100, 80, 160, 255},
		20, 16,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			shadow := bg.AddEffect(effect.DropShadow, 4, 4, 8, 0.6, 0.6, 0.7, 0.35)
			inShadow := bg.AddEffect(effect.InnerShadow, -3, -3, 6, 1, 1, 1, 0.7)
			effect.SetOwner(shadow, bg)
			effect.SetOwner(inShadow, bg)
			return func() {
					anim(shadow, "blur", 8, 18, 250*time.Millisecond)
					anim(shadow, "offsetX", 4, 7, 250*time.Millisecond)
					anim(shadow, "offsetY", 4, 7, 250*time.Millisecond)
				}, func() {
					anim(shadow, "blur", 18, 8, 200*time.Millisecond)
					anim(shadow, "offsetX", 7, 4, 200*time.Millisecond)
					anim(shadow, "offsetY", 7, 4, 200*time.Millisecond)
				}, func() {
					shadow.SetFloat("blur", 2)
					shadow.SetFloat("offsetX", 1)
					shadow.SetFloat("offsetY", 1)
					inShadow.SetFloat("blur", 10)
					canvas.Refresh(bg)
				}, func() {
					anim(shadow, "blur", 2, 8, 150*time.Millisecond)
					anim(shadow, "offsetX", 1, 4, 150*time.Millisecond)
					anim(shadow, "offsetY", 1, 4, 150*time.Millisecond)
					anim(inShadow, "blur", 10, 6, 150*time.Millisecond)
				}
		},
	)
}

func extrude3D() fyne.CanvasObject {
	return makeButton("3D Extrude",
		color.NRGBA{18, 18, 22, 255},
		color.NRGBA{255, 215, 0, 255},
		8, 16,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			shadow := bg.AddEffect(effect.DropShadow, 0, 6, 0, 0.75, 0.55, 0, 1.0)
			effect.SetOwner(shadow, bg)
			return func() {
					anim(shadow, "offsetY", 6, 10, 150*time.Millisecond)
					anim(shadow, "blur", 0, 3, 150*time.Millisecond)
				}, func() {
					anim(shadow, "offsetY", 10, 6, 150*time.Millisecond)
					anim(shadow, "blur", 3, 0, 150*time.Millisecond)
				}, func() {
					shadow.SetFloat("offsetY", 1)
					shadow.SetFloat("blur", 0)
					canvas.Refresh(bg)
				}, func() {
					anim(shadow, "offsetY", 1, 6, 100*time.Millisecond)
				}
		},
	)
}

func gradientGlow() fyne.CanvasObject {
	return makeButton("Gradient Glow",
		color.NRGBA{100, 40, 220, 255},
		color.White,
		16, 16,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			bg.AddEffect(effect.LinearGradientOverlay,
				0.6, 0.2, 0.9, 0.5, // startColor RGBA
				0.3, 0.0, 0.8, 0.5, // endColor RGBA
				45)                  // angle
			glow := bg.AddEffect(effect.OuterGlow, 0, 0.5, 0.2, 1.0, 0.8)
			bright := bg.AddEffect(effect.Brightness, 1.0)
			effect.SetOwner(glow, bg)
			effect.SetOwner(bright, bg)
			return func() {
					anim(glow, "radius", 0, 14, 300*time.Millisecond)
					anim(bright, "brightness", 1.0, 1.25, 300*time.Millisecond)
				}, func() {
					anim(glow, "radius", 14, 0, 200*time.Millisecond)
					anim(bright, "brightness", 1.25, 1.0, 200*time.Millisecond)
				}, nil, nil
		},
	)
}

func floatingIsland() fyne.CanvasObject {
	return makeButton("Floating Island",
		color.NRGBA{37, 99, 235, 255},
		color.White,
		24, 20,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			shadow := bg.AddEffect(effect.DropShadow, 0, 6, 10, 0.14, 0.38, 0.9, 0.4)
			effect.SetOwner(shadow, bg)
			return func() {
					anim(shadow, "blur", 10, 24, 300*time.Millisecond)
					anim(shadow, "offsetY", 6, 16, 300*time.Millisecond)
				}, func() {
					anim(shadow, "blur", 24, 10, 200*time.Millisecond)
					anim(shadow, "offsetY", 16, 6, 200*time.Millisecond)
				}, nil, nil
		},
	)
}

func outlineFill() fyne.CanvasObject {
	bgCol := color.NRGBA{100, 40, 220, 255}
	return makeButton("Outline Fill",
		color.NRGBA{100, 40, 220, 30},
		color.NRGBA{100, 40, 220, 255},
		8, 8,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			sat := bg.AddEffect(effect.Saturate, 1.0)
			bright := bg.AddEffect(effect.Brightness, 1.0)
			effect.SetOwner(sat, bg)
			effect.SetOwner(bright, bg)
			return func() {
					anim(bright, "brightness", 1.0, 8.0, 200*time.Millisecond)
					lbl.Color = color.White
					canvas.Refresh(lbl)
				}, func() {
					anim(bright, "brightness", 8.0, 1.0, 160*time.Millisecond)
					lbl.Color = bgCol
					canvas.Refresh(lbl)
				}, nil, nil
		},
	)
}

func neonGlow() fyne.CanvasObject {
	return makeButton("Neon Glow",
		color.NRGBA{10, 10, 25, 255},
		color.NRGBA{0, 255, 180, 255},
		10, 20,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			glow := bg.AddEffect(effect.OuterGlow, 5, 0.0, 1.0, 0.7, 0.9)
			bright := bg.AddEffect(effect.Brightness, 1.0)
			effect.SetOwner(glow, bg)
			effect.SetOwner(bright, bg)
			return func() {
					anim(glow, "radius", 5, 20, 300*time.Millisecond)
					anim(bright, "brightness", 1.0, 1.4, 300*time.Millisecond)
				}, func() {
					anim(glow, "radius", 20, 5, 200*time.Millisecond)
					anim(bright, "brightness", 1.4, 1.0, 200*time.Millisecond)
				}, nil, nil
		},
	)
}

func retroCRT() fyne.CanvasObject {
	return makeButton("RETRO CRT",
		color.NRGBA{8, 24, 8, 255},
		color.NRGBA{0, 255, 0, 255},
		4, 8,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			scan := bg.AddEffect(effect.Scanlines, 1.0, 0.15)
			crt := bg.AddEffect(effect.CRT, 0.0, 0.0, 0.0)
			grain := bg.AddEffect(effect.FilmGrain, 0.0, 0)
			bright := bg.AddEffect(effect.Brightness, 1.0)
			effect.SetOwner(scan, bg)
			effect.SetOwner(crt, bg)
			effect.SetOwner(grain, bg)
			effect.SetOwner(bright, bg)
			return func() {
					anim(scan, "opacity", 0.15, 0.35, 200*time.Millisecond)
					anim(crt, "curvature", 0.0, 0.02, 200*time.Millisecond)
					anim(crt, "scanlineIntensity", 0.0, 0.35, 200*time.Millisecond)
					anim(crt, "vignetteStrength", 0.0, 0.3, 200*time.Millisecond)
					anim(grain, "intensity", 0.0, 0.12, 200*time.Millisecond)
					anim(bright, "brightness", 1.0, 1.5, 200*time.Millisecond)
				}, func() {
					anim(scan, "opacity", 0.35, 0.15, 160*time.Millisecond)
					anim(crt, "curvature", 0.02, 0.0, 160*time.Millisecond)
					anim(crt, "scanlineIntensity", 0.35, 0.0, 160*time.Millisecond)
					anim(crt, "vignetteStrength", 0.3, 0.0, 160*time.Millisecond)
					anim(grain, "intensity", 0.12, 0.0, 160*time.Millisecond)
					anim(bright, "brightness", 1.5, 1.0, 160*time.Millisecond)
				}, nil, nil
		},
	)
}

func glitchBtn() fyne.CanvasObject {
	return makeButton("Glitch",
		color.NRGBA{20, 0, 40, 255},
		color.NRGBA{255, 60, 180, 255},
		6, 16,
		func(bg *canvas.Rectangle, lbl *canvas.Text) (func(), func(), func(), func()) {
			chroma := bg.AddEffect(effect.ChromaticAberration, 0)
			scan := bg.AddEffect(effect.Scanlines, 1.0, 0.1)
			effect.SetOwner(chroma, bg)
			effect.SetOwner(scan, bg)
			return func() {
					anim(chroma, "offset", 0, 8, 80*time.Millisecond)
					anim(scan, "opacity", 0.1, 0.4, 80*time.Millisecond)
				}, func() {
					anim(chroma, "offset", 8, 0, 120*time.Millisecond)
					anim(scan, "opacity", 0.4, 0.1, 120*time.Millisecond)
				}, nil, nil
		},
	)
}

// ── main ─────────────────────────────────────────────────────────────────────

func main() {
	a := app.New()
	w := a.NewWindow("Button Effects — Fyne Showcase")
	w.Resize(fyne.NewSize(900, 780))

	titleText := canvas.NewText("Button Effects — Hover to interact", color.White)
	titleText.TextSize = 20
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	hint := canvas.NewText("Each button uses raw canvas.Rectangle + AddEffect — no custom widget", color.NRGBA{150, 150, 150, 255})
	hint.TextSize = 11

	header := container.NewVBox(
		container.NewCenter(titleText),
		container.NewCenter(hint),
		widget.NewSeparator(),
	)

	content := container.NewVBox(
		header,
		row(floatingGlass(),
			`bg.AddEffect(effect.GaussianBlur, 0)`,
			"Blur + deepening drop shadow on hover\n(CSS: backdrop-filter + box-shadow)"),
		widget.NewSeparator(),
		row(neumorphic(),
			`bg.AddEffect(effect.DropShadow, ...)\nbg.AddEffect(effect.InnerShadow, ...)`,
			"Outer shadow lifts, inner shadow presses\n(CSS: neumorphism box-shadow pattern)"),
		widget.NewSeparator(),
		row(extrude3D(),
			`bg.AddEffect(effect.DropShadow, 0, 6, 0, ...)`,
			"Hard zero-blur shadow = 3D depth\nShrinks on press\n(CSS: box-shadow + translateY)"),
		widget.NewSeparator(),
		row(gradientGlow(),
			`bg.AddEffect(effect.LinearGradientOverlay, ...)\nbg.AddEffect(effect.OuterGlow, ...)`,
			"Gradient overlay + expanding outer glow\n(CSS: linear-gradient + filter:brightness)"),
		widget.NewSeparator(),
		row(floatingIsland(),
			`bg.AddEffect(effect.DropShadow, 0, 6, 10, ...)`,
			"Coloured shadow deepens and rises on hover\n(CSS: box-shadow + translateY)"),
		widget.NewSeparator(),
		row(outlineFill(),
			`bg.AddEffect(effect.Brightness, 1.0)`,
			"Near-transparent bg floods with brightness\n(CSS: background-color transition)"),
		widget.NewSeparator(),
		row(neonGlow(),
			`bg.AddEffect(effect.OuterGlow, 5, ...)`,
			"Expanding neon outer glow on hover\n(beyond CSS — needs WebGL in browser)"),
		widget.NewSeparator(),
		row(retroCRT(),
			`bg.AddEffect(effect.CRT, ...)\nbg.AddEffect(effect.Scanlines, ...)`,
			"Scanlines + CRT curve + film grain on hover\n(beyond CSS — needs WebGL in browser)"),
		widget.NewSeparator(),
		row(glitchBtn(),
			`bg.AddEffect(effect.ChromaticAberration, 0)`,
			"Chromatic aberration snaps in on hover\n(beyond CSS — GPU only)"),
	)

	w.SetContent(container.NewScroll(content))
	w.ShowAndRun()
}
