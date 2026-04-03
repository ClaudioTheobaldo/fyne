//go:build windows || darwin || linux || openbsd || freebsd

package gl

import (
	"image"
	"unsafe"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

// glFormatBytesPerTexel returns the number of bytes per texel for common GL
// internal formats used in streaming uploads.
func glFormatBytesPerTexel(glFormat uint32) int {
	switch glFormat {
	case fdLumA: // GL_LUMINANCE_ALPHA = 0x190A: 2 bytes/texel
		return 2
	case fdRGB: // GL_RGB = 0x1907: 3 bytes/texel
		return 3
	case fdRGBA: // GL_RGBA = 0x1908: 4 bytes/texel
		return 4
	default: // fdLum (GL_LUMINANCE = 0x1909) and unknown: 1 byte/texel
		return 1
	}
}

// uploadGenericPlane writes plane data into the write PBO, then uploads from
// the read PBO into the texture at the specified texture unit.
// w and h are the texture dimensions in texels; stride is the raw byte row length.
// glFormat and glType are the GL format and data type constants.
func (p *painter) uploadGenericPlane(
	pbo planePBOPair, tex Texture, writeIdx, readIdx int,
	w, h, stride, bytesPerTexel int,
	data []byte, ready bool,
	glFormat, glType, texUnit uint32,
) Texture {
	pixelSize := stride * h

	p.ctx.BindBuffer(pixelUnpackBuffer, pbo.buffers[writeIdx])
	p.ctx.BufferDataBytes(pixelUnpackBuffer, pixelSize, nil, streamDraw)
	ptr := p.ctx.MapBuffer(pixelUnpackBuffer, writeOnly)
	if ptr != nil {
		dst := unsafe.Slice((*byte)(ptr), pixelSize)
		copy(dst, data[:pixelSize])
		p.ctx.UnmapBuffer(pixelUnpackBuffer)
	}

	rowLenTexels := stride / bytesPerTexel
	if rowLenTexels != w {
		p.ctx.PixelStorei(unpackRowLength, int32(rowLenTexels))
	}

	p.ctx.ActiveTexture(texUnit)
	p.ctx.BindTexture(texture2D, tex)
	if ready {
		p.ctx.BindBuffer(pixelUnpackBuffer, pbo.buffers[readIdx])
		p.ctx.TexSubImage2DPBO(texture2D, 0, 0, 0, w, h, glFormat, glType)
	} else {
		p.ctx.BindBuffer(pixelUnpackBuffer, pbo.buffers[writeIdx])
		p.ctx.TexImage2DPBO(texture2D, 0, w, h, glFormat, glType)
	}

	if rowLenTexels != w {
		p.ctx.PixelStorei(unpackRowLength, 0)
	}

	p.ctx.BindBuffer(pixelUnpackBuffer, noBuffer)
	p.logError()
	return tex
}

// uploadStreamingRawFrame uploads all planes of a RawFrame via PBO double-buffering
// according to the format descriptor. Returns the uploaded textures (up to 4).
func (p *painter) uploadStreamingRawFrame(img *canvas.StreamingImage, frame *canvas.RawFrame, desc formatDescriptor) [4]Texture {
	w := frame.Width
	h := frame.Height

	// Compute per-plane sizes and create/reuse PBO state
	var planeSizes [4]int
	for pi := 0; pi < desc.planeCount; pi++ {
		planeW := w / desc.chromaWDiv[pi]
		planeH := h / desc.chromaHDiv[pi]
		bpt := glFormatBytesPerTexel(desc.planeFormats[pi])
		if frame.Strides[pi] > 0 {
			planeSizes[pi] = frame.Strides[pi] * planeH
		} else {
			planeSizes[pi] = planeW * planeH * bpt
		}
	}

	state := p.getOrCreateStreamPBO(img, desc.planeCount, w, h, planeSizes)
	writeIdx := state.index
	readIdx := 1 - writeIdx

	for pi := 0; pi < desc.planeCount; pi++ {
		planeW := w / desc.chromaWDiv[pi]
		planeH := h / desc.chromaHDiv[pi]
		bpt := glFormatBytesPerTexel(desc.planeFormats[pi])
		stride := frame.Strides[pi]
		if stride == 0 {
			stride = planeW * bpt
		}
		texUnits := [4]uint32{texture0, texture1, texture2, texture3}
		state.textures[pi] = p.uploadGenericPlane(
			state.planes[pi], state.textures[pi], writeIdx, readIdx,
			planeW, planeH, stride, bpt,
			frame.Data[pi], state.ready,
			desc.planeFormats[pi], desc.planeDataTypes[pi], texUnits[pi],
		)
	}

	state.ready = true
	state.index = 1 - state.index
	img.SetTextureSize(w, h)
	return state.textures
}

// uploadSinglePlane writes a single-channel plane into the write PBO and uploads
// from the read PBO to the given texture. Returns the texture.
func (p *painter) uploadSinglePlane(
	pbo [2]Buffer, tex Texture, writeIdx, readIdx int,
	w, h, stride int, data []byte, ready bool, scaleMode canvas.ImageScale,
) Texture {
	pixelSize := stride * h

	// Write plane into current PBO
	p.ctx.BindBuffer(pixelUnpackBuffer, pbo[writeIdx])
	p.ctx.BufferDataBytes(pixelUnpackBuffer, pixelSize, nil, streamDraw)
	ptr := p.ctx.MapBuffer(pixelUnpackBuffer, writeOnly)
	if ptr != nil {
		dst := unsafe.Slice((*byte)(ptr), pixelSize)
		copy(dst, data[:pixelSize])
		p.ctx.UnmapBuffer(pixelUnpackBuffer)
	}

	// Set row length for stride handling (reset after upload)
	if stride != w {
		p.ctx.PixelStorei(unpackRowLength, int32(stride))
	}

	// Upload from the other PBO to texture
	if ready {
		p.ctx.BindBuffer(pixelUnpackBuffer, pbo[readIdx])
		p.ctx.ActiveTexture(texture0)
		p.ctx.BindTexture(texture2D, tex)
		p.ctx.TexSubImage2DPBO(texture2D, 0, 0, 0, w, h, colorFormatLuminance, unsignedByte)
	} else {
		// First frame — allocate texture from write PBO
		p.ctx.BindBuffer(pixelUnpackBuffer, pbo[writeIdx])
		p.ctx.ActiveTexture(texture0)
		p.ctx.BindTexture(texture2D, tex)
		p.ctx.TexImage2DPBO(texture2D, 0, w, h, colorFormatLuminance, unsignedByte)
	}

	if stride != w {
		p.ctx.PixelStorei(unpackRowLength, 0)
	}

	p.ctx.BindBuffer(pixelUnpackBuffer, noBuffer)
	p.logError()
	return tex
}


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
