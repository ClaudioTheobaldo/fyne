//go:build ((gles || arm || arm64) && !android && !ios && !mobile && !darwin && !wasm && !test_web_driver) || ((android || ios || mobile) && (!wasm || !test_web_driver)) || wasm || test_web_driver

package gl

import _ "embed"

var (
	//go:embed shaders/line_es.frag
	shaderLineesFrag []byte

	//go:embed shaders/line_es.vert
	shaderLineesVert []byte

	//go:embed shaders/rectangle_es.frag
	shaderRectangleesFrag []byte

	//go:embed shaders/rectangle_es.vert
	shaderRectangleesVert []byte

	//go:embed shaders/round_rectangle_es.frag
	shaderRoundrectangleesFrag []byte

	//go:embed shaders/simple_es.frag
	shaderSimpleesFrag []byte

	//go:embed shaders/simple_es.vert
	shaderSimpleesVert []byte

	//go:embed shaders/polygon_es.frag
	shaderPolygonesFrag []byte

	//go:embed shaders/arc_es.frag
	shaderArcesFrag []byte

	//go:embed shaders/yuv_planar_es.frag
	shaderYuvPlanaresFrag []byte

	//go:embed shaders/yuva_planar_es.frag
	shaderYuvaPlanaresFrag []byte

	//go:embed shaders/nv_semiplanar_es.frag
	shaderNvSemiplanaresFrag []byte

	//go:embed shaders/packed_yuv422_es.frag
	shaderPackedYuv422esFrag []byte

	//go:embed shaders/grayscale_es.frag
	shaderGrayscaleesFrag []byte

	//go:embed shaders/yuv_planar_hibit_es.frag
	shaderYuvPlanarHibitesFrag []byte

	//go:embed shaders/nv_semiplanar_hibit_es.frag
	shaderNvSemiplanarHibitesFrag []byte

	//go:embed shaders/effect_passthrough_es.vert
	shaderEffectPassthroughesVert []byte

	//go:embed shaders/effect_passthrough_es.frag
	shaderEffectPassthroughesFrag []byte
)

func isGLES() bool { return true }

func shaderRectVertexSrc() []byte { return shaderRectangleesVert }

func effectPassthroughVertSrc() []byte { return shaderEffectPassthroughesVert }
func effectPassthroughFragSrc() []byte { return shaderEffectPassthroughesFrag }

func shaderSourceNamed(name string) ([]byte, []byte) {
	switch name {
	case "line_es":
		return shaderLineesVert, shaderLineesFrag
	case "simple_es":
		return shaderSimpleesVert, shaderSimpleesFrag
	case "rectangle_es":
		return shaderRectangleesVert, shaderRectangleesFrag
	case "round_rectangle_es":
		return shaderRectangleesVert, shaderRoundrectangleesFrag
	case "polygon_es":
		return shaderRectangleesVert, shaderPolygonesFrag
	case "arc_es":
		return shaderRectangleesVert, shaderArcesFrag
	case "yuv_planar_es":
		return shaderSimpleesVert, shaderYuvPlanaresFrag
	case "yuva_planar_es":
		return shaderSimpleesVert, shaderYuvaPlanaresFrag
	case "nv_semiplanar_es":
		return shaderSimpleesVert, shaderNvSemiplanaresFrag
	case "packed_yuv422_es":
		return shaderSimpleesVert, shaderPackedYuv422esFrag
	case "grayscale_es":
		return shaderSimpleesVert, shaderGrayscaleesFrag
	case "yuv_planar_hibit_es":
		return shaderSimpleesVert, shaderYuvPlanarHibitesFrag
	case "nv_semiplanar_hibit_es":
		return shaderSimpleesVert, shaderNvSemiplanarHibitesFrag
	}
	return nil, nil
}
