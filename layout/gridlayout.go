package layout

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*gridLayout)(nil)

type gridLayout struct {
	Cols            int
	vertical, adapt bool
}

// NewAdaptiveGridLayout returns a new grid layout which uses columns when horizontal but rows when vertical.
func NewAdaptiveGridLayout(rowcols int) fyne.Layout {
	return &gridLayout{Cols: rowcols, adapt: true}
}

// NewGridLayout returns a grid layout arranged in a specified number of columns.
// The number of rows will depend on how many children are in the container that uses this layout.
func NewGridLayout(cols int) fyne.Layout {
	return NewGridLayoutWithColumns(cols)
}

// NewGridLayoutWithColumns returns a new grid layout that specifies a column count and wrap to new rows when needed.
func NewGridLayoutWithColumns(cols int) fyne.Layout {
	return &gridLayout{Cols: cols}
}

// NewGridLayoutWithRows returns a new grid layout that specifies a row count that creates new rows as required.
func NewGridLayoutWithRows(rows int) fyne.Layout {
	return &gridLayout{Cols: rows, vertical: true}
}

// NewWeightedGridLayout returns a grid layout with proportional column widths
// specified by the given weight values. A weight of 2 means a column is twice as
// wide as a column with weight 1. Children are laid out left-to-right and wrap
// to a new row when all columns are filled.
// Children can carry ColSpan/RowSpan hints via WithSpan to span multiple cells.
//
// Since: 2.6
func NewWeightedGridLayout(weights ...float32) fyne.Layout {
	return &weightedGridLayout{weights: weights}
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*weightedGridLayout)(nil)

type weightedGridLayout struct {
	weights []float32
}

func (wg *weightedGridLayout) cols() int {
	if len(wg.weights) == 0 {
		return 1
	}
	return len(wg.weights)
}

func (wg *weightedGridLayout) resolveColWidths(containerWidth float32) []float32 {
	cols := wg.cols()
	padding := theme.Padding()
	totalWeight := float32(0)
	for _, wt := range wg.weights {
		totalWeight += wt
	}
	available := containerWidth - padding*float32(cols-1)
	widths := make([]float32, cols)
	for i, wt := range wg.weights {
		widths[i] = available * wt / totalWeight
	}
	return widths
}

func (wg *weightedGridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	cols := wg.cols()
	padding := theme.Padding()
	colWidths := wg.resolveColWidths(size.Width)

	type weightedCell struct {
		obj     fyne.CanvasObject
		colSpan int
		rowSpan int
	}
	var cells []weightedCell
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
		cells = append(cells, weightedCell{obj: obj, colSpan: cs, rowSpan: rs})
	}
	if len(cells) == 0 {
		return
	}

	// Compute row heights.
	rows := int(math.Ceil(float64(len(cells)) / float64(cols)))
	rowHeights := make([]float32, rows)
	col, row := 0, 0
	for _, c := range cells {
		if row < len(rowHeights) {
			if h := c.obj.MinSize().Height; h > rowHeights[row] {
				rowHeights[row] = h
			}
		}
		col += c.colSpan
		if col >= cols {
			col = 0
			row++
		}
	}

	// Place cells.
	col, row = 0, 0
	for _, c := range cells {
		x := float32(0)
		for i := 0; i < col; i++ {
			x += colWidths[i] + padding
		}
		y := float32(0)
		for i := 0; i < row; i++ {
			y += rowHeights[i] + padding
		}
		cellW := float32(0)
		for i := col; i < col+c.colSpan && i < cols; i++ {
			cellW += colWidths[i]
			if i > col {
				cellW += padding
			}
		}
		cellH := float32(0)
		for i := row; i < row+c.rowSpan && i < rows; i++ {
			cellH += rowHeights[i]
			if i > row {
				cellH += padding
			}
		}
		c.obj.Move(fyne.NewPos(x, y))
		c.obj.Resize(fyne.NewSize(cellW, cellH))

		col += c.colSpan
		if col >= cols {
			col = 0
			row++
		}
	}
}

func (wg *weightedGridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	cols := wg.cols()
	padding := theme.Padding()

	visible := 0
	maxMin := fyne.NewSize(0, 0)
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		visible++
		maxMin = maxMin.Max(obj.MinSize())
	}
	if visible == 0 {
		return fyne.NewSize(0, 0)
	}

	rows := int(math.Ceil(float64(visible) / float64(cols)))

	// Minimum container width: the smallest-weight column must fit maxMin.Width,
	// so total width = maxMin.Width * (totalWeight/minWeight) + padding*(cols-1).
	totalWeight := float32(0)
	minWeight := wg.weights[0]
	for _, wt := range wg.weights {
		totalWeight += wt
		if wt < minWeight {
			minWeight = wt
		}
	}
	minWidth := maxMin.Width*totalWeight/minWeight + padding*float32(cols-1)
	minHeight := maxMin.Height*float32(rows) + padding*float32(rows-1)
	return fyne.NewSize(minWidth, minHeight)
}

func (g *gridLayout) horizontal() bool {
	if g.adapt {
		return fyne.IsHorizontal(fyne.CurrentDevice().Orientation())
	}

	return !g.vertical
}

func (g *gridLayout) countRows(objects []fyne.CanvasObject) int {
	if g.Cols < 1 {
		g.Cols = 1
	}
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(g.Cols)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getLeading(size float64, offset int) float32 {
	ret := (size + float64(theme.Padding())) * float64(offset)
	return float32(ret)
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getTrailing(size float64, offset int) float32 {
	return getLeading(size, offset+1) - theme.Padding()
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *gridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	rows := g.countRows(objects)

	padding := theme.Padding()

	primaryObjects := rows
	secondaryObjects := g.Cols
	if g.horizontal() {
		primaryObjects, secondaryObjects = secondaryObjects, primaryObjects
	}

	padWidth := float32(primaryObjects-1) * padding
	padHeight := float32(secondaryObjects-1) * padding
	cellWidth := float64(size.Width-padWidth) / float64(primaryObjects)
	cellHeight := float64(size.Height-padHeight) / float64(secondaryObjects)

	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		x1 := getLeading(cellWidth, col)
		y1 := getLeading(cellHeight, row)
		x2 := getTrailing(cellWidth, col)
		y2 := getTrailing(cellHeight, row)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if g.horizontal() {
			if (i+1)%g.Cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%g.Cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
		i++
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a GridLayout this is the size of the largest child object multiplied by
// the required number of columns and rows, with appropriate padding between
// children.
func (g *gridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := g.countRows(objects)
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	padding := theme.Padding()

	primaryObjects := rows
	secondaryObjects := g.Cols
	if g.horizontal() {
		primaryObjects, secondaryObjects = secondaryObjects, primaryObjects
	}

	width := minSize.Width * float32(primaryObjects)
	height := minSize.Height * float32(secondaryObjects)
	xpad := padding * fyne.Max(float32(primaryObjects-1), 0)
	ypad := padding * fyne.Max(float32(secondaryObjects-1), 0)

	return fyne.NewSize(width+xpad, height+ypad)
}
