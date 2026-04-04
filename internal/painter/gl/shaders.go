//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

package gl

import _ "embed"

var (
	//go:embed shaders/line.frag
	shaderLineFrag []byte

	//go:embed shaders/line.vert
	shaderLineVert []byte

	//go:embed shaders/rectangle.frag
	shaderRectangleFrag []byte

	//go:embed shaders/rectangle.vert
	shaderRectangleVert []byte

	//go:embed shaders/round_rectangle.frag
	shaderRoundrectangleFrag []byte

	//go:embed shaders/simple.frag
	shaderSimpleFrag []byte

	//go:embed shaders/simple.vert
	shaderSimpleVert []byte

	//go:embed shaders/polygon.frag
	shaderPolygonFrag []byte

	//go:embed shaders/arc.frag
	shaderArcFrag []byte

	//go:embed shaders/yuv_planar.frag
	shaderYuvPlanarFrag []byte

	//go:embed shaders/yuva_planar.frag
	shaderYuvaPlanarFrag []byte

	//go:embed shaders/nv_semiplanar.frag
	shaderNvSemiplanarFrag []byte

	//go:embed shaders/packed_yuv422.frag
	shaderPackedYuv422Frag []byte

	//go:embed shaders/grayscale.frag
	shaderGrayscaleFrag []byte

	//go:embed shaders/yuv_planar_hibit.frag
	shaderYuvPlanarHibitFrag []byte

	//go:embed shaders/nv_semiplanar_hibit.frag
	shaderNvSemiplanarHibitFrag []byte

	//go:embed shaders/effect_passthrough.vert
	shaderEffectPassthroughVert []byte

	//go:embed shaders/effect_passthrough.frag
	shaderEffectPassthroughFrag []byte
)

func isGLES() bool { return false }

func shaderRectVertexSrc() []byte { return shaderRectangleVert }

func effectPassthroughVertSrc() []byte { return shaderEffectPassthroughVert }
func effectPassthroughFragSrc() []byte { return shaderEffectPassthroughFrag }

func shaderSourceNamed(name string) ([]byte, []byte) {
	switch name {
	case "line":
		return shaderLineVert, shaderLineFrag
	case "simple":
		return shaderSimpleVert, shaderSimpleFrag
	case "rectangle":
		return shaderRectangleVert, shaderRectangleFrag
	case "round_rectangle":
		return shaderRectangleVert, shaderRoundrectangleFrag
	case "polygon":
		return shaderRectangleVert, shaderPolygonFrag
	case "arc":
		return shaderRectangleVert, shaderArcFrag
	case "yuv_planar":
		return shaderSimpleVert, shaderYuvPlanarFrag
	case "yuva_planar":
		return shaderSimpleVert, shaderYuvaPlanarFrag
	case "nv_semiplanar":
		return shaderSimpleVert, shaderNvSemiplanarFrag
	case "packed_yuv422":
		return shaderSimpleVert, shaderPackedYuv422Frag
	case "grayscale":
		return shaderSimpleVert, shaderGrayscaleFrag
	case "yuv_planar_hibit":
		return shaderSimpleVert, shaderYuvPlanarHibitFrag
	case "nv_semiplanar_hibit":
		return shaderSimpleVert, shaderNvSemiplanarHibitFrag
	}
	return nil, nil
}
