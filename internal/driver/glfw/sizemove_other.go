//go:build !windows

package glfw

// The modal move/size repaint pump is a Windows-only concern (the Win32 modal loop blocks
// the run loop). Other platforms keep their run loop during a move/resize, so these are
// no-ops there.

func (w *window) installSizeMoveRepaint()   {}
func (w *window) uninstallSizeMoveRepaint() {}
