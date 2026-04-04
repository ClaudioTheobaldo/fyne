// Package main is a visual showcase of the new layout-hint system.
// Run with: go run ./cmd/layout_showcase
package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ── palette ───────────────────────────────────────────────────────────────────

var (
	cBlue   = color.NRGBA{0x42, 0x85, 0xF4, 0xFF}
	cGreen  = color.NRGBA{0x34, 0xA8, 0x53, 0xFF}
	cYellow = color.NRGBA{0xF9, 0xA8, 0x25, 0xFF}
	cRed    = color.NRGBA{0xEA, 0x43, 0x35, 0xFF}
	cPurple = color.NRGBA{0x9C, 0x27, 0xB0, 0xFF}
	cTeal   = color.NRGBA{0x00, 0x96, 0x88, 0xFF}
	cOrange = color.NRGBA{0xFF, 0x57, 0x22, 0xFF}
	cIndigo = color.NRGBA{0x3F, 0x51, 0xB5, 0xFF}
)

// ── helpers ───────────────────────────────────────────────────────────────────

// box returns a colored block with a white bold centered label.
func box(label string, col color.NRGBA, minW, minH float32) fyne.CanvasObject {
	bg := canvas.NewRectangle(col)
	bg.SetMinSize(fyne.NewSize(minW, minH))
	lbl := canvas.NewText(label, color.White)
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	lbl.Alignment = fyne.TextAlignCenter
	return container.NewStack(bg, container.NewCenter(lbl))
}

// note returns a word-wrapped explanation label.
func note(text string) *widget.Label {
	l := widget.NewLabel(text)
	l.Wrapping = fyne.TextWrapWord
	return l
}

// head returns a bold section heading.
func head(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
}

// withMinH enforces a minimum height on content via a transparent spacer in a Stack.
func withMinH(h float32, obj fyne.CanvasObject) fyne.CanvasObject {
	spacer := canvas.NewRectangle(color.Alpha{})
	spacer.SetMinSize(fyne.NewSize(0, h))
	return container.NewStack(spacer, obj)
}

// viewportLayout constrains its single child to a specific width, letting its
// height grow to the child's MinSize. Used to demonstrate wrap behaviour at
// controlled widths without relying on window resize.
type viewportLayout struct{ width float32 }

func (v *viewportLayout) Layout(objs []fyne.CanvasObject, _ fyne.Size) {
	for _, o := range objs {
		h := o.MinSize().Height
		o.Resize(fyne.NewSize(v.width, h))
		o.Move(fyne.NewPos(0, 0))
	}
}
func (v *viewportLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	h := float32(0)
	for _, o := range objs {
		if oh := o.MinSize().Height; oh > h {
			h = oh
		}
	}
	return fyne.NewSize(v.width, h)
}

// ── Tab 1: Flex Grow ──────────────────────────────────────────────────────────

func makeGrowTab() fyne.CanvasObject {
	// --- Interactive section ---
	bGrowLabel := widget.NewLabel("B  grow = 1.0")

	var flexRow *fyne.Container

	rebuild := func(bGrow float32) {
		if flexRow == nil {
			return
		}
		bGrowLabel.SetText(fmt.Sprintf("B  grow = %.1f", bGrow))
		flexRow.Objects = []fyne.CanvasObject{
			layout.WithGrow(box("A  ×1", cBlue, 70, 50), 1),
			layout.WithGrow(box(fmt.Sprintf("B  ×%.1f", bGrow), cGreen, 70, 50), bGrow),
			layout.WithGrow(box("C  ×1", cRed, 70, 50), 1),
		}
		flexRow.Refresh()
	}

	slider := widget.NewSlider(0, 5)
	slider.Step = 0.1
	slider.Value = 1
	slider.OnChanged = func(v float64) { rebuild(float32(v)) }

	flexRow = container.New(
		layout.NewFlexHBoxLayout(),
		layout.WithGrow(box("A  ×1", cBlue, 70, 50), 1),
		layout.WithGrow(box("B  ×1.0", cGreen, 70, 50), 1),
		layout.WithGrow(box("C  ×1", cRed, 70, 50), 1),
	)

	return container.NewVBox(
		note("Old HBox: children stay at MinSize — free space is wasted at the right edge.\n"+
			"FlexHBox + WithGrow: remaining space distributed proportionally by grow factor."),
		widget.NewSeparator(),

		head("Old  NewHBox  —  wasted space at right:"),
		container.NewHBox(
			box("A", cBlue, 70, 50),
			box("B", cGreen, 70, 50),
			box("C", cRed, 70, 50),
		),
		widget.NewSeparator(),

		head("New  FlexHBox  —  grow = 1 : 2 : 1, boxes fill 100% of the width:"),
		container.New(
			layout.NewFlexHBoxLayout(),
			layout.WithGrow(box("A  ×1", cBlue, 70, 50), 1),
			layout.WithGrow(box("B  ×2", cGreen, 70, 50), 2),
			layout.WithGrow(box("C  ×1", cRed, 70, 50), 1),
		),
		widget.NewSeparator(),

		head("New  FlexVBox  —  grow = 1 : 3 : 1 fills height:"),
		withMinH(190,
			container.New(
				layout.NewFlexVBoxLayout(),
				layout.WithGrow(box("A  ×1", cBlue, 0, 30), 1),
				layout.WithGrow(box("B  ×3", cTeal, 0, 30), 3),
				layout.WithGrow(box("C  ×1", cRed, 0, 30), 1),
			),
		),
		widget.NewSeparator(),

		head("Interactive — drag the slider to change B's grow factor live:"),
		container.NewHBox(bGrowLabel, slider),
		flexRow,
	)
}

// ── Tab 2: Justify Content ────────────────────────────────────────────────────

func makeJustifyTab() fyne.CanvasObject {
	justifyRow := func(label string, j layout.Justify) fyne.CanvasObject {
		row := container.New(
			layout.NewFlexHBoxLayout(layout.WithJustify(j)),
			box("1", cBlue, 55, 40),
			box("2", cGreen, 55, 40),
			box("3", cRed, 55, 40),
		)
		lbl := widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
		return container.NewVBox(lbl, row)
	}

	return container.NewVBox(
		note("Old HBox only supports JustifyStart. To fake other modes you needed layout.NewSpacer() hacks — verbose and fragile.\n"+
			"FlexHBox + WithJustify supports all 6 CSS justify-content modes."),
		widget.NewSeparator(),

		head("Old spacer hacks — cluttered and inflexible:"),
		widget.NewLabelWithStyle(
			"Center:        NewHBox( Spacer, 1, 2, 3, Spacer )",
			fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		container.NewHBox(
			layout.NewSpacer(),
			box("1", cBlue, 55, 40), box("2", cGreen, 55, 40), box("3", cRed, 55, 40),
			layout.NewSpacer(),
		),
		widget.NewLabelWithStyle(
			"SpaceBetween: NewHBox( 1, Spacer, 2, Spacer, 3 )",
			fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		container.NewHBox(
			box("1", cBlue, 55, 40), layout.NewSpacer(),
			box("2", cGreen, 55, 40), layout.NewSpacer(),
			box("3", cRed, 55, 40),
		),
		widget.NewLabelWithStyle(
			"SpaceAround / SpaceEvenly: ✗ not achievable with Spacers",
			fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		widget.NewSeparator(),

		head("New  FlexHBox  —  one option, any mode, no Spacer clutter:"),
		justifyRow("JustifyStart  (same as old HBox)", layout.JustifyStart),
		justifyRow("JustifyEnd  —  items pushed to the right", layout.JustifyEnd),
		justifyRow("JustifyCenter  —  centered as a group", layout.JustifyCenter),
		justifyRow("JustifySpaceBetween  —  equal gaps between, none at edges", layout.JustifySpaceBetween),
		justifyRow("JustifySpaceAround  —  half-gap at edges, full gap between", layout.JustifySpaceAround),
		justifyRow("JustifySpaceEvenly  —  equal gap everywhere including edges", layout.JustifySpaceEvenly),
	)
}

// ── Tab 3: Cross-axis Alignment ───────────────────────────────────────────────

func makeAlignTab() fyne.CanvasObject {
	alignRow := func(label string, opt layout.FlexOption) fyne.CanvasObject {
		row := withMinH(110,
			container.New(
				layout.NewFlexHBoxLayout(opt),
				box("60", cBlue, 60, 60),
				box("35", cGreen, 60, 35),
				box("85", cRed, 60, 85),
				box("50", cPurple, 60, 50),
			),
		)
		lbl := widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
		return container.NewVBox(lbl, row)
	}

	return container.NewVBox(
		note("Old HBox always stretches all children to the container height — no choice.\n"+
			"To get any other alignment you had to wrap every child in container.NewCenter(), hardcode positions, or write a custom layout.\n"+
			"Numbers on boxes = MinSize height."),
		widget.NewSeparator(),

		head("Old: faking vertical center — wrap every child in container.NewCenter():"),
		note("Verbose, and only works for AlignCenter. AlignStart and AlignEnd were impossible."),
		withMinH(110,
			container.NewHBox(
				container.NewCenter(box("60", cBlue, 60, 60)),
				container.NewCenter(box("35", cGreen, 60, 35)),
				container.NewCenter(box("85", cRed, 60, 85)),
				container.NewCenter(box("50", cPurple, 60, 50)),
			),
		),
		widget.NewSeparator(),

		head("New  FlexHBox  —  one option, any alignment mode:"),
		alignRow("AlignStretch  (default — same as old HBox)", layout.WithAlignItems(layout.AlignStretch)),
		alignRow("AlignStart  —  items top-aligned at their own height  ← impossible before", layout.WithAlignItems(layout.AlignStart)),
		alignRow("AlignCenter  —  items vertically centered", layout.WithAlignItems(layout.AlignCenter)),
		alignRow("AlignEnd  —  items bottom-aligned  ← impossible before", layout.WithAlignItems(layout.AlignEnd)),
		widget.NewSeparator(),

		note("WithAlign per child — container is AlignStart, the tallest box overridden to AlignEnd:"),
		withMinH(110,
			container.New(
				layout.NewFlexHBoxLayout(layout.WithAlignItems(layout.AlignStart)),
				box("60", cBlue, 60, 60),
				layout.WithAlign(box("85 end", cRed, 60, 85), layout.AlignEnd),
				box("50", cPurple, 60, 50),
			),
		),
	)
}

// ── Tab 4: Weighted Grid ──────────────────────────────────────────────────────

func makeWeightedGridTab() fyne.CanvasObject {
	return container.NewVBox(
		note("Old GridLayout always makes every column the same width.\n"+
			"WeightedGridLayout assigns proportional widths via weights."),
		widget.NewSeparator(),

		head("Old  NewGridWithColumns(3)  —  three equal columns:"),
		container.New(
			layout.NewGridLayoutWithColumns(3),
			box("Nav", cIndigo, 0, 80),
			box("Content", cTeal, 0, 80),
			box("Panel", cOrange, 0, 80),
		),
		widget.NewSeparator(),

		head("New  WeightedGridLayout(1, 3, 1)  —  sidebar : main : sidebar:"),
		container.New(
			layout.NewWeightedGridLayout(1, 3, 1),
			box("Nav", cIndigo, 0, 80),
			box("Content (3×)", cTeal, 0, 80),
			box("Panel", cOrange, 0, 80),
		),
		widget.NewSeparator(),

		head("8-column weighted dashboard — edges 2×, inner columns 1×:"),
		container.New(
			layout.NewWeightedGridLayout(2, 1, 1, 1, 1, 1, 1, 2),
			box("A", cBlue, 0, 60), box("B", cGreen, 0, 60), box("C", cYellow, 0, 60),
			box("D", cRed, 0, 60), box("E", cPurple, 0, 60), box("F", cTeal, 0, 60),
			box("G", cOrange, 0, 60), box("H", cIndigo, 0, 60),
		),
		widget.NewSeparator(),

		head("WeightedGridLayout(1, 2, 1) with ColSpan — first item spans 2 cols:"),
		container.New(
			layout.NewWeightedGridLayout(1, 2, 1),
			layout.WithSpan(box("Span 2 cols", cPurple, 0, 50), 2, 1),
			box("Single", cOrange, 0, 50),
			box("A", cBlue, 0, 50), box("B", cGreen, 0, 50), box("C", cRed, 0, 50),
		),
	)
}

// ── Tab 5: FlexGrid (CSS Grid) ────────────────────────────────────────────────

func makeFlexGridTab() fyne.CanvasObject {
	return container.NewVBox(
		note("FlexGridLayout brings CSS Grid to Fyne: fr fractional units, fixed tracks,\n"+
			"explicit row sizes, ColSpan and RowSpan."),
		widget.NewSeparator(),

		head("Fixed(120) + Fr(1) + Fr(2)  with a header spanning all 3 columns:"),
		container.New(
			layout.NewFlexGridLayout([]layout.TrackSize{
				layout.Fixed(120), layout.Fr(1), layout.Fr(2),
			}),
			layout.WithSpan(box("Header — spans all 3 columns", cPurple, 0, 36), 3, 1),
			box("Sidebar 120dp", cIndigo, 120, 60),
			box("Col B  1fr", cTeal, 0, 60),
			box("Col C  2fr", cOrange, 0, 60),
			box("Sidebar", cIndigo, 120, 60),
			box("Footer B", cTeal, 0, 60),
			box("Footer C", cOrange, 0, 60),
		),
		widget.NewSeparator(),

		head("3 × Fr(1)  —  first cell spans 2 rows:"),
		container.New(
			layout.NewFlexGridLayout([]layout.TrackSize{layout.Fr(1), layout.Fr(1), layout.Fr(1)}),
			layout.WithSpan(box("RowSpan 2", cRed, 0, 0), 1, 2),
			box("B", cGreen, 0, 50), box("C", cBlue, 0, 50),
			box("D", cYellow, 0, 50), box("E", cPurple, 0, 50),
		),
		widget.NewSeparator(),

		head("Explicit rows: 1 column, rows = Fixed(40) + Fr(1) + Fixed(40):"),
		withMinH(160,
			container.New(
				layout.NewFlexGridLayout(
					[]layout.TrackSize{layout.Fr(1)},
					layout.WithFlexRows([]layout.TrackSize{
						layout.Fixed(40), layout.Fr(1), layout.Fixed(40),
					}),
				),
				box("Header  40dp", cPurple, 0, 0),
				box("Content  1fr", cTeal, 0, 0),
				box("Footer  40dp", cIndigo, 0, 0),
			),
		),
	)
}

// ── Tab 6: Basis & MaxSize ────────────────────────────────────────────────────

func makeBasisMaxSizeTab() fyne.CanvasObject {
	return container.NewVBox(
		note("Basis sets the initial main-axis size before grow/shrink is applied.\n"+
			"MaxWidth / MaxHeight cap how far an item can grow — essential for keeping UI sane on wide screens."),
		widget.NewSeparator(),

		head("Basis only (no grow) — items start at their basis, ignoring MinSize:"),
		note("All four boxes have the same MinSize(60×50) but different Basis values."),
		container.New(
			layout.NewFlexHBoxLayout(),
			layout.WithHint(box("B:80", cBlue, 60, 50), layout.LayoutHint{Basis: 80}),
			layout.WithHint(box("B:160", cGreen, 60, 50), layout.LayoutHint{Basis: 160}),
			layout.WithHint(box("B:120", cRed, 60, 50), layout.LayoutHint{Basis: 120}),
			layout.WithHint(box("B:200", cPurple, 60, 50), layout.LayoutHint{Basis: 200}),
		),
		widget.NewSeparator(),

		head("Basis + Grow — start at basis, then share remaining space:"),
		note("Blue: Basis=80 Grow=1 | Green: Basis=160 Grow=1 | Red: Basis=80 Grow=2\n"+
			"Red gets twice as much of the leftover space despite the same starting basis."),
		container.New(
			layout.NewFlexHBoxLayout(),
			layout.WithHint(box("B:80 G:1", cBlue, 60, 50), layout.LayoutHint{Basis: 80, Grow: 1}),
			layout.WithHint(box("B:160 G:1", cGreen, 60, 50), layout.LayoutHint{Basis: 160, Grow: 1}),
			layout.WithHint(box("B:80 G:2", cRed, 60, 50), layout.LayoutHint{Basis: 80, Grow: 2}),
		),
		widget.NewSeparator(),

		head("MaxWidth — grow=1 on all, but Blue caps at 150dp, Green at 240dp, Red is free:"),
		note("Useful for keeping buttons or inputs from stretching absurdly wide on large screens."),
		container.New(
			layout.NewFlexHBoxLayout(),
			layout.WithHint(box("max 150", cBlue, 60, 50), layout.LayoutHint{Grow: 1, MaxWidth: 150}),
			layout.WithHint(box("max 240", cGreen, 60, 50), layout.LayoutHint{Grow: 1, MaxWidth: 240}),
			layout.WithHint(box("no max", cRed, 60, 50), layout.LayoutHint{Grow: 1}),
		),
		widget.NewSeparator(),

		head("Practical MaxWidth: form action buttons — grow but stay button-sized:"),
		note("Each button has grow=1 so they fill the row proportionally, but are capped at 140dp."),
		container.New(
			layout.NewFlexHBoxLayout(layout.WithJustify(layout.JustifyEnd)),
			layout.WithHint(widget.NewButton("Cancel", nil), layout.LayoutHint{Grow: 1, MaxWidth: 140}),
			layout.WithHint(widget.NewButton("Save Draft", nil), layout.LayoutHint{Grow: 1, MaxWidth: 140}),
			layout.WithHint(widget.NewButton("Publish", nil), layout.LayoutHint{Grow: 1, MaxWidth: 140}),
		),
		widget.NewSeparator(),

		head("MaxHeight — FlexVBox where middle item is capped at 50dp:"),
		withMinH(180,
			container.New(
				layout.NewFlexVBoxLayout(),
				layout.WithGrow(box("Grow free", cBlue, 0, 20), 1),
				layout.WithHint(box("Max 50dp", cRed, 0, 20), layout.LayoutHint{Grow: 1, MaxHeight: 50}),
				layout.WithGrow(box("Grow free", cGreen, 0, 20), 1),
			),
		),
	)
}

// ── Tab 7: Margins ────────────────────────────────────────────────────────────

func makeMarginTab() fyne.CanvasObject {
	return container.NewVBox(
		note("Old layouts have no per-child margin — you must nest Padded containers.\n"+
			"WithMargin adds asymmetric outer spacing to any individual child without a wrapper."),
		widget.NewSeparator(),

		head("Old: must wrap each spaced item in container.NewPadded():"),
		container.NewHBox(
			container.NewPadded(box("A", cBlue, 60, 50)),
			box("B", cGreen, 60, 50),
			container.NewPadded(box("C", cRed, 60, 50)),
			container.NewPadded(box("D", cPurple, 60, 50)),
		),
		widget.NewSeparator(),

		head("New: WithMargin — asymmetric spacing per child, no wrapper needed:"),
		note("A: uniform 12  |  B: none  |  C: top 20, left 8  |  D: right 20, bottom 4"),
		container.New(
			layout.NewFlexHBoxLayout(),
			layout.WithMargin(box("A", cBlue, 60, 50), layout.NewUniformInsets(12)),
			box("B", cGreen, 60, 50),
			layout.WithMargin(box("C", cRed, 60, 50), layout.Insets{Top: 20, Left: 8}),
			layout.WithMargin(box("D", cPurple, 60, 50), layout.Insets{Right: 20, Bottom: 4}),
		),
		widget.NewSeparator(),

		head("Margins in a FlexGrid — cells with varying inner breathing room:"),
		container.New(
			layout.NewFlexGridLayout([]layout.TrackSize{layout.Fr(1), layout.Fr(1), layout.Fr(1)}),
			layout.WithMargin(box("A m:8", cBlue, 0, 50), layout.NewUniformInsets(8)),
			layout.WithMargin(box("B m:4", cGreen, 0, 50), layout.NewUniformInsets(4)),
			layout.WithMargin(box("C m:16", cRed, 0, 50), layout.NewUniformInsets(16)),
			layout.WithMargin(box("D m:8", cPurple, 0, 50), layout.NewUniformInsets(8)),
			layout.WithMargin(box("E m:0", cTeal, 0, 50), layout.NewUniformInsets(0)),
			layout.WithMargin(box("F m:12", cOrange, 0, 50), layout.NewUniformInsets(12)),
		),
	)
}

// ── Tab 8: Reverse & Wrap ─────────────────────────────────────────────────────

func makeReverseWrapTab() fyne.CanvasObject {
	// Small set for the reverse demo (fits in any window width).
	mkSmall := func() []fyne.CanvasObject {
		cols := []color.NRGBA{cBlue, cGreen, cRed, cPurple, cTeal, cOrange, cYellow}
		lbls := []string{"A", "B", "C", "D", "E", "F", "G"}
		objs := make([]fyne.CanvasObject, len(cols))
		for i := range cols {
			objs[i] = box(lbls[i], cols[i], 65, 48)
		}
		return objs
	}

	// Larger set for wrap demos: 12 boxes of 90px → total ~1164px, wider than
	// the 880px starting window, so wrapping is visible immediately and becomes
	// more pronounced as the window narrows.
	cols12 := []color.NRGBA{cBlue, cGreen, cRed, cPurple, cTeal, cOrange, cYellow, cIndigo,
		cBlue, cGreen, cRed, cPurple}
	lbls12 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	mkWrap := func() []fyne.CanvasObject {
		objs := make([]fyne.CanvasObject, len(cols12))
		for i := range cols12 {
			objs[i] = box(lbls12[i], cols12[i], 90, 48)
		}
		return objs
	}

	return container.NewVBox(
		note("WithReverse renders items in reverse order without modifying the slice.\n"+
			"WithWrap flows items onto the next line when they overflow — resize the window to see it reflow."),
		widget.NewSeparator(),

		head("Normal order  A → G:"),
		container.New(layout.NewFlexHBoxLayout(), mkSmall()...),

		head("WithReverse(true)  —  G → A, the underlying slice is unchanged:"),
		container.New(layout.NewFlexHBoxLayout(layout.WithReverse(true)), mkSmall()...),
		widget.NewSeparator(),

		head("HBox WithWrap(true)  —  12 boxes, wraps at current window width, reflows on resize:"),
		container.New(layout.NewFlexHBoxLayout(layout.WithWrap(true)), mkWrap()...),
		widget.NewSeparator(),

		head("VBox WithWrap(true)  —  overflows to next column:"),
		withMinH(130,
			container.New(layout.NewFlexVBoxLayout(layout.WithWrap(true)), mkSmall()...),
		),
		widget.NewSeparator(),

		head("Reverse + Wrap combined  —  12 → 1 wrapping right-to-left:"),
		container.New(
			layout.NewFlexHBoxLayout(layout.WithReverse(true), layout.WithWrap(true)),
			mkWrap()...,
		),
	)
}

// ── Tab 9: Real World ─────────────────────────────────────────────────────────

func makeRealWorldTab() fyne.CanvasObject {
	// ── App shell ──
	navItems := container.NewVBox(
		widget.NewButton("Dashboard", nil),
		widget.NewButton("Projects", nil),
		widget.NewButton("Reports", nil),
		widget.NewButton("Settings", nil),
		layout.NewSpacer(),
		widget.NewButton("Logout", nil),
	)
	sidebar := widget.NewCard("Navigation", "", navItems)

	mainItems := container.NewVBox(
		widget.NewLabelWithStyle("Main Content Area", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		note("This region has Grow=1 and fills all remaining horizontal space.\n\n"+
			"Old Fyne required hardcoding sidebar widths (e.g. BorderLayout with a fixed-width left widget).\n"+
			"With the new system you declare: sidebar Basis=160, inspector Basis=140, content Grow=1.\n"+
			"The layout is fully resize-aware — try dragging the window wider or narrower."),
	)
	mainCard := widget.NewCard("Content", "", mainItems)

	inspectorItems := container.NewVBox(
		widget.NewLabel("X:      100"),
		widget.NewLabel("Y:      200"),
		widget.NewLabel("Width:  350"),
		widget.NewLabel("Height: 200"),
		widget.NewSeparator(),
		widget.NewLabel("Fill: #4285F4"),
		widget.NewLabel("Opacity: 100%"),
	)
	inspector := widget.NewCard("Inspector", "", inspectorItems)

	appShell := withMinH(230,
		container.New(
			layout.NewFlexHBoxLayout(),
			layout.WithHint(sidebar, layout.LayoutHint{Basis: 160}),
			layout.WithGrow(mainCard, 1),
			layout.WithHint(inspector, layout.LayoutHint{Basis: 150}),
		),
	)

	// ── Proportional form ──
	form := container.New(
		layout.NewWeightedGridLayout(1, 2),
		widget.NewLabel("First Name"), widget.NewEntry(),
		widget.NewLabel("Last Name"), widget.NewEntry(),
		widget.NewLabel("Email Address"), widget.NewEntry(),
		widget.NewLabel("Role"), widget.NewEntry(),
		widget.NewLabel("Department"), widget.NewEntry(),
	)

	formButtons := container.New(
		layout.NewFlexHBoxLayout(layout.WithJustify(layout.JustifyEnd)),
		layout.WithHint(widget.NewButton("Reset", nil), layout.LayoutHint{Grow: 1, MaxWidth: 120}),
		layout.WithHint(widget.NewButton("Save", nil), layout.LayoutHint{Grow: 1, MaxWidth: 120}),
	)

	// ── Toolbar ──
	toolbar := container.New(
		layout.NewFlexHBoxLayout(layout.WithAlignItems(layout.AlignCenter)),
		widget.NewButtonWithIcon("", fyne.NewStaticResource("", nil), nil), // placeholder icon btn
		layout.WithGrow(widget.NewEntry(), 1),                              // search bar grows
		layout.WithHint(widget.NewButton("Filter", nil), layout.LayoutHint{Basis: 80}),
		layout.WithHint(widget.NewButton("Export", nil), layout.LayoutHint{Basis: 80}),
	)

	return container.NewVBox(
		head("App shell — sidebar + fluid content + inspector:"),
		note("sidebar Basis=160  ·  inspector Basis=150  ·  content Grow=1. Resize the window to see it in action."),
		appShell,
		widget.NewSeparator(),

		head("Proportional form — labels take 1fr, inputs take 2fr:"),
		note("WeightedGridLayout(1, 2) aligns all labels in one column and all inputs in another. Buttons are right-aligned with MaxWidth caps."),
		widget.NewCard("Edit Profile", "", container.NewVBox(form, formButtons)),
		widget.NewSeparator(),

		head("Toolbar — search bar grows, action buttons have fixed basis:"),
		note("The search Entry has Grow=1, filter/export buttons have Basis=80. The whole row is AlignCenter so buttons stay vertically centred."),
		toolbar,
	)
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	a := app.New()
	w := a.NewWindow("Layout Hints Showcase")
	w.Resize(fyne.NewSize(880, 640))

	tabs := container.NewAppTabs(
		container.NewTabItem("Flex Grow", container.NewScroll(makeGrowTab())),
		container.NewTabItem("Justify", container.NewScroll(makeJustifyTab())),
		container.NewTabItem("Align", container.NewScroll(makeAlignTab())),
		container.NewTabItem("Weighted Grid", container.NewScroll(makeWeightedGridTab())),
		container.NewTabItem("FlexGrid", container.NewScroll(makeFlexGridTab())),
		container.NewTabItem("Basis & MaxSize", container.NewScroll(makeBasisMaxSizeTab())),
		container.NewTabItem("Margins", container.NewScroll(makeMarginTab())),
		container.NewTabItem("Reverse & Wrap", container.NewVScroll(makeReverseWrapTab())),
		container.NewTabItem("Real World", container.NewScroll(makeRealWorldTab())),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(tabs)
	w.ShowAndRun()
}
