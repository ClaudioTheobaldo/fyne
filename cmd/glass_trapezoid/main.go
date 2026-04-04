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

// Custom perspective+glass fragment shader.
// Creates a 3D trapezoid by applying perspective UV distortion,
// then adds a gloss stripe and edge highlight for glass appearance.
const glassTrapezoidFrag = `#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform vec2 resolution;

uniform float perspective;    // 0=flat, 0.3=moderate 3D, 0.6=extreme
uniform float glossY;         // vertical position of gloss highlight (0-1)
uniform float glossIntensity; // brightness of the gloss stripe
uniform float edgeFade;       // how much edges darken (depth cue)
uniform float refraction;     // chromatic split amount for glass refraction

varying vec2 fragTexCoord;

void main() {
    vec2 uv = fragTexCoord;

    // --- Perspective trapezoid distortion ---
    // Top edge stays wide, bottom narrows (as if tilted away from viewer)
    float narrowing = mix(1.0, 1.0 - perspective, uv.y);
    float wideningTop = mix(1.0, 1.0 + perspective * 0.15, 1.0 - uv.y);
    float scale = narrowing * wideningTop;
    float newX = 0.5 + (uv.x - 0.5) / max(scale, 0.01);

    // Discard pixels outside the trapezoid
    if (newX < 0.0 || newX > 1.0) {
        gl_FragColor = vec4(0.0);
        return;
    }

    vec2 perspUV = vec2(newX, uv.y);

    // --- Glass refraction: chromatic split through the "glass" ---
    vec2 refractDir = (perspUV - 0.5) * texelSize * refraction;
    float r = texture2D(tex, perspUV + refractDir).r;
    float g = texture2D(tex, perspUV).g;
    float b = texture2D(tex, perspUV - refractDir).b;
    float a = texture2D(tex, perspUV).a;
    vec4 color = vec4(r, g, b, a);

    // --- Gloss highlight stripe (simulates light reflection on glass surface) ---
    float glossDist = abs(uv.y - glossY);
    float gloss = exp(-glossDist * glossDist * 80.0) * glossIntensity;
    color.rgb += gloss;

    // --- Subtle top-edge highlight (glass rim light) ---
    float rimTop = exp(-uv.y * 12.0) * 0.25;
    color.rgb += rimTop;

    // --- Edge darkening for 3D depth ---
    float edgeX = 1.0 - pow(abs(newX - 0.5) * 2.0, 2.0);
    float edgeY = 1.0 - pow(abs(uv.y - 0.3) * 1.5, 2.0);
    float edgeMask = clamp(edgeX * edgeY, 0.0, 1.0);
    color.rgb *= mix(1.0 - edgeFade, 1.0, edgeMask);

    // --- Bottom shadow gradient (ground shadow visible through glass) ---
    float bottomShadow = smoothstep(0.7, 1.0, uv.y) * 0.15;
    color.rgb -= bottomShadow;

    gl_FragColor = clamp(color, 0.0, 1.0);
}
`

func main() {
	a := app.New()
	w := a.NewWindow("3D Glass Trapezoid Button")
	w.Resize(fyne.NewSize(600, 500))

	btn := &trapezoidButton{text: "GLASS BUTTON"}
	btn.ExtendBaseWidget(btn)

	note := canvas.NewText(
		"Custom shader: perspective UV + chromatic refraction + gloss stripe + rim light + edge fade",
		color.NRGBA{130, 130, 130, 255},
	)
	note.TextSize = 10
	note.TextStyle = fyne.TextStyle{Monospace: true}
	note.Alignment = fyne.TextAlignCenter

	w.SetContent(container.NewVBox(
		container.NewCenter(btn),
		widget.NewSeparator(),
		container.NewCenter(note),
	))
	w.ShowAndRun()
}

type trapezoidButton struct {
	widget.BaseWidget
	text string
	bg   *canvas.Rectangle
	lbl  *canvas.Text

	glassEff  *effect.Effect
	shadowEff *effect.Effect
	blurEff   *effect.Effect
}

var _ desktop.Hoverable = (*trapezoidButton)(nil)
var _ desktop.Mouseable = (*trapezoidButton)(nil)
var _ fyne.Tappable = (*trapezoidButton)(nil)

func (b *trapezoidButton) CreateRenderer() fyne.WidgetRenderer {
	// Rich dark blue glass base — opaque enough to see the trapezoid clearly
	b.bg = canvas.NewRectangle(color.NRGBA{15, 50, 120, 255})
	b.bg.CornerRadius = 4

	b.lbl = canvas.NewText(b.text, color.NRGBA{200, 230, 255, 255})
	b.lbl.TextSize = 17
	b.lbl.TextStyle = fyne.TextStyle{Bold: true}
	b.lbl.Alignment = fyne.TextAlignCenter

	// Glass trapezoid via custom shader
	b.glassEff = b.bg.AddCustomEffect(glassTrapezoidFrag, map[string][]float32{
		"perspective":    {0.35},
		"glossY":         {0.25},
		"glossIntensity": {0.6},
		"edgeFade":       {0.3},
		"refraction":     {3.0},
	})
	effect.SetOwner(b.glassEff, b.bg)

	// Shadow beneath for depth
	b.shadowEff = b.bg.AddEffect(effect.DropShadow, 0, 6, 8, 0.05, 0.1, 0.2, 0.5)
	effect.SetOwner(b.shadowEff, b.bg)

	// Subtle blur for frosted glass
	b.blurEff = b.bg.AddEffect(effect.GaussianBlur, 0.5)
	effect.SetOwner(b.blurEff, b.bg)

	return &trapezoidRenderer{btn: b}
}

func (b *trapezoidButton) MinSize() fyne.Size {
	return fyne.NewSize(280, 72)
}

func (b *trapezoidButton) Tapped(_ *fyne.PointEvent) {}

func (b *trapezoidButton) MouseIn(_ *desktop.MouseEvent) {
	dur := 400 * time.Millisecond
	// Tilt more, increase glass effects
	anim(b.glassEff, "perspective", 0.35, 0.50, dur)
	anim(b.glassEff, "glossIntensity", 0.45, 0.75, dur)
	anim(b.glassEff, "refraction", 2.0, 5.0, dur)
	anim(b.glassEff, "edgeFade", 0.25, 0.35, dur)
	// Shadow deepens
	anim(b.shadowEff, "blur", 8, 16, dur)
	anim(b.shadowEff, "offsetY", 6, 12, dur)
	// Slight frosted blur
	anim(b.blurEff, "radius", 0.5, 1.5, dur)
}

func (b *trapezoidButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *trapezoidButton) MouseOut() {
	dur := 300 * time.Millisecond
	anim(b.glassEff, "perspective", 0.50, 0.35, dur)
	anim(b.glassEff, "glossIntensity", 0.75, 0.45, dur)
	anim(b.glassEff, "refraction", 5.0, 2.0, dur)
	anim(b.glassEff, "edgeFade", 0.35, 0.25, dur)
	anim(b.shadowEff, "blur", 16, 8, dur)
	anim(b.shadowEff, "offsetY", 12, 6, dur)
	anim(b.blurEff, "radius", 1.5, 0.5, dur)
}

func (b *trapezoidButton) MouseDown(_ *desktop.MouseEvent) {
	// Press: flatten back toward viewer
	b.glassEff.SetFloat("perspective", 0.15)
	b.glassEff.SetFloat("glossIntensity", 0.3)
	b.shadowEff.SetFloat("blur", 3)
	b.shadowEff.SetFloat("offsetY", 2)
}

func (b *trapezoidButton) MouseUp(_ *desktop.MouseEvent) {
	dur := 200 * time.Millisecond
	anim(b.glassEff, "perspective", 0.15, 0.50, dur)
	anim(b.glassEff, "glossIntensity", 0.3, 0.75, dur)
	anim(b.shadowEff, "blur", 3, 16, dur)
	anim(b.shadowEff, "offsetY", 2, 12, dur)
}

type trapezoidRenderer struct{ btn *trapezoidButton }

func (r *trapezoidRenderer) Layout(size fyne.Size) {
	r.btn.bg.Resize(size)
	textH := r.btn.lbl.TextSize * 1.4
	r.btn.lbl.Resize(fyne.NewSize(size.Width, textH))
	r.btn.lbl.Move(fyne.NewPos(0, (size.Height-textH)/2))
}

func (r *trapezoidRenderer) MinSize() fyne.Size          { return r.btn.MinSize() }
func (r *trapezoidRenderer) Refresh()                     {}
func (r *trapezoidRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.btn.bg, r.btn.lbl} }
func (r *trapezoidRenderer) Destroy()                     {}

func anim(eff *effect.Effect, param string, from, to float32, dur time.Duration) {
	a := fyne.NewAnimation(dur, func(t float32) {
		eff.SetFloat(param, from+(to-from)*t)
	})
	a.Curve = fyne.AnimationEaseOut
	a.Start()
}
