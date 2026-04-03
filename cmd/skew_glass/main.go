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

func main() {
	a := app.New()
	w := a.NewWindow("Skew Glass Hover Effect")
	w.Resize(fyne.NewSize(500, 400))

	btn := &glassButton{}
	btn.text = "Hover Me"
	btn.ExtendBaseWidget(btn)

	note := canvas.NewText("Hover: Skew + Blur + ChromaticAberration + Brightness — animated in/out", color.NRGBA{140, 255, 140, 255})
	note.TextSize = 11
	note.TextStyle = fyne.TextStyle{Monospace: true}

	w.SetContent(container.NewVBox(
		container.NewCenter(btn),
		widget.NewSeparator(),
		container.NewCenter(note),
	))
	w.ShowAndRun()
}

// --- Glass Button Widget ---

type glassButton struct {
	widget.BaseWidget
	text string
	bg   *canvas.Rectangle
	lbl  *canvas.Text

	skewEff   *effect.Effect
	blurEff   *effect.Effect
	chromaEff *effect.Effect
	brightEff *effect.Effect
	shadowEff *effect.Effect
}

var _ desktop.Hoverable = (*glassButton)(nil)
var _ desktop.Mouseable = (*glassButton)(nil)
var _ fyne.Tappable = (*glassButton)(nil)

func (b *glassButton) CreateRenderer() fyne.WidgetRenderer {
	b.bg = canvas.NewRectangle(color.NRGBA{30, 100, 220, 255})
	b.bg.CornerRadius = 14

	b.lbl = canvas.NewText(b.text, color.White)
	b.lbl.TextSize = 18
	b.lbl.TextStyle = fyne.TextStyle{Bold: true}
	b.lbl.Alignment = fyne.TextAlignCenter

	// Idle state: subtle shadow
	b.shadowEff = b.bg.AddEffect(effect.DropShadow, 0, 3, 5, 0.1, 0.2, 0.5, 0.3)
	effect.SetOwner(b.shadowEff, b.bg)

	// Pre-create all effects at 0 intensity so we can animate them
	b.skewEff = b.bg.AddEffect(effect.Skew, 0, 0)
	effect.SetOwner(b.skewEff, b.bg)

	b.blurEff = b.bg.AddEffect(effect.GaussianBlur, 0)
	effect.SetOwner(b.blurEff, b.bg)

	b.chromaEff = b.bg.AddEffect(effect.ChromaticAberration, 0)
	effect.SetOwner(b.chromaEff, b.bg)

	b.brightEff = b.bg.AddEffect(effect.Brightness, 1.0)
	effect.SetOwner(b.brightEff, b.bg)

	return &glassRenderer{btn: b}
}

func (b *glassButton) MinSize() fyne.Size {
	return fyne.NewSize(240, 64)
}

func (b *glassButton) Tapped(_ *fyne.PointEvent) {}

func (b *glassButton) MouseIn(_ *desktop.MouseEvent) {
	dur := 350 * time.Millisecond

	// Skew rightwards-upwards: positive skewX tilts right, negative skewY lifts top-right
	// tan(30°) ≈ 0.577
	anim(b.skewEff, "skewX", 0, 0.55, dur)
	anim(b.skewEff, "skewY", 0, -0.15, dur)

	// Glass distortion: light blur + chromatic split
	anim(b.blurEff, "radius", 0, 1.5, dur)
	anim(b.chromaEff, "offset", 0, 4.0, dur)

	// Brighten as if light is hitting the surface
	anim(b.brightEff, "brightness", 1.0, 1.25, dur)

	// Shadow deepens and shifts with the skew
	anim(b.shadowEff, "blur", 5, 14, dur)
	anim(b.shadowEff, "offsetX", 0, -4, dur)
	anim(b.shadowEff, "offsetY", 3, 8, dur)
}

func (b *glassButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *glassButton) MouseOut() {
	dur := 250 * time.Millisecond

	anim(b.skewEff, "skewX", 0.55, 0, dur)
	anim(b.skewEff, "skewY", -0.15, 0, dur)

	anim(b.blurEff, "radius", 1.5, 0, dur)
	anim(b.chromaEff, "offset", 4.0, 0, dur)

	anim(b.brightEff, "brightness", 1.25, 1.0, dur)

	anim(b.shadowEff, "blur", 14, 5, dur)
	anim(b.shadowEff, "offsetX", -4, 0, dur)
	anim(b.shadowEff, "offsetY", 8, 3, dur)
}

func (b *glassButton) MouseDown(_ *desktop.MouseEvent) {
	// Press: snap flat
	b.skewEff.SetFloat("skewX", 0.1)
	b.skewEff.SetFloat("skewY", 0)
	b.brightEff.SetFloat("brightness", 0.9)
	b.shadowEff.SetFloat("blur", 2)
	b.shadowEff.SetFloat("offsetY", 1)
}

func (b *glassButton) MouseUp(_ *desktop.MouseEvent) {
	dur := 200 * time.Millisecond
	anim(b.skewEff, "skewX", 0.1, 0.55, dur)
	anim(b.skewEff, "skewY", 0, -0.15, dur)
	anim(b.brightEff, "brightness", 0.9, 1.25, dur)
	anim(b.shadowEff, "blur", 2, 14, dur)
	anim(b.shadowEff, "offsetY", 1, 8, dur)
}

type glassRenderer struct {
	btn *glassButton
}

func (r *glassRenderer) Layout(size fyne.Size) {
	r.btn.bg.Resize(size)
	r.btn.lbl.Resize(size)
	h := r.btn.lbl.MinSize().Height
	r.btn.lbl.Move(fyne.NewPos(0, (size.Height-h)/2))
}

func (r *glassRenderer) MinSize() fyne.Size        { return r.btn.MinSize() }
func (r *glassRenderer) Refresh()                   {}
func (r *glassRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.btn.bg, r.btn.lbl} }
func (r *glassRenderer) Destroy()                   {}

func anim(eff *effect.Effect, param string, from, to float32, dur time.Duration) {
	a := fyne.NewAnimation(dur, func(t float32) {
		eff.SetFloat(param, from+(to-from)*t)
	})
	a.Curve = fyne.AnimationEaseOut
	a.Start()
}
