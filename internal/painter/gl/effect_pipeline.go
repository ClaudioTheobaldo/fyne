package gl

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/canvas/effect"
)

// effectPipeline manages the FBO ping-pong rendering pipeline for shader effects.
type effectPipeline struct {
	quadBuf     Buffer                              // fullscreen quad vertex buffer
	passthrough ProgramState                        // passthrough shader for compositing final result
	progCache   map[effect.EffectType]*ProgramState // compiled effect shader programs
	customCache map[string]*ProgramState            // compiled custom shader programs keyed by source
	bilateral   map[int]*ProgramState               // BilateralBlur programs keyed by window radius
	initialized bool

	// Render targets, one set per object size in use. See targetsFor.
	targets    []*effectTargets
	targetTick uint64 // bumped on every lookup; drives least-recently-used eviction
}

// effectTargets is the ping-pong FBO pair (plus optional MSAA buffer) for one
// pixel size.
type effectTargets struct {
	fboA, fboB    Framebuffer
	texA, texB    Texture
	width, height int
	bytes         int    // approximate VRAM held, for the cache budget
	lastUsed      uint64 // effectPipeline.targetTick at the last lookup

	// MSAA support
	msaaFBO     Framebuffer
	msaaRBO     Renderbuffer // multisample color renderbuffer
	msaaSamples int
}

// Limits for the render target cache. A 2560x1440 pair is ~29 MB, so the byte
// budget keeps one full-resolution video tile plus the small widget targets
// resident, without letting several large sizes pile up. A single target larger
// than the budget is still allowed; it just evicts everything else.
const (
	maxEffectTargets      = 4
	maxEffectTargetsBytes = 64 << 20
)

// fullscreen quad: 2 triangles covering [-1,1] with matching UV [0,1]
// Layout: x, y, u, v (4 floats per vertex, 4 vertices as triangle strip)
var effectQuadVerts = []float32{
	-1, -1, 0, 0, // bottom-left
	1, -1, 1, 0, // bottom-right
	-1, 1, 0, 1, // top-left
	1, 1, 1, 1, // top-right
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

	linkOK := p.ctx.GetProgrami(prog, linkStatus) != glFalse

	// The program owns its compiled code after linking; release the shader
	// objects so they are not held for the process lifetime.
	p.ctx.DeleteShader(vertShader)
	p.ctx.DeleteShader(fragShader)

	if !linkOK {
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
	ep.bilateral = make(map[int]*ProgramState)
	ep.initialized = true
}

// targetsFor returns render targets of exactly w x h pixels, creating them if
// needed.
//
// There is one pipeline per painter, shared by every object that has effects.
// It used to hold a single FBO pair and reallocate it whenever the requested
// size changed, so two differently-sized objects drawn in the same frame
// destroyed and recreated both textures twice per frame (measured at 123
// reallocations a second from two ordinary widgets). Keeping a small cache keyed
// by size lets each object reuse its own pair. Exact sizes are kept rather than
// rendering into a sub-rectangle of a larger target so the shaders, which
// sample in [0,1] UV space, stay untouched.
func (p *painter) targetsFor(w, h int) *effectTargets {
	ep := &p.effectPipe
	ep.targetTick++

	for _, t := range ep.targets {
		if t.width == w && t.height == h {
			t.lastUsed = ep.targetTick
			return t
		}
	}

	t := p.newEffectTargets(w, h)
	t.lastUsed = ep.targetTick
	ep.targets = append(ep.targets, t)
	p.evictEffectTargets(t)
	return t
}

// evictEffectTargets frees least-recently-used targets until the cache is
// within its count and byte limits. keep is never evicted.
func (p *painter) evictEffectTargets(keep *effectTargets) {
	ep := &p.effectPipe
	for {
		total := 0
		for _, t := range ep.targets {
			total += t.bytes
		}
		if len(ep.targets) <= maxEffectTargets && total <= maxEffectTargetsBytes {
			return
		}

		oldest := -1
		for i, t := range ep.targets {
			if t == keep {
				continue
			}
			if oldest < 0 || t.lastUsed < ep.targets[oldest].lastUsed {
				oldest = i
			}
		}
		if oldest < 0 {
			return // only keep is left
		}

		p.deleteEffectTargets(ep.targets[oldest])
		ep.targets = append(ep.targets[:oldest], ep.targets[oldest+1:]...)
	}
}

func (p *painter) deleteEffectTargets(t *effectTargets) {
	p.ctx.DeleteFramebuffer(t.fboA)
	p.ctx.DeleteFramebuffer(t.fboB)
	p.ctx.DeleteTexture(t.texA)
	p.ctx.DeleteTexture(t.texB)
	if t.msaaSamples > 0 {
		p.ctx.DeleteFramebuffer(t.msaaFBO)
		p.ctx.DeleteRenderbuffer(t.msaaRBO)
	}
}

// newEffectTargets creates the two FBOs and their color textures.
func (p *painter) newEffectTargets(w, h int) *effectTargets {
	t := &effectTargets{width: w, height: h, bytes: 2 * w * h * 4}

	// Create two RGBA textures
	t.texA = p.newFBOTexture(w, h)
	t.texB = p.newFBOTexture(w, h)

	// Create and setup FBO A
	t.fboA = p.ctx.CreateFramebuffer()
	p.ctx.BindFramebuffer(glFramebuffer, t.fboA)
	p.ctx.FramebufferTexture2D(glFramebuffer, colorAttachment0, texture2D, t.texA, 0)
	status := p.ctx.CheckFramebufferStatus(glFramebuffer)
	if status != framebufferComplete {
		fyne.LogError("EffectPipeline: FBO A incomplete", fmt.Errorf("status: 0x%x", status))
	}

	// Create and setup FBO B
	t.fboB = p.ctx.CreateFramebuffer()
	p.ctx.BindFramebuffer(glFramebuffer, t.fboB)
	p.ctx.FramebufferTexture2D(glFramebuffer, colorAttachment0, texture2D, t.texB, 0)
	status = p.ctx.CheckFramebufferStatus(glFramebuffer)
	if status != framebufferComplete {
		fyne.LogError("EffectPipeline: FBO B incomplete", fmt.Errorf("status: 0x%x", status))
	}

	// Create MSAA FBO if hardware multisampling is active
	if p.aaMode.NeedsMSAA() {
		samples := int32(p.aaMode.Samples())
		t.msaaRBO = p.ctx.CreateRenderbuffer()
		p.ctx.BindRenderbuffer(glRenderbuffer, t.msaaRBO)
		p.ctx.RenderbufferStorageMultisample(glRenderbuffer, samples, colorFormatRGBA, int32(w), int32(h))

		t.msaaFBO = p.ctx.CreateFramebuffer()
		p.ctx.BindFramebuffer(glFramebuffer, t.msaaFBO)
		p.ctx.FramebufferRenderbuffer(glFramebuffer, colorAttachment0, glRenderbuffer, t.msaaRBO)
		status = p.ctx.CheckFramebufferStatus(glFramebuffer)
		if status != framebufferComplete {
			fyne.LogError("EffectPipeline: MSAA FBO incomplete", fmt.Errorf("status: 0x%x", status))
		}
		t.msaaSamples = int(samples)
		t.bytes += w * h * 4 * int(samples)
	}

	// Unbind FBO (bind default framebuffer 0)
	p.bindDefaultFramebuffer()
	return t
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
	if kind == effect.BilateralBlur {
		return p.getBilateralProgram(eff)
	}
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

// bilateralMaxRadius matches the window the shader supported before it was
// specialised per radius.
const bilateralMaxRadius = 8

// getBilateralProgram returns the BilateralBlur program whose loop covers
// exactly ceil(radius) taps either side of the centre.
//
// The shader needs a constant loop bound, so a single program had to loop over
// the largest window and mask out the unused taps; every pixel then paid the
// loop overhead of a 17x17 window even at radius 1. Compiling one program per
// whole-number radius (at most nine, built on first use) keeps the cost in
// proportion to the strength. The radius is read on every draw, so a radius
// changed in place with SetFloat picks up the right program too.
func (p *painter) getBilateralProgram(eff *effect.Effect) *ProgramState {
	ep := &p.effectPipe

	var radius float32
	if v, ok := eff.Uniform("radius").(float32); ok {
		radius = v
	}
	r := int(math.Ceil(float64(radius)))
	if r < 0 {
		r = 0
	} else if r > bilateralMaxRadius {
		r = bilateralMaxRadius
	}

	if ps, ok := ep.bilateral[r]; ok {
		return ps
	}

	const marker = "#define MAX_RADIUS 8"
	src := string(effectShaderSource(effect.BilateralBlur))
	if strings.Count(src, marker) != 1 {
		fyne.LogError("EffectPipeline: bilateral shader must contain its radius define exactly once", nil)
		return &ep.passthrough
	}
	src = strings.Replace(src, marker, "#define MAX_RADIUS "+strconv.Itoa(r), 1)

	ps := p.compileEffectProgram(src)
	if ps == nil {
		return &ep.passthrough
	}
	ep.bilateral[r] = ps
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

	linkOK := p.ctx.GetProgrami(prog, linkStatus) != glFalse

	// The linked program keeps its own copy of the compiled code, so the shader
	// objects are dead weight from here on. Deleting them flags them for removal;
	// the driver frees each once it is no longer attached to a live program.
	// Without this every compile left two shader objects allocated for the
	// lifetime of the process.
	p.ctx.DeleteShader(vertShader)
	p.ctx.DeleteShader(fragShader)

	if !linkOK {
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

	t := p.targetsFor(w, h)

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
	if t.msaaSamples > 0 {
		// Render into multisample FBO
		p.ctx.BindFramebuffer(glFramebuffer, t.msaaFBO)
		p.ctx.Viewport(0, 0, w, h)
		p.ctx.ClearColor(0, 0, 0, 0)
		p.ctx.Clear(bitColorBuffer)
		p.drawObjectDirect(o, fyne.NewPos(0, 0), fboFrame)

		// Resolve MSAA into fboA (texture-backed) via blit
		p.ctx.BindFramebuffer(glReadFramebuffer, t.msaaFBO)
		p.ctx.BindFramebuffer(glDrawFramebuffer, t.fboA)
		p.ctx.BlitFramebuffer(0, 0, int32(w), int32(h), 0, 0, int32(w), int32(h), bitColorBuffer, glLinear)
		p.ctx.BindFramebuffer(glFramebuffer, noFramebuffer)
	} else {
		// No MSAA: render directly into fboA
		p.ctx.BindFramebuffer(glFramebuffer, t.fboA)
		p.ctx.Viewport(0, 0, w, h)
		p.ctx.ClearColor(0, 0, 0, 0)
		p.ctx.Clear(bitColorBuffer)
		p.drawObjectDirect(o, fyne.NewPos(0, 0), fboFrame)
	}

	// Step 2: Ping-pong through effect chain
	srcTex := t.texA
	dstTex := t.texB
	// useB tracks which FBO is the current destination (true=fboB, false=fboA)
	useB := true

	resW := float32(w)
	resH := float32(h)
	texelW := 1.0 / resW
	texelH := 1.0 / resH

	// Viewport, clear colour, blend mode and the quad VBO are identical for
	// every pass, so set them once here rather than per pass. The quad's
	// contents never change after init, so re-uploading it each pass was a
	// full glBufferData reallocation of static data.
	//
	// The blend mode is deliberately left as src-alpha over the transparent
	// clear. For a source with partial alpha that multiplies RGB by alpha once
	// per pass, which compounds across a chain — switching to no blending, or
	// dropping the per-pass clear, would avoid it and save a full-tile clear
	// per pass. It is not done here because the change is only safe to judge
	// against a translucent source: the effect chain is also used by widgets
	// whose antialiased edges carry real partial alpha, and there is no
	// coverage for those. Left as-is deliberately rather than overlooked.
	p.ctx.Viewport(0, 0, w, h)
	p.ctx.ClearColor(0, 0, 0, 0)
	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.ctx.BindBuffer(arrayBuffer, ep.quadBuf)

	for _, eff := range active {
		ps := p.getEffectProgram(eff)
		passes := effect.PassCount(eff.Type())

		for pass := 0; pass < passes; pass++ {
			dstFBO := t.fboA
			if useB {
				dstFBO = t.fboB
			}

			p.ctx.BindFramebuffer(glFramebuffer, dstFBO)
			p.ctx.Clear(bitColorBuffer)

			p.ctx.UseProgram(ps.ref)

			// Bind source texture
			p.ctx.ActiveTexture(texture0)
			p.ctx.BindTexture(texture2D, srcTex)

			// Standard uniforms
			// tex is a sampler: it takes Uniform1i. Setting it through
			// SetUniform1f would also poison that uniform's cached previous
			// value, defeating the redundant-upload check for it.
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
			// RangeUniforms rather than Uniforms: this runs for every pass of
			// every effect of every frame, and Uniforms allocates a fresh map
			// each call.
			eff.RangeUniforms(func(name string, val any) {
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
			})

			// Draw fullscreen quad. The VBO is already bound and filled; only
			// the attribute pointers need rebinding, since they are program state.
			p.UpdateVertexArray(*ps, "vert", 2, 4, 0)
			p.UpdateVertexArray(*ps, "vertTexCoord", 2, 4, 2)

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
