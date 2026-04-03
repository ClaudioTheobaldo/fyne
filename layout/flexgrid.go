package layout

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// TrackSize describes the size of a single grid track (column or row).
// If Fixed > 0 the track has a fixed size; otherwise Fr defines a fractional
// share of the remaining space (similar to CSS fr units).
//
// Since: 2.6
type TrackSize struct {
	Fr    float32 // fractional share (1fr, 2fr, …)
	Fixed float32 // fixed dp size; takes precedence over Fr when > 0
	Min   float32 // minimum size for fr tracks (0 = no minimum)
}

// Fr is a shorthand constructor for a fractional TrackSize.
func Fr(fr float32) TrackSize { return TrackSize{Fr: fr} }

// Fixed is a shorthand constructor for a fixed TrackSize.
func Fixed(px float32) TrackSize { return TrackSize{Fixed: px} }

// FlexGridOption configures a FlexGridLayout.
//
// Since: 2.6
type FlexGridOption func(*flexGridConfig)

type flexGridConfig struct {
	rows   []TrackSize
	colGap float32
	rowGap float32
}

// WithFlexRows sets explicit row track sizes. When not provided rows are
// auto-sized to the maximum MinSize.Height of the children in each row.
func WithFlexRows(rows []TrackSize) FlexGridOption {
	return func(c *flexGridConfig) { c.rows = rows }
}

// WithColGap overrides the column gap (default: theme.Padding()).
func WithColGap(gap float32) FlexGridOption {
	return func(c *flexGridConfig) { c.colGap = gap }
}

// WithRowGap overrides the row gap (default: theme.Padding()).
func WithRowGap(gap float32) FlexGridOption {
	return func(c *flexGridConfig) { c.rowGap = gap }
}

// Declare conformity with Layout interface.
var _ fyne.Layout = (*flexGridLayout)(nil)

type flexGridLayout struct {
	columns []TrackSize
	rows    []TrackSize // nil = auto
	colGap  float32     // 0 = use theme
	rowGap  float32     // 0 = use theme
}

// NewFlexGridLayout returns a CSS Grid-like layout with columns defined by
// the given TrackSize slice. Use Fr(n) for fractional columns and Fixed(n) for
// fixed-width columns. Rows are auto-sized by default.
//
// Children can carry ColSpan/RowSpan hints via WithSpan.
//
//	// 3 columns: 1fr, 2fr, 100dp fixed
//	layout.NewFlexGridLayout([]layout.TrackSize{
//	    layout.Fr(1), layout.Fr(2), layout.Fixed(100),
//	})
//
// Since: 2.6
func NewFlexGridLayout(columns []TrackSize, opts ...FlexGridOption) fyne.Layout {
	cfg := &flexGridConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return &flexGridLayout{
		columns: columns,
		rows:    cfg.rows,
		colGap:  cfg.colGap,
		rowGap:  cfg.rowGap,
	}
}

func (f *flexGridLayout) colGapSize() float32 {
	if f.colGap > 0 {
		return f.colGap
	}
	return theme.Padding()
}

func (f *flexGridLayout) rowGapSize() float32 {
	if f.rowGap > 0 {
		return f.rowGap
	}
	return theme.Padding()
}

// resolveTrackWidths returns the resolved pixel width of each column.
func (f *flexGridLayout) resolveTrackWidths(containerWidth float32) []float32 {
	cols := len(f.columns)
	gap := f.colGapSize()
	totalFixed := float32(0)
	totalFr := float32(0)
	for _, t := range f.columns {
		if t.Fixed > 0 {
			totalFixed += t.Fixed
		} else {
			totalFr += t.Fr
		}
	}
	remaining := containerWidth - totalFixed - gap*float32(cols-1)
	if remaining < 0 {
		remaining = 0
	}
	widths := make([]float32, cols)
	for i, t := range f.columns {
		if t.Fixed > 0 {
			widths[i] = t.Fixed
		} else {
			if totalFr > 0 {
				w := remaining * t.Fr / totalFr
				if t.Min > 0 && w < t.Min {
					w = t.Min
				}
				widths[i] = w
			}
		}
	}
	return widths
}

type flexCell struct {
	obj     fyne.CanvasObject
	colSpan int
	rowSpan int
}

func (f *flexGridLayout) collectCells(objects []fyne.CanvasObject) []flexCell {
	var cells []flexCell
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		hint, _ := GetHint(obj)
		cs := hint.ColSpan
		if cs < 1 {
			cs = 1
		}
		rs := hint.RowSpan
		if rs < 1 {
			rs = 1
		}
		cells = append(cells, flexCell{obj: obj, colSpan: cs, rowSpan: rs})
	}
	return cells
}

// placeResult holds the computed grid placement for all cells.
type placeResult struct {
	colStart []int
	rowStart []int
	numRows  int
}

// place computes the grid placement using a simple occupancy matrix.
func (f *flexGridLayout) place(cells []flexCell) placeResult {
	cols := len(f.columns)
	if cols == 0 {
		return placeResult{}
	}

	colStart := make([]int, len(cells))
	rowStart := make([]int, len(cells))

	// occupancy[row][col] = true if occupied
	occupied := [][]bool{}
	ensureRow := func(r int) {
		for len(occupied) <= r {
			occupied = append(occupied, make([]bool, cols))
		}
	}
	canPlace := func(r, c, cs, rs int) bool {
		if c+cs > cols {
			return false
		}
		for dr := 0; dr < rs; dr++ {
			ensureRow(r + dr)
			for dc := 0; dc < cs; dc++ {
				if occupied[r+dr][c+dc] {
					return false
				}
			}
		}
		return true
	}
	markOccupied := func(r, c, cs, rs int) {
		for dr := 0; dr < rs; dr++ {
			ensureRow(r + dr)
			for dc := 0; dc < cs; dc++ {
				occupied[r+dr][c+dc] = true
			}
		}
	}

	curRow, curCol := 0, 0
	for i, c := range cells {
		// Advance cursor to find a free cell.
		for {
			ensureRow(curRow)
			if canPlace(curRow, curCol, c.colSpan, c.rowSpan) {
				break
			}
			curCol++
			if curCol >= cols {
				curCol = 0
				curRow++
			}
		}
		colStart[i] = curCol
		rowStart[i] = curRow
		markOccupied(curRow, curCol, c.colSpan, c.rowSpan)
		// Advance cursor past this cell.
		curCol += c.colSpan
		if curCol >= cols {
			curCol = 0
			curRow++
		}
	}

	numRows := len(occupied)
	if numRows == 0 {
		numRows = 1
	}
	return placeResult{colStart: colStart, rowStart: rowStart, numRows: numRows}
}

// resolveRowHeights computes the pixel height for each row.
func (f *flexGridLayout) resolveRowHeights(cells []flexCell, placement placeResult, containerHeight float32) []float32 {
	numRows := placement.numRows
	if len(f.rows) > 0 {
		return f.resolveTrackHeights(containerHeight, numRows)
	}
	// Auto-size: each row is the max MinSize.Height of cells in that row.
	heights := make([]float32, numRows)
	for i, c := range cells {
		r := placement.rowStart[i]
		if r < numRows {
			if h := c.obj.MinSize().Height; h > heights[r] {
				heights[r] = h
			}
		}
	}
	return heights
}

func (f *flexGridLayout) resolveTrackHeights(containerHeight float32, numRows int) []float32 {
	rows := f.rows
	if len(rows) == 0 {
		return make([]float32, numRows)
	}
	gap := f.rowGapSize()
	totalFixed := float32(0)
	totalFr := float32(0)
	for _, t := range rows {
		if t.Fixed > 0 {
			totalFixed += t.Fixed
		} else {
			totalFr += t.Fr
		}
	}
	remaining := containerHeight - totalFixed - gap*float32(len(rows)-1)
	if remaining < 0 {
		remaining = 0
	}
	heights := make([]float32, len(rows))
	for i, t := range rows {
		if t.Fixed > 0 {
			heights[i] = t.Fixed
		} else if totalFr > 0 {
			h := remaining * t.Fr / totalFr
			if t.Min > 0 && h < t.Min {
				h = t.Min
			}
			heights[i] = h
		}
	}
	return heights
}

func (f *flexGridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(f.columns) == 0 {
		return
	}
	cells := f.collectCells(objects)
	if len(cells) == 0 {
		return
	}

	placement := f.place(cells)
	colWidths := f.resolveTrackWidths(size.Width)
	rowHeights := f.resolveRowHeights(cells, placement, size.Height)
	colGap := f.colGapSize()
	rowGap := f.rowGapSize()

	// Compute x positions for each column.
	colX := make([]float32, len(colWidths))
	x := float32(0)
	for i, w := range colWidths {
		colX[i] = x
		x += w + colGap
	}
	// Compute y positions for each row.
	rowY := make([]float32, len(rowHeights))
	y := float32(0)
	for i, h := range rowHeights {
		rowY[i] = y
		y += h + rowGap
	}

	for i, c := range cells {
		col := placement.colStart[i]
		row := placement.rowStart[i]

		cx := colX[col]
		cy := float32(0)
		if row < len(rowY) {
			cy = rowY[row]
		}

		// Width: span across colSpan columns.
		cw := float32(0)
		for dc := 0; dc < c.colSpan && col+dc < len(colWidths); dc++ {
			cw += colWidths[col+dc]
			if dc > 0 {
				cw += colGap
			}
		}
		// Height: span across rowSpan rows.
		ch := float32(0)
		for dr := 0; dr < c.rowSpan && row+dr < len(rowHeights); dr++ {
			ch += rowHeights[row+dr]
			if dr > 0 {
				ch += rowGap
			}
		}

		c.obj.Move(fyne.NewPos(cx, cy))
		c.obj.Resize(fyne.NewSize(cw, ch))
	}
}

func (f *flexGridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(f.columns) == 0 {
		return fyne.NewSize(0, 0)
	}
	cells := f.collectCells(objects)
	if len(cells) == 0 {
		return fyne.NewSize(0, 0)
	}

	cols := len(f.columns)
	colGap := f.colGapSize()
	rowGap := f.rowGapSize()

	// Min width: sum of fixed columns + fr minimums + gaps.
	minW := float32(0)
	for _, t := range f.columns {
		if t.Fixed > 0 {
			minW += t.Fixed
		} else {
			minW += t.Min
		}
	}
	minW += colGap * float32(cols-1)

	// Min height: auto-size rows using max child MinSize.Height per row.
	numRows := int(math.Ceil(float64(len(cells)) / float64(cols)))
	rowHeights := make([]float32, numRows)
	for i, c := range cells {
		row := i / cols
		if row < numRows {
			if h := c.obj.MinSize().Height; h > rowHeights[row] {
				rowHeights[row] = h
			}
		}
	}
	minH := float32(0)
	for _, h := range rowHeights {
		minH += h
	}
	minH += rowGap * float32(numRows-1)

	return fyne.NewSize(minW, minH)
}
