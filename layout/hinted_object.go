package layout

import "fyne.io/fyne/v2"

// Declare conformity with interfaces.
var _ fyne.Widget = (*HintedObject)(nil)
var _ Hinted = (*HintedObject)(nil)

// HintedObject wraps a CanvasObject, attaches a LayoutHint, and is itself a
// proper fyne.Widget so Fyne's canvas renders it via its renderer.
//
// The hint's Margin is applied in the renderer: the inner object is placed at
// LOCAL position (Margin.Left, Margin.Top) within the widget bounds and sized
// to (widgetW − margins) × (widgetH − margins). MaxWidth/MaxHeight further cap
// growth. Because positions are LOCAL (relative to the widget's own origin),
// Fyne's canvas correctly adds the widget's canvas position when rendering.
//
// Since: 2.6
type HintedObject struct {
	child fyne.CanvasObject
	hint  LayoutHint
	pos   fyne.Position
	size  fyne.Size
}

// newHintedObject creates a HintedObject wrapping child with the given hint.
func newHintedObject(child fyne.CanvasObject, hint LayoutHint) *HintedObject {
	return &HintedObject{child: child, hint: hint}
}

// ── Hinted ────────────────────────────────────────────────────────────────────

// LayoutHint returns the hint attached to this object.
func (h *HintedObject) LayoutHint() LayoutHint { return h.hint }

// Unwrap returns the inner CanvasObject.
func (h *HintedObject) Unwrap() fyne.CanvasObject { return h.child }

// ── fyne.CanvasObject ─────────────────────────────────────────────────────────

// MinSize returns the inner object's MinSize expanded by the hint's Margin.
func (h *HintedObject) MinSize() fyne.Size {
	inner := h.child.MinSize()
	m := h.hint.Margin
	return fyne.NewSize(inner.Width+m.Left+m.Right, inner.Height+m.Top+m.Bottom)
}

// Move records the widget's canvas position. The inner object's LOCAL position
// (margin offset) is not updated here — it is set by Resize → doLayout.
// Fyne's canvas adds the widget's position to renderer objects during rendering.
func (h *HintedObject) Move(pos fyne.Position) { h.pos = pos }

// Position returns the widget's canvas position.
func (h *HintedObject) Position() fyne.Position { return h.pos }

// Resize records the widget size and immediately repositions and resizes the
// inner object according to the margin and MaxSize constraints.
func (h *HintedObject) Resize(size fyne.Size) {
	h.size = size
	h.doLayout(size)
}

// Size returns the widget's current size.
func (h *HintedObject) Size() fyne.Size { return h.size }

// Visible delegates to the inner object — hiding the child also hides the wrapper.
func (h *HintedObject) Visible() bool { return h.child.Visible() }

// Hide hides the inner object (and therefore this wrapper).
func (h *HintedObject) Hide() { h.child.Hide() }

// Show shows the inner object (and therefore this wrapper).
func (h *HintedObject) Show() { h.child.Show() }

// Refresh requests a repaint of the inner object.
func (h *HintedObject) Refresh() { h.child.Refresh() }

// ── fyne.Widget ───────────────────────────────────────────────────────────────

// CreateRenderer returns the renderer for this widget.
// The renderer's Objects() exposes the inner child to Fyne's canvas so it is
// actually painted; its Layout() mirrors doLayout for Fyne-triggered reflows.
func (h *HintedObject) CreateRenderer() fyne.WidgetRenderer {
	return &hintedRenderer{child: h.child, hint: &h.hint}
}

// ── internal ──────────────────────────────────────────────────────────────────

// doLayout places the inner child at LOCAL (Margin.Left, Margin.Top) within
// the widget bounds and sizes it to size − margins, clamped by MaxWidth/MaxHeight.
func (h *HintedObject) doLayout(size fyne.Size) {
	m := h.hint.Margin
	w := size.Width - m.Left - m.Right
	ht := size.Height - m.Top - m.Bottom
	if h.hint.MaxWidth > 0 && w > h.hint.MaxWidth {
		w = h.hint.MaxWidth
	}
	if h.hint.MaxHeight > 0 && ht > h.hint.MaxHeight {
		ht = h.hint.MaxHeight
	}
	if w < 0 {
		w = 0
	}
	if ht < 0 {
		ht = 0
	}
	h.child.Move(fyne.NewPos(m.Left, m.Top))
	h.child.Resize(fyne.NewSize(w, ht))
}

// ── renderer ──────────────────────────────────────────────────────────────────

type hintedRenderer struct {
	child fyne.CanvasObject
	hint  *LayoutHint
}

// Layout mirrors doLayout so Fyne-triggered reflows (e.g. window resize) also
// apply the margin and MaxSize constraints.
func (r *hintedRenderer) Layout(size fyne.Size) {
	m := r.hint.Margin
	w := size.Width - m.Left - m.Right
	h := size.Height - m.Top - m.Bottom
	if r.hint.MaxWidth > 0 && w > r.hint.MaxWidth {
		w = r.hint.MaxWidth
	}
	if r.hint.MaxHeight > 0 && h > r.hint.MaxHeight {
		h = r.hint.MaxHeight
	}
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	r.child.Move(fyne.NewPos(m.Left, m.Top))
	r.child.Resize(fyne.NewSize(w, h))
}

// MinSize returns the inner object's MinSize expanded by the hint's Margin.
func (r *hintedRenderer) MinSize() fyne.Size {
	inner := r.child.MinSize()
	m := r.hint.Margin
	return fyne.NewSize(inner.Width+m.Left+m.Right, inner.Height+m.Top+m.Bottom)
}

func (r *hintedRenderer) Refresh() { r.child.Refresh() }
func (r *hintedRenderer) Destroy() {}
func (r *hintedRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.child} }
