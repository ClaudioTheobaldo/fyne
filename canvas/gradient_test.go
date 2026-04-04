package canvas_test

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"fyne.io/fyne/v2/canvas"
	internalTest "fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestNewHorizontalGradient(t *testing.T) {
	horizontal := canvas.NewHorizontalGradient(color.Black, color.Transparent)

	smallImg := horizontal.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0xe6, 0xb3, 0x7f, 0x4c, 0x19},
		{0xe6, 0xb3, 0x7f, 0x4c, 0x19},
		{0xe6, 0xb3, 0x7f, 0x4c, 0x19},
		{0xe6, 0xb3, 0x7f, 0x4c, 0x19},
		{0xe6, 0xb3, 0x7f, 0x4c, 0x19},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	img := horizontal.Generate(51, 5)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, img.At(0, 0))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(25, i))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, img.At(50, 0))
}

func TestNewHorizontalGradient_Flipped(t *testing.T) {
	horizontal := canvas.NewHorizontalGradient(color.Black, color.Transparent)
	horizontal.Angle -= 180

	smallImg := horizontal.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0x19, 0x4c, 0x7f, 0xb3, 0xe6},
		{0x19, 0x4c, 0x7f, 0xb3, 0xe6},
		{0x19, 0x4c, 0x7f, 0xb3, 0xe6},
		{0x19, 0x4c, 0x7f, 0xb3, 0xe6},
		{0x19, 0x4c, 0x7f, 0xb3, 0xe6},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	img := horizontal.Generate(51, 5)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, img.At(0, 0))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(25, i))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, img.At(50, 0))
}

func TestNewVerticalGradient(t *testing.T) {
	vertical := canvas.NewVerticalGradient(color.Black, color.Transparent)

	smallImg := vertical.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0xe6, 0xe6, 0xe6, 0xe6, 0xe6},
		{0xb3, 0xb3, 0xb3, 0xb3, 0xb3},
		{0x7f, 0x7f, 0x7f, 0x7f, 0x7f},
		{0x4c, 0x4c, 0x4c, 0x4c, 0x4c},
		{0x19, 0x19, 0x19, 0x19, 0x19},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	imgVert := vertical.Generate(5, 51)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, imgVert.At(0, 0))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, imgVert.At(i, 25))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, imgVert.At(0, 50))
}

func TestNewVerticalGradient_Flipped(t *testing.T) {
	vertical := canvas.NewVerticalGradient(color.Black, color.Transparent)
	vertical.Angle += 180

	smallImg := vertical.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0x19, 0x19, 0x19, 0x19, 0x19},
		{0x4c, 0x4c, 0x4c, 0x4c, 0x4c},
		{0x7f, 0x7f, 0x7f, 0x7f, 0x7f},
		{0xb3, 0xb3, 0xb3, 0xb3, 0xb3},
		{0xe6, 0xe6, 0xe6, 0xe6, 0xe6},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	imgVert := vertical.Generate(5, 51)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, imgVert.At(0, 0))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, imgVert.At(i, 25))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, imgVert.At(0, 50))
}

func TestNewLinearGradient_45(t *testing.T) {
	negative := canvas.NewLinearGradient(color.Black, color.Transparent, 45.0)
	negativeBackward := canvas.NewLinearGradient(color.Black, color.Transparent, -315.0)

	smallImg := negative.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0x7f, 0x99, 0xb3, 0xcc, 0xe6},
		{0x66, 0x7f, 0x99, 0xb3, 0xcc},
		{0x4c, 0x66, 0x7f, 0x99, 0xb3},
		{0x33, 0x4c, 0x66, 0x7f, 0x99},
		{0x19, 0x33, 0x4c, 0x66, 0x7f},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	smallBackwardImg := negativeBackward.Generate(5, 5)
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallBackwardImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	img := negative.Generate(51, 51)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, img.At(50, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, img.At(0, 50))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(i, i))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(0, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(50, 50))
}

func TestNewLinearGradient_225(t *testing.T) {
	negative := canvas.NewLinearGradient(color.Black, color.Transparent, 225.0)

	smallImg := negative.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0x7f, 0x66, 0x4c, 0x33, 0x19},
		{0x99, 0x7f, 0x66, 0x4c, 0x33},
		{0xb3, 0x99, 0x7f, 0x66, 0x4c},
		{0xcc, 0xb3, 0x99, 0x7f, 0x66},
		{0xe6, 0xcc, 0xb3, 0x99, 0x7f},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	img := negative.Generate(51, 51)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, img.At(50, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, img.At(0, 50))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(i, i))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(0, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(50, 50))
}

func TestNewLinearGradient_135(t *testing.T) {
	positive := canvas.NewLinearGradient(color.Black, color.Transparent, 135.0)

	smallImg := positive.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0x19, 0x33, 0x4c, 0x66, 0x7f},
		{0x33, 0x4c, 0x66, 0x7f, 0x99},
		{0x4c, 0x66, 0x7f, 0x99, 0xb3},
		{0x66, 0x7f, 0x99, 0xb3, 0xcc},
		{0x7f, 0x99, 0xb3, 0xcc, 0xe6},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	img := positive.Generate(51, 51)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, img.At(0, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, img.At(50, 50))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(50-i, i))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(50, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(0, 50))
}

func TestNewLinearGradient_315(t *testing.T) {
	positive := canvas.NewLinearGradient(color.Black, color.Transparent, 315.0)

	smallImg := positive.Generate(5, 5)
	expectedAlphaValues := [][]uint8{
		{0xe6, 0xcc, 0xb3, 0x99, 0x7f},
		{0xcc, 0xb3, 0x99, 0x7f, 0x66},
		{0xb3, 0x99, 0x7f, 0x66, 0x4c},
		{0x99, 0x7f, 0x66, 0x4c, 0x33},
		{0x7f, 0x66, 0x4c, 0x33, 0x19},
	}
	for y, xv := range expectedAlphaValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{0, 0, 0, v}, smallImg.At(x, y), "alpha value at %d,%d", x, y)
		}
	}

	img := positive.Generate(51, 51)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xfd}, img.At(0, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x02}, img.At(50, 50))
	for i := 0; i < 5; i++ {
		assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(50-i, i))
	}
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(50, 0))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7f}, img.At(0, 50))
}

func TestNewRadialGradient(t *testing.T) {
	circle := canvas.NewRadialGradient(color.Black, color.Transparent)

	{
		imgOddDiameter := circle.Generate(5, 5)
		expectedAlphaValues := [][]uint8{
			{0x00, 0x1b, 0x33, 0x1b, 0x00},
			{0x1b, 0x6f, 0x99, 0x6f, 0x1b},
			{0x33, 0x99, 0xff, 0x99, 0x33},
			{0x1b, 0x6f, 0x99, 0x6f, 0x1b},
			{0x00, 0x1b, 0x33, 0x1b, 0x00},
		}
		for y, xv := range expectedAlphaValues {
			for x, v := range xv {
				assert.Equal(t, color.NRGBA{0, 0, 0, v}, imgOddDiameter.At(x, y), "alpha value at %d,%d", x, y)
			}
		}
	}

	{
		imgEvenDiameter := circle.Generate(6, 6)
		expectedAlphaValues := [][]uint8{
			{0x00, 0x07, 0x26, 0x26, 0x07, 0x00},
			{0x07, 0x4a, 0x79, 0x79, 0x4a, 0x07},
			{0x26, 0x79, 0xc3, 0xc3, 0x79, 0x26},
			{0x26, 0x79, 0xc3, 0xc3, 0x79, 0x26},
			{0x07, 0x4a, 0x79, 0x79, 0x4a, 0x07},
			{0x00, 0x07, 0x26, 0x26, 0x07, 0x00},
		}
		for y, xv := range expectedAlphaValues {
			for x, v := range xv {
				assert.Equal(t, color.NRGBA{0, 0, 0, v}, imgEvenDiameter.At(x, y), "alpha value at %d,%d", x, y)
			}
		}
	}

	imgCircle := circle.Generate(10, 10)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x18}, imgCircle.At(9, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x4a}, imgCircle.At(8, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7d}, imgCircle.At(7, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xaf}, imgCircle.At(6, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xdb}, imgCircle.At(5, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xdb}, imgCircle.At(4, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xaf}, imgCircle.At(3, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x7d}, imgCircle.At(2, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x4a}, imgCircle.At(1, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x18}, imgCircle.At(0, 5))

	circle.CenterOffsetX = 0.1
	circle.CenterOffsetY = 0.1
	imgCircleOffset := circle.Generate(10, 10)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xe1}, imgCircleOffset.At(5, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xbc}, imgCircleOffset.At(4, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x93}, imgCircleOffset.At(3, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x69}, imgCircleOffset.At(2, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x3e}, imgCircleOffset.At(1, 5))

	circle.CenterOffsetX = -0.1
	circle.CenterOffsetY = -0.1
	imgCircleOffset = circle.Generate(10, 10)
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xa5}, imgCircleOffset.At(5, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xbc}, imgCircleOffset.At(4, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xbc}, imgCircleOffset.At(3, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0xa5}, imgCircleOffset.At(2, 5))
	assert.Equal(t, color.NRGBA{0, 0, 0, 0x83}, imgCircleOffset.At(1, 5))
}

func TestRadialGradient_Elliptical(t *testing.T) {
	// Elliptical with ScaleX=2.0 should stretch horizontally
	g := canvas.NewRadialGradient(color.Black, color.Transparent)
	g.ScaleX = 2.0
	g.ScaleY = 1.0
	img := g.Generate(10, 10)

	// Center pixel should be opaque
	centerColor := img.At(5, 5).(color.NRGBA)
	assert.Greater(t, centerColor.A, uint8(0x80), "center should be dark")

	// Compare horizontal vs vertical falloff: horizontal should be slower due to ScaleX=2
	rightColor := img.At(8, 5).(color.NRGBA)
	bottomColor := img.At(5, 8).(color.NRGBA)
	// With ScaleX=2.0, the gradient extends further horizontally,
	// so at the same pixel distance, horizontal should have higher alpha (darker)
	assert.Greater(t, rightColor.A, bottomColor.A,
		"horizontal falloff should be slower than vertical with ScaleX=2")
}

func TestRadialGradient_DefaultScaleBackwardCompat(t *testing.T) {
	// ScaleX=0, ScaleY=0 should behave identically to the old circular gradient
	g := canvas.NewRadialGradient(color.Black, color.Transparent)
	imgDefault := g.Generate(10, 10)

	g2 := canvas.NewRadialGradient(color.Black, color.Transparent)
	g2.ScaleX = 1.0
	g2.ScaleY = 1.0
	imgExplicit := g2.Generate(10, 10)

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			assert.Equal(t, imgDefault.At(x, y), imgExplicit.At(x, y),
				"default scale should match explicit 1.0 at %d,%d", x, y)
		}
	}
}

func TestLinearGradient_MultiStop(t *testing.T) {
	g := &canvas.LinearGradient{
		Angle: 0, // vertical
		Stops: []canvas.GradientStop{
			{Color: color.NRGBA{R: 255, A: 255}, Position: 0.0},  // red at top
			{Color: color.NRGBA{G: 255, A: 255}, Position: 0.5},  // green at middle
			{Color: color.NRGBA{B: 255, A: 255}, Position: 1.0},  // blue at bottom
		},
	}
	img := g.Generate(1, 100)

	// Top should be red
	topColor := img.At(0, 0).(color.NRGBA)
	assert.Greater(t, topColor.R, uint8(250), "top red channel should be near 255")
	assert.Less(t, topColor.G, uint8(5), "top green channel should be near 0")

	// Middle should be green
	midColor := img.At(0, 50).(color.NRGBA)
	assert.Greater(t, midColor.G, uint8(200), "middle should be mostly green")

	// Bottom should be blue
	bottomColor := img.At(0, 99).(color.NRGBA)
	assert.Greater(t, bottomColor.B, uint8(200), "bottom should be mostly blue")
}

func TestRadialGradient_MultiStop(t *testing.T) {
	g := canvas.NewRadialGradient(color.Transparent, color.Transparent)
	g.Stops = []canvas.GradientStop{
		{Color: color.NRGBA{R: 255, A: 255}, Position: 0.0},
		{Color: color.Transparent, Position: 0.5},
	}
	img := g.Generate(10, 10)

	// Center should be red
	centerColor := img.At(5, 5).(color.NRGBA)
	assert.Greater(t, centerColor.R, uint8(100), "center should be reddish")
	assert.Greater(t, centerColor.A, uint8(100), "center should be opaque-ish")

	// Edge should be transparent (past the 0.5 stop)
	edgeColor := img.At(0, 0).(color.NRGBA)
	assert.Less(t, edgeColor.A, uint8(50), "edge should be mostly transparent")
}

func TestGradient_colorComputation(t *testing.T) {
	bg := internalTest.NewCheckedImage(50, 50, 1, 2)
	bounds := image.Rect(0, 0, 49, 49)
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, bg, image.Pt(0, 0), draw.Src)
	g := canvas.NewHorizontalGradient(color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, color.Transparent)
	ovl := g.Generate(50, 50)
	draw.Draw(dst, bounds, ovl, image.Pt(0, 0), draw.Over)
	test.AssertImageMatches(t, "gradient_colors.png", dst)
}
