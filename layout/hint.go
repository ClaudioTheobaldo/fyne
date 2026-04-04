package layout

import "fyne.io/fyne/v2"

// Insets describes spacing on all four sides of an object.
type Insets struct {
	Top, Bottom, Left, Right float32
}

// NewUniformInsets returns Insets with the same value on all four sides.
func NewUniformInsets(v float32) Insets {
	return Insets{Top: v, Bottom: v, Left: v, Right: v}
}

// NewSymmetricInsets returns Insets with given horizontal and vertical values.
func NewSymmetricInsets(horizontal, vertical float32) Insets {
	return Insets{Top: vertical, Bottom: vertical, Left: horizontal, Right: horizontal}
}

// Alignment controls per-child cross-axis positioning within a layout.
type Alignment int

const (
	// AlignDefault defers to the container's align-items setting.
	AlignDefault Alignment = iota
	// AlignStart positions the child at the start of the cross axis.
	AlignStart
	// AlignCenter centers the child on the cross axis.
	AlignCenter
	// AlignEnd positions the child at the end of the cross axis.
	AlignEnd
	// AlignStretch stretches the child to fill the cross axis (default for VBox/HBox).
	AlignStretch
)

// Justify controls how children are distributed along the main axis.
type Justify int

const (
	// JustifyStart packs children toward the start (default).
	JustifyStart Justify = iota
	// JustifyEnd packs children toward the end.
	JustifyEnd
	// JustifyCenter centers children along the main axis.
	JustifyCenter
	// JustifySpaceBetween places equal space between children (no space at edges).
	JustifySpaceBetween
	// JustifySpaceAround places equal space around each child (half-size at edges).
	JustifySpaceAround
	// JustifySpaceEvenly places equal space before, between, and after children.
	JustifySpaceEvenly
)

// LayoutHint carries per-child layout metadata. Attach it to a CanvasObject
// using WithHint or one of the convenience wrappers (WithGrow, WithAlign, etc.).
// All zero values produce the same behaviour as the current layout algorithms.
//
// Since: 2.6
type LayoutHint struct {
	// Grow is the flex-grow factor. 0 means the child does not grow beyond its basis.
	Grow float32
	// Shrink is the flex-shrink factor. 0 means the child does not shrink below MinSize.
	Shrink float32
	// Basis is the initial main-axis size before grow/shrink is applied.
	// 0 means use the child's MinSize in the main axis.
	Basis float32

	// Align is the per-child cross-axis alignment (align-self).
	// AlignDefault defers to the container's AlignItems setting.
	Align Alignment

	// Margin adds outer spacing around the child on all four sides.
	Margin Insets

	// MaxWidth caps the child's width. 0 means unconstrained.
	MaxWidth float32
	// MaxHeight caps the child's height. 0 means unconstrained.
	MaxHeight float32

	// ColSpan is the number of grid columns the child spans. 0 and 1 both mean one cell.
	ColSpan int
	// RowSpan is the number of grid rows the child spans. 0 and 1 both mean one cell.
	RowSpan int

	// Order overrides visual placement order. 0 means use source order.
	Order int
}

// Hinted is implemented by objects that carry a LayoutHint.
// Layouts use a type assertion to detect and read hints.
//
// Since: 2.6
type Hinted interface {
	LayoutHint() LayoutHint
	Unwrap() fyne.CanvasObject
}

// MaxSizable is an optional interface for objects that declare a maximum size.
// Layouts that support it will not make the child larger than its MaxSize.
//
// Since: 2.6
type MaxSizable interface {
	MaxSize() fyne.Size
}
