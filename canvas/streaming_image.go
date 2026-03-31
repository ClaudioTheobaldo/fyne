package canvas

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*StreamingImage)(nil)

// StreamingImage is a canvas object designed for efficient rendering of rapidly
// updating image data such as video frames. Unlike Image, it reuses its GPU
// texture and updates pixel data in-place via glTexSubImage2D, avoiding the
// overhead of texture allocation and deallocation on every frame.
//
// Since: 2.8
type StreamingImage struct {
	baseObject

	// Set a translucency value > 0.0 to fade the image
	Translucency float64
	// Specify the fill mode for the image within its container
	FillMode ImageFill
	// Specify the type of scaling interpolation applied to the image
	ScaleMode ImageScale

	// DisablePBO forces the use of direct TexSubImage2D instead of PBO
	// double-buffering. Useful for benchmarking or platforms where PBOs
	// cause issues.
	DisablePBO bool
	// DisableTexReuse forces texture reallocation each frame instead of
	// updating in-place. Useful for benchmarking to isolate PBO benefits
	// from texture reuse benefits.
	DisableTexReuse bool

	mu           sync.Mutex
	pendingFrame *image.RGBA
	texWidth     int
	texHeight    int
}

// Alpha is a convenience function that returns the alpha value for a streaming image
// based on its Translucency value. The result is 1.0 - Translucency.
func (s *StreamingImage) Alpha() float64 {
	return 1.0 - s.Translucency
}

// Hide will set this streaming image to not be visible.
func (s *StreamingImage) Hide() {
	s.baseObject.Hide()

	repaint(s)
}

// Move the streaming image to a new position, relative to its parent / canvas.
func (s *StreamingImage) Move(pos fyne.Position) {
	if s.Position() == pos {
		return
	}

	s.baseObject.Move(pos)

	repaint(s)
}

// Resize on a streaming image causes the new display size to be set and then calls Refresh.
// This does not affect the underlying texture dimensions, which are determined by the frame data.
func (s *StreamingImage) Resize(size fyne.Size) {
	if size == s.Size() {
		return
	}

	s.baseObject.Resize(size)
	Refresh(s)
}

// Refresh causes this streaming image to be redrawn.
func (s *StreamingImage) Refresh() {
	Refresh(s)
}

// UpdateFrame provides new pixel data to be displayed. The frame will be
// uploaded to the GPU on the next paint cycle. This method is safe to call
// from any goroutine.
//
// The provided image.RGBA should not be modified after calling this method
// until the next call to UpdateFrame.
func (s *StreamingImage) UpdateFrame(frame *image.RGBA) {
	s.mu.Lock()
	s.pendingFrame = frame
	s.mu.Unlock()

	repaint(s)
}

// ConsumePendingFrame returns the most recently provided frame and clears the
// pending state. Returns nil if no new frame is available. This is called by the
// painter on the GL thread.
func (s *StreamingImage) ConsumePendingFrame() *image.RGBA {
	s.mu.Lock()
	frame := s.pendingFrame
	s.pendingFrame = nil
	s.mu.Unlock()

	return frame
}

// TextureSize returns the current texture dimensions.
func (s *StreamingImage) TextureSize() (int, int) {
	s.mu.Lock()
	w, h := s.texWidth, s.texHeight
	s.mu.Unlock()
	return w, h
}

// SetTextureSize updates the stored texture dimensions. Called by the painter
// after allocating or reallocating the GPU texture.
func (s *StreamingImage) SetTextureSize(w, h int) {
	s.mu.Lock()
	s.texWidth = w
	s.texHeight = h
	s.mu.Unlock()
}

// NewStreamingImage returns a new StreamingImage instance optimized for
// displaying rapidly updating image data such as video frames.
//
// Since: 2.8
func NewStreamingImage() *StreamingImage {
	return &StreamingImage{
		ScaleMode: ImageScaleFastest,
		FillMode:  ImageFillStretch,
	}
}
