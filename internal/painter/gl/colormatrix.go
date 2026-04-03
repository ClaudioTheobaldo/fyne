package gl

import "fyne.io/fyne/v2/canvas"

// ColorMatrix holds a 3x3 YUV-to-RGB conversion matrix and a pre-subtraction
// offset vector. The shader applies:
//
//	yuv_shifted = vec3(y, u, v) - Offset
//	rgb = vec3(dot(Row0, yuv_shifted), dot(Row1, yuv_shifted), dot(Row2, yuv_shifted))
type ColorMatrix struct {
	Row0, Row1, Row2 [3]float32
	Offset           [3]float32
}

// ComputeColorMatrix derives the YUV-to-RGB conversion matrix for the given
// color space, range, and bit depth. For 8-bit content pass bitDepth=8.
func ComputeColorMatrix(colorSpace, colorRange, bitDepth int) ColorMatrix {
	switch colorSpace {
	case canvas.ColorSpaceBT709:
		return buildYUVMatrix(0.2126, 0.0722, colorRange, bitDepth)
	case canvas.ColorSpaceBT2020:
		return buildYUVMatrix(0.2627, 0.0593, colorRange, bitDepth)
	case canvas.ColorSpaceSMPTE240M:
		return buildYUVMatrix(0.212, 0.087, colorRange, bitDepth)
	case canvas.ColorSpaceFCC:
		return buildYUVMatrix(0.30, 0.11, colorRange, bitDepth)
	case canvas.ColorSpaceYCgCo:
		return ycocgMatrix(colorRange, bitDepth)
	case canvas.ColorSpaceICtCp:
		return ictcpMatrix(colorRange, bitDepth)
	default: // BT.601 (ColorSpaceBT601 = 0)
		return buildYUVMatrix(0.299, 0.114, colorRange, bitDepth)
	}
}

// buildYUVMatrix constructs a YUV-to-RGB matrix from luma coefficients Kr and
// Kb (standard ITU derivation), adjusted for the given range and bit depth.
func buildYUVMatrix(kr, kb float64, colorRange, bitDepth int) ColorMatrix {
	kg := 1.0 - kr - kb
	bits := 1 << uint(bitDepth)
	maxVal := float64(bits - 1)

	var yOff, uvOff, yScale, uvScale float64
	if colorRange == canvas.ColorRangeLimited {
		// Limited range: Y in [16,235] for 8-bit, UV in [16,240] for 8-bit.
		// Scale the reference values proportionally for higher bit depths.
		yMin := 16.0 * maxVal / 255.0
		yMax := 235.0 * maxVal / 255.0
		uvMin := 16.0 * maxVal / 255.0
		uvMax := 240.0 * maxVal / 255.0
		yOff = yMin / maxVal
		uvOff = (uvMin + uvMax) / 2.0 / maxVal // = 128/255 ≈ 0.50196 for 8-bit
		yScale = maxVal / (yMax - yMin)
		uvScale = maxVal / (uvMax - uvMin)
	} else {
		// Full range: Y, U, V all in [0, maxVal].
		yOff = 0
		uvOff = 0.5
		yScale = 1.0
		uvScale = 1.0
	}

	// Standard ITU YCbCr-to-RGB derivation:
	//   R = Y + 2(1-Kr)*V
	//   G = Y - 2*Kr*(1-Kr)/Kg*V - 2*Kb*(1-Kb)/Kg*U
	//   B = Y + 2(1-Kb)*U
	rv := 2.0 * (1.0 - kr)
	gv := -2.0 * kr * (1.0 - kr) / kg
	gu := -2.0 * kb * (1.0 - kb) / kg
	bu := 2.0 * (1.0 - kb)

	return ColorMatrix{
		// Each row: [Y_coeff, U_coeff, V_coeff]
		Row0: [3]float32{float32(yScale), 0, float32(rv * uvScale)},
		Row1: [3]float32{float32(yScale), float32(gu * uvScale), float32(gv * uvScale)},
		Row2: [3]float32{float32(yScale), float32(bu * uvScale), 0},
		Offset: [3]float32{float32(yOff), float32(uvOff), float32(uvOff)},
	}
}

// ycocgMatrix returns the YCgCo-to-RGB matrix.
// YCgCo: R = Y - Cg + Co, G = Y + Cg, B = Y - Cg - Co.
func ycocgMatrix(colorRange, bitDepth int) ColorMatrix {
	_ = colorRange
	_ = bitDepth
	return ColorMatrix{
		Row0:   [3]float32{1, -1, 1},
		Row1:   [3]float32{1, 1, 0},
		Row2:   [3]float32{1, -1, -1},
		Offset: [3]float32{0, 0.5, 0.5},
	}
}

// ictcpMatrix returns an approximate ICtCp-to-RGB linear matrix
// (ITU-R BT.2100). This is a display-referred approximation; full PQ/HLG
// inverse EOTF is deferred to a future transfer-function pass.
func ictcpMatrix(colorRange, bitDepth int) ColorMatrix {
	_ = colorRange
	_ = bitDepth
	return ColorMatrix{
		Row0:   [3]float32{1, 0.00860514, 0.11102860},
		Row1:   [3]float32{1, -0.00860514, -0.11102860},
		Row2:   [3]float32{1, 0.56003134, -0.32062717},
		Offset: [3]float32{0, 0.5, 0.5},
	}
}
