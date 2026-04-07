//go:build windows

package gl

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	wglAccessReadOnlyNV = 0x0000
)

// wglInterop manages WGL_NV_DX_interop2 function pointers and state.
// All methods must be called on the GL render thread.
type wglInterop struct {
	dxDevice uintptr // HANDLE from wglDXOpenDeviceNV

	// Function pointers loaded via wglGetProcAddress
	fnDXOpenDevice      uintptr
	fnDXCloseDevice     uintptr
	fnDXRegisterObject  uintptr
	fnDXUnregisterObject uintptr
	fnDXLockObjects     uintptr
	fnDXUnlockObjects   uintptr

	available   bool
	deviceOpen  bool
}

var (
	opengl32           = syscall.NewLazyDLL("opengl32.dll")
	procWglGetProcAddr = opengl32.NewProc("wglGetProcAddress")
)

func wglGetProcAddress(name string) uintptr {
	namePtr, _ := syscall.BytePtrFromString(name)
	ret, _, _ := procWglGetProcAddr.Call(uintptr(unsafe.Pointer(namePtr)))
	return ret
}

// initWGLInterop loads all WGL_NV_DX_interop2 function pointers.
// Must be called with a current GL context (i.e., on the render thread).
func (w *wglInterop) initWGLInterop() {
	w.fnDXOpenDevice = wglGetProcAddress("wglDXOpenDeviceNV")
	w.fnDXCloseDevice = wglGetProcAddress("wglDXCloseDeviceNV")
	w.fnDXRegisterObject = wglGetProcAddress("wglDXRegisterObjectNV")
	w.fnDXUnregisterObject = wglGetProcAddress("wglDXUnregisterObjectNV")
	w.fnDXLockObjects = wglGetProcAddress("wglDXLockObjectsNV")
	w.fnDXUnlockObjects = wglGetProcAddress("wglDXUnlockObjectsNV")

	w.available = w.fnDXOpenDevice != 0 &&
		w.fnDXCloseDevice != 0 &&
		w.fnDXRegisterObject != 0 &&
		w.fnDXUnregisterObject != 0 &&
		w.fnDXLockObjects != 0 &&
		w.fnDXUnlockObjects != 0

	if w.available {
		fmt.Println("[WGL_INTEROP] WGL_NV_DX_interop2 functions loaded successfully")
	} else {
		fmt.Println("[WGL_INTEROP] WGL_NV_DX_interop2 NOT available")
	}
}

// openDevice registers a D3D11 device with OpenGL for interop.
// d3d11Device is an ID3D11Device* pointer.
func (w *wglInterop) openDevice(d3d11Device uintptr) error {
	if !w.available {
		return fmt.Errorf("WGL_NV_DX_interop2 not available")
	}
	if w.deviceOpen {
		return nil // Already open
	}

	ret, _, _ := syscall.SyscallN(w.fnDXOpenDevice, d3d11Device)
	if ret == 0 {
		return fmt.Errorf("wglDXOpenDeviceNV failed")
	}
	w.dxDevice = ret
	w.deviceOpen = true
	fmt.Printf("[WGL_INTEROP] D3D11 device registered with GL (handle=%x)\n", ret)
	return nil
}

// closeDevice unregisters the D3D11 device from OpenGL.
func (w *wglInterop) closeDevice() {
	if !w.deviceOpen {
		return
	}
	syscall.SyscallN(w.fnDXCloseDevice, w.dxDevice)
	w.dxDevice = 0
	w.deviceOpen = false
}

// registerObject registers a D3D11 texture with OpenGL.
// dxObject is an ID3D11Texture2D*, glName is the GL texture name,
// glTarget is GL_TEXTURE_2D, access is WGL_ACCESS_READ_ONLY_NV.
// Returns the interop handle.
func (w *wglInterop) registerObject(dxObject uintptr, glName uint32, glTarget, access uint32) (uintptr, error) {
	if !w.deviceOpen {
		return 0, fmt.Errorf("device not open")
	}
	ret, _, _ := syscall.SyscallN(w.fnDXRegisterObject,
		w.dxDevice,
		dxObject,
		uintptr(glName),
		uintptr(glTarget),
		uintptr(access))
	if ret == 0 {
		return 0, fmt.Errorf("wglDXRegisterObjectNV failed")
	}
	return ret, nil
}

// unregisterObject unregisters a previously registered D3D11 object.
func (w *wglInterop) unregisterObject(handle uintptr) {
	if !w.deviceOpen || handle == 0 {
		return
	}
	syscall.SyscallN(w.fnDXUnregisterObject, w.dxDevice, handle)
}

// lockObjects locks interop objects for GL access. Must be called before
// using the GL textures that were registered via registerObject.
func (w *wglInterop) lockObjects(handles []uintptr) error {
	if len(handles) == 0 {
		return nil
	}
	ret, _, _ := syscall.SyscallN(w.fnDXLockObjects,
		w.dxDevice,
		uintptr(len(handles)),
		uintptr(unsafe.Pointer(&handles[0])))
	if ret == 0 {
		return fmt.Errorf("wglDXLockObjectsNV failed")
	}
	return nil
}

// unlockObjects unlocks interop objects after GL rendering is done.
func (w *wglInterop) unlockObjects(handles []uintptr) {
	if len(handles) == 0 {
		return
	}
	syscall.SyscallN(w.fnDXUnlockObjects,
		w.dxDevice,
		uintptr(len(handles)),
		uintptr(unsafe.Pointer(&handles[0])))
}
