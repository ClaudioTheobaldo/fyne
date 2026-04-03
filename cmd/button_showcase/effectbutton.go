package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Compile-time interface checks
var _ fyne.Tappable = (*EffectButton)(nil)
var _ desktop.Hoverable = (*EffectButton)(nil)
var _ desktop.Mouseable = (*EffectButton)(nil)

// EffectButton is a custom button widget that supports shader effects
// on its background rectangle, with hover and press state transitions.
type EffectButton struct {
	widget.BaseWidget

	Text      string
	TextColor color.Color
	TextSize  float32
	BgColor   color.Color
	MinWidth  float32
	MinHeight float32
	Radius    float32

	OnTapped   func()
	OnHoverIn  func(bg *canvas.Rectangle, label *canvas.Text)
	OnHoverOut func(bg *canvas.Rectangle, label *canvas.Text)
	OnPress    func(bg *canvas.Rectangle, label *canvas.Text)
	OnRelease  func(bg *canvas.Rectangle, label *canvas.Text)
	// OnInit is called once after CreateRenderer to apply idle-state effects
	OnInit func(bg *canvas.Rectangle, label *canvas.Text)

	hovered bool
	pressed bool
	bg      *canvas.Rectangle
	label   *canvas.Text
}

func NewEffectButton(text string, bgColor, textColor color.Color, tapped func()) *EffectButton {
	b := &EffectButton{
		Text:      text,
		TextColor: textColor,
		BgColor:   bgColor,
		TextSize:  14,
		MinWidth:  180,
		MinHeight: 50,
		Radius:    12,
		OnTapped:  tapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *EffectButton) CreateRenderer() fyne.WidgetRenderer {
	b.bg = canvas.NewRectangle(b.BgColor)
	b.bg.CornerRadius = b.Radius

	b.label = canvas.NewText(b.Text, b.TextColor)
	b.label.TextSize = b.TextSize
	b.label.TextStyle = fyne.TextStyle{Bold: true}
	b.label.Alignment = fyne.TextAlignCenter

	if b.OnInit != nil {
		b.OnInit(b.bg, b.label)
	}

	return &effectButtonRenderer{btn: b, bg: b.bg, label: b.label}
}

func (b *EffectButton) MinSize() fyne.Size {
	return fyne.NewSize(b.MinWidth, b.MinHeight)
}

func (b *EffectButton) Tapped(_ *fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *EffectButton) MouseIn(_ *desktop.MouseEvent) {
	b.hovered = true
	if b.OnHoverIn != nil && b.bg != nil {
		b.OnHoverIn(b.bg, b.label)
	}
}

func (b *EffectButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *EffectButton) MouseOut() {
	b.hovered = false
	if b.OnHoverOut != nil && b.bg != nil {
		b.OnHoverOut(b.bg, b.label)
	}
}

func (b *EffectButton) MouseDown(_ *desktop.MouseEvent) {
	b.pressed = true
	if b.OnPress != nil && b.bg != nil {
		b.OnPress(b.bg, b.label)
	}
}

func (b *EffectButton) MouseUp(_ *desktop.MouseEvent) {
	b.pressed = false
	if b.OnRelease != nil && b.bg != nil {
		b.OnRelease(b.bg, b.label)
	}
}

type effectButtonRenderer struct {
	btn   *EffectButton
	bg    *canvas.Rectangle
	label *canvas.Text
}

func (r *effectButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	textH := r.label.TextSize * 1.4
	r.label.Resize(fyne.NewSize(size.Width, textH))
	r.label.Move(fyne.NewPos(0, (size.Height-textH)/2))
}

func (r *effectButtonRenderer) MinSize() fyne.Size {
	return r.btn.MinSize()
}

func (r *effectButtonRenderer) Refresh() {
	r.bg.FillColor = r.btn.BgColor
	r.bg.CornerRadius = r.btn.Radius
	r.label.Text = r.btn.Text
	r.label.Color = r.btn.TextColor
	r.label.TextSize = r.btn.TextSize
	canvas.Refresh(r.bg)
	canvas.Refresh(r.label)
}

func (r *effectButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.label}
}

func (r *effectButtonRenderer) Destroy() {}

// animateEffect smoothly transitions a float uniform
func animateEffect(eff *effect.Effect, param string, from, to float32, duration time.Duration) {
	anim := fyne.NewAnimation(duration, func(t float32) {
		eff.SetFloat(param, from+(to-from)*t)
	})
	anim.Curve = fyne.AnimationEaseOut
	anim.Start()
}
