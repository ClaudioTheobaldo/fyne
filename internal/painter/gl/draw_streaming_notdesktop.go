//go:build !(windows || darwin || linux || openbsd || freebsd)

package gl

import (
	"image"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

// glFormatBytesPerTexel returns the number of bytes per texel for common GL formats.
func glFormatBytesPerTexel(glFormat uint32) int {
	switch glFormat {
	case fdLumA:
		return 2
	case fdRGB:
		return 3
	case fdRGBA:
		return 4
	default:
		return 1
	}
}

// uploadStreamingRawFrame uploads all planes of a RawFrame directly via TexSubImage2D.
// This is the non-desktop (mobile/WASM) path that does not use PBOs.
func (p *painter) uploadStreamingRawFrame(img *canvas.StreamingImage, frame *canvas.RawFrame, desc formatDescriptor) [4]Texture {
	w := frame.Width
	h := frame.Height
	texUnits := [4]uint32{texture0, texture1, texture2, texture3}
	var textures [4]Texture

	for pi := 0; pi < desc.planeCount; pi++ {
		planeW := w / desc.chromaWDiv[pi]
		planeH := h / desc.chromaHDiv[pi]
		bpt := glFormatBytesPerTexel(desc.planeFormats[pi])
		stride := frame.Strides[pi]
		if stride == 0 {
			stride = planeW * bpt
		}
		data := frame.Data[pi]
		glFmt := desc.planeFormats[pi]
		glType := desc.planeDataTypes[pi]

		p.ctx.ActiveTexture(texUnits[pi])
		// Reuse or create texture
		tex := p.newTexture(img.ScaleMode)
		p.ctx.BindTexture(texture2D, tex)
		if stride != planeW*bpt {
			p.ctx.PixelStorei(unpackRowLength, int32(stride/bpt))
		}
		p.ctx.TexImage2D(texture2D, 0, planeW, planeH, glFmt, glType, data)
		if stride != planeW*bpt {
			p.ctx.PixelStorei(unpackRowLength, 0)
		}
		p.logError()
		textures[pi] = tex
	}

	img.SetTextureSize(w, h)
	return textures
}

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
