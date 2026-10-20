//go:build !wasm && !test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/scale"
)

var _ desktop.InputInjector = (*window)(nil)

// injectOnMain runs f on the main thread and waits for it. Events from GLFW are
// processed there, and the process* handlers rely on it.
func injectOnMain(f func()) {
	if async.IsMainGoroutine() {
		f()
		return
	}
	runOnMain(f)
}

func (w *window) InjectMouseMove(pos fyne.Position) {
	injectOnMain(func() {
		if w.canvas == nil {
			return
		}
		// processMouseMoved takes screen pixels, exactly as GLFW reports them.
		x := scale.ToScreenCoordinate(w.canvas, pos.X)
		y := scale.ToScreenCoordinate(w.canvas, pos.Y)
		w.processMouseMoved(float64(x), float64(y))
	})
}

func (w *window) InjectMouseButton(button desktop.MouseButton, pressed bool, modifiers fyne.KeyModifier) {
	act := release
	if pressed {
		act = press
	}
	injectOnMain(func() {
		if w.canvas == nil {
			return
		}
		w.processMouseClicked(button, act, modifiers)
	})
}

func (w *window) InjectScroll(dx, dy float32) {
	injectOnMain(func() {
		if w.canvas == nil {
			return
		}
		w.processMouseScrolled(float64(dx), float64(dy))
	})
}

func (w *window) InjectKey(key fyne.KeyName, pressed bool, modifiers fyne.KeyModifier) {
	act := release
	if pressed {
		act = press
	}
	// Mirror keyPressed: CurrentKeyModifiers includes a modifier key while it is
	// held, and the ASCII name is only set for letters.
	current := modifiers
	if pressed {
		current |= modifierForKey(key)
	} else {
		current &^= modifierForKey(key)
	}
	var ascii fyne.KeyName
	if len(key) == 1 && key[0] >= 'A' && key[0] <= 'Z' {
		ascii = key
	}
	injectOnMain(func() {
		if w.canvas == nil {
			return
		}
		w.driver.currentKeyModifiers = current
		w.processKeyPressed(key, ascii, 0, act, modifiers)
	})
}

func (w *window) InjectChar(char rune) {
	injectOnMain(func() {
		if w.canvas == nil {
			return
		}
		w.processCharInput(char)
	})
}

func (w *window) ObjectAt(pos fyne.Position) (obj fyne.CanvasObject, rel fyne.Position) {
	injectOnMain(func() {
		if w.canvas == nil {
			return
		}
		obj, rel, _ = w.findObjectAtPositionMatching(w.canvas, pos, takesPointerInput)
	})
	return obj, rel
}

func (w *window) WalkVisibleObjects(fn func(obj fyne.CanvasObject, pos fyne.Position, layer int) (stop bool)) {
	injectOnMain(func() {
		if w.canvas == nil {
			return
		}
		walk := func(root fyne.CanvasObject, layer int) bool {
			if root == nil {
				return false
			}
			return driver.WalkVisibleObjectTree(root, func(o fyne.CanvasObject, p, _ fyne.Position, _ fyne.Size) bool {
				return fn(o, p, layer)
			}, nil)
		}
		if walk(w.canvas.menu, 0) || walk(w.canvas.Content(), 0) {
			return
		}
		for i, ov := range w.canvas.Overlays().List() {
			if walk(ov, i+1) {
				return
			}
		}
	})
}

func takesPointerInput(o fyne.CanvasObject) bool {
	switch o.(type) {
	case fyne.Tappable, fyne.SecondaryTappable, fyne.DoubleTappable, fyne.Draggable,
		fyne.Scrollable, fyne.Focusable, desktop.Hoverable, desktop.Mouseable:
		return true
	}
	return false
}

func modifierForKey(key fyne.KeyName) fyne.KeyModifier {
	switch key {
	case desktop.KeyShiftLeft, desktop.KeyShiftRight:
		return fyne.KeyModifierShift
	case desktop.KeyControlLeft, desktop.KeyControlRight:
		return fyne.KeyModifierControl
	case desktop.KeyAltLeft, desktop.KeyAltRight:
		return fyne.KeyModifierAlt
	case desktop.KeySuperLeft, desktop.KeySuperRight:
		return fyne.KeyModifierSuper
	}
	return 0
}
