package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// FlexOption configures a flex box layout.
//
// Since: 2.6
type FlexOption func(*flexConfig)

type flexConfig struct {
	paddingFunc func() float32
	justify     Justify
	alignItems  Alignment
	wrap        bool
	reverse     bool
}

// WithJustify sets the justify-content behaviour of a flex box layout.
//
// Since: 2.6
func WithJustify(j Justify) FlexOption {
	return func(c *flexConfig) { c.justify = j }
}

// WithAlignItems sets the default cross-axis alignment for all children.
//
// Since: 2.6
func WithAlignItems(a Alignment) FlexOption {
	return func(c *flexConfig) { c.alignItems = a }
}

// WithWrap enables flex-wrap so children that overflow the main axis flow
// into additional lines.
//
// Since: 2.6
func WithWrap(wrap bool) FlexOption {
	return func(c *flexConfig) { c.wrap = wrap }
}

// WithReverse reverses the placement order of children.
//
// Since: 2.6
func WithReverse(reverse bool) FlexOption {
	return func(c *flexConfig) { c.reverse = reverse }
}

// WithGapPadding overrides the inter-child gap with a fixed value.
//
// Since: 2.6
func WithGapPadding(gap float32) FlexOption {
	return func(c *flexConfig) { c.paddingFunc = func() float32 { return gap } }
}

// NewVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom. The objects are always displayed
// at their vertical MinSize. Use a different layout if the objects are intended
// to be larger than their vertical MinSize.
func NewVBoxLayout() fyne.Layout {
	return vBoxLayout{paddingFunc: theme.Padding}
}

// NewHBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right. The objects are always displayed
// at their horizontal MinSize. Use a different layout if the objects are intended
// to be larger than their horizontal MinSize.
func NewHBoxLayout() fyne.Layout {
	return hBoxLayout{paddingFunc: theme.Padding}
}

// NewCustomPaddedHBoxLayout returns a layout similar to HBoxLayout that uses a custom
// amount of padding in between objects instead of the theme.Padding value.
//
// Since: 2.5
func NewCustomPaddedHBoxLayout(padding float32) fyne.Layout {
	return hBoxLayout{paddingFunc: func() float32 { return padding }}
}

// NewCustomPaddedVBoxLayout returns a layout similar to VBoxLayout that uses a custom
// amount of padding in between objects instead of the theme.Padding value.
//
// Since: 2.5
func NewCustomPaddedVBoxLayout(padding float32) fyne.Layout {
	return vBoxLayout{paddingFunc: func() float32 { return padding }}
}

// NewFlexVBoxLayout returns a vertical flex box layout with the given options.
// Use WithJustify, WithAlignItems, WithWrap, WithReverse, WithGapPadding to
// configure flex behaviour. Children can carry LayoutHints via WithGrow,
// WithFlex, WithAlign, etc. for per-child control.
//
// Since: 2.6
func NewFlexVBoxLayout(opts ...FlexOption) fyne.Layout {
	cfg := &flexConfig{paddingFunc: theme.Padding, alignItems: AlignStretch}
	for _, o := range opts {
		o(cfg)
	}
	return vBoxLayout{
		paddingFunc: cfg.paddingFunc,
		justify:     cfg.justify,
		alignItems:  cfg.alignItems,
		wrap:        cfg.wrap,
		reverse:     cfg.reverse,
	}
}

// NewFlexHBoxLayout returns a horizontal flex box layout with the given options.
//
// Since: 2.6
func NewFlexHBoxLayout(opts ...FlexOption) fyne.Layout {
	cfg := &flexConfig{paddingFunc: theme.Padding, alignItems: AlignStretch}
	for _, o := range opts {
		o(cfg)
	}
	return hBoxLayout{
		paddingFunc: cfg.paddingFunc,
		justify:     cfg.justify,
		alignItems:  cfg.alignItems,
		wrap:        cfg.wrap,
		reverse:     cfg.reverse,
	}
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*vBoxLayout)(nil)

type vBoxLayout struct {
	paddingFunc func() float32
	justify     Justify
	alignItems  Alignment
	wrap        bool
	reverse     bool
}

// boxItem is an internal helper for the flex algorithm.
type boxItem struct {
	obj   fyne.CanvasObject
	hint  LayoutHint
	basis float32
	isSpc bool
}

// Layout is called to pack all child objects into a specified size.
// When no flex hints are present it behaves identically to the previous
// implementation (backward compatible).
func (v vBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if v.wrap {
		v.layoutWrap(objects, size)
		return
	}
	v.layoutLine(objects, size, 0)
}

func (v vBoxLayout) collectItems(objects []fyne.CanvasObject) []boxItem {
	src := objects
	if v.reverse {
		src = reversedObjs(objects)
	}
	items := make([]boxItem, 0, len(src))
	for _, obj := range src {
		if !obj.Visible() {
			continue
		}
		if isVerticalSpacer(obj) {
			items = append(items, boxItem{obj: obj, isSpc: true})
			continue
		}
		hint, _ := GetHint(obj)
		basis := hint.Basis
		if basis <= 0 {
			basis = obj.MinSize().Height
		}
		items = append(items, boxItem{obj: obj, hint: hint, basis: basis})
	}
	return items
}

func (v vBoxLayout) layoutLine(objects []fyne.CanvasObject, size fyne.Size, mainOffset float32) float32 {
	padding := v.paddingFunc()
	items := v.collectItems(objects)
	if len(items) == 0 {
		return 0
	}

	// Let items that know their intrinsic height (e.g. wrap containers)
	// report it given the available width, so VBox allocates the right space.
	for i, it := range items {
		if !it.isSpc {
			if is, ok := it.obj.(fyne.IntrinsicSizer); ok {
				items[i].basis = is.IntrinsicSize(size.Width).Height
			}
		}
	}

	hasGrow, hasSpacers := false, false
	nonSpcCount := 0
	for _, it := range items {
		if it.isSpc {
			hasSpacers = true
		} else {
			nonSpcCount++
			if it.hint.Grow > 0 {
				hasGrow = true
			}
		}
	}

	// Legacy spacer path – fully backward compatible.
	if !hasGrow && hasSpacers {
		return vLegacyLayout(items, size, padding, mainOffset)
	}

	// Flex path.
	totalBasis := float32(0)
	for _, it := range items {
		if !it.isSpc {
			totalBasis += it.basis
		}
	}
	freeSpace := size.Height - totalBasis - padding*float32(imax(nonSpcCount-1, 0))

	sizes := vDistribute(items, freeSpace)

	// When grow distributed the space, use JustifyStart (pack from top).
	appliedJustify := v.justify
	totalGrow := float32(0)
	for _, it := range items {
		if !it.isSpc {
			totalGrow += it.hint.Grow
		}
	}
	if totalGrow > 0 {
		appliedJustify = JustifyStart
	}

	j := vComputeJustify(appliedJustify, totalBasis, size.Height, padding, nonSpcCount)
	y := mainOffset + j.startY

	for i, it := range items {
		if it.isSpc {
			it.obj.Move(fyne.NewPos(0, y))
			it.obj.Resize(fyne.NewSize(0, 0))
			continue
		}
		x, w := vCrossAlign(v.alignItems, it, size.Width)
		h := sizes[i]
		it.obj.Move(fyne.NewPos(x, y))
		it.obj.Resize(fyne.NewSize(w, h))
		y += h + j.itemGap
	}
	return y - mainOffset
}

// vJustifyResult holds the starting offset and item gap for a line.
type vJustifyResult struct {
	startY  float32
	itemGap float32 // full inter-item gap (replaces padding for distribute modes)
}

func vComputeJustify(justify Justify, totalBasis, containerH, padding float32, nonSpcCount int) vJustifyResult {
	// For distribute modes (SpaceBetween/Around/Evenly) the justify gap fully
	// replaces the default padding so that spacing is CSS-correct.
	switch justify {
	case JustifyEnd:
		freeSpace := containerH - totalBasis - padding*float32(imax(nonSpcCount-1, 0))
		return vJustifyResult{startY: freeSpace, itemGap: padding}
	case JustifyCenter:
		freeSpace := containerH - totalBasis - padding*float32(imax(nonSpcCount-1, 0))
		return vJustifyResult{startY: freeSpace / 2, itemGap: padding}
	case JustifySpaceBetween:
		total := containerH - totalBasis
		gap := float32(0)
		if nonSpcCount > 1 {
			gap = total / float32(nonSpcCount-1)
		}
		return vJustifyResult{startY: 0, itemGap: gap}
	case JustifySpaceAround:
		total := containerH - totalBasis
		gap := total / float32(nonSpcCount)
		return vJustifyResult{startY: gap / 2, itemGap: gap}
	case JustifySpaceEvenly:
		total := containerH - totalBasis
		gap := total / float32(nonSpcCount+1)
		return vJustifyResult{startY: gap, itemGap: gap}
	default: // JustifyStart
		return vJustifyResult{startY: 0, itemGap: padding}
	}
}

func vDistribute(items []boxItem, freeSpace float32) []float32 {
	sizes := make([]float32, len(items))
	if freeSpace > 0 {
		totalGrow := float32(0)
		for _, it := range items {
			if !it.isSpc {
				totalGrow += it.hint.Grow
			}
		}
		if totalGrow > 0 {
			for i, it := range items {
				if !it.isSpc {
					sizes[i] = it.basis + freeSpace*it.hint.Grow/totalGrow
				}
			}
			return sizes
		}
	} else if freeSpace < 0 {
		totalShrinkScaled := float32(0)
		for _, it := range items {
			if !it.isSpc {
				totalShrinkScaled += it.hint.Shrink * it.basis
			}
		}
		if totalShrinkScaled > 0 {
			for i, it := range items {
				if !it.isSpc {
					shrinkAmt := (-freeSpace) * it.hint.Shrink * it.basis / totalShrinkScaled
					sizes[i] = fyne.Max(it.obj.MinSize().Height, it.basis-shrinkAmt)
				}
			}
			return sizes
		}
	}
	for i, it := range items {
		if !it.isSpc {
			sizes[i] = it.basis
		}
	}
	return sizes
}


func vCrossAlign(containerAlign Alignment, it boxItem, containerWidth float32) (x, w float32) {
	align := it.hint.Align
	if align == AlignDefault {
		align = containerAlign
		if align == AlignDefault {
			align = AlignStretch
		}
	}
	childMinW := it.obj.MinSize().Width
	switch align {
	case AlignStart:
		return 0, childMinW
	case AlignCenter:
		w = childMinW
		return (containerWidth - w) / 2, w
	case AlignEnd:
		w = childMinW
		return containerWidth - w, w
	default: // AlignStretch
		return 0, containerWidth
	}
}

func vLegacyLayout(items []boxItem, size fyne.Size, padding, mainOffset float32) float32 {
	spacers, visibleObjects := 0, 0
	total := float32(0)
	for _, it := range items {
		if it.isSpc {
			spacers++
		} else {
			visibleObjects++
			total += it.obj.MinSize().Height
		}
	}
	extra := size.Height - total - padding*float32(imax(visibleObjects-1, 0))
	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}
	y := mainOffset
	for _, it := range items {
		if it.isSpc {
			it.obj.Move(fyne.NewPos(0, y))
			it.obj.Resize(fyne.NewSize(size.Width, spacerSize))
			y += spacerSize
			continue
		}
		it.obj.Move(fyne.NewPos(0, y))
		height := it.obj.MinSize().Height
		y += padding + height
		it.obj.Resize(fyne.NewSize(size.Width, height))
	}
	return y - mainOffset
}

func (v vBoxLayout) layoutWrap(objects []fyne.CanvasObject, size fyne.Size) {
	padding := v.paddingFunc()
	src := objects
	if v.reverse {
		src = reversedObjs(objects)
	}

	// Group visible non-spacer objects into columns that fit within size.Height.
	type col struct{ objs []fyne.CanvasObject }
	var cols []col
	var current col
	currentH, first := float32(0), true

	for _, obj := range src {
		if !obj.Visible() || isVerticalSpacer(obj) {
			continue
		}
		h := obj.MinSize().Height
		if !first && currentH+padding+h > size.Height {
			cols = append(cols, current)
			current = col{}
			currentH = 0
			first = true
		}
		current.objs = append(current.objs, obj)
		if first {
			currentH = h
			first = false
		} else {
			currentH += padding + h
		}
	}
	if len(current.objs) > 0 {
		cols = append(cols, current)
	}

	x := float32(0)
	for _, col := range cols {
		colW := float32(0)
		for _, obj := range col.objs {
			if w := obj.MinSize().Width; w > colW {
				colW = w
			}
		}
		subLayout := vBoxLayout{paddingFunc: v.paddingFunc, justify: v.justify, alignItems: v.alignItems}
		subLayout.layoutLine(col.objs, fyne.NewSize(colW, size.Height), 0)
		for _, obj := range col.objs {
			pos := obj.Position()
			obj.Move(fyne.NewPos(x+pos.X, pos.Y))
		}
		x += colW + padding
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
func (v vBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	if v.wrap {
		// When wrapping, the minimum height is just the tallest single item
		// (everything can stack into one row). Width is sum of all items.
		padding := v.paddingFunc()
		addPadding := false
		for _, child := range objects {
			if !child.Visible() || isVerticalSpacer(child) {
				continue
			}
			childMin := child.MinSize()
			minSize.Height = fyne.Max(childMin.Height, minSize.Height)
			minSize.Width += childMin.Width
			if addPadding {
				minSize.Width += padding
			}
			addPadding = true
		}
		return minSize
	}
	addPadding := false
	padding := v.paddingFunc()
	for _, child := range objects {
		if !child.Visible() || isVerticalSpacer(child) {
			continue
		}
		childMin := child.MinSize()
		minSize.Width = fyne.Max(childMin.Width, minSize.Width)
		minSize.Height += childMin.Height
		if addPadding {
			minSize.Height += padding
		}
		addPadding = true
	}
	return minSize
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*hBoxLayout)(nil)

type hBoxLayout struct {
	paddingFunc func() float32
	justify     Justify
	alignItems  Alignment
	wrap        bool
	reverse     bool
}

// Layout is called to pack all child objects into a specified size.
// When no flex hints are present it behaves identically to the previous
// implementation (backward compatible).
func (g hBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if g.wrap {
		g.layoutWrap(objects, size)
		return
	}
	g.layoutLine(objects, size, 0)
}

func (g hBoxLayout) collectItems(objects []fyne.CanvasObject) []boxItem {
	src := objects
	if g.reverse {
		src = reversedObjs(objects)
	}
	items := make([]boxItem, 0, len(src))
	for _, obj := range src {
		if !obj.Visible() {
			continue
		}
		if isHorizontalSpacer(obj) {
			items = append(items, boxItem{obj: obj, isSpc: true})
			continue
		}
		hint, _ := GetHint(obj)
		basis := hint.Basis
		if basis <= 0 {
			basis = obj.MinSize().Width
		}
		items = append(items, boxItem{obj: obj, hint: hint, basis: basis})
	}
	return items
}

func (g hBoxLayout) layoutLine(objects []fyne.CanvasObject, size fyne.Size, mainOffset float32) float32 {
	padding := g.paddingFunc()
	items := g.collectItems(objects)
	if len(items) == 0 {
		return 0
	}

	hasGrow, hasSpacers := false, false
	nonSpcCount := 0
	for _, it := range items {
		if it.isSpc {
			hasSpacers = true
		} else {
			nonSpcCount++
			if it.hint.Grow > 0 {
				hasGrow = true
			}
		}
	}

	if !hasGrow && hasSpacers {
		return hLegacyLayout(items, size, padding, mainOffset)
	}

	totalBasis := float32(0)
	for _, it := range items {
		if !it.isSpc {
			totalBasis += it.basis
		}
	}
	freeSpace := size.Width - totalBasis - padding*float32(imax(nonSpcCount-1, 0))

	sizes := hDistribute(items, freeSpace)

	appliedJustifyH := g.justify
	totalGrowH := float32(0)
	for _, it := range items {
		if !it.isSpc {
			totalGrowH += it.hint.Grow
		}
	}
	if totalGrowH > 0 {
		appliedJustifyH = JustifyStart
	}

	j := hComputeJustify(appliedJustifyH, totalBasis, size.Width, padding, nonSpcCount)
	x := mainOffset + j.startX

	for i, it := range items {
		if it.isSpc {
			it.obj.Move(fyne.NewPos(x, 0))
			it.obj.Resize(fyne.NewSize(0, 0))
			continue
		}
		y, h := hCrossAlign(g.alignItems, it, size.Height)
		w := sizes[i]
		it.obj.Move(fyne.NewPos(x, y))
		it.obj.Resize(fyne.NewSize(w, h))
		x += w + j.itemGap
	}
	return x - mainOffset
}

type hJustifyResult struct {
	startX  float32
	itemGap float32
}

func hComputeJustify(justify Justify, totalBasis, containerW, padding float32, nonSpcCount int) hJustifyResult {
	switch justify {
	case JustifyEnd:
		freeSpace := containerW - totalBasis - padding*float32(imax(nonSpcCount-1, 0))
		return hJustifyResult{startX: freeSpace, itemGap: padding}
	case JustifyCenter:
		freeSpace := containerW - totalBasis - padding*float32(imax(nonSpcCount-1, 0))
		return hJustifyResult{startX: freeSpace / 2, itemGap: padding}
	case JustifySpaceBetween:
		total := containerW - totalBasis
		gap := float32(0)
		if nonSpcCount > 1 {
			gap = total / float32(nonSpcCount-1)
		}
		return hJustifyResult{startX: 0, itemGap: gap}
	case JustifySpaceAround:
		total := containerW - totalBasis
		gap := total / float32(nonSpcCount)
		return hJustifyResult{startX: gap / 2, itemGap: gap}
	case JustifySpaceEvenly:
		total := containerW - totalBasis
		gap := total / float32(nonSpcCount+1)
		return hJustifyResult{startX: gap, itemGap: gap}
	default:
		return hJustifyResult{startX: 0, itemGap: padding}
	}
}

func hDistribute(items []boxItem, freeSpace float32) []float32 {
	sizes := make([]float32, len(items))
	if freeSpace > 0 {
		totalGrow := float32(0)
		for _, it := range items {
			if !it.isSpc {
				totalGrow += it.hint.Grow
			}
		}
		if totalGrow > 0 {
			for i, it := range items {
				if !it.isSpc {
					sizes[i] = it.basis + freeSpace*it.hint.Grow/totalGrow
				}
			}
			return sizes
		}
	} else if freeSpace < 0 {
		totalShrinkScaled := float32(0)
		for _, it := range items {
			if !it.isSpc {
				totalShrinkScaled += it.hint.Shrink * it.basis
			}
		}
		if totalShrinkScaled > 0 {
			for i, it := range items {
				if !it.isSpc {
					shrinkAmt := (-freeSpace) * it.hint.Shrink * it.basis / totalShrinkScaled
					sizes[i] = fyne.Max(it.obj.MinSize().Width, it.basis-shrinkAmt)
				}
			}
			return sizes
		}
	}
	for i, it := range items {
		if !it.isSpc {
			sizes[i] = it.basis
		}
	}
	return sizes
}


func hCrossAlign(containerAlign Alignment, it boxItem, containerHeight float32) (y, h float32) {
	align := it.hint.Align
	if align == AlignDefault {
		align = containerAlign
		if align == AlignDefault {
			align = AlignStretch
		}
	}
	childMinH := it.obj.MinSize().Height
	switch align {
	case AlignStart:
		return 0, childMinH
	case AlignCenter:
		h = childMinH
		return (containerHeight - h) / 2, h
	case AlignEnd:
		h = childMinH
		return containerHeight - h, h
	default: // AlignStretch
		return 0, containerHeight
	}
}

func hLegacyLayout(items []boxItem, size fyne.Size, padding, mainOffset float32) float32 {
	spacers, visibleObjects := 0, 0
	total := float32(0)
	for _, it := range items {
		if it.isSpc {
			spacers++
		} else {
			visibleObjects++
			total += it.obj.MinSize().Width
		}
	}
	extra := size.Width - total - padding*float32(imax(visibleObjects-1, 0))
	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}
	x := mainOffset
	for _, it := range items {
		if it.isSpc {
			it.obj.Move(fyne.NewPos(x, 0))
			it.obj.Resize(fyne.NewSize(spacerSize, size.Height))
			x += spacerSize
			continue
		}
		it.obj.Move(fyne.NewPos(x, 0))
		width := it.obj.MinSize().Width
		x += padding + width
		it.obj.Resize(fyne.NewSize(width, size.Height))
	}
	return x - mainOffset
}

func (g hBoxLayout) layoutWrap(objects []fyne.CanvasObject, size fyne.Size) {
	padding := g.paddingFunc()
	src := objects
	if g.reverse {
		src = reversedObjs(objects)
	}

	type row struct{ objs []fyne.CanvasObject }
	var rows []row
	var current row
	currentW, first := float32(0), true

	for _, obj := range src {
		if !obj.Visible() || isHorizontalSpacer(obj) {
			continue
		}
		w := obj.MinSize().Width
		if !first && currentW+padding+w > size.Width {
			rows = append(rows, current)
			current = row{}
			currentW = 0
			first = true
		}
		current.objs = append(current.objs, obj)
		if first {
			currentW = w
			first = false
		} else {
			currentW += padding + w
		}
	}
	if len(current.objs) > 0 {
		rows = append(rows, current)
	}

	y := float32(0)
	for _, row := range rows {
		rowH := float32(0)
		for _, obj := range row.objs {
			if h := obj.MinSize().Height; h > rowH {
				rowH = h
			}
		}
		subLayout := hBoxLayout{paddingFunc: g.paddingFunc, justify: g.justify, alignItems: g.alignItems}
		subLayout.layoutLine(row.objs, fyne.NewSize(size.Width, rowH), 0)
		for _, obj := range row.objs {
			pos := obj.Position()
			obj.Move(fyne.NewPos(pos.X, y+pos.Y))
		}
		y += rowH + padding
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
func (g hBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	if g.wrap {
		// When wrapping, the minimum width is just the widest single item
		// (everything can stack into one column). Height is sum of all items.
		padding := g.paddingFunc()
		addPadding := false
		for _, child := range objects {
			if !child.Visible() || isHorizontalSpacer(child) {
				continue
			}
			childMin := child.MinSize()
			minSize.Width = fyne.Max(childMin.Width, minSize.Width)
			minSize.Height += childMin.Height
			if addPadding {
				minSize.Height += padding
			}
			addPadding = true
		}
		return minSize
	}
	addPadding := false
	padding := g.paddingFunc()
	for _, child := range objects {
		if !child.Visible() || isHorizontalSpacer(child) {
			continue
		}
		childMin := child.MinSize()
		minSize.Height = fyne.Max(childMin.Height, minSize.Height)
		minSize.Width += childMin.Width
		if addPadding {
			minSize.Width += padding
		}
		addPadding = true
	}
	return minSize
}

// LayoutConstrained implements [fyne.ConstrainedLayout] for hBoxLayout.
// When wrap is enabled it runs the wrap algorithm against the real available
// width and returns the height actually consumed, letting the parent
// container resize itself and trigger a re-layout if the height changed.
func (g hBoxLayout) LayoutConstrained(objects []fyne.CanvasObject, c fyne.Constraints) fyne.Size {
	if !g.wrap {
		// Non-wrapping: lay out at natural height (MinSize), not the potentially
		// unbounded MaxSize — avoids inflating parent VBox allocations to inf.
		natH := g.MinSize(objects).Height
		g.layoutLine(objects, fyne.NewSize(c.MaxSize.Width, natH), 0)
		return fyne.NewSize(c.MaxSize.Width, natH)
	}
	// First pass: layout at the available width to determine needed height.
	availW := c.MaxSize.Width
	padding := g.paddingFunc()

	// Build rows at availW.
	type row struct{ objs []fyne.CanvasObject }
	src := objects
	if g.reverse {
		src = reversedObjs(objects)
	}
	var rows []row
	var current row
	currentW, first := float32(0), true
	for _, obj := range src {
		if !obj.Visible() || isHorizontalSpacer(obj) {
			continue
		}
		w := obj.MinSize().Width
		if !first && currentW+padding+w > availW {
			rows = append(rows, current)
			current = row{}
			currentW = 0
			first = true
		}
		current.objs = append(current.objs, obj)
		if first {
			currentW = w
			first = false
		} else {
			currentW += padding + w
		}
	}
	if len(current.objs) > 0 {
		rows = append(rows, current)
	}

	// Compute total height and position rows.
	y := float32(0)
	for _, row := range rows {
		rowH := float32(0)
		for _, obj := range row.objs {
			if h := obj.MinSize().Height; h > rowH {
				rowH = h
			}
		}
		subLayout := hBoxLayout{paddingFunc: g.paddingFunc, justify: g.justify, alignItems: g.alignItems}
		subLayout.layoutLine(row.objs, fyne.NewSize(availW, rowH), 0)
		for _, obj := range row.objs {
			pos := obj.Position()
			obj.Move(fyne.NewPos(pos.X, y+pos.Y))
		}
		y += rowH + padding
	}
	if y > 0 {
		y -= padding // remove trailing gap
	}
	return fyne.NewSize(availW, y)
}

// Declare conformity with ConstrainedLayout when wrap is used.
var _ fyne.ConstrainedLayout = hBoxLayout{wrap: true}

func isVerticalSpacer(obj fyne.CanvasObject) bool {
	spacer, ok := obj.(SpacerObject)
	return ok && spacer.ExpandVertical()
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	spacer, ok := obj.(SpacerObject)
	return ok && spacer.ExpandHorizontal()
}

func reversedObjs(objects []fyne.CanvasObject) []fyne.CanvasObject {
	n := len(objects)
	out := make([]fyne.CanvasObject, n)
	for i, o := range objects {
		out[n-1-i] = o
	}
	return out
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
