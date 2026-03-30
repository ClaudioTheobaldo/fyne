//go:build windows || darwin || linux || openbsd || freebsd

package gl

import (
	"image"
	"unsafe"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

// uploadStreamingFrameDirect uploads frame pixels via TexSubImage2D without PBOs.
// Used when PBOs are disabled or as a fallback.
func (p *painter) uploadStreamingFrameDirect(img *canvas.StreamingImage, frame *image.RGBA) Texture {
	w := frame.Rect.Dx()
	h := frame.Rect.Dy()

	existingTex, cached := cache.GetTexture(img)
	texW, texH := img.TextureSize()

	if !img.DisableTexReuse && cached && w == texW && h == texH {
		// Fast path: reuse texture, update pixels in-place
		texture := Texture(existingTex)
		p.ctx.ActiveTexture(texture0)
		p.ctx.BindTexture(texture2D, texture)
		p.ctx.TexSubImage2D(texture2D, 0, 0, 0, w, h, colorFormatRGBA, unsignedByte, frame.Pix)
		p.logError()
		return texture
	}

	// Slow path: (re)allocate texture
	if cached {
		p.ctx.DeleteTexture(Texture(existingTex))
		cache.DeleteTexture(img)
	}
	texture := p.imgToTexture(frame, img.ScaleMode)
	cache.SetTexture(img, cache.TextureType(texture), p.canvas)
	img.SetTextureSize(w, h)
	return texture
}

// uploadStreamingFrame uses PBO double-buffering to asynchronously upload
// frame pixels to the GPU. Falls back to direct TexSubImage2D if PBOs are disabled.
func (p *painter) uploadStreamingFrame(img *canvas.StreamingImage, frame *image.RGBA) Texture {
	if img.DisablePBO {
		return p.uploadStreamingFrameDirect(img, frame)
	}

	w := frame.Rect.Dx()
	h := frame.Rect.Dy()
	pbo := p.getOrCreatePBO(img, w, h)
	pixelSize := w * h * 4

	// --- Write pixels into the current PBO via MapBuffer ---
	writeIdx := pbo.index
	p.ctx.BindBuffer(pixelUnpackBuffer, pbo.buffers[writeIdx])
	// Orphan the buffer to avoid GPU sync stall (driver can allocate a new backing store)
	p.ctx.BufferDataBytes(pixelUnpackBuffer, pixelSize, nil, streamDraw)
	ptr := p.ctx.MapBuffer(pixelUnpackBuffer, writeOnly)
	if ptr != nil {
		dst := unsafe.Slice((*byte)(ptr), len(frame.Pix))
		copy(dst, frame.Pix)
		p.ctx.UnmapBuffer(pixelUnpackBuffer)
	}
	p.logError()

	// --- Upload to texture from the OTHER PBO (filled last frame) ---
	var texture Texture
	readIdx := 1 - writeIdx

	existingTex, cached := cache.GetTexture(img)
	texW, texH := img.TextureSize()
	canReuse := !img.DisableTexReuse && pbo.ready && cached && w == texW && h == texH

	if canReuse {
		// Fast path: async DMA from read PBO to existing texture
		texture = Texture(existingTex)
		p.ctx.BindBuffer(pixelUnpackBuffer, pbo.buffers[readIdx])
		p.ctx.ActiveTexture(texture0)
		p.ctx.BindTexture(texture2D, texture)
		p.ctx.TexSubImage2DPBO(texture2D, 0, 0, 0, w, h, colorFormatRGBA, unsignedByte)
		p.logError()
	} else {
		// New texture from PBO (first frame, resolution change, or tex reuse disabled)
		if cached {
			p.ctx.DeleteTexture(Texture(existingTex))
			cache.DeleteTexture(img)
		}

		srcIdx := writeIdx
		if pbo.ready {
			srcIdx = readIdx // use the previously filled PBO if available
		}

		tex := p.newTexture(img.ScaleMode)
		p.ctx.BindBuffer(pixelUnpackBuffer, pbo.buffers[srcIdx])
		p.ctx.TexImage2DPBO(texture2D, 0, w, h, colorFormatRGBA, unsignedByte)
		p.logError()

		cache.SetTexture(img, cache.TextureType(tex), p.canvas)
		img.SetTextureSize(w, h)
		texture = tex
		pbo.ready = true
	}

	// Unbind PBO so other TexImage2D calls aren't affected
	p.ctx.BindBuffer(pixelUnpackBuffer, noBuffer)
	p.logError()

	// Swap for next frame
	pbo.index = 1 - pbo.index

	return texture
}
