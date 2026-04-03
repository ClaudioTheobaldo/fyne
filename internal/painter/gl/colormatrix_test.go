package gl

import (
	"math"
	"testing"

	"fyne.io/fyne/v2/canvas"
)

func TestComputeColorMatrix_BT601_FullRange_8bit(t *testing.T) {
	// BT.601 full-range 8-bit must reproduce the coefficients previously
	// hardcoded in yuv420p.frag:
	//   R = Y + 1.402*(V-0.5)
	//   G = Y - 0.344136*(U-0.5) - 0.714136*(V-0.5)
	//   B = Y + 1.772*(U-0.5)
	m := ComputeColorMatrix(canvas.ColorSpaceBT601, canvas.ColorRangeFull, 8)

	// Row0: [1, 0, 1.402]
	assertNear(t, "Row0[0]", float64(m.Row0[0]), 1.0, 1e-4)
	assertNear(t, "Row0[1]", float64(m.Row0[1]), 0.0, 1e-4)
	assertNear(t, "Row0[2]", float64(m.Row0[2]), 1.402, 1e-3)

	// Row1: [1, -0.344136, -0.714136]
	assertNear(t, "Row1[0]", float64(m.Row1[0]), 1.0, 1e-4)
	assertNear(t, "Row1[1]", float64(m.Row1[1]), -0.344136, 1e-3)
	assertNear(t, "Row1[2]", float64(m.Row1[2]), -0.714136, 1e-3)

	// Row2: [1, 1.772, 0]
	assertNear(t, "Row2[0]", float64(m.Row2[0]), 1.0, 1e-4)
	assertNear(t, "Row2[1]", float64(m.Row2[1]), 1.772, 1e-3)
	assertNear(t, "Row2[2]", float64(m.Row2[2]), 0.0, 1e-4)

	// Offset: [0, 0.5, 0.5]
	assertNear(t, "Offset[0]", float64(m.Offset[0]), 0.0, 1e-4)
	assertNear(t, "Offset[1]", float64(m.Offset[1]), 0.5, 1e-4)
	assertNear(t, "Offset[2]", float64(m.Offset[2]), 0.5, 1e-4)
}

func TestComputeColorMatrix_BT709_FullRange_8bit(t *testing.T) {
	m := ComputeColorMatrix(canvas.ColorSpaceBT709, canvas.ColorRangeFull, 8)

	// BT.709 full-range: R = Y + 1.5748*(V-0.5)
	assertNear(t, "Row0[2] (Rv)", float64(m.Row0[2]), 2.0*(1.0-0.2126), 1e-3)
	assertNear(t, "Row2[1] (Bu)", float64(m.Row2[1]), 2.0*(1.0-0.0722), 1e-3)
	assertNear(t, "Offset[1]", float64(m.Offset[1]), 0.5, 1e-4)
}

func TestComputeColorMatrix_BT601_LimitedRange_8bit(t *testing.T) {
	m := ComputeColorMatrix(canvas.ColorSpaceBT601, canvas.ColorRangeLimited, 8)

	// Y scale = 255/(235-16) = 255/219 ≈ 1.1644
	yScale := 255.0 / 219.0
	assertNear(t, "Row0[0] (yScale)", float64(m.Row0[0]), yScale, 1e-3)

	// Y offset = 16/255 ≈ 0.0627
	assertNear(t, "Offset[0] (yOff)", float64(m.Offset[0]), 16.0/255.0, 1e-4)

	// UV offset = 128/255 ≈ 0.5020
	assertNear(t, "Offset[1] (uvOff)", float64(m.Offset[1]), 128.0/255.0, 1e-4)
}

func TestComputeColorMatrix_BT2020_FullRange_10bit(t *testing.T) {
	m := ComputeColorMatrix(canvas.ColorSpaceBT2020, canvas.ColorRangeFull, 10)

	// BT.2020 Kr=0.2627, Kb=0.0593
	rv := 2.0 * (1.0 - 0.2627)
	bu := 2.0 * (1.0 - 0.0593)
	assertNear(t, "Row0[2] (Rv)", float64(m.Row0[2]), rv, 1e-3)
	assertNear(t, "Row2[1] (Bu)", float64(m.Row2[1]), bu, 1e-3)
	// Bit depth doesn't affect full-range coefficients, only offsets
	assertNear(t, "Row0[0] (yScale)", float64(m.Row0[0]), 1.0, 1e-4)
}

func TestComputeColorMatrix_YCgCo(t *testing.T) {
	m := ComputeColorMatrix(canvas.ColorSpaceYCgCo, canvas.ColorRangeFull, 8)

	// YCgCo: R=Y-Cg+Co, G=Y+Cg, B=Y-Cg-Co
	// Row0: [1, -1, 1]
	assertNear(t, "Row0[0]", float64(m.Row0[0]), 1.0, 1e-4)
	assertNear(t, "Row0[1]", float64(m.Row0[1]), -1.0, 1e-4)
	assertNear(t, "Row0[2]", float64(m.Row0[2]), 1.0, 1e-4)

	// Row1: [1, 1, 0]
	assertNear(t, "Row1[0]", float64(m.Row1[0]), 1.0, 1e-4)
	assertNear(t, "Row1[1]", float64(m.Row1[1]), 1.0, 1e-4)
	assertNear(t, "Row1[2]", float64(m.Row1[2]), 0.0, 1e-4)
}

func assertNear(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s: got %f, want %f (tolerance %f)", name, got, want, tol)
	}
}
