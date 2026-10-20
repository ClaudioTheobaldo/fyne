package desktop

import "fyne.io/fyne/v2"

// InputInjector is implemented by desktop windows that can receive synthetic
// input. Injected events enter the driver at the same point real OS events do,
// so everything downstream runs unchanged: hit-testing against overlays, menus
// and content, hover in/out, drag thresholds, tap vs double-tap timing, focus
// and shortcut resolution. Only the OS-to-driver layer is skipped.
//
// Intended for in-process test automation. Methods may be called from any
// goroutine; they block until the event has been processed on the main thread.
// Calling them from inside a callback that already runs on the main thread is
// also safe.
//
// Obtain one with a type assertion on a fyne.Window:
//
//	if inj, ok := w.(desktop.InputInjector); ok { ... }
type InputInjector interface {
	// InjectMouseMove moves the pointer to pos, in canvas coordinates.
	InjectMouseMove(pos fyne.Position)

	// InjectMouseButton presses or releases button at the current pointer
	// position. Move the pointer first: the driver falls back to reading the OS
	// cursor when the pointer has never been positioned.
	InjectMouseButton(button MouseButton, pressed bool, modifiers fyne.KeyModifier)

	// InjectScroll scrolls by dx, dy wheel steps at the current pointer position.
	InjectScroll(dx, dy float32)

	// InjectKey presses or releases the named key with the given modifiers held.
	InjectKey(key fyne.KeyName, pressed bool, modifiers fyne.KeyModifier)

	// InjectChar delivers one typed character, as a keyboard would after a key press.
	InjectChar(char rune)

	// ObjectAt reports what a pointer event at pos would reach, using the
	// driver's own hit-testing: the deepest visible object under pos that takes
	// pointer input (tap, secondary or double tap, hover, drag, mouse, scroll or
	// focus), and pos relative to it. While an overlay is open only the top
	// overlay is searched, as for real input. Returns nil when nothing matches.
	ObjectAt(pos fyne.Position) (fyne.CanvasObject, fyne.Position)

	// WalkVisibleObjects visits every visible object in the window: the main
	// menu and content (layer 0), then each overlay from bottom to top (layers
	// 1..n). pos is the object's absolute canvas position. Returning true from
	// fn stops the walk.
	WalkVisibleObjects(fn func(obj fyne.CanvasObject, pos fyne.Position, layer int) (stop bool))
}
