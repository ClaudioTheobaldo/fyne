// Package main is a visual showcase of every pixel format supported by
// canvas.StreamingImage. It converts a 64×64 color-wheel test image into each
// format, round-trips it back to RGBA, and displays all results in a
// scrollable grid so you can visually compare fidelity and artifacts.
//
// Run with:
//
//	go run ./cmd/streaming_showcase
package main

import (
	"encoding/binary"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/painter/gl"
	incanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

const imgW, imgH = 64, 64

// ─────────────────────────────────────────────────────────────────────────────
// Test-image generator
// ─────────────────────────────────────────────────────────────────────────────

// makeTestImage returns a 64×64 color-wheel RGBA image with a white-to-gray
// radial fade, giving strong signals in all channels.
func makeTestImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	cx, cy := float64(imgW)/2-0.5, float64(imgH)/2-0.5
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			dx, dy := float64(x)-cx, float64(y)-cy
			dist := math.Sqrt(dx*dx+dy*dy) / (float64(imgW) / 2)
			hue := (math.Atan2(dy, dx)/math.Pi + 1) / 2 // 0…1
			sat := math.Min(dist, 1)
			val := 1 - 0.25*dist
			r, g, b := hsvToRGB(hue, sat, val)
			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	if s == 0 {
		c := uint8(v * 255)
		return c, c, c
	}
	h6 := h * 6
	i := math.Floor(h6)
	f := h6 - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))
	var r, g, b float64
	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}
	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

// ─────────────────────────────────────────────────────────────────────────────
// BT.601 full-range helpers (shared by all YUV encode / decode paths)
// ─────────────────────────────────────────────────────────────────────────────

func rgbToYCbCr(r, g, b uint8) (y, cb, cr uint8) {
	rf, gf, bf := float64(r), float64(g), float64(b)
	yf := 0.299*rf + 0.587*gf + 0.114*bf
	cbf := -0.168736*rf - 0.331264*gf + 0.5*bf + 128
	crf := 0.5*rf - 0.418688*gf - 0.081312*bf + 128
	clamp := func(v float64) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(math.Round(v))
	}
	return clamp(yf), clamp(cbf), clamp(crf)
}

func yCbCrToRGB(yv, cb, cr uint8) (r, g, b uint8) {
	y, u, v := float64(yv), float64(cb)-128, float64(cr)-128
	rf := y + 1.402*v
	gf := y - 0.344136*u - 0.714136*v
	bf := y + 1.772*u
	clamp := func(v float64) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(math.Round(v))
	}
	return clamp(rf), clamp(gf), clamp(bf)
}

func lumaOf(r, g, b uint8) uint8 {
	y, _, _ := rgbToYCbCr(r, g, b)
	return y
}

// ─────────────────────────────────────────────────────────────────────────────
// Format entry
// ─────────────────────────────────────────────────────────────────────────────

type entry struct {
	label   string
	convert func(*image.RGBA) *image.RGBA // full round-trip: RGBA → format → RGBA
}

// ─────────────────────────────────────────────────────────────────────────────
// CPU-convert round-trips  (encode here, decode via gl.ConvertRawFrameToRGBA)
// ─────────────────────────────────────────────────────────────────────────────

func cpuRoundTrip(pixFmt int, frame *incanvas.RawFrame) *image.RGBA {
	out := gl.ConvertRawFrameToRGBA(pixFmt, frame)
	if out == nil {
		return image.NewRGBA(image.Rect(0, 0, imgW, imgH)) // blank on failure
	}
	return out
}

func encodeRGB24(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*3)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			i := (y*imgW + x) * 3
			data[i], data[i+1], data[i+2] = c.R, c.G, c.B
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeBGR24(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*3)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			i := (y*imgW + x) * 3
			data[i], data[i+1], data[i+2] = c.B, c.G, c.R
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeBGRA(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*4)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			i := (y*imgW + x) * 4
			data[i], data[i+1], data[i+2], data[i+3] = c.B, c.G, c.R, c.A
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeARGB(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*4)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			i := (y*imgW + x) * 4
			data[i], data[i+1], data[i+2], data[i+3] = c.A, c.R, c.G, c.B
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeABGR(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*4)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			i := (y*imgW + x) * 4
			data[i], data[i+1], data[i+2], data[i+3] = c.A, c.B, c.G, c.R
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeRGB565(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*2)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			r5 := uint16(c.R >> 3)
			g6 := uint16(c.G >> 2)
			b5 := uint16(c.B >> 3)
			v := (r5 << 11) | (g6 << 5) | b5
			binary.LittleEndian.PutUint16(data[(y*imgW+x)*2:], v)
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeBGR565(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*2)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			b5 := uint16(c.B >> 3)
			g6 := uint16(c.G >> 2)
			r5 := uint16(c.R >> 3)
			v := (b5 << 11) | (g6 << 5) | r5
			binary.LittleEndian.PutUint16(data[(y*imgW+x)*2:], v)
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeGBRP(src *image.RGBA, bits int) *incanvas.RawFrame {
	n := imgW * imgH
	bpt := 1
	if bits > 8 {
		bpt = 2
	}
	g, b, r := make([]byte, n*bpt), make([]byte, n*bpt), make([]byte, n*bpt)
	maxValInt := (1 << uint(bits)) - 1
	maxVal := float64(maxValInt)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			idx := y*imgW + x
			if bits == 8 {
				g[idx], b[idx], r[idx] = c.G, c.B, c.R
			} else {
				scale := maxVal / 255.0
				gv := uint16(math.Round(float64(c.G) * scale))
				bv := uint16(math.Round(float64(c.B) * scale))
				rv := uint16(math.Round(float64(c.R) * scale))
				binary.LittleEndian.PutUint16(g[idx*2:], gv)
				binary.LittleEndian.PutUint16(b[idx*2:], bv)
				binary.LittleEndian.PutUint16(r[idx*2:], rv)
			}
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{g, b, r}}
}

func encodeGBRAP(src *image.RGBA) *incanvas.RawFrame {
	n := imgW * imgH
	g, b, r, a := make([]byte, n), make([]byte, n), make([]byte, n), make([]byte, n)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			i := y*imgW + x
			g[i], b[i], r[i], a[i] = c.G, c.B, c.R, c.A
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{g, b, r, a}}
}

func encodeMono(src *image.RGBA, whiteIsZero bool) *incanvas.RawFrame {
	stride := (imgW + 7) / 8
	data := make([]byte, stride*imgH)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			luma := lumaOf(c.R, c.G, c.B)
			isLight := luma >= 128
			var bit byte
			if whiteIsZero {
				if !isLight {
					bit = 1 // black → bit=1 in MONOWHITE
				}
			} else {
				if isLight {
					bit = 1 // white → bit=1 in MONOBLACK
				}
			}
			if bit == 1 {
				data[y*stride+x/8] |= 1 << uint(7-x%8)
			}
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

// encodeBayer samples only the channel appropriate for each Bayer grid position.
// pat: 0=BGGR, 1=RGGB, 2=GBRG, 3=GRBG
func encodeBayer(src *image.RGBA, pat int) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH)
	// colorAtGrid: returns (r,g,b channel index: 0=R,1=G,2=B) for Bayer cell (cx,cy)
	colorAtGrid := func(cx, cy int) int {
		idx := cy*2 + cx
		switch pat {
		case 0: // BGGR: B G / G R
			switch idx {
			case 0:
				return 2 // B
			case 1, 2:
				return 1 // G
			default:
				return 0 // R
			}
		case 1: // RGGB: R G / G B
			switch idx {
			case 0:
				return 0 // R
			case 1, 2:
				return 1 // G
			default:
				return 2 // B
			}
		case 2: // GBRG: G B / R G
			switch idx {
			case 0, 3:
				return 1 // G
			case 1:
				return 2 // B
			default:
				return 0 // R
			}
		default: // GRBG: G R / B G
			switch idx {
			case 0, 3:
				return 1 // G
			case 1:
				return 0 // R
			default:
				return 2 // B
			}
		}
	}
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			ch := colorAtGrid(x%2, y%2)
			var val uint8
			switch ch {
			case 0:
				val = c.R
			case 1:
				val = c.G
			case 2:
				val = c.B
			}
			data[y*imgW+x] = val
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeGrayF32(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*4)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			luma := float32(lumaOf(c.R, c.G, c.B)) / 255.0
			binary.LittleEndian.PutUint32(data[(y*imgW+x)*4:], math.Float32bits(luma))
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeRGBF32(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*12)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			base := (y*imgW + x) * 12
			binary.LittleEndian.PutUint32(data[base+0:], math.Float32bits(float32(c.R)/255))
			binary.LittleEndian.PutUint32(data[base+4:], math.Float32bits(float32(c.G)/255))
			binary.LittleEndian.PutUint32(data[base+8:], math.Float32bits(float32(c.B)/255))
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

func encodeRGBAF32(src *image.RGBA) *incanvas.RawFrame {
	data := make([]byte, imgW*imgH*16)
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			base := (y*imgW + x) * 16
			binary.LittleEndian.PutUint32(data[base+0:], math.Float32bits(float32(c.R)/255))
			binary.LittleEndian.PutUint32(data[base+4:], math.Float32bits(float32(c.G)/255))
			binary.LittleEndian.PutUint32(data[base+8:], math.Float32bits(float32(c.B)/255))
			binary.LittleEndian.PutUint32(data[base+12:], math.Float32bits(float32(c.A)/255))
		}
	}
	return &incanvas.RawFrame{Width: imgW, Height: imgH, Data: [4][]byte{data}}
}

// ─────────────────────────────────────────────────────────────────────────────
// GPU-shader-category encode + CPU decode (for visual showcase only)
// ─────────────────────────────────────────────────────────────────────────────

// encodeYUVPlanar encodes src into Y + Cb + Cr planes.
// chromaWDiv and chromaHDiv control chroma subsampling (e.g. 2,2 for 420).
func encodeYUVPlanar(src *image.RGBA, chromaWDiv, chromaHDiv int) (yPlane, uPlane, vPlane []byte) {
	yPlane = make([]byte, imgW*imgH)
	cW, cH := imgW/chromaWDiv, imgH/chromaHDiv
	uPlane = make([]byte, cW*cH)
	vPlane = make([]byte, cW*cH)

	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			c := src.RGBAAt(x, y)
			yv, _, _ := rgbToYCbCr(c.R, c.G, c.B)
			yPlane[y*imgW+x] = yv
		}
	}
	// Average Cb/Cr over each chroma block
	for cy := 0; cy < cH; cy++ {
		for cx := 0; cx < cW; cx++ {
			var sumCb, sumCr float64
			count := 0
			for dy := 0; dy < chromaHDiv; dy++ {
				for dx := 0; dx < chromaWDiv; dx++ {
					px, py := cx*chromaWDiv+dx, cy*chromaHDiv+dy
					if px < imgW && py < imgH {
						c := src.RGBAAt(px, py)
						_, cb, cr := rgbToYCbCr(c.R, c.G, c.B)
						sumCb += float64(cb)
						sumCr += float64(cr)
						count++
					}
				}
			}
			if count > 0 {
				uPlane[cy*cW+cx] = uint8(math.Round(sumCb / float64(count)))
				vPlane[cy*cW+cx] = uint8(math.Round(sumCr / float64(count)))
			}
		}
	}
	return
}

func decodeYUVPlanar(y, u, v []byte, chromaWDiv, chromaHDiv int) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	cW := imgW / chromaWDiv
	for py := 0; py < imgH; py++ {
		for px := 0; px < imgW; px++ {
			yv := y[py*imgW+px]
			cb := u[(py/chromaHDiv)*cW+(px/chromaWDiv)]
			cr := v[(py/chromaHDiv)*cW+(px/chromaWDiv)]
			r, g, b := yCbCrToRGB(yv, cb, cr)
			out.SetRGBA(px, py, color.RGBA{r, g, b, 255})
		}
	}
	return out
}

func yuvPlanarEntry(label string, wDiv, hDiv int) entry {
	return entry{label, func(src *image.RGBA) *image.RGBA {
		y, u, v := encodeYUVPlanar(src, wDiv, hDiv)
		return decodeYUVPlanar(y, u, v, wDiv, hDiv)
	}}
}

// YUVA planar: same as YUV + alpha plane
func yuvaDecodeEntry(label string, wDiv, hDiv int) entry {
	return entry{label, func(src *image.RGBA) *image.RGBA {
		y, u, v := encodeYUVPlanar(src, wDiv, hDiv)
		out := decodeYUVPlanar(y, u, v, wDiv, hDiv)
		// Restore alpha from source (alpha plane is lossless in YUVA)
		for py := 0; py < imgH; py++ {
			for px := 0; px < imgW; px++ {
				c := out.RGBAAt(px, py)
				c.A = src.RGBAAt(px, py).A
				out.SetRGBA(px, py, c)
			}
		}
		return out
	}}
}

// Hi-bit planar: encode to 16-bit LE per channel, decode same way
func hibitPlanarEntry(label string, bits, wDiv, hDiv int) entry {
	return entry{label, func(src *image.RGBA) *image.RGBA {
		// Encode to hi-bit (round-trip through GBRP{bits} format)
		yp, up, vp := encodeYUVPlanarHiBit(src, bits, wDiv, hDiv)
		return decodeYUVPlanarHiBit(yp, up, vp, bits, wDiv, hDiv)
	}}
}

func encodeYUVPlanarHiBit(src *image.RGBA, bits, wDiv, hDiv int) (y, u, v []byte) {
	maxValInt := (1 << uint(bits)) - 1
	maxVal := float64(maxValInt)
	scale := maxVal / 255.0
	y = make([]byte, imgW*imgH*2)
	cW, cH := imgW/wDiv, imgH/hDiv
	u = make([]byte, cW*cH*2)
	v = make([]byte, cW*cH*2)
	for py := 0; py < imgH; py++ {
		for px := 0; px < imgW; px++ {
			c := src.RGBAAt(px, py)
			yv, _, _ := rgbToYCbCr(c.R, c.G, c.B)
			binary.LittleEndian.PutUint16(y[(py*imgW+px)*2:], uint16(math.Round(float64(yv)*scale)))
		}
	}
	for cy := 0; cy < cH; cy++ {
		for cx := 0; cx < cW; cx++ {
			var sumCb, sumCr float64
			count := 0
			for dy := 0; dy < hDiv; dy++ {
				for dx := 0; dx < wDiv; dx++ {
					px2, py2 := cx*wDiv+dx, cy*hDiv+dy
					if px2 < imgW && py2 < imgH {
						c := src.RGBAAt(px2, py2)
						_, cb, cr := rgbToYCbCr(c.R, c.G, c.B)
						sumCb += float64(cb)
						sumCr += float64(cr)
						count++
					}
				}
			}
			if count > 0 {
				cbHi := uint16(math.Round(sumCb / float64(count) * scale))
				crHi := uint16(math.Round(sumCr / float64(count) * scale))
				binary.LittleEndian.PutUint16(u[(cy*cW+cx)*2:], cbHi)
				binary.LittleEndian.PutUint16(v[(cy*cW+cx)*2:], crHi)
			}
		}
	}
	return
}

func decodeYUVPlanarHiBit(y, u, v []byte, bits, wDiv, hDiv int) *image.RGBA {
	maxValInt := (1 << uint(bits)) - 1
	maxVal := float64(maxValInt)
	cW := imgW / wDiv
	out := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	for py := 0; py < imgH; py++ {
		for px := 0; px < imgW; px++ {
			yRaw := binary.LittleEndian.Uint16(y[(py*imgW+px)*2:])
			cbRaw := binary.LittleEndian.Uint16(u[((py/hDiv)*cW+px/wDiv)*2:])
			crRaw := binary.LittleEndian.Uint16(v[((py/hDiv)*cW+px/wDiv)*2:])
			yv := uint8(math.Round(float64(yRaw) / maxVal * 255))
			cb := uint8(math.Round(float64(cbRaw) / maxVal * 255))
			cr := uint8(math.Round(float64(crRaw) / maxVal * 255))
			r, g, b := yCbCrToRGB(yv, cb, cr)
			out.SetRGBA(px, py, color.RGBA{r, g, b, 255})
		}
	}
	return out
}

// NV semiplanar: Y plane + interleaved UV (or VU) plane
func encodeNV(src *image.RGBA, swapUV bool, hDiv int) (yPlane, uvPlane []byte) {
	yPlane = make([]byte, imgW*imgH)
	cW, cH := imgW/2, imgH/hDiv
	uvPlane = make([]byte, cW*cH*2)
	for py := 0; py < imgH; py++ {
		for px := 0; px < imgW; px++ {
			c := src.RGBAAt(px, py)
			yv, _, _ := rgbToYCbCr(c.R, c.G, c.B)
			yPlane[py*imgW+px] = yv
		}
	}
	for cy := 0; cy < cH; cy++ {
		for cx := 0; cx < cW; cx++ {
			var sumCb, sumCr float64
			count := 0
			for dy := 0; dy < hDiv; dy++ {
				for dx := 0; dx < 2; dx++ {
					px2, py2 := cx*2+dx, cy*hDiv+dy
					if px2 < imgW && py2 < imgH {
						c := src.RGBAAt(px2, py2)
						_, cb, cr := rgbToYCbCr(c.R, c.G, c.B)
						sumCb += float64(cb)
						sumCr += float64(cr)
						count++
					}
				}
			}
			if count > 0 {
				cb := uint8(math.Round(sumCb / float64(count)))
				cr := uint8(math.Round(sumCr / float64(count)))
				idx := (cy*cW + cx) * 2
				if swapUV {
					uvPlane[idx], uvPlane[idx+1] = cr, cb
				} else {
					uvPlane[idx], uvPlane[idx+1] = cb, cr
				}
			}
		}
	}
	return
}

func decodeNV(yPlane, uvPlane []byte, swapUV bool, hDiv int) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	cW := imgW / 2
	for py := 0; py < imgH; py++ {
		for px := 0; px < imgW; px++ {
			yv := yPlane[py*imgW+px]
			idx := ((py/hDiv)*cW + px/2) * 2
			u0, u1 := uvPlane[idx], uvPlane[idx+1]
			var cb, cr uint8
			if swapUV {
				cr, cb = u0, u1
			} else {
				cb, cr = u0, u1
			}
			r, g, b := yCbCrToRGB(yv, cb, cr)
			out.SetRGBA(px, py, color.RGBA{r, g, b, 255})
		}
	}
	return out
}

func nvEntry(label string, swapUV bool, hDiv int) entry {
	return entry{label, func(src *image.RGBA) *image.RGBA {
		y, uv := encodeNV(src, swapUV, hDiv)
		return decodeNV(y, uv, swapUV, hDiv)
	}}
}

// Hi-bit NV semiplanar (P010/P016 style): Y as LUMINANCE_ALPHA, UV as RGBA
func nvHiBitEntry(label string, bits int) entry {
	return entry{label, func(src *image.RGBA) *image.RGBA {
		maxValInt := (1 << uint(bits)) - 1
		maxVal := float64(maxValInt)
		scale := maxVal / 255.0
		// Y plane: 2 bytes per pixel (LE uint16), LUMINANCE_ALPHA layout
		yp := make([]byte, imgW*imgH*2)
		// UV plane: 4 bytes per UV pair (U_lo,U_hi,V_lo,V_hi), at half-width
		cW := imgW / 2
		uvp := make([]byte, cW*(imgH/2)*4)
		for py := 0; py < imgH; py++ {
			for px := 0; px < imgW; px++ {
				c := src.RGBAAt(px, py)
				yv, _, _ := rgbToYCbCr(c.R, c.G, c.B)
				binary.LittleEndian.PutUint16(yp[(py*imgW+px)*2:], uint16(math.Round(float64(yv)*scale)))
			}
		}
		for cy := 0; cy < imgH/2; cy++ {
			for cx := 0; cx < cW; cx++ {
				var sumCb, sumCr float64
				for dy := 0; dy < 2; dy++ {
					for dx := 0; dx < 2; dx++ {
						c := src.RGBAAt(cx*2+dx, cy*2+dy)
						_, cb, cr := rgbToYCbCr(c.R, c.G, c.B)
						sumCb += float64(cb)
						sumCr += float64(cr)
					}
				}
				cbHi := uint16(math.Round(sumCb / 4 * scale))
				crHi := uint16(math.Round(sumCr / 4 * scale))
				idx := (cy*cW + cx) * 4
				binary.LittleEndian.PutUint16(uvp[idx+0:], cbHi)
				binary.LittleEndian.PutUint16(uvp[idx+2:], crHi)
			}
		}
		// Decode back
		out := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
		for py := 0; py < imgH; py++ {
			for px := 0; px < imgW; px++ {
				yRaw := binary.LittleEndian.Uint16(yp[(py*imgW+px)*2:])
				uvIdx := ((py/2)*cW + px/2) * 4
				cbRaw := binary.LittleEndian.Uint16(uvp[uvIdx+0:])
				crRaw := binary.LittleEndian.Uint16(uvp[uvIdx+2:])
				yv := uint8(math.Round(float64(yRaw) / maxVal * 255))
				cb := uint8(math.Round(float64(cbRaw) / maxVal * 255))
				cr := uint8(math.Round(float64(crRaw) / maxVal * 255))
				r, g, b := yCbCrToRGB(yv, cb, cr)
				out.SetRGBA(px, py, color.RGBA{r, g, b, 255})
			}
		}
		return out
	}}
}

// Packed YUV422: mode 0=YUYV, 1=UYVY, 2=YVYU
func encodePackedYUV422(src *image.RGBA, mode int) *incanvas.RawFrame {
	// Width must be even; RGBA stride → RGBA at half-width
	data := make([]byte, imgW/2*imgH*4)
	for py := 0; py < imgH; py++ {
		for bx := 0; bx < imgW/2; bx++ {
			c0 := src.RGBAAt(bx*2, py)
			c1 := src.RGBAAt(bx*2+1, py)
			y0, _, _ := rgbToYCbCr(c0.R, c0.G, c0.B)
			y1, _, _ := rgbToYCbCr(c1.R, c1.G, c1.B)
			_, cb, cr := rgbToYCbCr(
				(uint8((int(c0.R)+int(c1.R))/2)),
				(uint8((int(c0.G)+int(c1.G))/2)),
				(uint8((int(c0.B)+int(c1.B))/2)),
			)
			idx := (py*(imgW/2) + bx) * 4
			switch mode {
			case 0: // YUYV
				data[idx+0], data[idx+1], data[idx+2], data[idx+3] = y0, cb, y1, cr
			case 1: // UYVY
				data[idx+0], data[idx+1], data[idx+2], data[idx+3] = cb, y0, cr, y1
			case 2: // YVYU
				data[idx+0], data[idx+1], data[idx+2], data[idx+3] = y0, cr, y1, cb
			}
		}
	}
	return &incanvas.RawFrame{Width: imgW / 2, Height: imgH, Data: [4][]byte{data}}
}

func decodePackedYUV422(data []byte, mode int) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	for py := 0; py < imgH; py++ {
		for bx := 0; bx < imgW/2; bx++ {
			idx := (py*(imgW/2) + bx) * 4
			var y0, y1, cb, cr uint8
			switch mode {
			case 0: // YUYV
				y0, cb, y1, cr = data[idx+0], data[idx+1], data[idx+2], data[idx+3]
			case 1: // UYVY
				cb, y0, cr, y1 = data[idx+0], data[idx+1], data[idx+2], data[idx+3]
			case 2: // YVYU
				y0, cr, y1, cb = data[idx+0], data[idx+1], data[idx+2], data[idx+3]
			}
			r0, g0, b0 := yCbCrToRGB(y0, cb, cr)
			r1, g1, b1 := yCbCrToRGB(y1, cb, cr)
			out.SetRGBA(bx*2, py, color.RGBA{r0, g0, b0, 255})
			out.SetRGBA(bx*2+1, py, color.RGBA{r1, g1, b1, 255})
		}
	}
	return out
}

func packedYUVEntry(label string, mode int) entry {
	return entry{label, func(src *image.RGBA) *image.RGBA {
		frame := encodePackedYUV422(src, mode)
		return decodePackedYUV422(frame.Data[0], mode)
	}}
}

// GRAY8 / GRAY16 / YA8
func grayEntry(label string, bits int, withAlpha bool) entry {
	return entry{label, func(src *image.RGBA) *image.RGBA {
		out := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
		for py := 0; py < imgH; py++ {
			for px := 0; px < imgW; px++ {
				c := src.RGBAAt(px, py)
				luma := lumaOf(c.R, c.G, c.B)
				if bits == 16 {
					// Round-trip through 16-bit: scale up then back down
					luma16 := uint16(luma) * 257 // 0…255 → 0…65535
					luma = uint8(luma16 / 257)
				}
				alpha := c.A
				if !withAlpha {
					alpha = 255
				}
				out.SetRGBA(px, py, color.RGBA{luma, luma, luma, alpha})
			}
		}
		return out
	}}
}

// ─────────────────────────────────────────────────────────────────────────────
// Full format table
// ─────────────────────────────────────────────────────────────────────────────

func allEntries() []entry {
	cpuRT := func(fmt int, enc func(*image.RGBA) *incanvas.RawFrame) func(*image.RGBA) *image.RGBA {
		return func(src *image.RGBA) *image.RGBA {
			return cpuRoundTrip(fmt, enc(src))
		}
	}

	return []entry{
		// ── Original ──────────────────────────────────────────────────────
		{"Original\n(RGBA)", func(src *image.RGBA) *image.RGBA { return src }},

		// ── Packed RGB ────────────────────────────────────────────────────
		{"RGB24", cpuRT(incanvas.PixelFormatRGB24, encodeRGB24)},
		{"BGR24", cpuRT(incanvas.PixelFormatBGR24, encodeBGR24)},
		{"BGRA", cpuRT(incanvas.PixelFormatBGRA, encodeBGRA)},
		{"ARGB", cpuRT(incanvas.PixelFormatARGB, encodeARGB)},
		{"ABGR", cpuRT(incanvas.PixelFormatABGR, encodeABGR)},

		// ── Packed RGB 16-bit ─────────────────────────────────────────────
		{"RGB565\n(5+6+5)", cpuRT(incanvas.PixelFormatRGB565, encodeRGB565)},
		{"BGR565\n(5+6+5)", cpuRT(incanvas.PixelFormatBGR565, encodeBGR565)},

		// ── Planar GBR ────────────────────────────────────────────────────
		{"GBRP\n8-bit", cpuRT(incanvas.PixelFormatGBRP, func(s *image.RGBA) *incanvas.RawFrame { return encodeGBRP(s, 8) })},
		{"GBRP\n10-bit", cpuRT(incanvas.PixelFormatGBRP10, func(s *image.RGBA) *incanvas.RawFrame { return encodeGBRP(s, 10) })},
		{"GBRP\n12-bit", cpuRT(incanvas.PixelFormatGBRP12, func(s *image.RGBA) *incanvas.RawFrame { return encodeGBRP(s, 12) })},
		{"GBRP\n16-bit", cpuRT(incanvas.PixelFormatGBRP16, func(s *image.RGBA) *incanvas.RawFrame { return encodeGBRP(s, 16) })},
		{"GBRAP\n(+alpha)", cpuRT(incanvas.PixelFormatGBRAP, encodeGBRAP)},

		// ── Monochrome ────────────────────────────────────────────────────
		{"MONOWHITE\n(1-bit)", cpuRT(incanvas.PixelFormatMONOWHITE, func(s *image.RGBA) *incanvas.RawFrame { return encodeMono(s, true) })},
		{"MONOBLACK\n(1-bit)", cpuRT(incanvas.PixelFormatMONOBLACK, func(s *image.RGBA) *incanvas.RawFrame { return encodeMono(s, false) })},

		// ── Bayer ─────────────────────────────────────────────────────────
		{"Bayer\nBGGR8", cpuRT(incanvas.PixelFormatBAYER_BGGR8, func(s *image.RGBA) *incanvas.RawFrame { return encodeBayer(s, 0) })},
		{"Bayer\nRGGB8", cpuRT(incanvas.PixelFormatBAYER_RGGB8, func(s *image.RGBA) *incanvas.RawFrame { return encodeBayer(s, 1) })},
		{"Bayer\nGBRG8", cpuRT(incanvas.PixelFormatBAYER_GBRG8, func(s *image.RGBA) *incanvas.RawFrame { return encodeBayer(s, 2) })},
		{"Bayer\nGRBG8", cpuRT(incanvas.PixelFormatBAYER_GRBG8, func(s *image.RGBA) *incanvas.RawFrame { return encodeBayer(s, 3) })},

		// ── Float ─────────────────────────────────────────────────────────
		{"GRAYF32\n(float)", cpuRT(incanvas.PixelFormatGRAYF32, encodeGrayF32)},
		{"RGBF32\n(float)", cpuRT(incanvas.PixelFormatRGBF32, encodeRGBF32)},
		{"RGBAF32\n(float)", cpuRT(incanvas.PixelFormatRGBAF32, encodeRGBAF32)},

		// ── YUV planar (GPU shader category) ─────────────────────────────
		yuvPlanarEntry("YUV420P\n(4:2:0)", 2, 2),
		yuvPlanarEntry("YUV422P\n(4:2:2)", 2, 1),
		yuvPlanarEntry("YUV444P\n(4:4:4)", 1, 1),
		yuvPlanarEntry("YUV440P\n(4:4:0)", 1, 2),
		yuvPlanarEntry("YUV411P\n(4:1:1)", 4, 1),
		yuvPlanarEntry("YUV410P\n(4:1:0)", 4, 4),

		// ── YUVA planar ───────────────────────────────────────────────────
		yuvaDecodeEntry("YUVA420P\n(+alpha)", 2, 2),
		yuvaDecodeEntry("YUVA422P\n(+alpha)", 2, 1),
		yuvaDecodeEntry("YUVA444P\n(+alpha)", 1, 1),

		// ── YUV hi-bit planar ─────────────────────────────────────────────
		hibitPlanarEntry("YUV420P10\n(10-bit)", 10, 2, 2),
		hibitPlanarEntry("YUV420P12\n(12-bit)", 12, 2, 2),
		hibitPlanarEntry("YUV420P16\n(16-bit)", 16, 2, 2),
		hibitPlanarEntry("YUV444P10\n(10-bit)", 10, 1, 1),
		hibitPlanarEntry("YUV444P16\n(16-bit)", 16, 1, 1),

		// ── NV semiplanar ─────────────────────────────────────────────────
		nvEntry("NV12\n(Y+UV)", false, 2),
		nvEntry("NV21\n(Y+VU)", true, 2),
		nvEntry("NV16\n(Y+UV 4:2:2)", false, 1),
		nvEntry("NV24\n(Y+UV 4:4:4)", false, 1), // cW=imgW/2 still, but hDiv=1

		// ── Hi-bit NV ─────────────────────────────────────────────────────
		nvHiBitEntry("P010\n(10-bit NV)", 10),
		nvHiBitEntry("P016\n(16-bit NV)", 16),

		// ── Packed YUV422 ─────────────────────────────────────────────────
		packedYUVEntry("YUYV422", 0),
		packedYUVEntry("UYVY422", 1),
		packedYUVEntry("YVYU422", 2),

		// ── Grayscale (GPU shader category) ──────────────────────────────
		grayEntry("GRAY8", 8, false),
		grayEntry("GRAY16", 16, false),
		grayEntry("YA8\n(gray+α)", 8, true),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// UI helpers
// ─────────────────────────────────────────────────────────────────────────────

const thumbSize = float32(72)
const cols = 7

func makeCell(e entry, base *image.RGBA) fyne.CanvasObject {
	converted := e.convert(base)

	img := canvas.NewImageFromImage(converted)
	img.FillMode = canvas.ImageFillOriginal
	img.SetMinSize(fyne.NewSize(thumbSize, thumbSize))

	lbl := widget.NewLabel(e.label)
	lbl.Alignment = fyne.TextAlignCenter
	lbl.Wrapping = fyne.TextWrapOff

	return container.NewVBox(
		container.NewCenter(img),
		lbl,
	)
}

// ─────────────────────────────────────────────────────────────────────────────
// main
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	a := app.New()
	w := a.NewWindow("StreamingImage – Pixel Format Showcase")

	base := makeTestImage()
	entries := allEntries()

	grid := container.NewAdaptiveGrid(cols)
	for _, e := range entries {
		grid.Add(makeCell(e, base))
	}

	scroll := container.NewScroll(grid)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(float32(cols)*110, 700))
	w.ShowAndRun()
}
