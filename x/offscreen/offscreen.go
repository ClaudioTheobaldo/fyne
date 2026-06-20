//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

// Package offscreen exposes a headless renderer that runs a single custom
// fragment-shader effect over an image with no window or event loop. It is a
// thin public wrapper over the internal GL implementation, intended for offline
// shader validation/debugging (e.g. rendering an AddCustomEffect shader to a PNG
// outside the running app).
package offscreen

import (
	"image"

	"fyne.io/fyne/v2/internal/painter/gl"
)

// RenderEffect renders src through the fragment shader fragSrc at outW x outH and
// returns the result. fragSrc follows the custom-effect contract
// (`uniform sampler2D tex`, `varying vec2 fragTexCoord`, writing gl_FragColor);
// uniforms maps a name to 1..4 floats (scalar / vec2 / vec3 / vec4), exactly as
// effect.Effect.Uniforms() supplies them. Output matches the on-GPU result of the
// live effect pipeline; an identity (passthrough) shader returns src unchanged,
// with row 0 at the top for both input and output.
//
// It opens and tears down a hidden GL context, so it needs a usable display and
// must not run concurrently with a live Fyne GL window in the same process.
func RenderEffect(src *image.RGBA, fragSrc string, uniforms map[string][]float32, outW, outH int) (*image.RGBA, error) {
	return gl.RenderEffectOffscreen(src, fragSrc, uniforms, outW, outH)
}
