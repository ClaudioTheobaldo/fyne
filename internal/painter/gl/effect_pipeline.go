package gl

import (
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
)

// effectPipeline manages the FBO ping-pong rendering pipeline for shader effects.
type effectPipeline struct {
	fboA, fboB    Framebuffer
	texA, texB    Texture
	width, height int
	quadBuf       Buffer                             // fullscreen quad vertex buffer
	passthrough   ProgramState                       // passthrough shader for compositing final result
	progCache     map[effect.EffectType]*ProgramState // compiled effect shader programs
	customCache   map[string]*ProgramState            // compiled custom shader programs keyed by source
	initialized   bool
	allocated     bool // whether FBOs have been allocated (avoids struct comparison on WASM)

	// MSAA support
	msaaFBO     Framebuffer
	msaaRBO     Renderbuffer // multisample color renderbuffer
	msaaSamples int
}

// fullscreen quad: 2 triangles covering [-1,1] with matching UV [0,1]
// Layout: x, y, u, v (4 floats per vertex, 4 vertices as triangle strip)
var effectQuadVerts = []float32{
	-1, -1, 0, 0, // bottom-left
	1, -1, 1, 0,  // bottom-right
	-1, 1, 0, 1,  // top-left
	1, 1, 1, 1,   // top-right
}

func (p *painter) initEffectPipeline() {
	ep := &p.effectPipe

	// Create fullscreen quad VBO
	ep.quadBuf = p.createBuffer(len(effectQuadVerts))
	p.updateBuffer(ep.quadBuf, effectQuadVerts)

	// Compile passthrough shader
	vertSrc := string(effectPassthroughVertSrc())
	fragSrc := string(effectPassthroughFragSrc())

	vertShader, err := p.compileShader(vertSrc, vertexShader)
	if err != nil {
		fyne.LogError("EffectPipeline: failed to compile passthrough vertex shader", err)
		return
	}
	fragShader, err := p.compileShader(fragSrc, fragmentShader)
	if err != nil {
		fyne.LogError("EffectPipeline: failed to compile passthrough fragment shader", err)
		return
	}

	prog := p.ctx.CreateProgram()
	p.ctx.AttachShader(prog, vertShader)
	p.ctx.AttachShader(prog, fragShader)
	p.ctx.LinkProgram(prog)

	if p.ctx.GetProgrami(prog, linkStatus) == glFalse {
		info := p.ctx.GetProgramInfoLog(prog)
		fyne.LogError("EffectPipeline: failed to link passthrough program", fmt.Errorf("%s", info))
		return
	}

	p.ctx.UseProgram(prog)

	ep.passthrough = ProgramState{
		ref:        prog,
		buff:       ep.quadBuf,
		uniforms:   make(map[string]*UniformState),
		attributes: make(map[string]Attribute),
	}

	// Discover passthrough uniforms
	p.discoverEffectUniforms(&ep.passthrough, prog, []string{"tex"})

	// Enable vertex attributes
	for _, name := range []string{"vert", "vertTexCoord"} {
		a := p.ctx.GetAttribLocation(prog, name)
		p.ctx.EnableVertexAttribArray(a)
		ep.passthrough.attributes[name] = a
	}

	ep.progCache = make(map[effect.EffectType]*ProgramState)
	ep.customCache = make(map[string]*ProgramState)
	ep.initialized = true
}

// ensureFBOSize creates or resizes the two FBOs and their color textures.
func (p *painter) ensureFBOSize(w, h int) {
	ep := &p.effectPipe
	if ep.width == w && ep.height == h && ep.allocated {
		return
	}

	// Clean up old resources
	if ep.allocated {
		p.ctx.DeleteFramebuffer(ep.fboA)
		p.ctx.DeleteFramebuffer(ep.fboB)
		p.ctx.DeleteTexture(ep.texA)
		p.ctx.DeleteTexture(ep.texB)
		if ep.msaaSamples > 0 {
			p.ctx.DeleteFramebuffer(ep.msaaFBO)
			p.ctx.DeleteRenderbuffer(ep.msaaRBO)
		}
	}

	// Create two RGBA textures
	ep.texA = p.newFBOTexture(w, h)
	ep.texB = p.newFBOTexture(w, h)

	// Create and setup FBO A
	ep.fboA = p.ctx.CreateFramebuffer()
	p.ctx.BindFramebuffer(glFramebuffer, ep.fboA)
	p.ctx.FramebufferTexture2D(glFramebuffer, colorAttachment0, texture2D, ep.texA, 0)
	status := p.ctx.CheckFramebufferStatus(glFramebuffer)
	if status != framebufferComplete {
		fyne.LogError("EffectPipeline: FBO A incomplete", fmt.Errorf("status: 0x%x", status))
	}

	// Create and setup FBO B
	ep.fboB = p.ctx.CreateFramebuffer()
	p.ctx.BindFramebuffer(glFramebuffer, ep.fboB)
	p.ctx.FramebufferTexture2D(glFramebuffer, colorAttachment0, texture2D, ep.texB, 0)
	status = p.ctx.CheckFramebufferStatus(glFramebuffer)
	if status != framebufferComplete {
		fyne.LogError("EffectPipeline: FBO B incomplete", fmt.Errorf("status: 0x%x", status))
	}

	// Create MSAA FBO if hardware multisampling is active
	if p.aaMode.NeedsMSAA() {
		samples := int32(p.aaMode.Samples())
		ep.msaaRBO = p.ctx.CreateRenderbuffer()
		p.ctx.BindRenderbuffer(glRenderbuffer, ep.msaaRBO)
		p.ctx.RenderbufferStorageMultisample(glRenderbuffer, samples, colorFormatRGBA, int32(w), int32(h))

		ep.msaaFBO = p.ctx.CreateFramebuffer()
		p.ctx.BindFramebuffer(glFramebuffer, ep.msaaFBO)
		p.ctx.FramebufferRenderbuffer(glFramebuffer, colorAttachment0, glRenderbuffer, ep.msaaRBO)
		status = p.ctx.CheckFramebufferStatus(glFramebuffer)
		if status != framebufferComplete {
			fyne.LogError("EffectPipeline: MSAA FBO incomplete", fmt.Errorf("status: 0x%x", status))
		}
		ep.msaaSamples = int(samples)
	} else {
		ep.msaaSamples = 0
	}

	// Unbind FBO (bind default framebuffer 0)
	p.bindDefaultFramebuffer()

	ep.width = w
	ep.height = h
	ep.allocated = true
}

// bindDefaultFramebuffer binds framebuffer 0 (the screen) in a cross-platform way.
func (p *painter) bindDefaultFramebuffer() {
	p.ctx.BindFramebuffer(glFramebuffer, noFramebuffer)
}

// discoverEffectUniforms finds uniform locations and adds them to the program state.
// This avoids direct comparison of Uniform values against 0 which doesn't compile on WASM.
func (p *painter) discoverEffectUniforms(ps *ProgramState, prog Program, names []string) {
	for _, name := range names {
		loc := p.ctx.GetUniformLocation(prog, name)
		// Store unconditionally — the SetUniform* methods use ok-check on the map
		ps.uniforms[name] = &UniformState{ref: loc}
	}
}

// newFBOTexture creates an empty RGBA texture suitable for FBO attachment.
func (p *painter) newFBOTexture(w, h int) Texture {
	tex := p.ctx.CreateTexture()
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, tex)
	p.ctx.TexParameteri(texture2D, textureMinFilter, textureFilterToGL[0]) // LINEAR
	p.ctx.TexParameteri(texture2D, textureMagFilter, textureFilterToGL[0])
	p.ctx.TexParameteri(texture2D, textureWrapS, clampToEdge)
	p.ctx.TexParameteri(texture2D, textureWrapT, clampToEdge)
	// Allocate empty texture (must pass allocated slice, not nil, for go-gl/gl Ptr())
	emptyData := make([]uint8, w*h*4)
	p.ctx.TexImage2D(texture2D, 0, w, h, colorFormatRGBA, unsignedByte, emptyData)
	p.logError()
	return tex
}

// getEffectProgram returns a compiled and cached ProgramState for the given effect type.
// Returns the passthrough program if no specific shader exists.
func (p *painter) getEffectProgram(eff *effect.Effect) *ProgramState {
	ep := &p.effectPipe

	if eff.Type() == effect.Custom {
		src := eff.CustomShaderSrc()
		if src == "" {
			return &ep.passthrough
		}
		if ps, ok := ep.customCache[src]; ok {
			return ps
		}
		ps := p.compileEffectProgram(src)
		if ps == nil {
			return &ep.passthrough
		}
		// Discover custom uniforms from the effect's own uniform map
		// so SetUniform* calls can find them at draw time.
		uniforms := eff.Uniforms()
		customNames := make([]string, 0, len(uniforms))
		for name := range uniforms {
			if _, exists := ps.uniforms[name]; !exists {
				customNames = append(customNames, name)
			}
		}
		if len(customNames) > 0 {
			p.ctx.UseProgram(ps.ref)
			p.discoverEffectUniforms(ps, ps.ref, customNames)
		}
		ep.customCache[src] = ps
		return ps
	}

	kind := eff.Type()
	if ps, ok := ep.progCache[kind]; ok {
		return ps
	}

	fragSrc := effectShaderSource(kind)
	if fragSrc == nil {
		return &ep.passthrough
	}

	ps := p.compileEffectProgram(string(fragSrc))
	if ps == nil {
		return &ep.passthrough
	}
	ep.progCache[kind] = ps
	return ps
}

// compileEffectProgram compiles a fragment shader with the effect passthrough vertex shader.
func (p *painter) compileEffectProgram(fragSrc string) *ProgramState {
	vertSrc := string(effectPassthroughVertSrc())

	vertShader, err := p.compileShader(vertSrc, vertexShader)
	if err != nil {
		fyne.LogError("EffectPipeline: failed to compile effect vertex shader", err)
		return nil
	}
	fragShader, err := p.compileShader(fragSrc, fragmentShader)
	if err != nil {
		fyne.LogError("EffectPipeline: failed to compile effect fragment shader", err)
		return nil
	}

	prog := p.ctx.CreateProgram()
	p.ctx.AttachShader(prog, vertShader)
	p.ctx.AttachShader(prog, fragShader)
	p.ctx.LinkProgram(prog)

	if p.ctx.GetProgrami(prog, linkStatus) == glFalse {
		info := p.ctx.GetProgramInfoLog(prog)
		fyne.LogError("EffectPipeline: failed to link effect program", fmt.Errorf("%s", info))
		return nil
	}

	p.ctx.UseProgram(prog)

	ps := &ProgramState{
		ref:        prog,
		buff:       p.effectPipe.quadBuf,
		uniforms:   make(map[string]*UniformState),
		attributes: make(map[string]Attribute),
	}

	// Discover all possible uniforms used by any built-in effect shader
	standardUniforms := []string{
		"tex", "resolution", "texelSize",
		// Color
		"brightness", "contrast", "saturation", "amount", "angle", "opacity",
		"colorMatrix", "matrixOffset", "levels", "gamma", "vibrance", "temperature",
		// Blur
		"radius", "direction", "sigmaSpace", "sigmaColor",
		// Distortion
		"pixelSize", "offset", "strength", "center", "distortion",
		"amplitude", "frequency", "speed", "noiseScale",
		"focusY", "focusWidth", "skewX", "skewY",
		"scaleX", "scaleY", "flipX", "flipY", "segments", "rotation",
		"curl", "phase",
		// Compositing
		"sourceColor", "targetColor", "tolerance",
		"darkColor", "lightColor", "shadowTint", "highlightTint", "balance",
		"redOut", "greenOut", "blueOut",
		"color0", "color1", "color2", "color3", "color4",
		"patternType", "time",
		"brightMin", "brightMax", "scaleMin", "scaleMax",
		"blockSize", "width",
		// Detail
		"threshold",
		// Shadow
		"offsetX", "offsetY", "blur", "spread", "shadowColor", "glowColor",
		// Stylization
		"intensity", "smoothness", "density", "dotSize", "spacing",
		"curvature", "scanlineIntensity", "vignetteStrength",
		// Lighting
		"lightPos", "diffuseConstant", "specularConstant", "specularExponent", "surfaceScale",
		"color",
		// Gradient overlay
		"startColor", "endColor",
		// Mask
		"startAlpha", "endAlpha",
		// Procedural
		"baseFrequency", "numOctaves", "seed",
		// Blend
		"blendOpacity",
		// Advanced
		"exposure",
	}
	p.discoverEffectUniforms(ps, prog, standardUniforms)

	// Enable vertex attributes
	for _, name := range []string{"vert", "vertTexCoord"} {
		a := p.ctx.GetAttribLocation(prog, name)
		p.ctx.EnableVertexAttribArray(a)
		ps.attributes[name] = a
	}

	return ps
}

// effectHolder is the interface we check for on canvas objects.
type effectHolder interface {
	HasEffects() bool
	Effects() []*effect.Effect
}

// drawObjectWithEffects renders a canvas object with its effect chain.
func (p *painter) drawObjectWithEffects(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	ep := &p.effectPipe
	if !ep.initialized {
		// Fallback: draw without effects
		p.drawObjectDirect(o, pos, frame)
		return
	}

	holder := o.(effectHolder)
	effects := holder.Effects()

	// Filter to enabled effects only
	var active []*effect.Effect
	for _, e := range effects {
		if e.Enabled() {
			active = append(active, e)
		}
	}
	if len(active) == 0 {
		p.drawObjectDirect(o, pos, frame)
		return
	}

	// Set owners on effects (lazy, done once per draw)
	for _, e := range active {
		effect.SetOwner(e, o)
	}

	// Compute FBO dimensions in pixels
	objSize := o.Size()
	w := int(math.Ceil(float64(objSize.Width * p.pixScale)))
	h := int(math.Ceil(float64(objSize.Height * p.pixScale)))
	if w <= 0 || h <= 0 {
		return
	}

	p.ensureFBOSize(w, h)

	// If this object sits inside a clip region (e.g. a Scrollable), a GL scissor
	// is currently enabled in WINDOW coordinates. The offscreen passes below draw
	// into FBOs with a (0,0)-origin viewport, so that window-space scissor rect
	// would fall entirely outside the FBO and clip the whole render away (black
	// output). Disable it for the offscreen work and restore it for the on-screen
	// composite (Step 3), which DOES want to be clipped to the tile.
	if p.clipEnabled {
		p.ctx.Disable(scissorTest)
	}

	// Save viewport state
	// (Fyne's viewport is set once per frame, we restore it after)

	// Step 1: Render object into FBO A (via MSAA FBO + resolve if active)
	fboFrame := fyne.NewSize(objSize.Width, objSize.Height)
	if ep.msaaSamples > 0 {
		// Render into multisample FBO
		p.ctx.BindFramebuffer(glFramebuffer, ep.msaaFBO)
		p.ctx.Viewport(0, 0, w, h)
		p.ctx.ClearColor(0, 0, 0, 0)
		p.ctx.Clear(bitColorBuffer)
		p.drawObjectDirect(o, fyne.NewPos(0, 0), fboFrame)

		// Resolve MSAA into fboA (texture-backed) via blit
		p.ctx.BindFramebuffer(glReadFramebuffer, ep.msaaFBO)
		p.ctx.BindFramebuffer(glDrawFramebuffer, ep.fboA)
		p.ctx.BlitFramebuffer(0, 0, int32(w), int32(h), 0, 0, int32(w), int32(h), bitColorBuffer, glLinear)
		p.ctx.BindFramebuffer(glFramebuffer, noFramebuffer)
	} else {
		// No MSAA: render directly into fboA
		p.ctx.BindFramebuffer(glFramebuffer, ep.fboA)
		p.ctx.Viewport(0, 0, w, h)
		p.ctx.ClearColor(0, 0, 0, 0)
		p.ctx.Clear(bitColorBuffer)
		p.drawObjectDirect(o, fyne.NewPos(0, 0), fboFrame)
	}

	// Step 2: Ping-pong through effect chain
	srcTex := ep.texA
	dstTex := ep.texB
	// useB tracks which FBO is the current destination (true=fboB, false=fboA)
	useB := true

	resW := float32(w)
	resH := float32(h)
	texelW := 1.0 / resW
	texelH := 1.0 / resH

	for _, eff := range active {
		ps := p.getEffectProgram(eff)
		passes := effect.PassCount(eff.Type())

		for pass := 0; pass < passes; pass++ {
			dstFBO := ep.fboA
			if useB {
				dstFBO = ep.fboB
			}

			p.ctx.BindFramebuffer(glFramebuffer, dstFBO)
			p.ctx.Viewport(0, 0, w, h)
			p.ctx.ClearColor(0, 0, 0, 0)
			p.ctx.Clear(bitColorBuffer)

			p.ctx.UseProgram(ps.ref)

			// Bind source texture
			p.ctx.ActiveTexture(texture0)
			p.ctx.BindTexture(texture2D, srcTex)

			// Standard uniforms
			p.SetUniform1f(*ps, "tex", 0) // sampler — will be ignored if it's not a float uniform, but Uniform1i is needed
			if u, ok := ps.uniforms["tex"]; ok {
				p.ctx.Uniform1i(u.ref, 0)
			}
			p.SetUniform2f(*ps, "resolution", resW, resH)
			p.SetUniform2f(*ps, "texelSize", texelW, texelH)

			// Set Gaussian blur direction for two-pass blur
			if eff.Type() == effect.GaussianBlur {
				if pass == 0 {
					p.SetUniform2f(*ps, "direction", 1.0, 0.0) // horizontal
				} else {
					p.SetUniform2f(*ps, "direction", 0.0, 1.0) // vertical
				}
			}

			// Set effect-specific uniforms from the Effect's uniform map
			uniforms := eff.Uniforms()
			for name, val := range uniforms {
				switch v := val.(type) {
				case float32:
					p.SetUniform1f(*ps, name, v)
				case [2]float32:
					p.SetUniform2f(*ps, name, v[0], v[1])
				case [3]float32:
					p.SetUniform3f(*ps, name, v[0], v[1], v[2])
				case [4]float32:
					p.SetUniform4f(*ps, name, v[0], v[1], v[2], v[3])
				}
			}

			// Draw fullscreen quad
			p.updateBuffer(ep.quadBuf, effectQuadVerts)
			p.UpdateVertexArray(*ps, "vert", 2, 4, 0)
			p.UpdateVertexArray(*ps, "vertTexCoord", 2, 4, 2)

			p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
			p.ctx.DrawArrays(triangleStrip, 0, 4)
			p.logError()

			// Swap source and destination
			srcTex, dstTex = dstTex, srcTex
			useB = !useB
		}
	}

	// Step 3: Composite final result to screen
	// Restore the clip scissor (disabled above) so the effect output is clipped to
	// the tile just like the un-effected object would have been.
	if p.clipEnabled {
		p.ctx.Enable(scissorTest)
	}
	p.ctx.BindFramebuffer(glFramebuffer, noFramebuffer)

	// Restore viewport to full frame
	frameW, frameH := p.scaleFrameSize(frame)
	p.ctx.Viewport(0, 0, int(frameW), int(frameH))

	// Draw the result texture at the object's actual position using the
	// existing simple texture program (p.program)
	p.compositeEffectResult(o, srcTex, pos, frame)
}

// compositeEffectResult draws the final FBO texture to the screen at the correct position.
func (p *painter) compositeEffectResult(o fyne.CanvasObject, tex Texture, pos fyne.Position, frame fyne.Size) {
	size := o.Size()

	// Compute screen-space quad coords (same as rectCoords but simplified)
	points := p.effectScreenQuad(pos, size, frame)

	p.ctx.UseProgram(p.program.ref)
	p.updateBuffer(p.program.buff, points)
	p.UpdateVertexArray(p.program, "vert", 3, 5, 0)
	p.UpdateVertexArray(p.program, "vertTexCoord", 2, 5, 3)

	// No corner radius for effect composite
	p.SetUniform1f(p.program, "cornerRadius", 0)
	p.SetUniform2f(p.program, "size", size.Width*p.pixScale, size.Height*p.pixScale)
	p.SetUniform4f(p.program, "inset", 0, 0, 1, 1) // full texture
	p.SetUniform1f(p.program, "alpha", 1.0)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, tex)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

// effectScreenQuad computes a 5-float-per-vertex quad (x,y,z, u,v) for
// positioning the effect result on screen, matching Fyne's NDC convention.
func (p *painter) effectScreenQuad(pos fyne.Position, size, frame fyne.Size) []float32 {
	x1 := pos.X
	y1 := pos.Y
	x2 := pos.X + size.Width
	y2 := pos.Y + size.Height

	// Convert to normalized device coordinates [-1, 1]
	x1Ndc := -1 + x1*2/frame.Width
	x2Ndc := -1 + x2*2/frame.Width
	y1Ndc := 1 - y1*2/frame.Height
	y2Ndc := 1 - y2*2/frame.Height

	return []float32{
		x1Ndc, y2Ndc, 0, 0, 0, // bottom-left: tex (0,0)
		x2Ndc, y2Ndc, 0, 1, 0, // bottom-right: tex (1,0)
		x1Ndc, y1Ndc, 0, 0, 1, // top-left: tex (0,1)
		x2Ndc, y1Ndc, 0, 1, 1, // top-right: tex (1,1)
	}
}

// drawObjectDirect calls the normal drawObject type-switch without the effect check.
// This is used to render the base object into an FBO.
func (p *painter) drawObjectDirect(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	switch obj := o.(type) {
	case *canvas.StreamingImage:
		p.drawStreamingImage(obj, pos, frame)
	case *canvas.Circle:
		p.drawCircle(obj, pos, frame)
	case *canvas.Line:
		p.drawLine(obj, pos, frame)
	case *canvas.Image:
		p.drawImage(obj, pos, frame)
	case *canvas.Raster:
		p.drawRaster(obj, pos, frame)
	case *canvas.Rectangle:
		p.drawRectangle(obj, pos, frame)
	case *canvas.Text:
		p.drawText(obj, pos, frame)
	case *canvas.LinearGradient:
		p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
	case *canvas.RadialGradient:
		p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
	case *canvas.Polygon:
		p.drawPolygon(obj, pos, frame)
	case *canvas.Arc:
		p.drawArc(obj, pos, frame)
	case *canvas.ShaderRect:
		p.drawShaderRect(obj, pos, frame)
	}
}
