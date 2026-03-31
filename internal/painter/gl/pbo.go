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

// yuvPBOState holds double-buffered PBOs and cached textures for YUV420P streaming.
// Three pairs of PBOs (one per plane: Y, U, V) plus three textures.
type yuvPBOState struct {
	yPBO, uPBO, vPBO [2]Buffer
	texY, texU, texV Texture
	index            int
	width, height    int
	ready            bool
}

// getOrCreateYUVPBO returns the YUV PBO state for the given StreamingImage,
// creating new PBOs and textures if needed or if dimensions changed.
func (p *painter) getOrCreateYUVPBO(img *canvas.StreamingImage, w, h int) *yuvPBOState {
	if p.yuvPBOStates == nil {
		p.yuvPBOStates = make(map[*canvas.StreamingImage]*yuvPBOState)
	}

	state, ok := p.yuvPBOStates[img]
	if ok && state.width == w && state.height == h {
		return state
	}

	if ok {
		p.destroyYUVPBOState(state)
	}

	state = &yuvPBOState{width: w, height: h}

	// Y plane: w * h bytes, U/V planes: (w/2) * (h/2) bytes each
	ySize := w * h
	uvSize := (w / 2) * (h / 2)

	// Create PBO pairs for each plane
	for i := 0; i < 2; i++ {
		state.yPBO[i] = p.ctx.CreateBuffer()
		p.ctx.BindBuffer(pixelUnpackBuffer, state.yPBO[i])
		p.ctx.BufferDataBytes(pixelUnpackBuffer, ySize, nil, streamDraw)

		state.uPBO[i] = p.ctx.CreateBuffer()
		p.ctx.BindBuffer(pixelUnpackBuffer, state.uPBO[i])
		p.ctx.BufferDataBytes(pixelUnpackBuffer, uvSize, nil, streamDraw)

		state.vPBO[i] = p.ctx.CreateBuffer()
		p.ctx.BindBuffer(pixelUnpackBuffer, state.vPBO[i])
		p.ctx.BufferDataBytes(pixelUnpackBuffer, uvSize, nil, streamDraw)
	}
	p.ctx.BindBuffer(pixelUnpackBuffer, noBuffer)
	p.logError()

	// Create textures for each plane
	state.texY = p.newTexture(img.ScaleMode)
	state.texU = p.newTexture(img.ScaleMode)
	state.texV = p.newTexture(img.ScaleMode)

	p.yuvPBOStates[img] = state
	return state
}

// destroyYUVPBO removes and deletes YUV PBOs and textures for a StreamingImage.
func (p *painter) destroyYUVPBO(img *canvas.StreamingImage) {
	if p.yuvPBOStates == nil {
		return
	}
	state, ok := p.yuvPBOStates[img]
	if !ok {
		return
	}
	p.destroyYUVPBOState(state)
	delete(p.yuvPBOStates, img)
}

// destroyYUVPBOState deletes the GL buffer objects and textures.
func (p *painter) destroyYUVPBOState(state *yuvPBOState) {
	for i := 0; i < 2; i++ {
		p.ctx.DeleteBuffer(state.yPBO[i])
		p.ctx.DeleteBuffer(state.uPBO[i])
		p.ctx.DeleteBuffer(state.vPBO[i])
	}
	p.ctx.DeleteTexture(state.texY)
	p.ctx.DeleteTexture(state.texU)
	p.ctx.DeleteTexture(state.texV)
}
