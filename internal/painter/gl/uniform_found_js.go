//go:build js

package gl

import gl "github.com/fyne-io/gl-js"

// uniformFound reports whether GetUniformLocation found the uniform. WebGL
// returns null rather than -1 for a uniform the program does not have.
func uniformFound(u Uniform) bool {
	v := gl.Uniform(u).Value
	return !v.IsNull() && !v.IsUndefined()
}
