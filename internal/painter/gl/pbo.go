package gl

import "fyne.io/fyne/v2/canvas"

// PBO-related GL constants (universal OpenGL spec values).
const (
	pixelUnpackBuffer uint32 = 0x88EC
	streamDraw        uint32 = 0x88E0
	writeOnly         uint32 = 0x88B9
)

// pboState holds the double-buffered Pixel Buffer Object state for a StreamingImage.
// Two PBOs are used in a ping-pong pattern: while the GPU reads from one PBO to
// update the texture, the CPU writes the next frame's pixels into the other.
type pboState struct {
	buffers [2]Buffer // ping-pong PBO pair
	index   int       // current write target (0 or 1)
	width   int       // allocated buffer dimensions
	height  int
	ready   bool // false until first frame has been uploaded via PBO
}

// getOrCreatePBO returns the PBO state for the given StreamingImage, creating
// new PBOs if needed or if the dimensions changed.
func (p *painter) getOrCreatePBO(img *canvas.StreamingImage, w, h int) *pboState {
	if p.pboStates == nil {
		p.pboStates = make(map[*canvas.StreamingImage]*pboState)
	}

	state, ok := p.pboStates[img]
	if ok && state.width == w && state.height == h {
		return state
	}

	// Dimensions changed or first creation — (re)allocate PBOs
	if ok {
		p.destroyPBOState(state)
	}

	size := w * h * 4 // RGBA
	state = &pboState{width: w, height: h}

	for i := 0; i < 2; i++ {
		state.buffers[i] = p.ctx.CreateBuffer()
		p.ctx.BindBuffer(pixelUnpackBuffer, state.buffers[i])
		p.ctx.BufferDataBytes(pixelUnpackBuffer, size, nil, streamDraw)
	}
	p.ctx.BindBuffer(pixelUnpackBuffer, noBuffer)
	p.logError()

	p.pboStates[img] = state
	return state
}

// destroyPBO removes and deletes PBOs for a StreamingImage.
func (p *painter) destroyPBO(img *canvas.StreamingImage) {
	if p.pboStates == nil {
		return
	}
	state, ok := p.pboStates[img]
	if !ok {
		return
	}
	p.destroyPBOState(state)
	delete(p.pboStates, img)
}

// destroyPBOState deletes the GL buffer objects in a pboState.
func (p *painter) destroyPBOState(state *pboState) {
	for i := 0; i < 2; i++ {
		p.ctx.DeleteBuffer(state.buffers[i])
	}
}

// planePBOPair holds a double-buffered PBO pair for one texture plane.
type planePBOPair struct {
	buffers [2]Buffer
}

// streamPBOState is the generic multi-plane PBO state for UpdateRawFrame uploads.
// Up to 4 planes are supported; planeCount tracks how many are in use.
type streamPBOState struct {
	planes     [4]planePBOPair
	textures   [4]Texture
	planeCount int
	index      int // current write-buffer index (0 or 1)
	width      int
	height     int
	ready      bool // false until the first complete frame is available
}

// getOrCreateStreamPBO returns the streamPBOState for img, re-creating it if
// dimensions or plane count changed.
func (p *painter) getOrCreateStreamPBO(img *canvas.StreamingImage, planeCount, w, h int, planeSizes [4]int) *streamPBOState {
	if p.rawPBOStates == nil {
		p.rawPBOStates = make(map[*canvas.StreamingImage]*streamPBOState)
	}

	state, ok := p.rawPBOStates[img]
	if ok && state.width == w && state.height == h && state.planeCount == planeCount {
		return state
	}

	if ok {
		p.destroyStreamPBOState(state)
	}

	state = &streamPBOState{planeCount: planeCount, width: w, height: h}
	for pi := 0; pi < planeCount; pi++ {
		for i := 0; i < 2; i++ {
			state.planes[pi].buffers[i] = p.ctx.CreateBuffer()
			p.ctx.BindBuffer(pixelUnpackBuffer, state.planes[pi].buffers[i])
			p.ctx.BufferDataBytes(pixelUnpackBuffer, planeSizes[pi], nil, streamDraw)
		}
		state.textures[pi] = p.newTexture(img.ScaleMode)
	}
	p.ctx.BindBuffer(pixelUnpackBuffer, noBuffer)
	p.logError()

	p.rawPBOStates[img] = state
	return state
}

// destroyStreamPBO removes streamPBOState for img.
func (p *painter) destroyStreamPBO(img *canvas.StreamingImage) {
	if p.rawPBOStates == nil {
		return
	}
	state, ok := p.rawPBOStates[img]
	if !ok {
		return
	}
	p.destroyStreamPBOState(state)
	delete(p.rawPBOStates, img)
}

// destroyStreamPBOState releases all GL resources in a streamPBOState.
func (p *painter) destroyStreamPBOState(state *streamPBOState) {
	for pi := 0; pi < state.planeCount; pi++ {
		for i := 0; i < 2; i++ {
			p.ctx.DeleteBuffer(state.planes[pi].buffers[i])
		}
		p.ctx.DeleteTexture(state.textures[pi])
	}
}
