//go:build windows

package glfw

import (
	"sync"
	"syscall"
	"unsafe"
)

// During a Win32 modal move/size loop (title-bar drag or border resize) the OS spins its
// own message loop inside DefWindowProc, so GLFW's PollEvents is blocked and the driver's
// run-loop draw never executes — live content (e.g. playing video) freezes until the drag
// ends, and holding the window still produces no events to repaint on at all. GLFW installs
// no timer for this, so we subclass the window's WndProc and run a WM_TIMER for the duration
// of the modal loop, repainting on every tick (the canonical Windows fix). Windows-only.

const (
	wmEnterSizeMove = 0x0231
	wmExitSizeMove  = 0x0232
	wmTimer         = 0x0113
	gwlpWndProc     = ^uintptr(3) // GWLP_WNDPROC (-4)
	sizeMoveTimerID = 0xF12E
	sizeMoveTimerMs = 15 // ~60 fps repaint while dragging
)

var (
	user32sm             = syscall.NewLazyDLL("user32.dll")
	procSetWindowLongPtr = user32sm.NewProc("SetWindowLongPtrW")
	procCallWindowProc   = user32sm.NewProc("CallWindowProcW")
	procDefWindowProc    = user32sm.NewProc("DefWindowProcW")
	procSetTimer         = user32sm.NewProc("SetTimer")
	procKillTimer        = user32sm.NewProc("KillTimer")

	sizeMoveMu    sync.Mutex
	sizeMoveState = map[uintptr]*sizeMoveSub{}             // hwnd -> subclass state
	sizeMoveProc  = syscall.NewCallback(sizeMoveWndProc) // shared C-callable WndProc
)

type sizeMoveSub struct {
	orig    uintptr // GLFW's original WndProc, to chain unhandled messages
	repaint func()  // repaint this window now (main thread)
}

// sizeMoveWndProc is the subclassed window procedure. It drives a repaint timer across the
// modal move/size loop and chains every message back to GLFW's original WndProc so normal
// behaviour is untouched.
func sizeMoveWndProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	sizeMoveMu.Lock()
	sub := sizeMoveState[hwnd]
	sizeMoveMu.Unlock()
	if sub == nil { // not (or no longer) subclassed
		r, _, _ := procDefWindowProc.Call(hwnd, msg, wparam, lparam)
		return r
	}

	switch msg {
	case wmEnterSizeMove:
		procSetTimer.Call(hwnd, sizeMoveTimerID, uintptr(sizeMoveTimerMs), 0)
	case wmExitSizeMove:
		procKillTimer.Call(hwnd, sizeMoveTimerID)
	case wmTimer:
		if wparam == sizeMoveTimerID {
			sub.repaint()
		}
	}

	r, _, _ := procCallWindowProc.Call(sub.orig, hwnd, msg, wparam, lparam)
	return r
}

// installSizeMoveRepaint subclasses the GLFW window so live content keeps painting during a
// move/resize drag. No-op if the native handle isn't available.
func (w *window) installSizeMoveRepaint() {
	if w.viewport == nil {
		return
	}
	hwnd := uintptr(unsafe.Pointer(w.viewport.GetWin32Window()))
	if hwnd == 0 {
		return
	}
	orig, _, _ := procSetWindowLongPtr.Call(hwnd, gwlpWndProc, sizeMoveProc)
	if orig == 0 {
		return
	}
	sizeMoveMu.Lock()
	sizeMoveState[hwnd] = &sizeMoveSub{
		orig: orig,
		repaint: func() {
			if w.visible && !w.isClosing() {
				w.RunWithContext(func() {
					w.driver.repaintWindow(w)
				})
			}
		},
	}
	sizeMoveMu.Unlock()
}

// uninstallSizeMoveRepaint restores GLFW's WndProc. Call while the HWND is still valid
// (before the viewport is destroyed).
func (w *window) uninstallSizeMoveRepaint() {
	if w.viewport == nil {
		return
	}
	hwnd := uintptr(unsafe.Pointer(w.viewport.GetWin32Window()))
	if hwnd == 0 {
		return
	}
	sizeMoveMu.Lock()
	sub := sizeMoveState[hwnd]
	delete(sizeMoveState, hwnd)
	sizeMoveMu.Unlock()
	if sub != nil {
		procKillTimer.Call(hwnd, sizeMoveTimerID)
		procSetWindowLongPtr.Call(hwnd, gwlpWndProc, sub.orig)
	}
}
