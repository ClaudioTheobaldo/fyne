package fyne

// Layout defines how [CanvasObject]s may be laid out in a specified Size.
//
// Deprecated: Implement [ConstrainedLayout] instead for correct intrinsic
// sizing — particularly for layouts whose height depends on their width
// (e.g. wrapping, auto-sizing). The old interface remains fully supported.
type Layout interface {
	// Layout will manipulate the listed [CanvasObject]s Size and Position
	// to fit within the specified size.
	Layout([]CanvasObject, Size)
	// MinSize calculates the smallest size that will fit the listed
	// [CanvasObject]s using this Layout algorithm.
	MinSize(objects []CanvasObject) Size
}

// Constraints describes the available space passed to a [ConstrainedLayout].
// A MaxSize component of +Inf means "no upper bound" on that axis.
type Constraints struct {
	// MinSize is the smallest size the container is allowed to occupy.
	MinSize Size
	// MaxSize is the largest size the container may occupy.
	// Use [NewInfSize] for an unconstrained axis.
	MaxSize Size
}

// NewInfSize returns a Size with both components set to +Inf, representing
// an unconstrained layout axis.
func NewInfSize() Size {
	return NewSize(inf, inf)
}

const inf = float32(1e9) // sentinel for "unbounded"

// IntrinsicSizer is an optional interface for [CanvasObject]s whose height
// depends on their available width (e.g. wrapping containers, text widgets).
// A parent layout should call IntrinsicSize(availW) instead of MinSize() when
// it knows the width it will allocate, so the child can report its true height.
type IntrinsicSizer interface {
	IntrinsicSize(availW float32) Size
}

// ConstrainedLayout extends [Layout] with a two-pass measurement protocol.
// LayoutConstrained is used for measurement only — it must be safe to call
// without side effects on positions/sizes when used as a measurement probe.
// The canonical use is [IntrinsicSizer]: a parent layout calls
// IntrinsicSize(availW) on a child container, which in turn calls
// LayoutConstrained to discover its true height at that width.
//
// Container.layout() still calls the plain Layout() method for the actual
// layout pass.  ConstrainedLayout is only consulted during measurement.
//
// Layouts that do not depend on the available width/height for their
// intrinsic size (e.g. simple horizontal/vertical stacks) need not implement
// this interface.
type ConstrainedLayout interface {
	Layout
	// LayoutConstrained measures and arranges objects within the given
	// constraints and returns the actual size consumed.
	// It is called both as a measurement probe (with MaxSize.Height = inf)
	// and as a full layout pass (with MaxSize = container size).
	LayoutConstrained(objects []CanvasObject, constraints Constraints) Size
}
