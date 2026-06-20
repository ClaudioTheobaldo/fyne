//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

package gl

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

// makeTopRedBottomBlue builds a w x h RGBA: top half red, bottom half blue.
// Used to detect both colour fidelity and vertical orientation.
func makeTopRedBottomBlue(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		c := color.RGBA{R: 255, A: 255}
		if y >= h/2 {
			c = color.RGBA{B: 255, A: 255}
		}
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func near(a, b uint8) bool {
	d := int(a) - int(b)
	return d < 8 && d > -8
}

// TestOffscreenPassthroughIdentity verifies the renderer round-trips a source
// unchanged through the passthrough shader — proving the GL context, sampling,
// readback, and the upright-in/upright-out orientation convention all hold.
func TestOffscreenPassthroughIdentity(t *testing.T) {
	src := makeTopRedBottomBlue(4, 4)
	out, err := RenderEffectOffscreen(src, string(effectPassthroughFragSrc()), nil, 4, 4)
	if err != nil {
		if strings.Contains(err.Error(), "glfw init") || strings.Contains(err.Error(), "hidden window") || strings.Contains(err.Error(), "gl init") {
			t.Skip("no usable GL context in this environment:", err)
		}
		t.Fatal(err)
	}

	top := out.RGBAAt(0, 0)       // expect red (orientation preserved)
	bottom := out.RGBAAt(0, 3)    // expect blue
	if !(near(top.R, 255) && near(top.G, 0) && near(top.B, 0)) {
		t.Errorf("top pixel = %+v, want red (orientation/identity broken)", top)
	}
	if !(near(bottom.R, 0) && near(bottom.G, 0) && near(bottom.B, 255)) {
		t.Errorf("bottom pixel = %+v, want blue", bottom)
	}
}

// TestOffscreenChannelSwap verifies a custom fragment shader actually runs:
// swapping R<->B turns the red top blue and the blue bottom red.
func TestOffscreenChannelSwap(t *testing.T) {
	const frag = `#version 110
uniform sampler2D tex;
varying vec2 fragTexCoord;
void main() {
    vec4 c = texture2D(tex, fragTexCoord);
    gl_FragColor = vec4(c.b, c.g, c.r, c.a);
}`
	src := makeTopRedBottomBlue(4, 4)
	out, err := RenderEffectOffscreen(src, frag, nil, 4, 4)
	if err != nil {
		if strings.Contains(err.Error(), "glfw init") || strings.Contains(err.Error(), "hidden window") || strings.Contains(err.Error(), "gl init") {
			t.Skip("no usable GL context in this environment:", err)
		}
		t.Fatal(err)
	}
	if top := out.RGBAAt(0, 0); !near(top.B, 255) || !near(top.R, 0) {
		t.Errorf("top after R<->B swap = %+v, want blue", top)
	}
	if bot := out.RGBAAt(0, 3); !near(bot.R, 255) || !near(bot.B, 0) {
		t.Errorf("bottom after R<->B swap = %+v, want red", bot)
	}
}

// TestOffscreenUniform verifies float uniforms reach the shader: a vec3 tint of
// (0.5,0,0) halves the red channel and zeroes the rest.
func TestOffscreenUniform(t *testing.T) {
	const frag = `#version 110
uniform sampler2D tex;
uniform vec3 tint;
varying vec2 fragTexCoord;
void main() {
    vec4 c = texture2D(tex, fragTexCoord);
    gl_FragColor = vec4(c.rgb * tint, c.a);
}`
	src := makeTopRedBottomBlue(4, 4)
	out, err := RenderEffectOffscreen(src, frag, map[string][]float32{"tint": {0.5, 0, 0}}, 4, 4)
	if err != nil {
		if strings.Contains(err.Error(), "glfw init") || strings.Contains(err.Error(), "hidden window") || strings.Contains(err.Error(), "gl init") {
			t.Skip("no usable GL context in this environment:", err)
		}
		t.Fatal(err)
	}
	top := out.RGBAAt(0, 0) // red 255 * 0.5 ~= 127, G/B zeroed
	if !near(top.R, 127) || !near(top.G, 0) || !near(top.B, 0) {
		t.Errorf("top after tint = %+v, want ~{127,0,0}", top)
	}
}
