//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

package gl

// offscreen.go is a self-contained, windowless renderer that runs a single
// fragment-shader effect over a source image and returns the result. It exists
// purely as an offline debug/validation oracle (render a custom effect shader to
// a PNG with no GUI, no event loop) and shares NO state with the live painter:
// it creates its own hidden GL context and its own FBO, and does not touch any
// painter method or the live effect pipeline.
//
// Faithfulness: it pairs the caller's fragment shader with the SAME passthrough
// vertex shader (effectPassthroughVertSrc) and SAME fullscreen quad
// (effectQuadVerts) the live effect pipeline uses, and feeds uniforms the same
// way, so a given (shader, uniforms) produces the same on-GPU result as
// AddCustomEffect. The source/passthrough convention is "upright in -> upright
// out": an identity (passthrough) shader returns the source image unchanged, and
// row 0 of both the source and the returned image is the top row.

import (
	"fmt"
	"image"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// RenderEffectOffscreen renders src through the fragment shader fragSrc (which
// must follow the custom-effect contract: `uniform sampler2D tex`,
// `varying vec2 fragTexCoord`, writing gl_FragColor) at outW x outH, and returns
// the result. uniforms maps a name to 1..4 floats (scalar / vec2 / vec3 / vec4),
// exactly as effect.Effect.Uniforms() would supply them.
//
// It is intended for offline tooling/tests; it opens (and tears down) a hidden
// GLFW window for the GL context, so it must run with a usable display and is
// not safe to call concurrently with a live Fyne GL window in the same process.
func RenderEffectOffscreen(src *image.RGBA, fragSrc string, uniforms map[string][]float32, outW, outH int) (out *image.RGBA, err error) {
	if src == nil {
		return nil, fmt.Errorf("offscreen: nil source image")
	}
	if outW <= 0 || outH <= 0 {
		return nil, fmt.Errorf("offscreen: invalid output size %dx%d", outW, outH)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("offscreen: glfw init: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	win, err := glfw.CreateWindow(1, 1, "fyne-offscreen", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("offscreen: create hidden window: %w", err)
	}
	defer win.Destroy()
	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("offscreen: gl init: %w", err)
	}

	prog, err := buildOffscreenProgram(string(effectPassthroughVertSrc()), fragSrc)
	if err != nil {
		return nil, err
	}
	defer gl.DeleteProgram(prog)

	// Source texture, uploaded as-is (row 0 = top). The quad maps source row 0 to
	// the framebuffer bottom and GL reads back bottom-first, so the two cancel:
	// sampling and readback are both upright and match the CPU reference model.
	srcTex := uploadSourceTexture(src)
	defer gl.DeleteTextures(1, &srcTex)

	// Output FBO + colour texture.
	var outTex uint32
	gl.GenTextures(1, &outTex)
	gl.BindTexture(gl.TEXTURE_2D, outTex)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(outW), int32(outH), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	defer gl.DeleteTextures(1, &outTex)

	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, outTex, 0)
	defer gl.DeleteFramebuffers(1, &fbo)
	if st := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); st != gl.FRAMEBUFFER_COMPLETE {
		return nil, fmt.Errorf("offscreen: framebuffer incomplete: 0x%x", st)
	}

	// Fullscreen quad (same vertices/layout as the live effect pipeline).
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(effectQuadVerts)*4, gl.Ptr(effectQuadVerts), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &vbo)

	gl.Viewport(0, 0, int32(outW), int32(outH))
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(prog)

	const stride = 4 * 4 // 4 floats per vertex
	vertLoc := uint32(gl.GetAttribLocation(prog, gl.Str("vert\x00")))
	texLoc := uint32(gl.GetAttribLocation(prog, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(vertLoc)
	gl.VertexAttribPointerWithOffset(vertLoc, 2, gl.FLOAT, false, stride, 0)
	gl.EnableVertexAttribArray(texLoc)
	gl.VertexAttribPointerWithOffset(texLoc, 2, gl.FLOAT, false, stride, 2*4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, srcTex)
	if loc := gl.GetUniformLocation(prog, gl.Str("tex\x00")); loc >= 0 {
		gl.Uniform1i(loc, 0)
	}

	for name, v := range uniforms {
		loc := gl.GetUniformLocation(prog, gl.Str(name+"\x00"))
		if loc < 0 {
			continue
		}
		switch len(v) {
		case 1:
			gl.Uniform1f(loc, v[0])
		case 2:
			gl.Uniform2f(loc, v[0], v[1])
		case 3:
			gl.Uniform3f(loc, v[0], v[1], v[2])
		case 4:
			gl.Uniform4f(loc, v[0], v[1], v[2], v[3])
		}
	}

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	pix := make([]byte, outW*outH*4)
	gl.ReadPixels(0, 0, int32(outW), int32(outH), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pix))

	// No flip on readback: the quad mapping + GL's bottom-first readback cancel,
	// so straight-through copy yields an upright, model-matching image.
	out = image.NewRGBA(image.Rect(0, 0, outW, outH))
	copy(out.Pix, pix)
	return out, nil
}

// uploadSourceTexture uploads an RGBA image as a 2D texture row 0 first (top),
// packed tight (honouring the image stride / sub-image origin).
func uploadSourceTexture(src *image.RGBA) uint32 {
	w := src.Rect.Dx()
	h := src.Rect.Dy()
	tight := make([]byte, w*h*4)
	row := w * 4
	for y := 0; y < h; y++ {
		srcOff := src.PixOffset(src.Rect.Min.X, src.Rect.Min.Y+y)
		copy(tight[y*row:(y+1)*row], src.Pix[srcOff:srcOff+row])
	}
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(tight))
	return tex
}

func buildOffscreenProgram(vertSrc, fragSrc string) (uint32, error) {
	vs, err := compileOffscreenShader(vertSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("offscreen: vertex shader: %w", err)
	}
	defer gl.DeleteShader(vs)
	fs, err := compileOffscreenShader(fragSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("offscreen: fragment shader: %w", err)
	}
	defer gl.DeleteShader(fs)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, fs)
	gl.LinkProgram(prog)
	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var n int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &n)
		logMsg := strings.Repeat("\x00", int(n)+1)
		gl.GetProgramInfoLog(prog, n, nil, gl.Str(logMsg))
		return 0, fmt.Errorf("offscreen: link program: %s", strings.TrimRight(logMsg, "\x00"))
	}
	return prog, nil
}

func compileOffscreenShader(src string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csrc, free := gl.Strs(src + "\x00")
	gl.ShaderSource(shader, 1, csrc, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var n int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &n)
		logMsg := strings.Repeat("\x00", int(n)+1)
		gl.GetShaderInfoLog(shader, n, nil, gl.Str(logMsg))
		gl.DeleteShader(shader)
		return 0, fmt.Errorf("%s", strings.TrimRight(logMsg, "\x00"))
	}
	return shader, nil
}
