package layout

import "fyne.io/fyne/v2"

// Declare conformity with Layout interface
var _ fyne.Layout = (*centerLayout)(nil)

type centerLayout struct{}

// NewCenterLayout creates a new CenterLayout instance
func NewCenterLayout() fyne.Layout {
	return &centerLayout{}
}

// Layout is called to pack all child objects into a specified size.
// For CenterLayout this sets all children to their minimum size, centered within the space.
// Children carrying an Alignment hint via WithAlign override the default centering:
//   - AlignStart  → left-aligned (horizontal), top-aligned (vertical) stays centered
//   - AlignEnd    → right-aligned (horizontal), bottom-aligned (vertical) stays centered
//   - AlignCenter → explicitly centered (same as default)
//   - AlignStretch → stretched to fill the full container size
func (c *centerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, child := range objects {
		hint, hasHint := GetHint(child)
		childMin := child.MinSize()

		var w, h float32
		if hasHint && hint.Align == AlignStretch {
			w, h = size.Width, size.Height
		} else {
			w, h = childMin.Width, childMin.Height
		}
		child.Resize(fyne.NewSize(w, h))

		var x, y float32
		if hasHint {
			switch hint.Align {
			case AlignStart:
				x = 0
			case AlignEnd:
				x = size.Width - w
			case AlignStretch:
				x = 0
			default: // AlignDefault, AlignCenter
				x = (size.Width - w) / 2
			}
		} else {
			x = (size.Width - w) / 2
		}
		y = (size.Height - h) / 2

		child.Move(fyne.NewPos(x, y))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For CenterLayout this is determined simply as the MinSize of the largest child.
func (c *centerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}
