package gl

import (
	"testing"

	"fyne.io/fyne/v2/canvas"
)

// TestFormatDescriptor_AllFormatsHaveEntry verifies every public PixelFormat
// constant resolves to a non-zero descriptor (i.e. is present in the map).
func TestFormatDescriptor_AllFormatsHaveEntry(t *testing.T) {
	formats := []struct {
		name   string
		format int
	}{
		{"RGBA", canvas.PixelFormatRGBA},
		{"YUV420P", canvas.PixelFormatYUV420P},
		{"YUV422P", canvas.PixelFormatYUV422P},
		{"YUV444P", canvas.PixelFormatYUV444P},
		{"YUV410P", canvas.PixelFormatYUV410P},
		{"YUV411P", canvas.PixelFormatYUV411P},
		{"YUV440P", canvas.PixelFormatYUV440P},
		{"YUVJ420P", canvas.PixelFormatYUVJ420P},
		{"YUVJ422P", canvas.PixelFormatYUVJ422P},
		{"YUVJ444P", canvas.PixelFormatYUVJ444P},
		{"YUVJ440P", canvas.PixelFormatYUVJ440P},
		{"YUVA420P", canvas.PixelFormatYUVA420P},
		{"YUVA422P", canvas.PixelFormatYUVA422P},
		{"YUVA444P", canvas.PixelFormatYUVA444P},
		{"YUV420P9", canvas.PixelFormatYUV420P9},
		{"YUV420P10", canvas.PixelFormatYUV420P10},
		{"YUV420P12", canvas.PixelFormatYUV420P12},
		{"YUV420P14", canvas.PixelFormatYUV420P14},
		{"YUV420P16", canvas.PixelFormatYUV420P16},
		{"YUV422P9", canvas.PixelFormatYUV422P9},
		{"YUV422P10", canvas.PixelFormatYUV422P10},
		{"YUV422P12", canvas.PixelFormatYUV422P12},
		{"YUV422P14", canvas.PixelFormatYUV422P14},
		{"YUV422P16", canvas.PixelFormatYUV422P16},
		{"YUV444P9", canvas.PixelFormatYUV444P9},
		{"YUV444P10", canvas.PixelFormatYUV444P10},
		{"YUV444P12", canvas.PixelFormatYUV444P12},
		{"YUV444P14", canvas.PixelFormatYUV444P14},
		{"YUV444P16", canvas.PixelFormatYUV444P16},
		{"YUVA420P9", canvas.PixelFormatYUVA420P9},
		{"YUVA420P10", canvas.PixelFormatYUVA420P10},
		{"YUVA420P16", canvas.PixelFormatYUVA420P16},
		{"YUVA422P9", canvas.PixelFormatYUVA422P9},
		{"YUVA422P10", canvas.PixelFormatYUVA422P10},
		{"YUVA422P16", canvas.PixelFormatYUVA422P16},
		{"YUVA444P9", canvas.PixelFormatYUVA444P9},
		{"YUVA444P10", canvas.PixelFormatYUVA444P10},
		{"YUVA444P16", canvas.PixelFormatYUVA444P16},
		{"NV12", canvas.PixelFormatNV12},
		{"NV21", canvas.PixelFormatNV21},
		{"NV16", canvas.PixelFormatNV16},
		{"NV24", canvas.PixelFormatNV24},
		{"NV42", canvas.PixelFormatNV42},
		{"P010", canvas.PixelFormatP010},
		{"P012", canvas.PixelFormatP012},
		{"P016", canvas.PixelFormatP016},
		{"YUYV422", canvas.PixelFormatYUYV422},
		{"UYVY422", canvas.PixelFormatUYVY422},
		{"YVYU422", canvas.PixelFormatYVYU422},
		{"RGB24", canvas.PixelFormatRGB24},
		{"BGR24", canvas.PixelFormatBGR24},
		{"BGRA", canvas.PixelFormatBGRA},
		{"ARGB", canvas.PixelFormatARGB},
		{"ABGR", canvas.PixelFormatABGR},
		{"RGB565", canvas.PixelFormatRGB565},
		{"BGR565", canvas.PixelFormatBGR565},
		{"GBRP", canvas.PixelFormatGBRP},
		{"GBRP10", canvas.PixelFormatGBRP10},
		{"GBRP12", canvas.PixelFormatGBRP12},
		{"GBRP16", canvas.PixelFormatGBRP16},
		{"GBRAP", canvas.PixelFormatGBRAP},
		{"GRAY8", canvas.PixelFormatGRAY8},
		{"GRAY16", canvas.PixelFormatGRAY16},
		{"YA8", canvas.PixelFormatYA8},
		{"MONOWHITE", canvas.PixelFormatMONOWHITE},
		{"MONOBLACK", canvas.PixelFormatMONOBLACK},
		{"BAYER_BGGR8", canvas.PixelFormatBAYER_BGGR8},
		{"BAYER_RGGB8", canvas.PixelFormatBAYER_RGGB8},
		{"BAYER_GBRG8", canvas.PixelFormatBAYER_GBRG8},
		{"BAYER_GRBG8", canvas.PixelFormatBAYER_GRBG8},
		{"BAYER_BGGR16", canvas.PixelFormatBAYER_BGGR16},
		{"BAYER_RGGB16", canvas.PixelFormatBAYER_RGGB16},
		{"BAYER_GBRG16", canvas.PixelFormatBAYER_GBRG16},
		{"BAYER_GRBG16", canvas.PixelFormatBAYER_GRBG16},
		{"GRAYF32", canvas.PixelFormatGRAYF32},
		{"RGBF32", canvas.PixelFormatRGBF32},
		{"RGBAF32", canvas.PixelFormatRGBAF32},
	}

	for _, tc := range formats {
		t.Run(tc.name, func(t *testing.T) {
			desc, ok := formatDescriptors[tc.format]
			if !ok {
				t.Fatalf("format %q (value %d) not found in formatDescriptors map", tc.name, tc.format)
			}
			// CPU-convert formats have planeCount=0 because they convert to RGBA on CPU
			// before any GL upload; only GPU-path formats need planeCount >= 1.
			if desc.category != categoryCPUConvert && desc.planeCount < 1 {
				t.Errorf("format %q (GPU path): planeCount=%d, expected >= 1", tc.name, desc.planeCount)
			}
		})
	}
}

// TestFormatDescriptor_CategoryConsistency checks invariants for each category.
func TestFormatDescriptor_CategoryConsistency(t *testing.T) {
	for fmtVal, desc := range formatDescriptors {
		switch desc.category {
		case categoryYUVPlanar:
			if desc.planeCount != 3 {
				t.Errorf("format %d (categoryYUVPlanar): planeCount=%d, expected 3", fmtVal, desc.planeCount)
			}
			if desc.bitDepth != 8 {
				t.Errorf("format %d (categoryYUVPlanar): bitDepth=%d, expected 8", fmtVal, desc.bitDepth)
			}
		case categoryYUVAPlanar:
			if desc.planeCount != 4 {
				t.Errorf("format %d (categoryYUVAPlanar): planeCount=%d, expected 4", fmtVal, desc.planeCount)
			}
			if !desc.hasAlpha {
				t.Errorf("format %d (categoryYUVAPlanar): hasAlpha=false, expected true", fmtVal)
			}
		case categoryNVSemiplanar:
			if desc.planeCount != 2 {
				t.Errorf("format %d (categoryNVSemiplanar): planeCount=%d, expected 2", fmtVal, desc.planeCount)
			}
		case categoryPackedYUV422:
			if desc.planeCount != 1 {
				t.Errorf("format %d (categoryPackedYUV422): planeCount=%d, expected 1", fmtVal, desc.planeCount)
			}
			if desc.packingMode < 0 || desc.packingMode > 2 {
				t.Errorf("format %d (categoryPackedYUV422): packingMode=%d, expected 0-2", fmtVal, desc.packingMode)
			}
		case categoryHiBitPlanar:
			if desc.planeCount < 3 {
				t.Errorf("format %d (categoryHiBitPlanar): planeCount=%d, expected >= 3", fmtVal, desc.planeCount)
			}
			if desc.bitDepth <= 8 {
				t.Errorf("format %d (categoryHiBitPlanar): bitDepth=%d, expected > 8", fmtVal, desc.bitDepth)
			}
		case categoryHiBitNV:
			if desc.planeCount != 2 {
				t.Errorf("format %d (categoryHiBitNV): planeCount=%d, expected 2", fmtVal, desc.planeCount)
			}
			if desc.bitDepth <= 8 {
				t.Errorf("format %d (categoryHiBitNV): bitDepth=%d, expected > 8", fmtVal, desc.bitDepth)
			}
		}
	}
}

// TestFormatDescriptor_NVSwapUV verifies NV21/NV42 have swapUV=true, NV12 has swapUV=false.
func TestFormatDescriptor_NVSwapUV(t *testing.T) {
	nv12 := formatDescriptors[canvas.PixelFormatNV12]
	if nv12.swapUV {
		t.Error("NV12 should have swapUV=false")
	}

	nv21 := formatDescriptors[canvas.PixelFormatNV21]
	if !nv21.swapUV {
		t.Error("NV21 should have swapUV=true")
	}

	nv42 := formatDescriptors[canvas.PixelFormatNV42]
	if !nv42.swapUV {
		t.Error("NV42 should have swapUV=true")
	}
}

// TestFormatDescriptor_PackedYUV422Modes verifies packing modes for YUYV/UYVY/YVYU.
func TestFormatDescriptor_PackedYUV422Modes(t *testing.T) {
	yuyv := formatDescriptors[canvas.PixelFormatYUYV422]
	if yuyv.packingMode != 0 {
		t.Errorf("YUYV422: packingMode=%d, expected 0", yuyv.packingMode)
	}

	uyvy := formatDescriptors[canvas.PixelFormatUYVY422]
	if uyvy.packingMode != 1 {
		t.Errorf("UYVY422: packingMode=%d, expected 1", uyvy.packingMode)
	}

	yvyu := formatDescriptors[canvas.PixelFormatYVYU422]
	if yvyu.packingMode != 2 {
		t.Errorf("YVYU422: packingMode=%d, expected 2", yvyu.packingMode)
	}
}

// TestFormatDescriptor_YUVPlanarChromaSubsampling checks chroma divisors for common YUV formats.
func TestFormatDescriptor_YUVPlanarChromaSubsampling(t *testing.T) {
	cases := []struct {
		name           string
		format         int
		chromaWDiv     [3]int // for planes 0,1,2
		chromaHDiv     [3]int
	}{
		{"YUV420P", canvas.PixelFormatYUV420P, [3]int{1, 2, 2}, [3]int{1, 2, 2}},
		{"YUV422P", canvas.PixelFormatYUV422P, [3]int{1, 2, 2}, [3]int{1, 1, 1}},
		{"YUV444P", canvas.PixelFormatYUV444P, [3]int{1, 1, 1}, [3]int{1, 1, 1}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			desc := formatDescriptors[tc.format]
			for i := 0; i < 3; i++ {
				if desc.chromaWDiv[i] != tc.chromaWDiv[i] {
					t.Errorf("plane %d chromaWDiv: got %d, want %d", i, desc.chromaWDiv[i], tc.chromaWDiv[i])
				}
				if desc.chromaHDiv[i] != tc.chromaHDiv[i] {
					t.Errorf("plane %d chromaHDiv: got %d, want %d", i, desc.chromaHDiv[i], tc.chromaHDiv[i])
				}
			}
		})
	}
}

// TestFormatDescriptor_HiBitDepths checks bit depths for hi-bit planar formats.
func TestFormatDescriptor_HiBitDepths(t *testing.T) {
	cases := []struct {
		name     string
		format   int
		bitDepth int
	}{
		{"YUV420P9", canvas.PixelFormatYUV420P9, 9},
		{"YUV420P10", canvas.PixelFormatYUV420P10, 10},
		{"YUV420P12", canvas.PixelFormatYUV420P12, 12},
		{"YUV420P14", canvas.PixelFormatYUV420P14, 14},
		{"YUV420P16", canvas.PixelFormatYUV420P16, 16},
		{"P010", canvas.PixelFormatP010, 10},
		{"P016", canvas.PixelFormatP016, 16},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			desc := formatDescriptors[tc.format]
			if desc.bitDepth != tc.bitDepth {
				t.Errorf("bitDepth: got %d, want %d", desc.bitDepth, tc.bitDepth)
			}
		})
	}
}

// TestFormatDescriptor_CPUConvertFormats checks that CPU-convert formats have correct category.
func TestFormatDescriptor_CPUConvertFormats(t *testing.T) {
	cpuFmts := []struct {
		name   string
		format int
	}{
		{"BGR24", canvas.PixelFormatBGR24},
		{"BGRA", canvas.PixelFormatBGRA},
		{"ARGB", canvas.PixelFormatARGB},
		{"ABGR", canvas.PixelFormatABGR},
		{"RGB565", canvas.PixelFormatRGB565},
		{"BGR565", canvas.PixelFormatBGR565},
		{"GBRP", canvas.PixelFormatGBRP},
		{"GBRP10", canvas.PixelFormatGBRP10},
		{"GBRAP", canvas.PixelFormatGBRAP},
		{"MONOWHITE", canvas.PixelFormatMONOWHITE},
		{"MONOBLACK", canvas.PixelFormatMONOBLACK},
		{"BAYER_BGGR8", canvas.PixelFormatBAYER_BGGR8},
		{"GRAYF32", canvas.PixelFormatGRAYF32},
		{"RGBF32", canvas.PixelFormatRGBF32},
		{"RGBAF32", canvas.PixelFormatRGBAF32},
	}

	for _, tc := range cpuFmts {
		t.Run(tc.name, func(t *testing.T) {
			desc := formatDescriptors[tc.format]
			if desc.category != categoryCPUConvert {
				t.Errorf("%s: category=%d, expected categoryCPUConvert (%d)",
					tc.name, desc.category, categoryCPUConvert)
			}
		})
	}
}

// TestFormatDescriptorFor_UnknownFormat verifies unknown format codes default to categoryCPUConvert
// (safe fallback: attempt CPU conversion to RGBA, which will return nil and skip draw).
func TestFormatDescriptorFor_UnknownFormat(t *testing.T) {
	desc := formatDescriptorFor(9999)
	if desc.category != categoryCPUConvert {
		t.Errorf("unknown format: category=%d, expected categoryCPUConvert (%d)", desc.category, categoryCPUConvert)
	}
}

// TestFormatDescriptor_GrayscaleHasAlpha checks YA8 has hasAlpha=true, GRAY8 does not.
func TestFormatDescriptor_GrayscaleHasAlpha(t *testing.T) {
	gray8 := formatDescriptors[canvas.PixelFormatGRAY8]
	if gray8.hasAlpha {
		t.Error("GRAY8 should have hasAlpha=false")
	}

	ya8 := formatDescriptors[canvas.PixelFormatYA8]
	if !ya8.hasAlpha {
		t.Error("YA8 should have hasAlpha=true")
	}
}
