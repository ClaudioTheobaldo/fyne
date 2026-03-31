package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/stretchr/testify/assert"
)

const testFragShader = `#version 110
uniform vec4 fill_color;
void main() {
    gl_FragColor = fill_color;
}`

func TestShaderRect_MinSize(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)
	min := rect.MinSize()

	assert.Positive(t, min.Width)
	assert.Positive(t, min.Height)
}

func TestShaderRect_FillColor(t *testing.T) {
	c := color.White
	rect := canvas.NewShaderRect(testFragShader, c)

	assert.Equal(t, c, rect.FillColor)
}

func TestShaderRect_FragmentShader(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)

	assert.Equal(t, testFragShader, rect.FragmentShader)
	assert.Empty(t, rect.FragmentShaderES)
}

func TestShaderRect_Uniforms(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)
	rect.Uniforms = map[string][]float32{
		"intensity": {0.5},
		"offset":    {1.0, 2.0},
		"tint":      {0.1, 0.2, 0.3, 0.4},
	}

	assert.Len(t, rect.Uniforms, 3)
	assert.Equal(t, []float32{0.5}, rect.Uniforms["intensity"])
	assert.Equal(t, []float32{1.0, 2.0}, rect.Uniforms["offset"])
	assert.Equal(t, []float32{0.1, 0.2, 0.3, 0.4}, rect.Uniforms["tint"])
}

func TestShaderRect_Resize(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)
	rect.Resize(fyne.NewSize(100, 50))

	assert.Equal(t, float32(100), rect.Size().Width)
	assert.Equal(t, float32(50), rect.Size().Height)
}

func TestShaderRect_Move(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)
	rect.Move(fyne.NewPos(10, 20))

	assert.Equal(t, float32(10), rect.Position().X)
	assert.Equal(t, float32(20), rect.Position().Y)
}

func TestShaderRect_Visibility(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)
	assert.True(t, rect.Visible())

	rect.Hide()
	assert.False(t, rect.Visible())

	rect.Show()
	assert.True(t, rect.Visible())
}

func TestShaderRect_CornerRadius(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)
	rect.CornerRadius = 8
	rect.TopRightCornerRadius = 4

	assert.Equal(t, float32(8), rect.CornerRadius)
	assert.Equal(t, float32(4), rect.TopRightCornerRadius)
	assert.Equal(t, float32(0), rect.TopLeftCornerRadius)
}

func TestShaderRect_StrokeProperties(t *testing.T) {
	rect := canvas.NewShaderRect(testFragShader, color.Black)
	rect.StrokeColor = color.White
	rect.StrokeWidth = 2.0

	assert.Equal(t, color.White, rect.StrokeColor)
	assert.Equal(t, float32(2.0), rect.StrokeWidth)
}
