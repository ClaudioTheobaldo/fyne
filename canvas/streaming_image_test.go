package canvas_test

import (
	"image"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// ── Construction ────────────────────────────────────────────────────────────

func TestNewStreamingImage_Defaults(t *testing.T) {
	img := canvas.NewStreamingImage()
	if img == nil {
		t.Fatal("NewStreamingImage returned nil")
	}
	if img.PixelFormat != canvas.PixelFormatRGBA {
		t.Errorf("PixelFormat: got %d, want PixelFormatRGBA (%d)", img.PixelFormat, canvas.PixelFormatRGBA)
	}
	if img.ColorSpace != canvas.ColorSpaceBT601 {
		t.Errorf("ColorSpace: got %d, want ColorSpaceBT601 (%d)", img.ColorSpace, canvas.ColorSpaceBT601)
	}
	if img.ColorRange != canvas.ColorRangeFull {
		t.Errorf("ColorRange: got %d, want ColorRangeFull (%d)", img.ColorRange, canvas.ColorRangeFull)
	}
	if img.Alpha() != 1.0 {
		t.Errorf("Alpha(): got %f, want 1.0", img.Alpha())
	}
}

func TestNewStreamingImageYUV420P_Format(t *testing.T) {
	img := canvas.NewStreamingImageYUV420P()
	if img.PixelFormat != canvas.PixelFormatYUV420P {
		t.Errorf("PixelFormat: got %d, want PixelFormatYUV420P (%d)", img.PixelFormat, canvas.PixelFormatYUV420P)
	}
}

func TestNewStreamingImageWithFormat(t *testing.T) {
	cases := []struct {
		name        string
		pixelFormat int
		colorSpace  int
		colorRange  int
	}{
		{"YUV420P BT601 full", canvas.PixelFormatYUV420P, canvas.ColorSpaceBT601, canvas.ColorRangeFull},
		{"YUV420P BT709 limited", canvas.PixelFormatYUV420P, canvas.ColorSpaceBT709, canvas.ColorRangeLimited},
		{"NV12 BT2020 full", canvas.PixelFormatNV12, canvas.ColorSpaceBT2020, canvas.ColorRangeFull},
		{"GRAY8 default", canvas.PixelFormatGRAY8, canvas.ColorSpaceBT601, canvas.ColorRangeFull},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			img := canvas.NewStreamingImageWithFormat(tc.pixelFormat, tc.colorSpace, tc.colorRange)
			if img.PixelFormat != tc.pixelFormat {
				t.Errorf("PixelFormat: got %d, want %d", img.PixelFormat, tc.pixelFormat)
			}
			if img.ColorSpace != tc.colorSpace {
				t.Errorf("ColorSpace: got %d, want %d", img.ColorSpace, tc.colorSpace)
			}
			if img.ColorRange != tc.colorRange {
				t.Errorf("ColorRange: got %d, want %d", img.ColorRange, tc.colorRange)
			}
		})
	}
}

// ── Alpha / Translucency ─────────────────────────────────────────────────────

func TestStreamingImage_Alpha(t *testing.T) {
	img := canvas.NewStreamingImage()
	if img.Alpha() != 1.0 {
		t.Errorf("default Alpha: got %f, want 1.0", img.Alpha())
	}

	img.Translucency = 0.25
	if img.Alpha() != 0.75 {
		t.Errorf("Alpha with Translucency=0.25: got %f, want 0.75", img.Alpha())
	}

	img.Translucency = 1.0
	if img.Alpha() != 0.0 {
		t.Errorf("fully transparent: got %f, want 0.0", img.Alpha())
	}
}

// ── UpdateFrame (RGBA legacy path) ───────────────────────────────────────────

func TestStreamingImage_UpdateFrame_StoredAndConsumed(t *testing.T) {
	img := canvas.NewStreamingImage()
	frame := image.NewRGBA(image.Rect(0, 0, 64, 64))

	// Before update: nothing pending.
	if got := img.ConsumePendingFrame(); got != nil {
		t.Error("ConsumePendingFrame: expected nil before UpdateFrame")
	}

	img.UpdateFrame(frame)

	got := img.ConsumePendingFrame()
	if got == nil {
		t.Fatal("ConsumePendingFrame: expected non-nil after UpdateFrame")
	}
	if got != frame {
		t.Error("ConsumePendingFrame: returned different pointer than supplied frame")
	}

	// Second consume: frame should be cleared.
	if got2 := img.ConsumePendingFrame(); got2 != nil {
		t.Error("ConsumePendingFrame: expected nil on second call")
	}
}

// ── UpdateYUVFrame (YUV420P legacy path) ─────────────────────────────────────

func TestStreamingImage_UpdateYUVFrame_StoredAndConsumed(t *testing.T) {
	img := canvas.NewStreamingImageYUV420P()

	if got := img.ConsumePendingYUVFrame(); got != nil {
		t.Error("ConsumePendingYUVFrame: expected nil before UpdateYUVFrame")
	}

	yuv := &canvas.YUV420PFrame{
		Y: make([]byte, 64*64), U: make([]byte, 32*32), V: make([]byte, 32*32),
		Width: 64, Height: 64,
	}
	img.UpdateYUVFrame(yuv)

	got := img.ConsumePendingYUVFrame()
	if got == nil {
		t.Fatal("ConsumePendingYUVFrame: expected non-nil after UpdateYUVFrame")
	}
	if got != yuv {
		t.Error("ConsumePendingYUVFrame: returned different pointer than supplied frame")
	}

	if got2 := img.ConsumePendingYUVFrame(); got2 != nil {
		t.Error("ConsumePendingYUVFrame: expected nil on second call")
	}
}

// ── UpdateRawFrame (generic path) ────────────────────────────────────────────

func TestStreamingImage_UpdateRawFrame_StoredAndConsumed(t *testing.T) {
	img := canvas.NewStreamingImageWithFormat(canvas.PixelFormatYUV422P, canvas.ColorSpaceBT709, canvas.ColorRangeFull)

	if got := img.ConsumePendingRawFrame(); got != nil {
		t.Error("ConsumePendingRawFrame: expected nil before UpdateRawFrame")
	}

	rawFrame := &canvas.RawFrame{
		Data:    [4][]byte{make([]byte, 64*64), make([]byte, 32*64), make([]byte, 32*64)},
		Strides: [4]int{64, 32, 32},
		Width:   64,
		Height:  64,
	}
	img.UpdateRawFrame(rawFrame)

	got := img.ConsumePendingRawFrame()
	if got == nil {
		t.Fatal("ConsumePendingRawFrame: expected non-nil after UpdateRawFrame")
	}
	if got != rawFrame {
		t.Error("ConsumePendingRawFrame: returned different pointer than supplied frame")
	}

	if got2 := img.ConsumePendingRawFrame(); got2 != nil {
		t.Error("ConsumePendingRawFrame: expected nil on second call")
	}
}

func TestStreamingImage_UpdateRawFrame_ReplacesOldFrame(t *testing.T) {
	img := canvas.NewStreamingImageWithFormat(canvas.PixelFormatNV12, canvas.ColorSpaceBT601, canvas.ColorRangeLimited)

	first := &canvas.RawFrame{Width: 64, Height: 64}
	second := &canvas.RawFrame{Width: 64, Height: 64}

	img.UpdateRawFrame(first)
	img.UpdateRawFrame(second) // second call replaces first (painter hasn't consumed yet)

	got := img.ConsumePendingRawFrame()
	if got != second {
		t.Error("expected second frame to replace first unreceived frame")
	}
}

// ── TextureSize ───────────────────────────────────────────────────────────────

func TestStreamingImage_TextureSize_InitialZero(t *testing.T) {
	img := canvas.NewStreamingImage()
	w, h := img.TextureSize()
	if w != 0 || h != 0 {
		t.Errorf("initial TextureSize: got (%d,%d), want (0,0)", w, h)
	}
}

func TestStreamingImage_SetTextureSize(t *testing.T) {
	img := canvas.NewStreamingImage()
	img.SetTextureSize(1920, 1080)
	w, h := img.TextureSize()
	if w != 1920 || h != 1080 {
		t.Errorf("TextureSize after Set: got (%d,%d), want (1920,1080)", w, h)
	}
}

// ── Resize / Move / Hide ─────────────────────────────────────────────────────

func TestStreamingImage_Resize(t *testing.T) {
	img := canvas.NewStreamingImage()
	img.Resize(fyne.NewSize(100, 50))
	sz := img.Size()
	if sz.Width != 100 || sz.Height != 50 {
		t.Errorf("Size after Resize: got %v, want {100 50}", sz)
	}
}

func TestStreamingImage_Move(t *testing.T) {
	img := canvas.NewStreamingImage()
	img.Move(fyne.NewPos(10, 20))
	pos := img.Position()
	if pos.X != 10 || pos.Y != 20 {
		t.Errorf("Position after Move: got %v, want {10 20}", pos)
	}
}

func TestStreamingImage_Hide_Show(t *testing.T) {
	img := canvas.NewStreamingImage()
	img.Hide()
	if img.Visible() {
		t.Error("after Hide(): Visible should be false")
	}
	img.Show()
	if !img.Visible() {
		t.Error("after Show(): Visible should be true")
	}
}

// ── Pixel format constants ────────────────────────────────────────────────────

func TestPixelFormatConstants_Unique(t *testing.T) {
	// Verify all named constants have distinct values.
	seen := make(map[int]string)
	check := func(name string, val int) {
		if prev, dup := seen[val]; dup {
			t.Errorf("PixelFormat value %d used by both %q and %q", val, prev, name)
		}
		seen[val] = name
	}

	check("RGBA", canvas.PixelFormatRGBA)
	check("YUV420P", canvas.PixelFormatYUV420P)
	check("YUV422P", canvas.PixelFormatYUV422P)
	check("YUV444P", canvas.PixelFormatYUV444P)
	check("YUV410P", canvas.PixelFormatYUV410P)
	check("YUV411P", canvas.PixelFormatYUV411P)
	check("YUV440P", canvas.PixelFormatYUV440P)
	check("YUVJ420P", canvas.PixelFormatYUVJ420P)
	check("YUVJ422P", canvas.PixelFormatYUVJ422P)
	check("YUVJ444P", canvas.PixelFormatYUVJ444P)
	check("YUVJ440P", canvas.PixelFormatYUVJ440P)
	check("YUVA420P", canvas.PixelFormatYUVA420P)
	check("YUVA422P", canvas.PixelFormatYUVA422P)
	check("YUVA444P", canvas.PixelFormatYUVA444P)
	check("YUV420P9", canvas.PixelFormatYUV420P9)
	check("YUV420P10", canvas.PixelFormatYUV420P10)
	check("YUV420P12", canvas.PixelFormatYUV420P12)
	check("YUV420P14", canvas.PixelFormatYUV420P14)
	check("YUV420P16", canvas.PixelFormatYUV420P16)
	check("YUV422P9", canvas.PixelFormatYUV422P9)
	check("YUV422P10", canvas.PixelFormatYUV422P10)
	check("YUV422P12", canvas.PixelFormatYUV422P12)
	check("YUV422P14", canvas.PixelFormatYUV422P14)
	check("YUV422P16", canvas.PixelFormatYUV422P16)
	check("YUV444P9", canvas.PixelFormatYUV444P9)
	check("YUV444P10", canvas.PixelFormatYUV444P10)
	check("YUV444P12", canvas.PixelFormatYUV444P12)
	check("YUV444P14", canvas.PixelFormatYUV444P14)
	check("YUV444P16", canvas.PixelFormatYUV444P16)
	check("YUVA420P9", canvas.PixelFormatYUVA420P9)
	check("YUVA420P10", canvas.PixelFormatYUVA420P10)
	check("YUVA420P16", canvas.PixelFormatYUVA420P16)
	check("YUVA422P9", canvas.PixelFormatYUVA422P9)
	check("YUVA422P10", canvas.PixelFormatYUVA422P10)
	check("YUVA422P16", canvas.PixelFormatYUVA422P16)
	check("YUVA444P9", canvas.PixelFormatYUVA444P9)
	check("YUVA444P10", canvas.PixelFormatYUVA444P10)
	check("YUVA444P16", canvas.PixelFormatYUVA444P16)
	check("NV12", canvas.PixelFormatNV12)
	check("NV21", canvas.PixelFormatNV21)
	check("NV16", canvas.PixelFormatNV16)
	check("NV24", canvas.PixelFormatNV24)
	check("NV42", canvas.PixelFormatNV42)
	check("P010", canvas.PixelFormatP010)
	check("P012", canvas.PixelFormatP012)
	check("P016", canvas.PixelFormatP016)
	check("YUYV422", canvas.PixelFormatYUYV422)
	check("UYVY422", canvas.PixelFormatUYVY422)
	check("YVYU422", canvas.PixelFormatYVYU422)
	check("RGB24", canvas.PixelFormatRGB24)
	check("BGR24", canvas.PixelFormatBGR24)
	check("BGRA", canvas.PixelFormatBGRA)
	check("ARGB", canvas.PixelFormatARGB)
	check("ABGR", canvas.PixelFormatABGR)
	check("RGB565", canvas.PixelFormatRGB565)
	check("BGR565", canvas.PixelFormatBGR565)
	check("GBRP", canvas.PixelFormatGBRP)
	check("GBRP10", canvas.PixelFormatGBRP10)
	check("GBRP12", canvas.PixelFormatGBRP12)
	check("GBRP16", canvas.PixelFormatGBRP16)
	check("GBRAP", canvas.PixelFormatGBRAP)
	check("GRAY8", canvas.PixelFormatGRAY8)
	check("GRAY16", canvas.PixelFormatGRAY16)
	check("YA8", canvas.PixelFormatYA8)
	check("MONOWHITE", canvas.PixelFormatMONOWHITE)
	check("MONOBLACK", canvas.PixelFormatMONOBLACK)
	check("BAYER_BGGR8", canvas.PixelFormatBAYER_BGGR8)
	check("BAYER_RGGB8", canvas.PixelFormatBAYER_RGGB8)
	check("BAYER_GBRG8", canvas.PixelFormatBAYER_GBRG8)
	check("BAYER_GRBG8", canvas.PixelFormatBAYER_GRBG8)
	check("BAYER_BGGR16", canvas.PixelFormatBAYER_BGGR16)
	check("BAYER_RGGB16", canvas.PixelFormatBAYER_RGGB16)
	check("BAYER_GBRG16", canvas.PixelFormatBAYER_GBRG16)
	check("BAYER_GRBG16", canvas.PixelFormatBAYER_GRBG16)
	check("GRAYF32", canvas.PixelFormatGRAYF32)
	check("RGBF32", canvas.PixelFormatRGBF32)
	check("RGBAF32", canvas.PixelFormatRGBAF32)
}

// ── Color space / range constants ─────────────────────────────────────────────

func TestColorSpaceConstants_Distinct(t *testing.T) {
	vals := []int{
		canvas.ColorSpaceBT601,
		canvas.ColorSpaceBT709,
		canvas.ColorSpaceBT2020,
		canvas.ColorSpaceSMPTE240M,
		canvas.ColorSpaceFCC,
		canvas.ColorSpaceYCgCo,
		canvas.ColorSpaceICtCp,
	}
	seen := make(map[int]bool)
	for _, v := range vals {
		if seen[v] {
			t.Errorf("duplicate ColorSpace value %d", v)
		}
		seen[v] = true
	}
}

func TestColorRangeConstants(t *testing.T) {
	if canvas.ColorRangeFull == canvas.ColorRangeLimited {
		t.Error("ColorRangeFull and ColorRangeLimited must be distinct")
	}
}

// ── RawFrame struct ───────────────────────────────────────────────────────────

func TestRawFrame_Fields(t *testing.T) {
	rf := &canvas.RawFrame{
		Data:    [4][]byte{{1, 2, 3}, nil, nil, nil},
		Strides: [4]int{3, 0, 0, 0},
		Width:   1,
		Height:  3,
	}
	if rf.Width != 1 || rf.Height != 3 {
		t.Errorf("RawFrame dimensions: got %dx%d, want 1x3", rf.Width, rf.Height)
	}
	if rf.Data[0][0] != 1 {
		t.Errorf("RawFrame Data[0][0]: got %d, want 1", rf.Data[0][0])
	}
}
