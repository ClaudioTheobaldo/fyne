package layout_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

// --- WeightedGridLayout ---

func TestWeightedGrid_EqualWeights(t *testing.T) {
	// 2 equal-weight columns should behave like a 2-col uniform grid.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewWeightedGridLayout(1, 1)
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 20))

	colW := (100 - padding) / 2
	assert.InDelta(t, float32(0), a.Position().X, 0.5)
	assert.InDelta(t, colW+padding, b.Position().X, 0.5)
	assert.InDelta(t, colW, a.Size().Width, 0.5)
	assert.InDelta(t, colW, b.Size().Width, 0.5)
}

func TestWeightedGrid_ProportionalWeights(t *testing.T) {
	// 1:2:1 columns.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	c := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewWeightedGridLayout(1, 2, 1)
	l.Layout([]fyne.CanvasObject{a, b, c}, fyne.NewSize(100, 20))

	available := float32(100) - padding*2 // 3 cols, 2 gaps
	col0W := available * 1 / 4
	col1W := available * 2 / 4
	col2W := available * 1 / 4

	assert.InDelta(t, float32(0), a.Position().X, 0.5)
	assert.InDelta(t, col0W+padding, b.Position().X, 0.5)
	assert.InDelta(t, col0W+padding+col1W+padding, c.Position().X, 0.5)

	assert.InDelta(t, col0W, a.Size().Width, 0.5)
	assert.InDelta(t, col1W, b.Size().Width, 0.5)
	assert.InDelta(t, col2W, c.Size().Width, 0.5)
}

func TestWeightedGrid_WrapToNextRow(t *testing.T) {
	// 2 columns; 3 items → second row.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	c := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewWeightedGridLayout(1, 1)
	l.Layout([]fyne.CanvasObject{a, b, c}, fyne.NewSize(100, 50))

	assert.Equal(t, float32(0), a.Position().Y, "a in row 0")
	assert.Equal(t, float32(0), b.Position().Y, "b in row 0")
	assert.InDelta(t, float32(20)+padding, c.Position().Y, 0.5, "c in row 1")
}

func TestWeightedGrid_ColSpan(t *testing.T) {
	// 3 columns, first item spans 2 columns.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	c := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewWeightedGridLayout(1, 1, 1)
	l.Layout([]fyne.CanvasObject{
		layout.WithSpan(a, 2, 1),
		b,
		c,
	}, fyne.NewSize(100, 20))

	available := float32(100) - padding*2
	colW := available / 3
	spanW := colW*2 + padding

	assert.InDelta(t, float32(0), a.Position().X, 0.5)
	assert.InDelta(t, spanW, a.Size().Width, 0.5, "a spans 2 cols")

	// b is in col 2 (after the 2-col span).
	assert.InDelta(t, spanW+padding, b.Position().X, 0.5)
}

func TestWeightedGrid_MinSize(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(30, 20))
	b := NewMinSizeRect(fyne.NewSize(30, 20))

	l := layout.NewWeightedGridLayout(1, 1)
	min := l.MinSize([]fyne.CanvasObject{a, b})

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

// --- FlexGridLayout ---

func TestFlexGrid_EqualFrColumns(t *testing.T) {
	// 2 equal fr columns.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewFlexGridLayout([]layout.TrackSize{layout.Fr(1), layout.Fr(1)})
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 20))

	colW := (100 - padding) / 2
	assert.InDelta(t, float32(0), a.Position().X, 0.5)
	assert.InDelta(t, colW+padding, b.Position().X, 0.5)
	assert.InDelta(t, colW, a.Size().Width, 0.5)
}

func TestFlexGrid_ProportionalFrColumns(t *testing.T) {
	// 1fr + 2fr + 1fr columns.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	c := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewFlexGridLayout([]layout.TrackSize{
		layout.Fr(1), layout.Fr(2), layout.Fr(1),
	})
	l.Layout([]fyne.CanvasObject{a, b, c}, fyne.NewSize(100, 20))

	available := float32(100) - padding*2
	col0W := available / 4
	col1W := available / 2

	assert.InDelta(t, col0W, a.Size().Width, 0.5)
	assert.InDelta(t, col1W, b.Size().Width, 0.5)
}

func TestFlexGrid_FixedAndFrColumns(t *testing.T) {
	// 100dp fixed + 1fr.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewFlexGridLayout([]layout.TrackSize{
		layout.Fixed(100), layout.Fr(1),
	})
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(200, 20))

	assert.InDelta(t, float32(100), a.Size().Width, 0.5, "fixed column is 100")
	assert.InDelta(t, float32(200)-100-padding, b.Size().Width, 0.5, "fr column fills rest")
}

func TestFlexGrid_AutoRows(t *testing.T) {
	// Two rows of 1 column each; row heights from content.
	a := NewMinSizeRect(fyne.NewSize(50, 30))
	b := NewMinSizeRect(fyne.NewSize(50, 40))
	padding := theme.Padding()

	l := layout.NewFlexGridLayout([]layout.TrackSize{layout.Fr(1)})
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 200))

	assert.Equal(t, float32(0), a.Position().Y)
	assert.InDelta(t, float32(30)+padding, b.Position().Y, 0.5)
	assert.InDelta(t, float32(30), a.Size().Height, 0.5)
	assert.InDelta(t, float32(40), b.Size().Height, 0.5)
}

func TestFlexGrid_ColSpan(t *testing.T) {
	// 3-column grid; first item spans 2 columns.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	c := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewFlexGridLayout([]layout.TrackSize{
		layout.Fr(1), layout.Fr(1), layout.Fr(1),
	})
	l.Layout([]fyne.CanvasObject{
		layout.WithSpan(a, 2, 1),
		b,
		c,
	}, fyne.NewSize(100, 20))

	available := float32(100) - padding*2
	colW := available / 3
	spanW := colW*2 + padding

	assert.InDelta(t, spanW, a.Size().Width, 0.5, "a spans 2 cols")
	assert.InDelta(t, spanW+padding, b.Position().X, 0.5, "b starts after a's span")
}

func TestFlexGrid_RowSpan(t *testing.T) {
	// 2-column grid; first item spans 2 rows.
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))
	c := NewMinSizeRect(fyne.NewSize(10, 20))
	padding := theme.Padding()

	l := layout.NewFlexGridLayout([]layout.TrackSize{layout.Fr(1), layout.Fr(1)})
	l.Layout([]fyne.CanvasObject{
		layout.WithSpan(a, 1, 2), // a spans 2 rows in col 0
		b,                         // b in col 1, row 0
		c,                         // c in col 1, row 1 (col 0 is occupied by a)
	}, fyne.NewSize(100, 50))

	// b and c should be at x = half of 100+padding ≈ col 1 position
	assert.True(t, b.Position().X > 0, "b in col 1")
	// a should span 2 rows in height
	assert.True(t, a.Size().Height > 20, "a spans 2 rows")
	_ = padding
}

func TestFlexGrid_MinSize(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(20, 30))
	b := NewMinSizeRect(fyne.NewSize(20, 30))

	l := layout.NewFlexGridLayout([]layout.TrackSize{layout.Fr(1), layout.Fr(1)})
	min := l.MinSize([]fyne.CanvasObject{a, b})

	assert.True(t, min.Width >= 0)
	assert.True(t, min.Height > 0)
}

func TestFlexGrid_WithExplicitRows(t *testing.T) {
	// Explicit rows: 1fr, 2fr.
	a := NewMinSizeRect(fyne.NewSize(10, 10))
	b := NewMinSizeRect(fyne.NewSize(10, 10))

	l := layout.NewFlexGridLayout(
		[]layout.TrackSize{layout.Fr(1)},
		layout.WithFlexRows([]layout.TrackSize{layout.Fr(1), layout.Fr(2)}),
	)
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 90))

	padding := theme.Padding()
	row0H := (float32(90) - padding) / 3 * 1
	assert.InDelta(t, row0H, a.Size().Height, 0.5)
	assert.InDelta(t, row0H*2, b.Size().Height, 0.5)
}

func TestFlexGrid_CustomGaps(t *testing.T) {
	a := NewMinSizeRect(fyne.NewSize(10, 20))
	b := NewMinSizeRect(fyne.NewSize(10, 20))

	l := layout.NewFlexGridLayout(
		[]layout.TrackSize{layout.Fr(1), layout.Fr(1)},
		layout.WithColGap(10),
	)
	l.Layout([]fyne.CanvasObject{a, b}, fyne.NewSize(100, 20))

	colW := (100 - 10) / 2
	assert.InDelta(t, colW, a.Size().Width, 0.5)
	assert.InDelta(t, colW+10, b.Position().X, 0.5)
}
