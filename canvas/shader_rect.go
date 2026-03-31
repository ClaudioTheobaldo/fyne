package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*ShaderRect)(nil)

// ShaderRect describes a rectangle primitive that is rendered using a user-supplied
// GLSL fragment shader. It receives the same standard uniforms as the built-in
// round_rectangle shader (frame_size, rect_coords, fill_color, stroke_color,
// stroke_width_half, rect_size_half, radius, edge_softness), and additionally
// allows custom float uniforms via the Uniforms map.
//
// The vertex shader is the standard rectangle vertex shader (vert + normal attributes).
//
// Since: 2.8
type ShaderRect struct {
	baseObject

	FillColor   color.Color // The rectangle fill color, passed as fill_color uniform
	StrokeColor color.Color // The rectangle stroke color, passed as stroke_color uniform
	StrokeWidth float32     // The stroke width of the rectangle

	// The radius of the rectangle corners
	CornerRadius float32

	// Per-corner radius overrides. When non-zero, these take precedence over CornerRadius.
	TopRightCornerRadius    float32
	TopLeftCornerRadius     float32
	BottomRightCornerRadius float32
	BottomLeftCornerRadius  float32

	// FragmentShader is the GLSL fragment shader source for desktop GL (#version 110).
	// The shader can use any of the standard uniforms listed in the type doc,
	// plus any custom uniforms declared in the Uniforms map.
	FragmentShader string

	// FragmentShaderES is the GLSL ES fragment shader source (#version 100 with precision).
	// If empty and the runtime requires ES, the painter will attempt to auto-convert
	// FragmentShader by prepending precision qualifiers and adjusting the version line.
	FragmentShaderES string

	// Uniforms holds user-defined custom uniform values. Each key is a uniform name
	// declared in the fragment shader, and the value is a slice of 1, 2, or 4 float32
	// values (mapped to uniform1f, uniform2f, or uniform4f respectively).
	Uniforms map[string][]float32
}

// Hide will set this shader rect to not be visible
func (s *ShaderRect) Hide() {
	s.baseObject.Hide()
	repaint(s)
}

// Move the shader rect to a new position, relative to its parent / canvas
func (s *ShaderRect) Move(pos fyne.Position) {
	if s.Position() == pos {
		return
	}
	s.baseObject.Move(pos)
	repaint(s)
}

// Refresh causes this shader rect to be redrawn with its configured state.
func (s *ShaderRect) Refresh() {
	Refresh(s)
}

// Resize updates the size of this shader rect.
func (s *ShaderRect) Resize(sz fyne.Size) {
	if sz == s.Size() {
		return
	}
	s.baseObject.Resize(sz)
	Refresh(s)
}

// NewShaderRect returns a new ShaderRect instance with the given fragment shader
// source and fill color.
func NewShaderRect(fragmentShader string, fill color.Color) *ShaderRect {
	return &ShaderRect{
		FillColor:      fill,
		FragmentShader: fragmentShader,
	}
}
