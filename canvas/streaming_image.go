package canvas

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
)

// Pixel format constants for StreamingImage, matching FFmpeg's AV_PIX_FMT_* values.
const (
	// Existing formats
	PixelFormatRGBA    = 0
	PixelFormatYUV420P = 1

	// Planar YUV 8-bit
	PixelFormatYUV422P  = 2
	PixelFormatYUV444P  = 3
	PixelFormatYUV410P  = 4
	PixelFormatYUV411P  = 5
	PixelFormatYUV440P  = 6
	PixelFormatYUVJ420P = 7 // full-range alias
	PixelFormatYUVJ422P = 8
	PixelFormatYUVJ444P = 9
	PixelFormatYUVJ440P = 10

	// Planar YUVA 8-bit (with alpha plane)
	PixelFormatYUVA420P = 11
	PixelFormatYUVA422P = 12
	PixelFormatYUVA444P = 13

	// High bit-depth planar YUV (little-endian, 9–16 bit)
	PixelFormatYUV420P9  = 14
	PixelFormatYUV420P10 = 15
	PixelFormatYUV420P12 = 16
	PixelFormatYUV420P14 = 17
	PixelFormatYUV420P16 = 18
	PixelFormatYUV422P9  = 19
	PixelFormatYUV422P10 = 20
	PixelFormatYUV422P12 = 21
	PixelFormatYUV422P14 = 22
	PixelFormatYUV422P16 = 23
	PixelFormatYUV444P9  = 24
	PixelFormatYUV444P10 = 25
	PixelFormatYUV444P12 = 26
	PixelFormatYUV444P14 = 27
	PixelFormatYUV444P16 = 28

	// High bit-depth YUVA
	PixelFormatYUVA420P9  = 29
	PixelFormatYUVA420P10 = 30
	PixelFormatYUVA420P16 = 31
	PixelFormatYUVA422P9  = 32
	PixelFormatYUVA422P10 = 33
	PixelFormatYUVA422P16 = 34
	PixelFormatYUVA444P9  = 35
	PixelFormatYUVA444P10 = 36
	PixelFormatYUVA444P16 = 37

	// Semi-planar NV (Y plane + interleaved UV)
	PixelFormatNV12 = 38
	PixelFormatNV21 = 39 // UV order reversed
	PixelFormatNV16 = 40 // 4:2:2 semi-planar
	PixelFormatNV24 = 41 // 4:4:4 semi-planar
	PixelFormatNV42 = 42 // 4:4:4 semi-planar, UV reversed
	PixelFormatP010 = 43 // 10-bit NV12
	PixelFormatP012 = 44 // 12-bit NV12
	PixelFormatP016 = 45 // 16-bit NV12

	// Packed YUV 4:2:2
	PixelFormatYUYV422 = 46
	PixelFormatUYVY422 = 47
	PixelFormatYVYU422 = 48

	// Packed RGB / BGR
	PixelFormatRGB24  = 49
	PixelFormatBGR24  = 50
	PixelFormatBGRA   = 51
	PixelFormatARGB   = 52
	PixelFormatABGR   = 53
	PixelFormatRGB565 = 54
	PixelFormatBGR565 = 55

	// Planar RGB (GBR plane order)
	PixelFormatGBRP   = 56
	PixelFormatGBRP10 = 57
	PixelFormatGBRP12 = 58
	PixelFormatGBRP16 = 59
	PixelFormatGBRAP  = 60 // planar RGBA

	// Grayscale
	PixelFormatGRAY8     = 61
	PixelFormatGRAY16    = 62
	PixelFormatYA8       = 63 // gray + alpha
	PixelFormatMONOWHITE = 64 // 1-bit, 0=white
	PixelFormatMONOBLACK = 65 // 1-bit, 0=black

	// Bayer raw sensor
	PixelFormatBAYER_BGGR8  = 66
	PixelFormatBAYER_RGGB8  = 67
	PixelFormatBAYER_GBRG8  = 68
	PixelFormatBAYER_GRBG8  = 69
	PixelFormatBAYER_BGGR16 = 70
	PixelFormatBAYER_RGGB16 = 71
	PixelFormatBAYER_GBRG16 = 72
	PixelFormatBAYER_GRBG16 = 73

	// Float formats
	PixelFormatGRAYF32 = 74
	PixelFormatRGBF32  = 75
	PixelFormatRGBAF32 = 76
)

// Color space constants, matching FFmpeg's AVCOL_SPC_* values.
const (
	ColorSpaceBT601    = 0 // ITU-R BT.601 (SD, NTSC/PAL)
	ColorSpaceBT709    = 1 // ITU-R BT.709 (HD)
	ColorSpaceBT2020   = 2 // ITU-R BT.2020 (UHD / 4K)
	ColorSpaceSMPTE240M = 3 // SMPTE 240M (obsolete HDTV)
	ColorSpaceFCC      = 4 // FCC (US SD broadcast)
	ColorSpaceYCgCo    = 5 // YCgCo (lossless codecs)
	ColorSpaceICtCp    = 6 // ICtCp (ITU-R BT.2100 HDR)
)

// Color range constants.
const (
	ColorRangeFull    = 0 // Full range: Y/U/V all [0, maxVal]
	ColorRangeLimited = 1 // Limited range: Y [16,235], UV [16,240] for 8-bit
)

// YUV420PFrame holds planar YUV 4:2:0 data for GPU upload.
// Y is full resolution (Width x Height), U and V are quarter resolution (Width/2 x Height/2).
// Stride values indicate the byte width of each plane row (may include padding).
type YUV420PFrame struct {
	Y, U, V                   []byte
	StrideY, StrideU, StrideV int
	Width, Height             int
}

// RawFrame holds planar or packed pixel data for any PixelFormat.
// Data holds up to 4 planes; unused planes are nil. The interpretation of
// each plane depends on the StreamingImage.PixelFormat setting.
//
// Strides are the byte widths of each plane row, including any padding.
// A zero stride means the plane is tightly packed: the painter will compute
// the correct stride from Width and the pixel format automatically. This
// covers the common case where data comes from a simple allocation.
//
// Set explicit strides when your source has row padding — for example,
// FFmpeg's AVFrame.linesize fields:
//
//	img.UpdateRawFrame(&canvas.RawFrame{
//	    Data:    [4][]byte{avFrame.Data[0], avFrame.Data[1], avFrame.Data[2]},
//	    Strides: [4]int{avFrame.Linesize[0], avFrame.Linesize[1], avFrame.Linesize[2]},
//	    Width:   avFrame.Width,
//	    Height:  avFrame.Height,
//	})
type RawFrame struct {
	Data    [4][]byte
	Strides [4]int
	Width   int
	Height  int
}

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

	// PixelFormat selects the frame data format. Default is PixelFormatRGBA.
	// Set at construction time; do not change after the image is visible.
	PixelFormat int

	// ColorSpace selects the YUV-to-RGB matrix coefficients. Default is
	// ColorSpaceBT601. Only meaningful for YUV pixel formats.
	ColorSpace int

	// ColorRange selects full or limited (broadcast) range. Default is
	// ColorRangeFull. Only meaningful for YUV pixel formats.
	ColorRange int

	mu           sync.Mutex
	pendingFrame *image.RGBA
	pendingYUV   *YUV420PFrame
	pendingRaw   *RawFrame
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

// UpdateYUVFrame provides new YUV420P planar data to be displayed. The frame
// will be uploaded to the GPU as three single-channel textures on the next
// paint cycle and converted to RGB by a fragment shader.
// This method is safe to call from any goroutine.
func (s *StreamingImage) UpdateYUVFrame(frame *YUV420PFrame) {
	s.mu.Lock()
	s.pendingYUV = frame
	s.mu.Unlock()

	repaint(s)
}

// ConsumePendingYUVFrame returns the most recently provided YUV frame and
// clears the pending state. Returns nil if no new frame is available.
// Called by the painter on the GL thread.
func (s *StreamingImage) ConsumePendingYUVFrame() *YUV420PFrame {
	s.mu.Lock()
	frame := s.pendingYUV
	s.pendingYUV = nil
	s.mu.Unlock()

	return frame
}

// UpdateRawFrame provides new frame data in the format specified by PixelFormat.
// The frame will be processed on the next paint cycle. This method is safe to
// call from any goroutine.
//
// The frame data should not be modified after calling this method until the
// next call to UpdateRawFrame.
func (s *StreamingImage) UpdateRawFrame(frame *RawFrame) {
	s.mu.Lock()
	s.pendingRaw = frame
	s.mu.Unlock()

	repaint(s)
}

// ConsumePendingRawFrame returns the most recently provided RawFrame and
// clears the pending state. Returns nil if no new frame is available.
// Called by the painter on the GL thread.
func (s *StreamingImage) ConsumePendingRawFrame() *RawFrame {
	s.mu.Lock()
	frame := s.pendingRaw
	s.pendingRaw = nil
	s.mu.Unlock()

	return frame
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

// NewStreamingImageYUV420P returns a StreamingImage configured for YUV420P
// planar input. Frames are uploaded as three single-channel textures and
// converted to RGB on the GPU via a fragment shader.
func NewStreamingImageYUV420P() *StreamingImage {
	return &StreamingImage{
		ScaleMode:   ImageScaleFastest,
		FillMode:    ImageFillStretch,
		PixelFormat: PixelFormatYUV420P,
	}
}

// NewStreamingImageWithFormat returns a StreamingImage configured for the
// specified pixel format, color space, and color range. See PixelFormat*,
// ColorSpace*, and ColorRange* constants.
func NewStreamingImageWithFormat(pixelFormat, colorSpace, colorRange int) *StreamingImage {
	return &StreamingImage{
		ScaleMode:   ImageScaleFastest,
		FillMode:    ImageFillStretch,
		PixelFormat: pixelFormat,
		ColorSpace:  colorSpace,
		ColorRange:  colorRange,
	}
}
