//go:build windows

package glfw

import (
	"sync"
	"syscall"
	"unsafe"
)

// During a Win32 modal move/size loop (title-bar drag or border resize) the OS spins its own
// message loop inside DefWindowProc, so GLFW's PollEvents is blocked and the driver's run-loop
// draw never executes — live content (e.g. playing video) freezes until the drag ends, and
// holding the window still produces no events to repaint on at all. GLFW installs no timer for
// this. We subclass the window's WndProc and, for the duration of the modal loop, run a
// WM_TIMER that on every tick:
//
//  1. drains the main-thread func queue, so queued UI work (new video frames, Refresh,
//     dirty-marking) actually applies during the drag instead of piling up behind the blocked
//     run loop — otherwise every window keeps painting the same stale frame;
//  2. relayouts the dragged window to its live size (see relayoutIfResized); and
//  3. repaints every visible window (see repaintAllDuringModalLoop), since all windows share
//     the one blocked main thread.
//
// Separately, some machines defer the actual client resize (WM_SIZE) until drag-release, so the
// content can't track the border during the drag; we force it live by applying the proposed
// rect on each WM_SIZING (see below). Every message is chained to GLFW's original WndProc so
// normal behaviour is untouched. Windows-only.

const (
	wmEnterSizeMove = 0x0231
	wmExitSizeMove  = 0x0232
	wmSizing        = 0x0214
	wmTimer         = 0x0113
	gwlpWndProc     = ^uintptr(3) // GWLP_WNDPROC (-4)
	sizeMoveTimerID = 0xF12E
	sizeMoveTimerMs = 15 // ~60 fps repaint while dragging

	swpNoZorder      = 0x0004
	swpNoActivate    = 0x0010
	swpNoOwnerZorder = 0x0200
	swpNoCopyBits    = 0x0100
)

// winRect mirrors Win32 RECT (LONGs). WM_SIZING's lParam points at the proposed window rect.
type winRect struct {
	left, top, right, bottom int32
}

var (
	user32sm             = syscall.NewLazyDLL("user32.dll")
	procSetWindowLongPtr = user32sm.NewProc("SetWindowLongPtrW")
	procCallWindowProc   = user32sm.NewProc("CallWindowProcW")
	procDefWindowProc    = user32sm.NewProc("DefWindowProcW")
	procSetTimer         = user32sm.NewProc("SetTimer")
	procKillTimer        = user32sm.NewProc("KillTimer")
	procSetWindowPos     = user32sm.NewProc("SetWindowPos")

	sizeMoveMu    sync.Mutex
	sizeMoveState = map[uintptr]*sizeMoveSub{}            // hwnd -> subclass state
	sizeMoveProc  = syscall.NewCallback(sizeMoveWndProc) // shared C-callable WndProc
)

type sizeMoveSub struct {
	orig    uintptr // GLFW's original WndProc, to chain unhandled messages
	win     *window // the subclassed window (to track its live size while dragging)
	repaint func()  // repaint all live windows now (main thread)
	lastW   int     // last window size we relaid out at, to skip no-op resizes
	lastH   int
	lastSzW int // last proposed WM_SIZING size we forced, to skip redundant SetWindowPos
	lastSzH int
}

// sizeMoveWndProc is the subclassed window procedure. It drives the modal-loop work off a
// WM_TIMER and chains every message back to GLFW's original WndProc so normal behaviour is
// untouched.
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
	case wmSizing:
		// Some machines defer the actual client resize (WM_SIZE) to drag-release, so the
		// content can't follow the border during the drag. Force it: apply the proposed window
		// rect immediately, which makes Windows commit the resize now → GLFW's normal WM_SIZE →
		// Fyne relayout fires each frame and the content tracks the border live.
		if r := (*winRect)(unsafe.Pointer(lparam)); r != nil {
			pw, ph := int(r.right-r.left), int(r.bottom-r.top)
			if pw > 0 && ph > 0 && (pw != sub.lastSzW || ph != sub.lastSzH) {
				sub.lastSzW, sub.lastSzH = pw, ph
				procSetWindowPos.Call(hwnd, 0,
					uintptr(r.left), uintptr(r.top), uintptr(pw), uintptr(ph),
					swpNoZorder|swpNoActivate|swpNoOwnerZorder|swpNoCopyBits)
			}
		}
	case wmTimer:
		if wparam == sizeMoveTimerID {
			drainMainFuncQueue()    // apply queued UI work (frame delivery, Refresh) so content stays live
			sub.relayoutIfResized() // track the border live during a resize drag
			sub.repaint()           // present every window
		}
	}

	r, _, _ := procCallWindowProc.Call(sub.orig, hwnd, msg, wparam, lparam)
	return r
}

// relayoutIfResized makes content track the window border live during a resize drag. We poll
// the live window size each timer tick and, if it changed, run the exact same processResized
// the normal size path uses (canvas resize + relayout). No-op for a pure move (size unchanged)
// and for windows other than the one being dragged. Runs on the main thread (in the modal loop).
func (s *sizeMoveSub) relayoutIfResized() {
	w := s.win
	if w == nil || w.viewport == nil || !w.visible || w.isClosing() {
		return
	}
	width, height := w.viewport.GetSize()
	if width == 0 || height == 0 || (width == s.lastW && height == s.lastH) {
		return
	}
	s.lastW, s.lastH = width, height
	w.processResized(width, height)
}

// repaintAllDuringModalLoop force-paints every visible window. The Win32 modal move/size loop
// runs on the shared main thread, so dragging ANY window blocks the run loop and stops ALL
// windows from drawing (most visibly: video in the main window freezes while a secondary window
// is dragged). We therefore repaint every window on each timer tick, not just the one being
// dragged. repaintWindow paints unconditionally (unlike drawSingleFrame, which is dirty-gated).
// Safe to read d.windows without a lock here: the run loop and window add/remove all run on this
// same main thread, which is parked inside the modal loop while this fires.
func (d *gLDriver) repaintAllDuringModalLoop() {
	for _, win := range d.windowList() {
		w, ok := win.(*window)
		if !ok || w.viewport == nil || !w.visible || w.isClosing() {
			continue
		}
		w.RunWithContext(func() {
			d.repaintWindow(w)
		})
	}
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
	width, height := w.viewport.GetSize()
	sizeMoveMu.Lock()
	sizeMoveState[hwnd] = &sizeMoveSub{
		orig:    orig,
		win:     w,
		repaint: w.driver.repaintAllDuringModalLoop,
		lastW:   width,
		lastH:   height,
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
