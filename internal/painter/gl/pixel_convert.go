package gl

import (
	"encoding/binary"
	"image"
	"math"

	"fyne.io/fyne/v2/canvas"
)

// ConvertRawFrameToRGBA is the exported entry point for CPU-side pixel format
// conversion. It wraps convertRawToRGBA and is intended for testing and showcase use.
func ConvertRawFrameToRGBA(pixelFormat int, frame *canvas.RawFrame) *image.RGBA {
	return convertRawToRGBA(pixelFormat, frame)
}

// convertRawToRGBA converts a RawFrame in a CPU-side pixel format to *image.RGBA.
// Returns nil if the pixel format is unsupported or the frame data is empty.
func convertRawToRGBA(pixelFormat int, frame *canvas.RawFrame) *image.RGBA {
	if len(frame.Data[0]) == 0 {
		return nil
	}
	out := image.NewRGBA(image.Rect(0, 0, frame.Width, frame.Height))

	switch pixelFormat {
	case canvas.PixelFormatRGB24:
		convertRGB24ToRGBA(frame, out)
	case canvas.PixelFormatBGR24:
		convertBGR24ToRGBA(frame, out)
	case canvas.PixelFormatBGRA:
		convertBGRAToRGBA(frame, out)
	case canvas.PixelFormatARGB:
		convertARGBToRGBA(frame, out)
	case canvas.PixelFormatABGR:
		convertABGRToRGBA(frame, out)
	case canvas.PixelFormatRGB565:
		convertRGB565ToRGBA(frame, out)
	case canvas.PixelFormatBGR565:
		convertBGR565ToRGBA(frame, out)
	case canvas.PixelFormatGBRP:
		convertGBRPToRGBA(frame, out, 8, false)
	case canvas.PixelFormatGBRP10:
		convertGBRPToRGBA(frame, out, 10, false)
	case canvas.PixelFormatGBRP12:
		convertGBRPToRGBA(frame, out, 12, false)
	case canvas.PixelFormatGBRP16:
		convertGBRPToRGBA(frame, out, 16, false)
	case canvas.PixelFormatGBRAP:
		convertGBRPToRGBA(frame, out, 8, true)
	case canvas.PixelFormatMONOWHITE:
		convertMonoToRGBA(frame, out, true)
	case canvas.PixelFormatMONOBLACK:
		convertMonoToRGBA(frame, out, false)
	case canvas.PixelFormatBAYER_BGGR8:
		convertBayerToRGBA(frame, out, bayerBGGR, 8)
	case canvas.PixelFormatBAYER_RGGB8:
		convertBayerToRGBA(frame, out, bayerRGGB, 8)
	case canvas.PixelFormatBAYER_GBRG8:
		convertBayerToRGBA(frame, out, bayerGBRG, 8)
	case canvas.PixelFormatBAYER_GRBG8:
		convertBayerToRGBA(frame, out, bayerGRBG, 8)
	case canvas.PixelFormatBAYER_BGGR16:
		convertBayerToRGBA(frame, out, bayerBGGR, 16)
	case canvas.PixelFormatBAYER_RGGB16:
		convertBayerToRGBA(frame, out, bayerRGGB, 16)
	case canvas.PixelFormatBAYER_GBRG16:
		convertBayerToRGBA(frame, out, bayerGBRG, 16)
	case canvas.PixelFormatBAYER_GRBG16:
		convertBayerToRGBA(frame, out, bayerGRBG, 16)
	case canvas.PixelFormatGRAYF32:
		convertGrayF32ToRGBA(frame, out)
	case canvas.PixelFormatRGBF32:
		convertRGBF32ToRGBA(frame, out)
	case canvas.PixelFormatRGBAF32:
		convertRGBAF32ToRGBA(frame, out)
	default:
		return nil
	}
	return out
}

func convertRGB24ToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 3
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			dRow[x*4+0] = sRow[x*3+0] // R
			dRow[x*4+1] = sRow[x*3+1] // G
			dRow[x*4+2] = sRow[x*3+2] // B
			dRow[x*4+3] = 255
		}
	}
}

func convertBGR24ToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 3
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			dRow[x*4+0] = sRow[x*3+2] // R (from B position)
			dRow[x*4+1] = sRow[x*3+1] // G
			dRow[x*4+2] = sRow[x*3+0] // B (from R position)
			dRow[x*4+3] = 255
		}
	}
}

func convertBGRAToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 4
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			dRow[x*4+0] = sRow[x*4+2] // R
			dRow[x*4+1] = sRow[x*4+1] // G
			dRow[x*4+2] = sRow[x*4+0] // B
			dRow[x*4+3] = sRow[x*4+3] // A
		}
	}
}

func convertARGBToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 4
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			dRow[x*4+0] = sRow[x*4+1] // R
			dRow[x*4+1] = sRow[x*4+2] // G
			dRow[x*4+2] = sRow[x*4+3] // B
			dRow[x*4+3] = sRow[x*4+0] // A
		}
	}
}

func convertABGRToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 4
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			dRow[x*4+0] = sRow[x*4+3] // R (ABGR: A=0,B=1,G=2,R=3)
			dRow[x*4+1] = sRow[x*4+2] // G
			dRow[x*4+2] = sRow[x*4+1] // B
			dRow[x*4+3] = sRow[x*4+0] // A
		}
	}
}

func convertRGB565ToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 2
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			v := binary.LittleEndian.Uint16(sRow[x*2:])
			r5 := uint8((v >> 11) & 0x1F)
			g6 := uint8((v >> 5) & 0x3F)
			b5 := uint8(v & 0x1F)
			dRow[x*4+0] = (r5 << 3) | (r5 >> 2)
			dRow[x*4+1] = (g6 << 2) | (g6 >> 4)
			dRow[x*4+2] = (b5 << 3) | (b5 >> 2)
			dRow[x*4+3] = 255
		}
	}
}

func convertBGR565ToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 2
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			v := binary.LittleEndian.Uint16(sRow[x*2:])
			b5 := uint8((v >> 11) & 0x1F) // MSBs are blue in BGR565
			g6 := uint8((v >> 5) & 0x3F)
			r5 := uint8(v & 0x1F) // LSBs are red in BGR565
			dRow[x*4+0] = (r5 << 3) | (r5 >> 2)
			dRow[x*4+1] = (g6 << 2) | (g6 >> 4)
			dRow[x*4+2] = (b5 << 3) | (b5 >> 2)
			dRow[x*4+3] = 255
		}
	}
}

// convertGBRPToRGBA converts planar GBR(A) to RGBA.
// Planes: Data[0]=G, Data[1]=B, Data[2]=R, Data[3]=A (if hasAlpha).
// bitDepth 8 uses 1 byte/sample; 10/12/14/16 use 2 bytes/sample (little-endian uint16).
func convertGBRPToRGBA(frame *canvas.RawFrame, out *image.RGBA, bitDepth int, hasAlpha bool) {
	g, b, r := frame.Data[0], frame.Data[1], frame.Data[2]
	gStr, bStr, rStr := frame.Strides[0], frame.Strides[1], frame.Strides[2]
	w, h := frame.Width, frame.Height

	if bitDepth == 8 {
		if gStr == 0 {
			gStr = w
		}
		if bStr == 0 {
			bStr = w
		}
		if rStr == 0 {
			rStr = w
		}
		var a []byte
		var aStr int
		if hasAlpha && len(frame.Data[3]) > 0 {
			a = frame.Data[3]
			aStr = frame.Strides[3]
			if aStr == 0 {
				aStr = w
			}
		}
		for y := 0; y < h; y++ {
			dRow := out.Pix[y*out.Stride:]
			gRow, bRow, rRow := g[y*gStr:], b[y*bStr:], r[y*rStr:]
			for x := 0; x < w; x++ {
				alpha := uint8(255)
				if a != nil {
					alpha = a[y*aStr+x]
				}
				dRow[x*4+0] = rRow[x]
				dRow[x*4+1] = gRow[x]
				dRow[x*4+2] = bRow[x]
				dRow[x*4+3] = alpha
			}
		}
	} else {
		bpt := 2 // bytes per sample
		if gStr == 0 {
			gStr = w * bpt
		}
		if bStr == 0 {
			bStr = w * bpt
		}
		if rStr == 0 {
			rStr = w * bpt
		}
		maxSampleInt := (1 << uint(bitDepth)) - 1
		maxSample := float64(maxSampleInt)
		for y := 0; y < h; y++ {
			dRow := out.Pix[y*out.Stride:]
			gRow, bRow, rRow := g[y*gStr:], b[y*bStr:], r[y*rStr:]
			for x := 0; x < w; x++ {
				rv := float64(binary.LittleEndian.Uint16(rRow[x*bpt:])) / maxSample
				gv := float64(binary.LittleEndian.Uint16(gRow[x*bpt:])) / maxSample
				bv := float64(binary.LittleEndian.Uint16(bRow[x*bpt:])) / maxSample
				dRow[x*4+0] = uint8(math.Round(rv * 255))
				dRow[x*4+1] = uint8(math.Round(gv * 255))
				dRow[x*4+2] = uint8(math.Round(bv * 255))
				dRow[x*4+3] = 255
			}
		}
	}
}

// convertMonoToRGBA converts packed 1-bit monochrome data to RGBA.
// Pixels are packed MSB-first, 8 pixels per byte.
// whiteIsZero=true: AV_PIX_FMT_MONOWHITE (0=white, 1=black).
// whiteIsZero=false: AV_PIX_FMT_MONOBLACK (0=black, 1=white).
func convertMonoToRGBA(frame *canvas.RawFrame, out *image.RGBA, whiteIsZero bool) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = (frame.Width + 7) / 8
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			bit := (sRow[x/8] >> uint(7-x%8)) & 1
			var gray uint8
			if (bit == 0) == whiteIsZero {
				gray = 255
			}
			dRow[x*4+0] = gray
			dRow[x*4+1] = gray
			dRow[x*4+2] = gray
			dRow[x*4+3] = 255
		}
	}
}

type bayerPattern int

const (
	bayerBGGR bayerPattern = iota
	bayerRGGB
	bayerGBRG
	bayerGRBG
)

// convertBayerToRGBA demosaics a Bayer-pattern image to RGBA using bilinear interpolation.
func convertBayerToRGBA(frame *canvas.RawFrame, out *image.RGBA, pat bayerPattern, bitDepth int) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	w, h := frame.Width, frame.Height
	bpt := 1
	if bitDepth > 8 {
		bpt = 2
	}
	if stride == 0 {
		stride = w * bpt
	}
	maxVal := float64(255)
	if bitDepth > 8 {
		maxValInt := (1 << uint(bitDepth)) - 1
		maxVal = float64(maxValInt)
	}

	// sample returns the normalized [0,1] value at (x,y), clamped at edges.
	sample := func(x, y int) float64 {
		if x < 0 {
			x = 0
		} else if x >= w {
			x = w - 1
		}
		if y < 0 {
			y = 0
		} else if y >= h {
			y = h - 1
		}
		if bpt == 1 {
			return float64(src[y*stride+x]) / maxVal
		}
		return float64(binary.LittleEndian.Uint16(src[y*stride+x*2:])) / maxVal
	}

	// colorAt returns which color (R/G/B) is at Bayer grid position (cx, cy) where cx=x%2, cy=y%2.
	colorAt := func(cx, cy int) (isR, isG, isB bool) {
		idx := cy*2 + cx
		switch pat {
		case bayerBGGR: // B G / G R
			switch idx {
			case 0:
				return false, false, true
			case 1, 2:
				return false, true, false
			default:
				return true, false, false
			}
		case bayerRGGB: // R G / G B
			switch idx {
			case 0:
				return true, false, false
			case 1, 2:
				return false, true, false
			default:
				return false, false, true
			}
		case bayerGBRG: // G B / R G
			switch idx {
			case 0, 3:
				return false, true, false
			case 1:
				return false, false, true
			default:
				return true, false, false
			}
		default: // bayerGRBG: G R / B G
			switch idx {
			case 0, 3:
				return false, true, false
			case 1:
				return true, false, false
			default:
				return false, false, true
			}
		}
	}

	for y := 0; y < h; y++ {
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < w; x++ {
			isR, _, isB := colorAt(x%2, y%2)
			s := sample(x, y)
			var rv, gv, bv float64
			switch {
			case isR:
				rv = s
				gv = (sample(x-1, y) + sample(x+1, y) + sample(x, y-1) + sample(x, y+1)) / 4
				bv = (sample(x-1, y-1) + sample(x+1, y-1) + sample(x-1, y+1) + sample(x+1, y+1)) / 4
			case isB:
				bv = s
				gv = (sample(x-1, y) + sample(x+1, y) + sample(x, y-1) + sample(x, y+1)) / 4
				rv = (sample(x-1, y-1) + sample(x+1, y-1) + sample(x-1, y+1) + sample(x+1, y+1)) / 4
			default: // isG
				gv = s
				// Determine which axis holds R vs B based on neighbors
				rNeighbor, _, _ := colorAt((x+1)%2, y%2)
				if rNeighbor {
					rv = (sample(x-1, y) + sample(x+1, y)) / 2
					bv = (sample(x, y-1) + sample(x, y+1)) / 2
				} else {
					bv = (sample(x-1, y) + sample(x+1, y)) / 2
					rv = (sample(x, y-1) + sample(x, y+1)) / 2
				}
			}
			dRow[x*4+0] = uint8(math.Round(math.Min(rv, 1) * 255))
			dRow[x*4+1] = uint8(math.Round(math.Min(gv, 1) * 255))
			dRow[x*4+2] = uint8(math.Round(math.Min(bv, 1) * 255))
			dRow[x*4+3] = 255
		}
	}
}

func convertGrayF32ToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 4
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			f := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*4:]))
			v := uint8(math.Round(math.Max(0, math.Min(float64(f), 1)) * 255))
			dRow[x*4+0] = v
			dRow[x*4+1] = v
			dRow[x*4+2] = v
			dRow[x*4+3] = 255
		}
	}
}

func convertRGBF32ToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 12 // 3 floats × 4 bytes
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			rf := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*12+0:]))
			gf := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*12+4:]))
			bf := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*12+8:]))
			dRow[x*4+0] = clampF32(rf)
			dRow[x*4+1] = clampF32(gf)
			dRow[x*4+2] = clampF32(bf)
			dRow[x*4+3] = 255
		}
	}
}

func convertRGBAF32ToRGBA(frame *canvas.RawFrame, out *image.RGBA) {
	src := frame.Data[0]
	stride := frame.Strides[0]
	if stride == 0 {
		stride = frame.Width * 16 // 4 floats × 4 bytes
	}
	for y := 0; y < frame.Height; y++ {
		sRow := src[y*stride:]
		dRow := out.Pix[y*out.Stride:]
		for x := 0; x < frame.Width; x++ {
			rf := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*16+0:]))
			gf := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*16+4:]))
			bf := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*16+8:]))
			af := math.Float32frombits(binary.LittleEndian.Uint32(sRow[x*16+12:]))
			dRow[x*4+0] = clampF32(rf)
			dRow[x*4+1] = clampF32(gf)
			dRow[x*4+2] = clampF32(bf)
			dRow[x*4+3] = clampF32(af)
		}
	}
}

// clampF32 converts a float32 in [0,1] to uint8, clamping out-of-range values.
func clampF32(v float32) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 255
	}
	return uint8(math.Round(float64(v) * 255))
}
