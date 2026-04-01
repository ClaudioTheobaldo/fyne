package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRectangle_AddEffect(t *testing.T) {
	rect := canvas.NewRectangle(color.White)

	e := rect.AddEffect(effect.GaussianBlur, 4.0)
	require.NotNil(t, e)
	assert.Equal(t, effect.GaussianBlur, e.Type())
	assert.True(t, rect.HasEffects())
}

func TestRectangle_StackEffects(t *testing.T) {
	rect := canvas.NewRectangle(color.White)

	blur := rect.AddEffect(effect.GaussianBlur, 4.0)
	brightness := rect.AddEffect(effect.Brightness, 1.2)
	opacity := rect.AddEffect(effect.Opacity, 0.5)

	effects := rect.Effects()
	require.Len(t, effects, 3)
	assert.Equal(t, blur, effects[0])
	assert.Equal(t, brightness, effects[1])
	assert.Equal(t, opacity, effects[2])
}

func TestRectangle_RemoveEffect(t *testing.T) {
	rect := canvas.NewRectangle(color.White)

	blur := rect.AddEffect(effect.GaussianBlur, 4.0)
	_ = rect.AddEffect(effect.Brightness, 1.2)

	rect.RemoveEffect(blur)

	effects := rect.Effects()
	require.Len(t, effects, 1)
	assert.Equal(t, effect.Brightness, effects[0].Type())
}

func TestRectangle_ClearEffects(t *testing.T) {
	rect := canvas.NewRectangle(color.White)

	rect.AddEffect(effect.GaussianBlur, 4.0)
	rect.AddEffect(effect.Brightness, 1.2)
	assert.True(t, rect.HasEffects())

	rect.ClearEffects()

	assert.False(t, rect.HasEffects())
	assert.Empty(t, rect.Effects())
}

func TestRectangle_HasEffects(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	assert.False(t, rect.HasEffects())

	e := rect.AddEffect(effect.Opacity, 0.5)
	assert.True(t, rect.HasEffects())

	rect.RemoveEffect(e)
	assert.False(t, rect.HasEffects())
}

func TestRectangle_DuplicateEffectType(t *testing.T) {
	rect := canvas.NewRectangle(color.White)

	blur1 := rect.AddEffect(effect.GaussianBlur, 4.0)
	blur2 := rect.AddEffect(effect.GaussianBlur, 12.0)

	assert.NotEqual(t, blur1, blur2)
	assert.Len(t, rect.Effects(), 2)

	// Each has its own params
	assert.Equal(t, float32(4.0), blur1.Uniform("radius"))
	assert.Equal(t, float32(12.0), blur2.Uniform("radius"))
}

func TestRectangle_AddCustomEffect(t *testing.T) {
	rect := canvas.NewRectangle(color.White)

	src := `void main() { gl_FragColor = vec4(1.0); }`
	e := rect.AddCustomEffect(src, map[string][]float32{
		"intensity": {0.5},
	})

	require.NotNil(t, e)
	assert.Equal(t, effect.Custom, e.Type())
	assert.Equal(t, src, e.CustomShaderSrc())
	assert.True(t, rect.HasEffects())
}

func TestRectangle_EffectsReturnsCopy(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	rect.AddEffect(effect.Brightness, 1.0)

	effects := rect.Effects()
	effects = append(effects, nil)

	// Original should be unaffected
	assert.Len(t, rect.Effects(), 1)
}

// Test that all primitives embedding baseObject support effects

func TestText_HasEffectSupport(t *testing.T) {
	text := canvas.NewText("Hello", color.White)
	assert.False(t, text.HasEffects())

	e := text.AddEffect(effect.Opacity, 0.5)
	assert.True(t, text.HasEffects())
	assert.NotNil(t, e)
}

func TestImage_HasEffectSupport(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/fyne.png")
	assert.False(t, img.HasEffects())

	e := img.AddEffect(effect.GaussianBlur, 4.0)
	assert.True(t, img.HasEffects())
	assert.NotNil(t, e)
}

func TestShaderRect_HasEffectSupport(t *testing.T) {
	sr := canvas.NewShaderRect(`void main() {}`, color.White)
	assert.False(t, sr.HasEffects())

	e := sr.AddEffect(effect.Brightness, 1.5)
	assert.True(t, sr.HasEffects())
	assert.NotNil(t, e)
}

// Test Circle and Line (which don't embed baseObject)

func TestCircle_HasEffectSupport(t *testing.T) {
	circle := canvas.NewCircle(color.White)
	assert.False(t, circle.HasEffects())

	e := circle.AddEffect(effect.GaussianBlur, 2.0)
	assert.True(t, circle.HasEffects())
	assert.NotNil(t, e)

	circle.RemoveEffect(e)
	assert.False(t, circle.HasEffects())
}

func TestLine_HasEffectSupport(t *testing.T) {
	line := canvas.NewLine(color.White)
	assert.False(t, line.HasEffects())

	e := line.AddEffect(effect.Opacity, 0.8)
	assert.True(t, line.HasEffects())
	assert.NotNil(t, e)

	line.ClearEffects()
	assert.False(t, line.HasEffects())
}

func TestPolygon_HasEffectSupport(t *testing.T) {
	poly := canvas.NewPolygon(6, color.White)
	assert.False(t, poly.HasEffects())

	e := poly.AddEffect(effect.Sepia, 1.0)
	assert.True(t, poly.HasEffects())
	assert.NotNil(t, e)
}

func TestArc_HasEffectSupport(t *testing.T) {
	arc := canvas.NewPieArc(0, 180, color.White)
	assert.False(t, arc.HasEffects())

	e := arc.AddEffect(effect.Invert, 1.0)
	assert.True(t, arc.HasEffects())
	assert.NotNil(t, e)
}

func TestLinearGradient_HasEffectSupport(t *testing.T) {
	grad := canvas.NewLinearGradient(color.White, color.Black, 0)
	assert.False(t, grad.HasEffects())

	e := grad.AddEffect(effect.GaussianBlur, 6.0)
	assert.True(t, grad.HasEffects())
	assert.NotNil(t, e)
}

func TestRadialGradient_HasEffectSupport(t *testing.T) {
	grad := canvas.NewRadialGradient(color.White, color.Black)
	assert.False(t, grad.HasEffects())

	e := grad.AddEffect(effect.Vignette, 0.8, 0.5)
	assert.True(t, grad.HasEffects())
	assert.NotNil(t, e)
}
