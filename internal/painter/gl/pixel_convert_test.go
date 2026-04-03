package gl

import (
	"encoding/binary"
	"image/color"
	"math"
	"testing"

	"fyne.io/fyne/v2/canvas"
)

// pxRGBA is a small helper to build a 1-pixel RawFrame with a flat byte slice.
func pxFrame(w, h int, planes ...[]byte) *canvas.RawFrame {
	f := &canvas.RawFrame{Width: w, Height: h}
	for i, p := range planes {
		if i < 4 {
			f.Data[i] = p
		}
	}
	return f
}

// expectRGBA checks that the pixel at (x,y) in img matches want within tolerance tol.
func expectRGBA(t *testing.T, label string, img interface{ At(int, int) color.Color }, x, y int, want color.RGBA, tol uint8) {
	t.Helper()
	got := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
	diff := func(a, b uint8) uint8 {
		if a >= b {
			return a - b
		}
		return b - a
	}
	if diff(got.R, want.R) > tol || diff(got.G, want.G) > tol ||
		diff(got.B, want.B) > tol || diff(got.A, want.A) > tol {
		t.Errorf("%s at (%d,%d): got RGBA(%d,%d,%d,%d) want RGBA(%d,%d,%d,%d) (tol %d)",
			label, x, y, got.R, got.G, got.B, got.A, want.R, want.G, want.B, want.A, tol)
	}
}

// f32LE encodes a float32 value to 4 little-endian bytes.
func f32LE(v float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(v))
	return b
}

// u16LE encodes a uint16 value to 2 little-endian bytes.
func u16LE(v uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	return b
}

// ---- Packed RGB conversions ------------------------------------------------

func TestConvertRGB24_SinglePixel(t *testing.T) {
	frame := pxFrame(1, 1, []byte{200, 100, 50}) // R,G,B
	out := convertRawToRGBA(canvas.PixelFormatRGB24, frame)
	expectRGBA(t, "RGB24", out, 0, 0, color.RGBA{200, 100, 50, 255}, 0)
}

func TestConvertBGR24_SinglePixel(t *testing.T) {
	// BGR order on disk: B=50, G=100, R=200
	frame := pxFrame(1, 1, []byte{50, 100, 200})
	out := convertRawToRGBA(canvas.PixelFormatBGR24, frame)
	expectRGBA(t, "BGR24", out, 0, 0, color.RGBA{200, 100, 50, 255}, 0)
}

func TestConvertBGRA_SinglePixel(t *testing.T) {
	// BGRA: B=50, G=100, R=200, A=128
	frame := pxFrame(1, 1, []byte{50, 100, 200, 128})
	out := convertRawToRGBA(canvas.PixelFormatBGRA, frame)
	expectRGBA(t, "BGRA", out, 0, 0, color.RGBA{200, 100, 50, 128}, 0)
}

func TestConvertARGB_SinglePixel(t *testing.T) {
	// ARGB: A=128, R=200, G=100, B=50
	frame := pxFrame(1, 1, []byte{128, 200, 100, 50})
	out := convertRawToRGBA(canvas.PixelFormatARGB, frame)
	expectRGBA(t, "ARGB", out, 0, 0, color.RGBA{200, 100, 50, 128}, 0)
}

func TestConvertABGR_SinglePixel(t *testing.T) {
	// ABGR: A=128, B=50, G=100, R=200
	frame := pxFrame(1, 1, []byte{128, 50, 100, 200})
	out := convertRawToRGBA(canvas.PixelFormatABGR, frame)
	expectRGBA(t, "ABGR", out, 0, 0, color.RGBA{200, 100, 50, 128}, 0)
}

func TestConvertRGB565_PrimaryColors(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  color.RGBA
	}{
		// Pure red: R=31,G=0,B=0 → value=0xF800 → LE=[0x00,0xF8]
		{"red", []byte{0x00, 0xF8}, color.RGBA{255, 0, 0, 255}},
		// Pure green: R=0,G=63,B=0 → value=0x07E0 → LE=[0xE0,0x07]
		{"green", []byte{0xE0, 0x07}, color.RGBA{0, 255, 0, 255}},
		// Pure blue: R=0,G=0,B=31 → value=0x001F → LE=[0x1F,0x00]
		{"blue", []byte{0x1F, 0x00}, color.RGBA{0, 0, 255, 255}},
		// Black
		{"black", []byte{0x00, 0x00}, color.RGBA{0, 0, 0, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := pxFrame(1, 1, tt.bytes)
			out := convertRawToRGBA(canvas.PixelFormatRGB565, frame)
			expectRGBA(t, "RGB565/"+tt.name, out, 0, 0, tt.want, 0)
		})
	}
}

func TestConvertBGR565_PrimaryColors(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  color.RGBA
	}{
		// BGR565: B in MSBits (15-11), G (10-5), R in LSBits (4-0)
		// Pure red: R=31 → value=0x001F → LE=[0x1F,0x00]
		{"red", []byte{0x1F, 0x00}, color.RGBA{255, 0, 0, 255}},
		// Pure green: G=63 → value=0x07E0 → LE=[0xE0,0x07]
		{"green", []byte{0xE0, 0x07}, color.RGBA{0, 255, 0, 255}},
		// Pure blue: B=31 → value=0xF800 → LE=[0x00,0xF8]
		{"blue", []byte{0x00, 0xF8}, color.RGBA{0, 0, 255, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := pxFrame(1, 1, tt.bytes)
			out := convertRawToRGBA(canvas.PixelFormatBGR565, frame)
			expectRGBA(t, "BGR565/"+tt.name, out, 0, 0, tt.want, 0)
		})
	}
}

// ---- Planar RGB -------------------------------------------------------

func TestConvertGBRP_8bit(t *testing.T) {
	// Planes: G=100, B=50, R=200
	frame := pxFrame(1, 1,
		[]byte{100}, // G plane
		[]byte{50},  // B plane
		[]byte{200}, // R plane
	)
	out := convertRawToRGBA(canvas.PixelFormatGBRP, frame)
	expectRGBA(t, "GBRP/8bit", out, 0, 0, color.RGBA{200, 100, 50, 255}, 0)
}

func TestConvertGBRP_10bit(t *testing.T) {
	// 10-bit: G=0 (plane0), B=0 (plane1), R=1023 (plane2) → pure red at full brightness
	r10 := u16LE(1023)
	frame := pxFrame(1, 1,
		[]byte{0x00, 0x00}, // G plane (0)
		[]byte{0x00, 0x00}, // B plane (0)
		r10,                // R plane (1023 → uint8 255)
	)
	out := convertRawToRGBA(canvas.PixelFormatGBRP10, frame)
	expectRGBA(t, "GBRP10/red", out, 0, 0, color.RGBA{255, 0, 0, 255}, 0)
}

func TestConvertGBRP_16bit_MidGray(t *testing.T) {
	// 16-bit mid gray: all channels at 32767 (≈ 50% of 65535)
	// 32767/65535 * 255 ≈ 127.5 → round to 128
	mid := u16LE(32767)
	frame := pxFrame(1, 1, mid, mid, mid)
	out := convertRawToRGBA(canvas.PixelFormatGBRP16, frame)
	// tolerance 1 because of rounding
	expectRGBA(t, "GBRP16/gray", out, 0, 0, color.RGBA{128, 128, 128, 255}, 1)
}

func TestConvertGBRAP_8bit(t *testing.T) {
	// 4 planes: G=100, B=50, R=200, A=77
	frame := pxFrame(1, 1,
		[]byte{100}, // G
		[]byte{50},  // B
		[]byte{200}, // R
		[]byte{77},  // A
	)
	out := convertRawToRGBA(canvas.PixelFormatGBRAP, frame)
	expectRGBA(t, "GBRAP/8bit", out, 0, 0, color.RGBA{200, 100, 50, 77}, 0)
}

// ---- Monochrome -------------------------------------------------------

func TestConvertMONOWHITE_WhiteAndBlack(t *testing.T) {
	// MONOWHITE: 0-bit = white (255), 1-bit = black (0), MSB first.
	// 1 byte = 8 pixels for an 8x1 frame.
	// 0x00 → all zero bits → all white
	// 0xFF → all one bits → all black
	t.Run("all_white", func(t *testing.T) {
		frame := pxFrame(8, 1, []byte{0x00})
		out := convertRawToRGBA(canvas.PixelFormatMONOWHITE, frame)
		for x := 0; x < 8; x++ {
			expectRGBA(t, "MONOWHITE/white", out, x, 0, color.RGBA{255, 255, 255, 255}, 0)
		}
	})
	t.Run("all_black", func(t *testing.T) {
		frame := pxFrame(8, 1, []byte{0xFF})
		out := convertRawToRGBA(canvas.PixelFormatMONOWHITE, frame)
		for x := 0; x < 8; x++ {
			expectRGBA(t, "MONOWHITE/black", out, x, 0, color.RGBA{0, 0, 0, 255}, 0)
		}
	})
	t.Run("alternating", func(t *testing.T) {
		// 0b10101010 = 0xAA, MSB first → bit values: 1,0,1,0,1,0,1,0
		// MONOWHITE: bit=1 → black(0), bit=0 → white(255)
		// So pixels: B,W,B,W,B,W,B,W
		frame := pxFrame(8, 1, []byte{0xAA})
		out := convertRawToRGBA(canvas.PixelFormatMONOWHITE, frame)
		for x := 0; x < 8; x++ {
			if x%2 == 0 {
				expectRGBA(t, "MONOWHITE/alt/black", out, x, 0, color.RGBA{0, 0, 0, 255}, 0)
			} else {
				expectRGBA(t, "MONOWHITE/alt/white", out, x, 0, color.RGBA{255, 255, 255, 255}, 0)
			}
		}
	})
}

func TestConvertMONOBLACK_WhiteAndBlack(t *testing.T) {
	// MONOBLACK: 0-bit = black (0), 1-bit = white (255), MSB first.
	t.Run("all_black", func(t *testing.T) {
		frame := pxFrame(8, 1, []byte{0x00})
		out := convertRawToRGBA(canvas.PixelFormatMONOBLACK, frame)
		for x := 0; x < 8; x++ {
			expectRGBA(t, "MONOBLACK/black", out, x, 0, color.RGBA{0, 0, 0, 255}, 0)
		}
	})
	t.Run("all_white", func(t *testing.T) {
		frame := pxFrame(8, 1, []byte{0xFF})
		out := convertRawToRGBA(canvas.PixelFormatMONOBLACK, frame)
		for x := 0; x < 8; x++ {
			expectRGBA(t, "MONOBLACK/white", out, x, 0, color.RGBA{255, 255, 255, 255}, 0)
		}
	})
}

// ---- Bayer demosaic ---------------------------------------------------

// bayerUniform verifies that a uniform Bayer frame (all pixels the same value)
// demosaics to a flat gray image, regardless of pattern.
func bayerUniform(t *testing.T, format int, label string) {
	t.Helper()
	const value = 128
	// 4×4 frame, all pixels = value
	data := make([]byte, 4*4)
	for i := range data {
		data[i] = value
	}
	frame := pxFrame(4, 4, data)
	out := convertRawToRGBA(format, frame)
	// Inner 2×2 (avoid border clamping artefacts at edges)
	for y := 1; y <= 2; y++ {
		for x := 1; x <= 2; x++ {
			expectRGBA(t, label, out, x, y, color.RGBA{value, value, value, 255}, 1)
		}
	}
}

func TestConvertBayer_UniformFrames(t *testing.T) {
	bayerUniform(t, canvas.PixelFormatBAYER_BGGR8, "BGGR")
	bayerUniform(t, canvas.PixelFormatBAYER_RGGB8, "RGGB")
	bayerUniform(t, canvas.PixelFormatBAYER_GBRG8, "GBRG")
	bayerUniform(t, canvas.PixelFormatBAYER_GRBG8, "GRBG")
}

func TestConvertBayer_RGGB_PureRed(t *testing.T) {
	// 4×4 RGGB: only R pixels (even col, even row) are lit; G and B = 0.
	// RGGB pattern: R G / G B
	// For a 4x4 frame with R pixels = 200, G = 0, B = 0:
	//   raw[y][x] = 200 if (x%2==0 && y%2==0), else 0
	data := make([]byte, 4*4)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if x%2 == 0 && y%2 == 0 {
				data[y*4+x] = 200
			}
		}
	}
	frame := pxFrame(4, 4, data)
	out := convertRawToRGBA(canvas.PixelFormatBAYER_RGGB8, frame)

	// The interior R pixel at (2,2) should have high R and low G,B.
	got := color.RGBAModel.Convert(out.At(2, 2)).(color.RGBA)
	if got.R < 150 {
		t.Errorf("RGGB pure-red: R pixel at (2,2) has R=%d, expected ≥150", got.R)
	}
	if got.G > 50 || got.B > 50 {
		t.Errorf("RGGB pure-red: R pixel at (2,2) has G=%d B=%d, expected near 0", got.G, got.B)
	}
}

func TestConvertBayer_16bit_Uniform(t *testing.T) {
	// 16-bit RGGB: all pixels = 500 (out of 65535 range, ~0.76% brightness)
	// Expected gray ≈ round(500/65535*255) = round(1.946) = 2
	const val = uint16(500)
	data := make([]byte, 4*4*2)
	for i := 0; i < 4*4; i++ {
		binary.LittleEndian.PutUint16(data[i*2:], val)
	}
	frame := pxFrame(4, 4, data)
	out := convertRawToRGBA(canvas.PixelFormatBAYER_RGGB16, frame)
	expected := uint8(math.Round(float64(val) / 65535.0 * 255))
	for y := 1; y <= 2; y++ {
		for x := 1; x <= 2; x++ {
			expectRGBA(t, "Bayer16/uniform", out, x, y, color.RGBA{expected, expected, expected, 255}, 1)
		}
	}
}

// ---- Float formats ----------------------------------------------------

func TestConvertGRAYF32_Values(t *testing.T) {
	tests := []struct {
		name  string
		val   float32
		wantV uint8
	}{
		{"white", 1.0, 255},
		{"black", 0.0, 0},
		{"mid", 0.5, 128},
		{"clamp_high", 2.0, 255},
		{"clamp_low", -1.0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := pxFrame(1, 1, f32LE(tt.val))
			out := convertRawToRGBA(canvas.PixelFormatGRAYF32, frame)
			want := color.RGBA{tt.wantV, tt.wantV, tt.wantV, 255}
			expectRGBA(t, "GRAYF32/"+tt.name, out, 0, 0, want, 1)
		})
	}
}

func TestConvertRGBF32_PrimaryColors(t *testing.T) {
	tests := []struct {
		name string
		r, g, b float32
		want color.RGBA
	}{
		{"red", 1.0, 0.0, 0.0, color.RGBA{255, 0, 0, 255}},
		{"green", 0.0, 1.0, 0.0, color.RGBA{0, 255, 0, 255}},
		{"blue", 0.0, 0.0, 1.0, color.RGBA{0, 0, 255, 255}},
		{"white", 1.0, 1.0, 1.0, color.RGBA{255, 255, 255, 255}},
		{"mid_gray", 0.5, 0.5, 0.5, color.RGBA{128, 128, 128, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := append(append(f32LE(tt.r), f32LE(tt.g)...), f32LE(tt.b)...)
			frame := pxFrame(1, 1, data)
			out := convertRawToRGBA(canvas.PixelFormatRGBF32, frame)
			expectRGBA(t, "RGBF32/"+tt.name, out, 0, 0, tt.want, 1)
		})
	}
}

func TestConvertRGBAF32_WithAlpha(t *testing.T) {
	// R=1.0, G=0.0, B=0.5, A=0.25
	data := append(append(append(f32LE(1.0), f32LE(0.0)...), f32LE(0.5)...), f32LE(0.25)...)
	frame := pxFrame(1, 1, data)
	out := convertRawToRGBA(canvas.PixelFormatRGBAF32, frame)
	expectRGBA(t, "RGBAF32", out, 0, 0,
		color.RGBA{255, 0, 128, 64}, // 0.5*255≈128, 0.25*255≈64
		1)
}

// ---- Stride handling --------------------------------------------------

func TestConvertRGB24_WithStride(t *testing.T) {
	// 2×1 frame with a stride of 8 bytes (padding after each row).
	// pixel(0,0)=red, pixel(1,0)=blue, 2 pad bytes after.
	data := []byte{
		255, 0, 0, // R=red
		0, 0, 255, // R=blue
		0, 0, // padding (stride=8)
	}
	frame := &canvas.RawFrame{
		Width: 2, Height: 1,
		Data:    [4][]byte{data},
		Strides: [4]int{8},
	}
	out := convertRawToRGBA(canvas.PixelFormatRGB24, frame)
	expectRGBA(t, "RGB24/stride/px0", out, 0, 0, color.RGBA{255, 0, 0, 255}, 0)
	expectRGBA(t, "RGB24/stride/px1", out, 1, 0, color.RGBA{0, 0, 255, 255}, 0)
}

// ---- Nil / unsupported format handling --------------------------------

func TestConvertRawToRGBA_NilOnEmpty(t *testing.T) {
	frame := pxFrame(1, 1) // no data
	out := convertRawToRGBA(canvas.PixelFormatRGB24, frame)
	if out != nil {
		t.Error("expected nil for empty frame data")
	}
}

func TestConvertRawToRGBA_NilOnUnknown(t *testing.T) {
	frame := pxFrame(1, 1, []byte{0xFF})
	// Use a format value that is not handled by convertRawToRGBA.
	out := convertRawToRGBA(-999, frame)
	if out != nil {
		t.Error("expected nil for unknown pixel format")
	}
}

// ---- Multi-pixel consistency: all formats produce correct dimensions --

func TestConvertAllFormats_OutputDimensions(t *testing.T) {
	const w, h = 4, 4

	type fmtCase struct {
		name   string
		format int
		data   func() [4][]byte
	}

	cases := []fmtCase{
		{"RGB24", canvas.PixelFormatRGB24, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*3)}
		}},
		{"BGR24", canvas.PixelFormatBGR24, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*3)}
		}},
		{"BGRA", canvas.PixelFormatBGRA, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*4)}
		}},
		{"ARGB", canvas.PixelFormatARGB, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*4)}
		}},
		{"ABGR", canvas.PixelFormatABGR, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*4)}
		}},
		{"RGB565", canvas.PixelFormatRGB565, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*2)}
		}},
		{"BGR565", canvas.PixelFormatBGR565, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*2)}
		}},
		{"GBRP", canvas.PixelFormatGBRP, func() [4][]byte {
			return [4][]byte{make([]byte, w*h), make([]byte, w*h), make([]byte, w*h)}
		}},
		{"GBRP10", canvas.PixelFormatGBRP10, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*2), make([]byte, w*h*2), make([]byte, w*h*2)}
		}},
		{"GBRAP", canvas.PixelFormatGBRAP, func() [4][]byte {
			return [4][]byte{make([]byte, w*h), make([]byte, w*h), make([]byte, w*h), make([]byte, w*h)}
		}},
		{"MONOWHITE", canvas.PixelFormatMONOWHITE, func() [4][]byte {
			rowBytes := (w + 7) / 8 // bytes per row (each row packed independently)
			return [4][]byte{make([]byte, rowBytes*h)}
		}},
		{"MONOBLACK", canvas.PixelFormatMONOBLACK, func() [4][]byte {
			rowBytes := (w + 7) / 8
			return [4][]byte{make([]byte, rowBytes*h)}
		}},
		{"BGGR8", canvas.PixelFormatBAYER_BGGR8, func() [4][]byte {
			return [4][]byte{make([]byte, w*h)}
		}},
		{"GRAYF32", canvas.PixelFormatGRAYF32, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*4)}
		}},
		{"RGBF32", canvas.PixelFormatRGBF32, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*12)}
		}},
		{"RGBAF32", canvas.PixelFormatRGBAF32, func() [4][]byte {
			return [4][]byte{make([]byte, w*h*16)}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := c.data()
			frame := &canvas.RawFrame{Width: w, Height: h, Data: d}
			out := convertRawToRGBA(c.format, frame)
			if out == nil {
				t.Fatalf("%s: convertRawToRGBA returned nil", c.name)
			}
			if out.Bounds().Dx() != w || out.Bounds().Dy() != h {
				t.Errorf("%s: got size %dx%d, want %dx%d",
					c.name, out.Bounds().Dx(), out.Bounds().Dy(), w, h)
			}
		})
	}
}
