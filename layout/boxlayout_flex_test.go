package layout_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

// --- VBox Flex Grow ---

func TestVBoxFlex_GrowEqual(t *testing.T) {
	// Two children with equal grow should split free space equally.
	a := NewMinSizeRect(fyne.NewSize(50, 20))
	b := NewMinSizeRect(fyne.NewSize(50, 20))
	padding := theme.Padding()

	l := layout.NewFlexVBoxLayout()
	containerSize := fyne.NewSize(100, 100)
	l.Layout([]fyne.CanvasObject{
		layout.WithGrow(a, 1),
		layout.WithGrow(b, 1),
	}, containerSize)

	// totalBasis=40, gaps=padding, freeSpace=100-40-padding
	freeSpace := float32(100) - 40 - padding
	expected := 20 + freeSpace/2

	assert.InDelta(t, expected, a.Size().Height, 0.5)
	assert.InDelta(t, expected, b.Size().Height, 0.5)
}

func TestVBoxFlex_GrowUnequal(t *testing.T) {
	// Grow 1:2 ratio.
	a := NewMinSizeRect(fyne.NewSize(50, 10))
	b := NewMinSizeRect(fyne.NewSize(50, 10))
	padding := theme.Padding()

	l := layout.NewFlexVBoxLayout()
	l.Layout([]fyne.CanvasObject{
		layout.WithGrow(a, 1),
		layout.WithGrow(b, 2),
	}, fyne.NewSize(100, 100))

	freeSpace := float32(100) - 20 - padding
	assert.InDelta(t, 10+freeSpace*1/3, a.Size().Height, 0.5)
	assert.InDelta(t, 10+freeSpace*2/3, b.Size().Height, 0.5)
}

func TestVBoxFlex_GrowSingleChild(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 20))

	l := layout.NewFlexVBoxLayout()
	l.Layout([]fyne.CanvasObject{
		layout.WithGrow(a, 1),
	}, fyne.NewSize(100, 100))

	// Single child absorbs all free space.
	assert.InDelta(t, float32(100), a.Size().Height, 0.5)
}

func TestVBoxFlex_NoGrowNoChange(t *testing.T) {
	// Without grow, children keep their MinSize height.
	a := NewMinSizeRect(fyne.NewSize(50, 30))
	b := NewMinSizeRect(fyne.NewSize(50, 40))

	l := layout.NewFlexVBoxLayout()
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 200))

	assert.Equal(t, float32(30), a.Size().Height)
	assert.Equal(t, float32(40), b.Size().Height)
}

// --- VBox Flex Shrink ---

func TestVBoxFlex_ShrinkEqual(t *testing.T) {
	// Container is smaller than sum of bases; shrink equally.
	a := NewMinSizeRect(fyne.NewSize(50, 10))
	b := NewMinSizeRect(fyne.NewSize(50, 10))
	padding := theme.Padding()

	l := layout.NewFlexVBoxLayout()
	l.Layout([]fyne.CanvasObject{
		layout.WithFlex(a, 0, 1, 60),
		layout.WithFlex(b, 0, 1, 60),
	}, fyne.NewSize(100, 50))

	// totalBasis=120, gap=padding, freeSpace=50-120-padding (negative)
	freeSpace := float32(50) - 120 - padding // very negative
	_ = freeSpace
	// Each shrinks by half the excess. But clamp to MinSize(10).
	// shrinkAmt per child = abs(freeSpace) * 1 * 60 / (1*60 + 1*60) = abs(freeSpace)/2
	// result = 60 - abs(freeSpace)/2 -- may be below MinSize
	// Both should be clamped to MinSize=10
	assert.True(t, a.Size().Height >= 10, "a should not shrink below MinSize")
	assert.True(t, b.Size().Height >= 10, "b should not shrink below MinSize")
}

func TestVBoxFlex_ShrinkClampToMinSize(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 30))

	l := layout.NewFlexVBoxLayout()
	l.Layout([]fyne.CanvasObject{
		layout.WithFlex(a, 0, 1, 100),
	}, fyne.NewSize(100, 20))

	// Cannot shrink below MinSize of 30.
	assert.True(t, a.Size().Height >= 30)
}

// --- VBox Justify ---

func TestVBoxFlex_JustifyEnd(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 30))

	l := layout.NewFlexVBoxLayout(layout.WithJustify(layout.JustifyEnd))
	l.Layout([]fyne.CanvasObject{a}, fyne.NewSize(100, 100))

	// Child should be at bottom: y = 100-30 = 70
	assert.InDelta(t, float32(70), a.Position().Y, 0.5)
}

func TestVBoxFlex_JustifyCenter(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 30))

	l := layout.NewFlexVBoxLayout(layout.WithJustify(layout.JustifyCenter))
	l.Layout([]fyne.CanvasObject{a}, fyne.NewSize(100, 100))

	// Child centered: y = (100-30)/2 = 35
	assert.InDelta(t, float32(35), a.Position().Y, 0.5)
}

func TestVBoxFlex_JustifySpaceBetween(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 20))
	b := NewMinSizeRect(fyne.NewSize(50, 20))
	c := NewMinSizeRect(fyne.NewSize(50, 20))

	l := layout.NewFlexVBoxLayout(layout.WithJustify(layout.JustifySpaceBetween))
	l.Layout([]fyne.CanvasObject{a, b, c}, fyne.NewSize(100, 100))

	// a at 0, c at 80 (100-20), b in middle at 40
	assert.InDelta(t, float32(0), a.Position().Y, 0.5)
	assert.InDelta(t, float32(80), c.Position().Y, 0.5)
}

func TestVBoxFlex_JustifySpaceEvenly(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 20))
	b := NewMinSizeRect(fyne.NewSize(50, 20))

	l := layout.NewFlexVBoxLayout(layout.WithJustify(layout.JustifySpaceEvenly))
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 100))

	// gap = (100-40) / (2+1) = 20. a at 20, b at 20+20+20 = 60.
	// (distribute modes own the full inter-item spacing, no extra theme padding)
	gap := (float32(100) - 40) / 3
	assert.InDelta(t, gap, a.Position().Y, 0.5)
	assert.InDelta(t, gap+20+gap, b.Position().Y, 0.5)
}

// --- VBox Cross-axis Alignment ---

func TestVBoxFlex_AlignItemsCenter(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 20))

	l := layout.NewFlexVBoxLayout(layout.WithAlignItems(layout.AlignCenter))
	l.Layout([]fyne.CanvasObject{a}, fyne.NewSize(100, 100))

	// Child width = MinSize.Width=30, centered x = (100-30)/2 = 35
	assert.InDelta(t, float32(35), a.Position().X, 0.5)
	assert.Equal(t, float32(30), a.Size().Width)
}

func TestVBoxFlex_AlignItemsEnd(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 20))

	l := layout.NewFlexVBoxLayout(layout.WithAlignItems(layout.AlignEnd))
	l.Layout([]fyne.CanvasObject{a}, fyne.NewSize(100, 100))

	// x = 100-30 = 70
	assert.InDelta(t, float32(70), a.Position().X, 0.5)
}

func TestVBoxFlex_AlignItemsStart(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 20))

	l := layout.NewFlexVBoxLayout(layout.WithAlignItems(layout.AlignStart))
	l.Layout([]fyne.CanvasObject{a}, fyne.NewSize(100, 100))

	assert.Equal(t, float32(0), a.Position().X)
	assert.Equal(t, float32(30), a.Size().Width)
}

func TestVBoxFlex_AlignSelfOverridesContainer(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 20))
	b := NewMinSizeRect(fyne.NewSize(30, 20))
	wB := layout.WithAlign(b, layout.AlignEnd) // keep wrapper reference

	l := layout.NewFlexVBoxLayout(layout.WithAlignItems(layout.AlignCenter))
	l.Layout([]fyne.CanvasObject{
		a,
		wB,
	}, fyne.NewSize(100, 100))

	// a is a plain rect (not wrapped) — centered at x = (100-30)/2 = 35
	assert.InDelta(t, float32(35), a.Position().X, 0.5)
	// wB is the HintedObject wrapper — its widget position should be at end: x = 100-30 = 70
	assert.InDelta(t, float32(70), wB.Position().X, 0.5)
}

// --- VBox Reverse ---

func TestVBoxFlex_Reverse(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 20))
	b := NewMinSizeRect(fyne.NewSize(50, 30))

	l := layout.NewFlexVBoxLayout(layout.WithReverse(true))
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 100))

	// Reversed: b is placed first (at top), a second.
	assert.True(t, b.Position().Y < a.Position().Y, "b should be above a when reversed")
}

// --- VBox Wrap ---

func TestVBoxFlex_Wrap(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 40))
	b := NewMinSizeRect(fyne.NewSize(30, 40))
	c := NewMinSizeRect(fyne.NewSize(30, 40))

	// Container height = 50; a and b together need 40+padding+40 > 50, so b wraps to new column.
	l := layout.NewFlexVBoxLayout(layout.WithWrap(true))
	l.Layout([]fyne.CanvasObject{a, b, c}, fyne.NewSize(200, 50))

	// a is in column 0, b and c are in column 1 (and possibly 2).
	assert.Equal(t, float32(0), a.Position().X, "a should be in column 0")
	assert.True(t, b.Position().X > 0, "b should wrap to next column")
}

// --- HBox Flex Grow ---

func TestHBoxFlex_GrowEqual(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(20, 50))
	b := NewMinSizeRect(fyne.NewSize(20, 50))
	padding := theme.Padding()

	l := layout.NewFlexHBoxLayout()
	l.Layout([]fyne.CanvasObject{
		layout.WithGrow(a, 1),
		layout.WithGrow(b, 1),
	}, fyne.NewSize(100, 100))

	freeSpace := float32(100) - 40 - padding
	expected := 20 + freeSpace/2
	assert.InDelta(t, expected, a.Size().Width, 0.5)
	assert.InDelta(t, expected, b.Size().Width, 0.5)
}

func TestHBoxFlex_JustifyCenter(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 50))

	l := layout.NewFlexHBoxLayout(layout.WithJustify(layout.JustifyCenter))
	l.Layout([]fyne.CanvasObject{a}, fyne.NewSize(100, 100))

	assert.InDelta(t, float32(35), a.Position().X, 0.5)
}

func TestHBoxFlex_AlignItemsCenter(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 20))

	l := layout.NewFlexHBoxLayout(layout.WithAlignItems(layout.AlignCenter))
	l.Layout([]fyne.CanvasObject{a}, fyne.NewSize(100, 100))

	// y = (100-20)/2 = 40
	assert.InDelta(t, float32(40), a.Position().Y, 0.5)
	assert.Equal(t, float32(20), a.Size().Height)
}

func TestHBoxFlex_Wrap(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(40, 30))
	b := NewMinSizeRect(fyne.NewSize(40, 30))
	c := NewMinSizeRect(fyne.NewSize(40, 30))

	l := layout.NewFlexHBoxLayout(layout.WithWrap(true))
	l.Layout([]fyne.CanvasObject{a, b, c}, fyne.NewSize(50, 200))

	// a in row 0, b wraps to row 1
	assert.Equal(t, float32(0), a.Position().Y, "a should be in row 0")
	assert.True(t, b.Position().Y > 0, "b should wrap to next row")
}

// --- Backward compatibility ---

func TestVBoxFlex_BackwardCompatSpacers(t *testing.T) {
	// When no flex hints are used, spacers should still work as before.
	a := NewMinSizeRect(fyne.NewSize(50, 30))
	spc := layout.NewSpacer()
	b := NewMinSizeRect(fyne.NewSize(50, 30))

	l := layout.NewFlexVBoxLayout()
	l.Layout([]fyne.CanvasObject{a, spc, b}, fyne.NewSize(100, 100))

	// a at top (y=0), spacer fills middle, b at bottom.
	// Legacy spacer: extra = 100-60-4 = 36, spacerSize=36, b at y=34+36=70.
	assert.Equal(t, float32(0), a.Position().Y)
	assert.InDelta(t, float32(70), b.Position().Y, 0.5)
}

func TestVBoxFlex_GrowTakesPrecedenceOverSpacer(t *testing.T) {
	// When a Grow hint is present, spacers become zero-size.
	a := NewMinSizeRect(fyne.NewSize(50, 20))
	spc := layout.NewSpacer()
	b := NewMinSizeRect(fyne.NewSize(50, 20))

	l := layout.NewFlexVBoxLayout()
	l.Layout([]fyne.CanvasObject{
		layout.WithGrow(a, 1),
		spc,
		layout.WithGrow(b, 1),
	}, fyne.NewSize(100, 100))

	// Spacer has zero size; grow distributes space.
	assert.Equal(t, float32(0), spc.Size().Height)
	assert.True(t, a.Size().Height > 20, "a should grow")
}

// --- Margin via WithMargin ---

func TestVBoxFlex_MarginOffsetPosition(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(50, 20))

	l := layout.NewFlexVBoxLayout()
	wrapped := layout.WithMargin(a, layout.Insets{Top: 10, Left: 5})
	l.Layout([]fyne.CanvasObject{wrapped}, fyne.NewSize(100, 100))

	// Inner object should be offset by margin.
	assert.InDelta(t, float32(5), a.Position().X, 0.5)
	assert.InDelta(t, float32(10), a.Position().Y, 0.5)
}
