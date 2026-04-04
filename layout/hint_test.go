package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/test"
)

func newRect(w, h float32) *canvas.Rectangle {
	r := canvas.NewRectangle(color.Black)
	r.SetMinSize(fyne.NewSize(w, h))
	return r
}

// --- HintedObject delegation ---

func TestHintedObject_DelegatesVisible(t *testing.T) {
	r := newRect(10, 10)
	h := layout.WithGrow(r, 1)
	if !h.Visible() {
		t.Fatal("should be visible by default")
	}
	r.Hide()
	if h.Visible() {
		t.Fatal("should reflect hidden state of inner object")
	}
}

func TestHintedObject_DelegatesHideShow(t *testing.T) {
	r := newRect(10, 10)
	h := layout.WithGrow(r, 1)
	h.Hide()
	if r.Visible() {
		t.Fatal("Hide() should hide inner object")
	}
	h.Show()
	if !r.Visible() {
		t.Fatal("Show() should show inner object")
	}
}

// --- MinSize includes margins ---

func TestHintedObject_MinSizeNoMargin(t *testing.T) {
	r := newRect(50, 30)
	h := layout.WithGrow(r, 1)
	got := h.MinSize()
	if got.Width != 50 || got.Height != 30 {
		t.Fatalf("expected 50x30, got %v", got)
	}
}

func TestHintedObject_MinSizeWithMargin(t *testing.T) {
	r := newRect(50, 30)
	h := layout.WithMargin(r, layout.Insets{Top: 5, Bottom: 10, Left: 3, Right: 7})
	got := h.MinSize()
	if got.Width != 60 || got.Height != 45 {
		t.Fatalf("expected 60x45, got %v", got)
	}
}

func TestHintedObject_MinSizeUniformMargin(t *testing.T) {
	r := newRect(20, 20)
	h := layout.WithMargin(r, layout.NewUniformInsets(5))
	got := h.MinSize()
	if got.Width != 30 || got.Height != 30 {
		t.Fatalf("expected 30x30, got %v", got)
	}
}

// --- Move offsets by margin ---

func TestHintedObject_MoveAppliesMarginOffset(t *testing.T) {
	r := newRect(50, 30)
	h := layout.WithMargin(r, layout.Insets{Top: 5, Left: 3})
	// Resize triggers renderer.Layout, which positions the inner object
	// at LOCAL (Margin.Left, Margin.Top) = (3, 5) within the widget bounds.
	h.Resize(fyne.NewSize(53, 35))
	innerPos := r.Position()
	if innerPos.X != 3 || innerPos.Y != 5 {
		t.Fatalf("expected inner local pos (3,5), got %v", innerPos)
	}
}

func TestHintedObject_MoveNoMargin(t *testing.T) {
	r := newRect(50, 30)
	h := layout.WithGrow(r, 1)
	// No margin: inner object at LOCAL (0, 0) within the widget.
	h.Resize(fyne.NewSize(50, 30))
	innerPos := r.Position()
	if innerPos.X != 0 || innerPos.Y != 0 {
		t.Fatalf("expected inner local pos (0,0), got %v", innerPos)
	}
}

// --- Position returns outer box position ---

func TestHintedObject_PositionReturnsOuterPos(t *testing.T) {
	r := newRect(50, 30)
	h := layout.WithMargin(r, layout.Insets{Top: 5, Left: 3})
	h.Move(fyne.NewPos(10, 20))
	outerPos := h.Position()
	if outerPos.X != 10 || outerPos.Y != 20 {
		t.Fatalf("expected outer pos (10,20), got %v", outerPos)
	}
}

// --- Resize accounts for margins ---

func TestHintedObject_ResizeAppliesMargin(t *testing.T) {
	r := newRect(10, 10)
	h := layout.WithMargin(r, layout.Insets{Top: 5, Bottom: 10, Left: 3, Right: 7})
	h.Resize(fyne.NewSize(100, 80))
	innerSize := r.Size()
	if innerSize.Width != 90 || innerSize.Height != 65 {
		t.Fatalf("expected inner 90x65, got %v", innerSize)
	}
}

// --- MaxSize clamping ---

func TestHintedObject_ResizeClampMaxWidth(t *testing.T) {
	r := newRect(10, 10)
	h := layout.WithMaxSize(r, 40, 0)
	h.Resize(fyne.NewSize(100, 80))
	innerSize := r.Size()
	if innerSize.Width != 40 {
		t.Fatalf("expected width clamped to 40, got %v", innerSize.Width)
	}
	if innerSize.Height != 80 {
		t.Fatalf("expected height 80 (unconstrained), got %v", innerSize.Height)
	}
}

func TestHintedObject_ResizeClampMaxHeight(t *testing.T) {
	r := newRect(10, 10)
	h := layout.WithMaxSize(r, 0, 50)
	h.Resize(fyne.NewSize(100, 80))
	innerSize := r.Size()
	if innerSize.Height != 50 {
		t.Fatalf("expected height clamped to 50, got %v", innerSize.Height)
	}
}

// --- GetHint ---

func TestGetHint_ReturnsHintForHintedObject(t *testing.T) {
	r := newRect(10, 10)
	h := layout.WithGrow(r, 2.5)
	hint, ok := layout.GetHint(h)
	if !ok {
		t.Fatal("expected ok=true for HintedObject")
	}
	if hint.Grow != 2.5 {
		t.Fatalf("expected Grow=2.5, got %v", hint.Grow)
	}
}

func TestGetHint_ReturnsZeroForPlainObject(t *testing.T) {
	r := newRect(10, 10)
	hint, ok := layout.GetHint(r)
	if ok {
		t.Fatal("expected ok=false for plain CanvasObject")
	}
	if hint.Grow != 0 || hint.Shrink != 0 || hint.Align != layout.AlignDefault {
		t.Fatal("expected zero-value hint for plain object")
	}
}

// --- UnwrapObject ---

func TestUnwrapObject_ReturnsInnerForHintedObject(t *testing.T) {
	r := newRect(10, 10)
	h := layout.WithGrow(r, 1)
	inner := layout.UnwrapObject(h)
	if inner != r {
		t.Fatal("expected UnwrapObject to return original rectangle")
	}
}

func TestUnwrapObject_ReturnsObjectForPlainObject(t *testing.T) {
	r := newRect(10, 10)
	result := layout.UnwrapObject(r)
	if result != r {
		t.Fatal("expected UnwrapObject to return object unchanged")
	}
}

// --- Insets helpers ---

func TestNewUniformInsets(t *testing.T) {
	ins := layout.NewUniformInsets(8)
	if ins.Top != 8 || ins.Bottom != 8 || ins.Left != 8 || ins.Right != 8 {
		t.Fatalf("expected all sides 8, got %+v", ins)
	}
}

func TestNewSymmetricInsets(t *testing.T) {
	ins := layout.NewSymmetricInsets(4, 6)
	if ins.Left != 4 || ins.Right != 4 || ins.Top != 6 || ins.Bottom != 6 {
		t.Fatalf("expected horizontal=4 vertical=6, got %+v", ins)
	}
}

// --- WithHint merges all fields ---

func TestWithHint_FullHint(t *testing.T) {
	r := newRect(10, 10)
	hint := layout.LayoutHint{
		Grow:      1,
		Shrink:    0.5,
		Basis:     20,
		Align:     layout.AlignCenter,
		MaxWidth:  100,
		MaxHeight: 200,
		ColSpan:   2,
		RowSpan:   3,
	}
	h := layout.WithHint(r, hint)
	got, ok := layout.GetHint(h)
	if !ok {
		t.Fatal("expected hint")
	}
	if got.Grow != 1 || got.Shrink != 0.5 || got.Basis != 20 ||
		got.Align != layout.AlignCenter || got.MaxWidth != 100 ||
		got.MaxHeight != 200 || got.ColSpan != 2 || got.RowSpan != 3 {
		t.Fatalf("hint fields not preserved: %+v", got)
	}
}

// Suppress unused import warning for theme.
var _ = theme.Padding
