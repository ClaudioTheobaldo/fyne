package gl

import "fyne.io/fyne/v2/canvas"

// shaderCategory groups pixel formats by GPU sampling pattern, determining
// which shader program and upload path to use.
type shaderCategory int

const (
	categorySimple        shaderCategory = iota // RGBA (existing simple program)
	categoryYUVPlanar                           // 3-plane 8-bit YUV
	categoryYUVAPlanar                          // 4-plane 8-bit YUVA (with alpha)
	categoryNVSemiplanar                        // Y + interleaved UV (NV12/NV21/etc.)
	categoryPackedYUV422                        // Packed YUYV/UYVY/YVYU
	categoryGrayscale                           // GRAY8, GRAY16, YA8, MONO
	categoryHiBitPlanar                         // 9–16 bit planar YUV (3 planes)
	categoryHiBitNV                             // 9–16 bit semi-planar (P010/P016)
	categoryCPUConvert                          // Exotic formats: convert to RGBA on CPU
)

// formatDescriptor describes how to upload and draw a pixel format.
type formatDescriptor struct {
	category       shaderCategory
	planeCount     int
	planeFormats   [4]uint32 // GL internal format per plane
	planeDataTypes [4]uint32 // GL data type per plane (fdByte or unsignedShort)
	chromaWDiv     [4]int    // width divisor per plane (1=full, 2=half, etc.)
	chromaHDiv     [4]int    // height divisor per plane
	hasAlpha       bool      // true if format includes an alpha channel
	swapUV         bool      // true for NV21/NV42 (UV byte order swapped)
	bitDepth       int       // 8, 9, 10, 12, 14, or 16
	packingMode    int       // for packed YUV422: 0=YUYV, 1=UYVY, 2=YVYU
	defaultRange   int       // canvas.ColorRangeFull or ColorRangeLimited
}

// Standard OpenGL format/type constants for use in format descriptors.
// Using explicit uint32 hex literals makes this file independent of the
// platform-specific gl_*.go constant definitions (which may be typed gl.Enum).
const (
	fdLum  uint32 = 0x1909 // GL_LUMINANCE
	fdLumA uint32 = 0x190A // GL_LUMINANCE_ALPHA
	fdRGBA uint32 = 0x1908 // GL_RGBA
	fdRGB  uint32 = 0x1907 // GL_RGB
	fdByte uint32 = 0x1401 // GL_UNSIGNED_BYTE
	fdU16  uint32 = 0x1403 // GL_UNSIGNED_SHORT (unused for now; hi-bit uses fdByte+packing)
)

// formatDescriptors maps each PixelFormat constant to its upload/draw descriptor.
// Formats not in this map are treated as categoryCPUConvert.
var formatDescriptors = map[int]formatDescriptor{
	// ── RGBA ─────────────────────────────────────────────────────────────────
	canvas.PixelFormatRGBA: {
		category:     categorySimple,
		planeCount:   1,
		planeFormats: [4]uint32{fdRGBA},
		planeDataTypes: [4]uint32{fdByte},
		chromaWDiv:   [4]int{1},
		chromaHDiv:   [4]int{1},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeFull,
	},

	// ── 8-bit planar YUV ─────────────────────────────────────────────────────
	canvas.PixelFormatYUV420P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVJ420P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeFull, // YUVJ = JPEG / full-range
	},
	canvas.PixelFormatYUV422P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVJ422P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeFull,
	},
	canvas.PixelFormatYUV444P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVJ444P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeFull,
	},
	canvas.PixelFormatYUV440P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 2, 2}, // vertical 2:1 subsampling
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVJ440P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeFull,
	},
	canvas.PixelFormatYUV410P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 4, 4},
		chromaHDiv:   [4]int{1, 4, 4},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV411P: {
		category:     categoryYUVPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 4, 4},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},

	// ── 8-bit YUVA (with separate alpha plane) ────────────────────────────────
	canvas.PixelFormatYUVA420P: {
		category:     categoryYUVAPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 2, 2, 1},
		hasAlpha:     true,
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA422P: {
		category:     categoryYUVAPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA444P: {
		category:     categoryYUVAPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLum, fdLum, fdLum, fdLum},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     8,
		defaultRange: canvas.ColorRangeLimited,
	},

	// ── High bit-depth planar YUV (10/12/14/16-bit, LE) ──────────────────────
	// Uploaded as GL_LUMINANCE_ALPHA (2 bytes per sample, LE packing).
	canvas.PixelFormatYUV420P9: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     9,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV420P10: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     10,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV420P12: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     12,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV420P14: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     14,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV420P16: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 2, 2},
		bitDepth:     16,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV422P9: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     9,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV422P10: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     10,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV422P12: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     12,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV422P14: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     14,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV422P16: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     16,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV444P9: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     9,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV444P10: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     10,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV444P12: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     12,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV444P14: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     14,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUV444P16: {
		category:     categoryHiBitPlanar,
		planeCount:   3,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1},
		bitDepth:     16,
		defaultRange: canvas.ColorRangeLimited,
	},

	// High bit-depth YUVA (same as above + alpha plane)
	canvas.PixelFormatYUVA420P9: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 2, 2, 1},
		hasAlpha:     true,
		bitDepth:     9,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA420P10: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 2, 2, 1},
		hasAlpha:     true,
		bitDepth:     10,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA420P16: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 2, 2, 1},
		hasAlpha:     true,
		bitDepth:     16,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA422P9: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     9,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA422P10: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     10,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA422P16: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2, 2, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     16,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA444P9: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     9,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA444P10: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     10,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYUVA444P16: {
		category:     categoryHiBitPlanar,
		planeCount:   4,
		planeFormats: [4]uint32{fdLumA, fdLumA, fdLumA, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte, fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1, 1, 1},
		chromaHDiv:   [4]int{1, 1, 1, 1},
		hasAlpha:     true,
		bitDepth:     16,
		defaultRange: canvas.ColorRangeLimited,
	},

	// ── Semi-planar NV ────────────────────────────────────────────────────────
	// UV plane is GL_LUMINANCE_ALPHA: .r = U (or V for NV21), .a = V (or U).
	canvas.PixelFormatNV12: {
		category:     categoryNVSemiplanar,
		planeCount:   2,
		planeFormats: [4]uint32{fdLum, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2},
		chromaHDiv:   [4]int{1, 2},
		bitDepth:     8,
		swapUV:       false,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatNV21: {
		category:     categoryNVSemiplanar,
		planeCount:   2,
		planeFormats: [4]uint32{fdLum, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2},
		chromaHDiv:   [4]int{1, 2},
		bitDepth:     8,
		swapUV:       true,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatNV16: {
		category:     categoryNVSemiplanar,
		planeCount:   2,
		planeFormats: [4]uint32{fdLum, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2},
		chromaHDiv:   [4]int{1, 1},
		bitDepth:     8,
		swapUV:       false,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatNV24: {
		category:     categoryNVSemiplanar,
		planeCount:   2,
		planeFormats: [4]uint32{fdLum, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1},
		chromaHDiv:   [4]int{1, 1},
		bitDepth:     8,
		swapUV:       false,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatNV42: {
		category:     categoryNVSemiplanar,
		planeCount:   2,
		planeFormats: [4]uint32{fdLum, fdLumA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 1},
		chromaHDiv:   [4]int{1, 1},
		bitDepth:     8,
		swapUV:       true,
		defaultRange: canvas.ColorRangeLimited,
	},

	// High bit-depth NV: Y uploaded as GL_LUMINANCE_ALPHA (2 bytes/sample, LE).
	// UV uploaded as GL_RGBA at half-width: .r=U_lo .g=U_hi .b=V_lo .a=V_hi per UV pair.
	// chromaWDiv[1]=2 → UV texture width = frame_width/2 texels, each covering one UV pair.
	canvas.PixelFormatP010: {
		category:     categoryHiBitNV,
		planeCount:   2,
		planeFormats: [4]uint32{fdLumA, fdRGBA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2},
		chromaHDiv:   [4]int{1, 2},
		bitDepth:     10,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatP012: {
		category:     categoryHiBitNV,
		planeCount:   2,
		planeFormats: [4]uint32{fdLumA, fdRGBA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2},
		chromaHDiv:   [4]int{1, 2},
		bitDepth:     12,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatP016: {
		category:     categoryHiBitNV,
		planeCount:   2,
		planeFormats: [4]uint32{fdLumA, fdRGBA},
		planeDataTypes: [4]uint32{fdByte, fdByte},
		chromaWDiv:   [4]int{1, 2},
		chromaHDiv:   [4]int{1, 2},
		bitDepth:     16,
		defaultRange: canvas.ColorRangeLimited,
	},

	// ── Packed YUV 4:2:2 ─────────────────────────────────────────────────────
	// Uploaded at half-width as RGBA (4 bytes = 2 pixels).
	// The shader decodes the Y/U/V from the packed texel.
	canvas.PixelFormatYUYV422: {
		category:     categoryPackedYUV422,
		planeCount:   1,
		planeFormats: [4]uint32{fdRGBA},
		planeDataTypes: [4]uint32{fdByte},
		chromaWDiv:   [4]int{2}, // upload width = frame_width / 2
		chromaHDiv:   [4]int{1},
		bitDepth:     8,
		packingMode:  0,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatUYVY422: {
		category:     categoryPackedYUV422,
		planeCount:   1,
		planeFormats: [4]uint32{fdRGBA},
		planeDataTypes: [4]uint32{fdByte},
		chromaWDiv:   [4]int{2},
		chromaHDiv:   [4]int{1},
		bitDepth:     8,
		packingMode:  1,
		defaultRange: canvas.ColorRangeLimited,
	},
	canvas.PixelFormatYVYU422: {
		category:     categoryPackedYUV422,
		planeCount:   1,
		planeFormats: [4]uint32{fdRGBA},
		planeDataTypes: [4]uint32{fdByte},
		chromaWDiv:   [4]int{2},
		chromaHDiv:   [4]int{1},
		bitDepth:     8,
		packingMode:  2,
		defaultRange: canvas.ColorRangeLimited,
	},

	// ── Grayscale ─────────────────────────────────────────────────────────────
	canvas.PixelFormatGRAY8: {
		category:     categoryGrayscale,
		planeCount:   1,
		planeFormats: [4]uint32{fdLum},
		planeDataTypes: [4]uint32{fdByte},
		chromaWDiv:   [4]int{1},
		chromaHDiv:   [4]int{1},
		bitDepth:     8,
		hasAlpha:     false,
		defaultRange: canvas.ColorRangeFull,
	},
	canvas.PixelFormatGRAY16: {
		// Uploaded as GL_LUMINANCE_ALPHA (2 bytes per pixel, LE).
		category:     categoryGrayscale,
		planeCount:   1,
		planeFormats: [4]uint32{fdLumA},
		planeDataTypes: [4]uint32{fdByte},
		chromaWDiv:   [4]int{1},
		chromaHDiv:   [4]int{1},
		bitDepth:     16,
		hasAlpha:     false,
		defaultRange: canvas.ColorRangeFull,
	},
	canvas.PixelFormatYA8: {
		// Gray + alpha, interleaved. GL_LUMINANCE_ALPHA.
		category:     categoryGrayscale,
		planeCount:   1,
		planeFormats: [4]uint32{fdLumA},
		planeDataTypes: [4]uint32{fdByte},
		chromaWDiv:   [4]int{1},
		chromaHDiv:   [4]int{1},
		bitDepth:     8,
		hasAlpha:     true,
		defaultRange: canvas.ColorRangeFull,
	},

	// ── CPU-converted formats ─────────────────────────────────────────────────
	// These are converted to RGBA on the CPU before upload.
	canvas.PixelFormatRGB24:    {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatBGR24:    {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatBGRA:     {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatARGB:     {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatABGR:     {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatRGB565:   {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatBGR565:   {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatGBRP:     {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatGBRP10:   {category: categoryCPUConvert, bitDepth: 10},
	canvas.PixelFormatGBRP12:   {category: categoryCPUConvert, bitDepth: 12},
	canvas.PixelFormatGBRP16:   {category: categoryCPUConvert, bitDepth: 16},
	canvas.PixelFormatGBRAP:    {category: categoryCPUConvert, bitDepth: 8, hasAlpha: true},
	canvas.PixelFormatMONOWHITE:{category: categoryCPUConvert, bitDepth: 1},
	canvas.PixelFormatMONOBLACK:{category: categoryCPUConvert, bitDepth: 1},
	canvas.PixelFormatBAYER_BGGR8:  {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatBAYER_RGGB8:  {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatBAYER_GBRG8:  {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatBAYER_GRBG8:  {category: categoryCPUConvert, bitDepth: 8},
	canvas.PixelFormatBAYER_BGGR16: {category: categoryCPUConvert, bitDepth: 16},
	canvas.PixelFormatBAYER_RGGB16: {category: categoryCPUConvert, bitDepth: 16},
	canvas.PixelFormatBAYER_GBRG16: {category: categoryCPUConvert, bitDepth: 16},
	canvas.PixelFormatBAYER_GRBG16: {category: categoryCPUConvert, bitDepth: 16},
	canvas.PixelFormatGRAYF32: {category: categoryCPUConvert, bitDepth: 32},
	canvas.PixelFormatRGBF32:  {category: categoryCPUConvert, bitDepth: 32},
	canvas.PixelFormatRGBAF32: {category: categoryCPUConvert, bitDepth: 32, hasAlpha: true},
}

// formatDescriptorFor returns the descriptor for the given pixel format.
// Formats not in the table fall back to categoryCPUConvert.
func formatDescriptorFor(pixelFormat int) formatDescriptor {
	if d, ok := formatDescriptors[pixelFormat]; ok {
		return d
	}
	return formatDescriptor{category: categoryCPUConvert, bitDepth: 8}
}
