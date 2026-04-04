//go:build (gles || arm || arm64) && !darwin && !wasm && !test_web_driver

package glfw

import "github.com/go-gl/glfw/v3.3/glfw"

func initWindowHints(msaaSamples int) {
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	// MSAA not supported on GLES2 — msaaSamples parameter ignored.
}
