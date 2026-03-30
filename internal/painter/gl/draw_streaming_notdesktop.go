//go:build !(windows || darwin || linux || openbsd || freebsd)

package gl

import (
	"image"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

// uploadStreamingFrame uploads frame pixels directly via TexSubImage2D (no PBO).
// This is the fallback for platforms where PBOs are not available (WASM, mobile, GLES2).
func (p *painter) uploadStreamingFrame(img *canvas.StreamingImage, frame *image.RGBA) Texture {
	w := frame.Rect.Dx()
	h := frame.Rect.Dy()

	existingTex, cached := cache.GetTexture(img)
	texW, texH := img.TextureSize()

	if !img.DisableTexReuse && cached && w == texW && h == texH {
		// Fast path: same dimensions, update pixels in-place
		texture := Texture(existingTex)
		p.ctx.ActiveTexture(texture0)
		p.ctx.BindTexture(texture2D, texture)
		p.ctx.TexSubImage2D(texture2D, 0, 0, 0, w, h, colorFormatRGBA, unsignedByte, frame.Pix)
		p.logError()
		return texture
	}

	// Slow path: dimensions changed or first frame — (re)allocate texture
	if cached {
		p.ctx.DeleteTexture(Texture(existingTex))
		cache.DeleteTexture(img)
	}
	texture := p.imgToTexture(frame, img.ScaleMode)
	cache.SetTexture(img, cache.TextureType(texture), p.canvas)
	img.SetTextureSize(w, h)
	return texture
}
