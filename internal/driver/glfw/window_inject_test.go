//go:build !wasm && !test_web_driver

package glfw

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/stretchr/testify/assert"
)

func TestWindow_InjectClickTaps(t *testing.T) {
	w := createWindow("Test")
	prop := canvas.NewRectangle(color.White)
	prop.SetMinSize(fyne.NewSize(100, 100))
	o := &tappableObject{Rectangle: canvas.NewRectangle(color.White)}
	w.SetContent(container.NewStack(prop, o))

	var inj desktop.InputInjector = w.window
	inj.InjectMouseMove(fyne.NewPos(50, 60))
	inj.InjectMouseButton(desktop.MouseButtonPrimary, true, 0)
	inj.InjectMouseButton(desktop.MouseButtonPrimary, false, 0)

	runOnMain(func() {
		if e, _ := o.popTapEvent().(*fyne.PointEvent); assert.NotNil(t, e, "tapped") {
			assert.Equal(t, fyne.NewPos(50, 60), e.AbsolutePosition)
			assert.Equal(t, fyne.NewPos(46, 56), e.Position)
		}
	})
}

func TestWindow_InjectSecondaryClick(t *testing.T) {
	w := createWindow("Test")
	prop := canvas.NewRectangle(color.White)
	prop.SetMinSize(fyne.NewSize(100, 100))
	o := &tappableObject{Rectangle: canvas.NewRectangle(color.White)}
	w.SetContent(container.NewStack(prop, o))

	var inj desktop.InputInjector = w.window
	inj.InjectMouseMove(fyne.NewPos(50, 60))
	inj.InjectMouseButton(desktop.MouseButtonSecondary, true, 0)
	inj.InjectMouseButton(desktop.MouseButtonSecondary, false, 0)

	runOnMain(func() {
		assert.Nil(t, o.popTapEvent(), "no primary tap")
		assert.NotNil(t, o.popSecondaryTapEvent(), "secondary tap")
	})
}

func TestWindow_InjectDragAndHover(t *testing.T) {
	w := createWindow("Test")
	dh := &draggableHoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	c := container.NewWithoutLayout(dh)
	dh.Resize(fyne.NewSize(20, 20))
	w.SetContent(c)
	repaintWindow(w)

	var inj desktop.InputInjector = w.window
	inj.InjectMouseMove(fyne.NewPos(10, 10))
	runOnMain(func() {
		assert.Equal(t,
			&desktop.MouseEvent{PointEvent: fyne.PointEvent{
				Position:         fyne.NewPos(6, 6),
				AbsolutePosition: fyne.NewPos(10, 10),
			}},
			dh.popMouseInEvent(),
		)
	})

	inj.InjectMouseButton(desktop.MouseButtonPrimary, true, 0)
	inj.InjectMouseMove(fyne.NewPos(12, 12))
	runOnMain(func() {
		assert.Equal(t,
			&fyne.DragEvent{
				PointEvent: fyne.PointEvent{
					Position:         fyne.NewPos(8, 8),
					AbsolutePosition: fyne.NewPos(12, 12),
				},
				Dragged: fyne.NewDelta(2, 2),
			},
			dh.popDragEvent(),
		)
	})

	inj.InjectMouseButton(desktop.MouseButtonPrimary, false, 0)
	runOnMain(func() {
		assert.NotNil(t, dh.popDragEndEvent(), "drag ended")
	})
}

func TestWindow_InjectKeysAndShortcut(t *testing.T) {
	w := createWindow("Test")
	w.SetContent(canvas.NewRectangle(color.White))

	var typedKeys []fyne.KeyName
	var typedRunes []rune
	shortcutHits := 0
	c := w.Canvas() // safeWindow hops to main itself, so not inside runOnMain
	c.SetOnTypedKey(func(ev *fyne.KeyEvent) { typedKeys = append(typedKeys, ev.Name) })
	c.SetOnTypedRune(func(r rune) { typedRunes = append(typedRunes, r) })
	c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: fyne.KeyModifierControl},
		func(fyne.Shortcut) { shortcutHits++ })

	var inj desktop.InputInjector = w.window
	inj.InjectKey(fyne.KeySpace, true, 0)
	inj.InjectKey(fyne.KeySpace, false, 0)
	inj.InjectChar('x')

	inj.InjectKey(desktop.KeyControlLeft, true, fyne.KeyModifierControl)
	inj.InjectKey(fyne.KeyD, true, fyne.KeyModifierControl)
	runOnMain(func() {
		assert.Equal(t, fyne.KeyModifierControl, d.CurrentKeyModifiers())
	})
	inj.InjectKey(fyne.KeyD, false, fyne.KeyModifierControl)
	inj.InjectKey(desktop.KeyControlLeft, false, 0)

	runOnMain(func() {
		// A bare modifier press is a typed key too, exactly as with a real keyboard;
		// the Ctrl+D press resolves to the shortcut instead.
		assert.Equal(t, []fyne.KeyName{fyne.KeySpace, desktop.KeyControlLeft}, typedKeys)
		assert.Equal(t, []rune{'x'}, typedRunes)
		assert.Equal(t, 1, shortcutHits)
		assert.Equal(t, fyne.KeyModifier(0), d.CurrentKeyModifiers())
	})
}

func TestWindow_ObjectAtFollowsOverlays(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	o := &tappableObject{Rectangle: canvas.NewRectangle(color.White)}
	o.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(o)
	repaintWindow(w)

	var inj desktop.InputInjector = w.window
	got, rel := inj.ObjectAt(fyne.NewPos(30, 40))
	assert.Equal(t, fyne.CanvasObject(o), got)
	assert.Equal(t, fyne.NewPos(30, 40), rel)

	// An open overlay captures the pointer, exactly as for real input.
	ov := &tappableObject{Rectangle: canvas.NewRectangle(color.Black)}
	ov.Resize(fyne.NewSize(10, 10))
	c := w.Canvas()
	c.Overlays().Add(ov)
	got, _ = inj.ObjectAt(fyne.NewPos(5, 5))
	assert.Equal(t, fyne.CanvasObject(ov), got)
	got, _ = inj.ObjectAt(fyne.NewPos(30, 40))
	assert.Nil(t, got, "content is not reachable under an overlay")
	c.Overlays().Remove(ov)
}

func TestWindow_WalkVisibleObjects(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	a := canvas.NewRectangle(color.White)
	b := canvas.NewRectangle(color.White)
	hidden := canvas.NewRectangle(color.White)
	hidden.Hide()
	content := container.NewWithoutLayout(a, b, hidden)
	b.Move(fyne.NewPos(20, 30))
	w.SetContent(content)
	ov := canvas.NewRectangle(color.Black)
	c := w.Canvas()
	c.Overlays().Add(ov)
	defer c.Overlays().Remove(ov)

	var inj desktop.InputInjector = w.window
	seen := map[fyne.CanvasObject]fyne.Position{}
	layers := map[fyne.CanvasObject]int{}
	inj.WalkVisibleObjects(func(o fyne.CanvasObject, p fyne.Position, layer int) bool {
		seen[o] = p
		layers[o] = layer
		return false
	})
	assert.Contains(t, seen, fyne.CanvasObject(a))
	assert.Equal(t, fyne.NewPos(20, 30), seen[b])
	assert.Contains(t, seen, fyne.CanvasObject(ov))
	assert.NotContains(t, seen, fyne.CanvasObject(hidden))
	assert.Equal(t, 0, layers[a])
	assert.Equal(t, 1, layers[ov])

	count := 0
	inj.WalkVisibleObjects(func(fyne.CanvasObject, fyne.Position, int) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count, "returning true stops the walk")
}
