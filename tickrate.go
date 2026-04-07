package fyne

// InitialTickRate is the render loop tick rate in Hz used when the GLFW driver
// starts. Set this before calling ShowAndRun. A value of 0 means the driver
// defaults to 60 Hz.
var InitialTickRate int

// SetRenderTickRate changes the render loop frame rate at runtime (Hz).
// It is wired by the GLFW driver immediately after the event ticker is created,
// so it is safe to call from any goroutine once the application is running.
// A nil value means the driver has not started yet.
var SetRenderTickRate func(rate int)
