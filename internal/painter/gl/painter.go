// Package gl provides a full Fyne render implementation using system OpenGL libraries.
package gl

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/aa"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/theme"
)

// Painter defines the functionality of our OpenGL based renderer
type Painter interface {
	// Init tell a new painter to initialise, usually called after a context is available
	Init()
	// Capture requests that the specified canvas be drawn to an in-memory image
	Capture(fyne.Canvas) image.Image
	// Clear tells our painter to prepare a fresh paint
	Clear()
	// Free is used to indicate that a certain canvas object is no longer needed
	Free(fyne.CanvasObject)
	// Paint a single fyne.CanvasObject but not its children.
	Paint(fyne.CanvasObject, fyne.Position, fyne.Size)
	// SetFrameBufferScale tells us when we have more than 1 framebuffer pixel for each output pixel
	SetFrameBufferScale(float32)
	// SetOutputSize is used to change the resolution of our output viewport
	SetOutputSize(int, int)
	// StartClipping tells us that the following paint actions should be clipped to the specified area.
	StartClipping(fyne.Position, fyne.Size)
	// StopClipping stops clipping paint actions.
	StopClipping()

	// SetAntiAliasingMode configures the active anti-aliasing technique.
	SetAntiAliasingMode(mode aa.Mode)
	// AntiAliasingMode returns the current anti-aliasing mode.
	AntiAliasingMode() aa.Mode
}

// NewPainter creates a new GL based renderer for the provided canvas.
// If it is a master painter it will also initialise OpenGL
func NewPainter(c fyne.Canvas, ctx driver.WithContext) Painter {
	p := &painter{canvas: c, contextProvider: ctx, aaMode: aa.SDF}
	p.SetFrameBufferScale(1.0)
	return p
}

type painter struct {
	canvas                fyne.Canvas
	ctx                   context
	contextProvider       driver.WithContext
	program               ProgramState
	lineProgram           ProgramState
	rectangleProgram      ProgramState
	roundRectangleProgram ProgramState
	polygonProgram        ProgramState
	arcProgram            ProgramState
	yuvPlanarProgram      ProgramState
	yuvaPlanarProgram     ProgramState
	nvSemiplanarProgram   ProgramState
	packedYUV422Program   ProgramState
	grayscaleProgram      ProgramState
	yuvPlanarHibitProgram ProgramState
	nvSemiplanarHibitProgram ProgramState
	texScale              float32
	pixScale              float32 // pre-calculate scale*texScale for each draw
	pboStates             map[*canvas.StreamingImage]*pboState
	rawPBOStates          map[*canvas.StreamingImage]*streamPBOState
	shaderCache           map[string]*ProgramState // cached user shader programs keyed by fragment source
	effectPipe            effectPipeline           // FBO-based effect rendering pipeline
	aaMode            aa.Mode              // active anti-aliasing technique
	clipEnabled       bool                 // true while a scissor clip is active (StartClipping..StopClipping)
}

type ProgramState struct {
	ref        Program
	buff       Buffer
	uniforms   map[string]*UniformState
	attributes map[string]Attribute
}

type UniformState struct {
	ref  Uniform
	prev [4]float32
}

func (p *painter) SetUniform1f(pState ProgramState, name string, v float32) {
	u, ok := pState.uniforms[name]
	if !ok {
		return
	}
	if u.prev[0] == v {
		return
	}
	u.prev[0] = v
	p.ctx.Uniform1f(u.ref, v)
}

func (p *painter) SetUniform2f(pState ProgramState, name string, v0, v1 float32) {
	u, ok := pState.uniforms[name]
	if !ok {
		return
	}
	if u.prev[0] == v0 && u.prev[1] == v1 {
		return
	}
	u.prev[0] = v0
	u.prev[1] = v1
	p.ctx.Uniform2f(u.ref, v0, v1)
}

func (p *painter) SetUniform3f(pState ProgramState, name string, v0, v1, v2 float32) {
	u, ok := pState.uniforms[name]
	if !ok {
		return
	}
	if u.prev[0] == v0 && u.prev[1] == v1 && u.prev[2] == v2 {
		return
	}
	u.prev[0] = v0
	u.prev[1] = v1
	u.prev[2] = v2
	p.ctx.Uniform3f(u.ref, v0, v1, v2)
}

func (p *painter) SetUniform4f(pState ProgramState, name string, v0, v1, v2, v3 float32) {
	u, ok := pState.uniforms[name]
	if !ok {
		return
	}
	if u.prev[0] == v0 && u.prev[1] == v1 && u.prev[2] == v2 && u.prev[3] == v3 {
		return
	}
	u.prev[0] = v0
	u.prev[1] = v1
	u.prev[2] = v2
	u.prev[3] = v3
	p.ctx.Uniform4f(u.ref, v0, v1, v2, v3)
}

func (p *painter) UpdateVertexArray(pState ProgramState, name string, size, stride, offset int) {
	a := pState.attributes[name]

	p.ctx.VertexAttribPointerWithOffset(a, size, float, false, stride*floatSize, offset*floatSize)
	p.logError()
}

// Declare conformity to Painter interface
var _ Painter = (*painter)(nil)

func (p *painter) Clear() {
	r, g, b, a := theme.Color(theme.ColorNameBackground).RGBA()
	p.ctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	p.ctx.Clear(bitColorBuffer | bitDepthBuffer)
	p.logError()
}

func (p *painter) Free(obj fyne.CanvasObject) {
	// Don't free StreamingImage textures/PBO on layout changes — they manage
	// their own frame lifecycle via UpdateFrame/UpdateRawFrame. Freeing them
	// during layout repositioning causes a white flash because the texture
	// and PBO state are destroyed before the next repaint.
	if _, isStreaming := obj.(*canvas.StreamingImage); isStreaming {
		return
	}
	p.freeTexture(obj)
}

func (p *painter) Paint(obj fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if obj.Visible() {
		p.drawObject(obj, pos, frame)
	}
}

// SetAntiAliasingMode configures the active anti-aliasing technique.
func (p *painter) SetAntiAliasingMode(mode aa.Mode) { p.aaMode = mode }

// AntiAliasingMode returns the current anti-aliasing mode.
func (p *painter) AntiAliasingMode() aa.Mode { return p.aaMode }

func (p *painter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
	p.pixScale = p.canvas.Scale() * p.texScale
}

func (p *painter) SetOutputSize(width, height int) {
	p.ctx.Viewport(0, 0, width, height)
	p.logError()
}

func (p *painter) StartClipping(pos fyne.Position, size fyne.Size) {
	x := p.textureScale(pos.X)
	y := p.textureScale(p.canvas.Size().Height - pos.Y - size.Height)
	w := p.textureScale(size.Width)
	h := p.textureScale(size.Height)
	p.ctx.Scissor(int32(x), int32(y), int32(w), int32(h))
	p.ctx.Enable(scissorTest)
	p.clipEnabled = true
	p.logError()
}

func (p *painter) StopClipping() {
	p.ctx.Disable(scissorTest)
	p.clipEnabled = false
	p.logError()
}

func (p *painter) compileShader(source string, shaderType uint32) (Shader, error) {
	shader := p.ctx.CreateShader(shaderType)

	p.ctx.ShaderSource(shader, source)
	p.logError()
	p.ctx.CompileShader(shader)
	p.logError()

	info := p.ctx.GetShaderInfoLog(shader)
	if p.ctx.GetShaderi(shader, compileStatus) == glFalse {
		return noShader, fmt.Errorf("failed to compile OpenGL shader:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE", info, source)
	}

	// The info is probably a null terminated string.
	// An empty info has been seen as "\x00" or "\x00\x00".
	if len(info) > 0 && info != "\x00" && info != "\x00\x00" {
		fmt.Printf("OpenGL shader compilation output:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE\n", info, source)
	}

	return shader, nil
}

func (p *painter) createProgram(shaderFilename string) Program {
	// Why a switch over a filename?
	// Because this allows for a minimal change, once we reach Go 1.16 and use go:embed instead of
	// fyne bundle.
	vertexSrc, fragmentSrc := shaderSourceNamed(shaderFilename)
	if vertexSrc == nil {
		panic("shader not found: " + shaderFilename)
	}

	vertShader, err := p.compileShader(string(vertexSrc), vertexShader)
	if err != nil {
		panic(err)
	}
	fragShader, err := p.compileShader(string(fragmentSrc), fragmentShader)
	if err != nil {
		panic(err)
	}

	prog := p.ctx.CreateProgram()
	p.ctx.AttachShader(prog, vertShader)
	p.ctx.AttachShader(prog, fragShader)
	p.ctx.LinkProgram(prog)

	info := p.ctx.GetProgramInfoLog(prog)
	if p.ctx.GetProgrami(prog, linkStatus) == glFalse {
		panic(fmt.Errorf("failed to link OpenGL program:\n%s", info))
	}

	// The info is probably a null terminated string.
	// An empty info has been seen as "\x00" or "\x00\x00".
	if len(info) > 0 && info != "\x00" && info != "\x00\x00" {
		fmt.Printf("OpenGL program linking output:\n%s\n", info)
	}

	if glErr := p.ctx.GetError(); glErr != 0 {
		panic(fmt.Sprintf("failed to link OpenGL program; error code: %x", glErr))
	}

	p.ctx.UseProgram(prog)

	return prog
}

func (p *painter) logError() {
	logGLError(p.ctx.GetError)
}

// getOrCompileShaderProgram returns a cached ProgramState for the given fragment shader source,
// compiling and linking it on first use. The standard rectangle vertex shader is always used.
// Returns nil if compilation or linking fails.
func (p *painter) getOrCompileShaderProgram(fragSrc string) *ProgramState {
	if ps, ok := p.shaderCache[fragSrc]; ok {
		return ps
	}

	vertSrc := string(shaderRectVertexSrc())

	vertShader, err := p.compileShader(vertSrc, vertexShader)
	if err != nil {
		fyne.LogError("ShaderRect: failed to compile vertex shader", err)
		return nil
	}
	fragShader, err := p.compileShader(fragSrc, fragmentShader)
	if err != nil {
		fyne.LogError("ShaderRect: failed to compile fragment shader", err)
		return nil
	}

	prog := p.ctx.CreateProgram()
	p.ctx.AttachShader(prog, vertShader)
	p.ctx.AttachShader(prog, fragShader)
	p.ctx.LinkProgram(prog)

	if p.ctx.GetProgrami(prog, linkStatus) == glFalse {
		info := p.ctx.GetProgramInfoLog(prog)
		fyne.LogError("ShaderRect: failed to link program", fmt.Errorf("%s", info))
		return nil
	}

	p.ctx.UseProgram(prog)

	ps := &ProgramState{
		ref:        prog,
		buff:       p.createBuffer(16),
		uniforms:   make(map[string]*UniformState),
		attributes: make(map[string]Attribute),
	}

	// Discover standard uniforms (skip those not present in the user's shader)
	standardUniforms := []string{
		"frame_size", "rect_coords",
		"stroke_width_half", "rect_size_half",
		"radius", "edge_softness",
		"fill_color", "stroke_color",
	}
	for _, name := range standardUniforms {
		loc := p.ctx.GetUniformLocation(prog, name)
		if loc >= 0 {
			ps.uniforms[name] = &UniformState{ref: loc}
		}
	}

	// Enable standard vertex attributes
	for _, name := range []string{"vert", "normal"} {
		a := p.ctx.GetAttribLocation(prog, name)
		p.ctx.EnableVertexAttribArray(a)
		ps.attributes[name] = a
	}

	p.shaderCache[fragSrc] = ps
	return ps
}

// discoverCustomUniforms discovers and sets custom uniform locations for a ShaderRect.
// This is called each draw to handle uniforms that weren't in the standard set.
func (p *painter) discoverCustomUniforms(ps *ProgramState, prog Program, uniforms map[string][]float32) {
	for name := range uniforms {
		if _, ok := ps.uniforms[name]; ok {
			continue // already known
		}
		loc := p.ctx.GetUniformLocation(prog, name)
		if loc >= 0 {
			ps.uniforms[name] = &UniformState{ref: loc}
		}
	}
}
