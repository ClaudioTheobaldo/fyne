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
