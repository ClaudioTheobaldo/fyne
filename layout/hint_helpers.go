package layout

import "fyne.io/fyne/v2"

// WithHint wraps obj with the given LayoutHint and returns a CanvasObject
// that carries the hint alongside the original object.
// Layouts that understand hints will read them via type assertion; other
// layouts see a plain CanvasObject and behave as before.
//
// Since: 2.6
func WithHint(obj fyne.CanvasObject, hint LayoutHint) fyne.CanvasObject {
	return newHintedObject(obj, hint)
}

// WithGrow wraps obj with a LayoutHint whose Grow field is set to grow.
// A grow value of 1 means the child absorbs its proportional share of free
// space in the main axis.
//
// Since: 2.6
func WithGrow(obj fyne.CanvasObject, grow float32) fyne.CanvasObject {
	return newHintedObject(obj, LayoutHint{Grow: grow})
}

// WithFlex wraps obj with explicit flex-grow, flex-shrink, and flex-basis values.
//
// Since: 2.6
func WithFlex(obj fyne.CanvasObject, grow, shrink, basis float32) fyne.CanvasObject {
	return newHintedObject(obj, LayoutHint{Grow: grow, Shrink: shrink, Basis: basis})
}

// WithAlign wraps obj with a per-child cross-axis alignment hint.
//
// Since: 2.6
func WithAlign(obj fyne.CanvasObject, align Alignment) fyne.CanvasObject {
	return newHintedObject(obj, LayoutHint{Align: align})
}

// WithMargin wraps obj with outer margin spacing on all four sides.
//
// Since: 2.6
func WithMargin(obj fyne.CanvasObject, margin Insets) fyne.CanvasObject {
	return newHintedObject(obj, LayoutHint{Margin: margin})
}

// WithMaxSize wraps obj with maximum width and height constraints.
// A value of 0 for either dimension means unconstrained.
//
// Since: 2.6
func WithMaxSize(obj fyne.CanvasObject, maxW, maxH float32) fyne.CanvasObject {
	return newHintedObject(obj, LayoutHint{MaxWidth: maxW, MaxHeight: maxH})
}

// WithSpan wraps obj with a grid column/row span hint.
// Values of 0 or 1 both mean a single cell.
//
// Since: 2.6
func WithSpan(obj fyne.CanvasObject, colSpan, rowSpan int) fyne.CanvasObject {
	return newHintedObject(obj, LayoutHint{ColSpan: colSpan, RowSpan: rowSpan})
}

// GetHint returns the LayoutHint for obj if it implements Hinted,
// plus a boolean indicating whether a hint was found.
// If obj does not carry a hint a zero-value LayoutHint is returned.
//
// Since: 2.6
func GetHint(obj fyne.CanvasObject) (LayoutHint, bool) {
	if h, ok := obj.(Hinted); ok {
		return h.LayoutHint(), true
	}
	return LayoutHint{}, false
}

// UnwrapObject returns the inner CanvasObject if obj is a HintedObject,
// otherwise returns obj unchanged.
//
// Since: 2.6
func UnwrapObject(obj fyne.CanvasObject) fyne.CanvasObject {
	if h, ok := obj.(Hinted); ok {
		return h.Unwrap()
	}
	return obj
}
