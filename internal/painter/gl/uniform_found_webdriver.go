//go:build test_web_driver && !js

package gl

import gl "github.com/fyne-io/gl-js"

// uniformFound reports whether GetUniformLocation found the uniform. Outside
// a browser gl-js wraps the desktop location, which is -1 when not found.
func uniformFound(u Uniform) bool {
	return gl.Uniform(u).Value >= 0
}
