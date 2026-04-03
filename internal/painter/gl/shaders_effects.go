//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

package gl

import (
	_ "embed"

	"fyne.io/fyne/v2/canvas/effect"
)

// --- Color effects ---

//go:embed shaders/effects/brightness.frag
var effectBrightnessFrag []byte

//go:embed shaders/effects/contrast.frag
var effectContrastFrag []byte

//go:embed shaders/effects/grayscale.frag
var effectGrayscaleFrag []byte

//go:embed shaders/effects/sepia.frag
var effectSepiaFrag []byte

//go:embed shaders/effects/invert.frag
var effectInvertFrag []byte

//go:embed shaders/effects/opacity.frag
var effectOpacityFrag []byte

//go:embed shaders/effects/hue_rotate.frag
var effectHueRotateFrag []byte

//go:embed shaders/effects/saturate.frag
var effectSaturateFrag []byte

//go:embed shaders/effects/color_matrix.frag
var effectColorMatrixFrag []byte

//go:embed shaders/effects/posterize.frag
var effectPosterizeFrag []byte

//go:embed shaders/effects/threshold.frag
var effectThresholdFrag []byte

//go:embed shaders/effects/gamma.frag
var effectGammaFrag []byte

//go:embed shaders/effects/vibrance.frag
var effectVibranceFrag []byte

//go:embed shaders/effects/temperature.frag
var effectTemperatureFrag []byte

// --- Blur effects ---

//go:embed shaders/effects/gaussian_blur.frag
var effectGaussianBlurFrag []byte

//go:embed shaders/effects/box_blur.frag
var effectBoxBlurFrag []byte

//go:embed shaders/effects/directional_blur.frag
var effectDirectionalBlurFrag []byte

//go:embed shaders/effects/zoom_blur.frag
var effectZoomBlurFrag []byte

//go:embed shaders/effects/bilateral_blur.frag
var effectBilateralBlurFrag []byte

// --- Distortion effects ---

//go:embed shaders/effects/pixelate.frag
var effectPixelateFrag []byte

//go:embed shaders/effects/ripple.frag
var effectRippleFrag []byte

//go:embed shaders/effects/swirl.frag
var effectSwirlFrag []byte

//go:embed shaders/effects/barrel.frag
var effectBarrelFrag []byte

//go:embed shaders/effects/spherize.frag
var effectSpherizeFrag []byte

//go:embed shaders/effects/fisheye.frag
var effectFisheyeFrag []byte

//go:embed shaders/effects/chromatic_aberration.frag
var effectChromaticAberrationFrag []byte

//go:embed shaders/effects/rgb_shift.frag
var effectRGBShiftFrag []byte

//go:embed shaders/effects/frosted.frag
var effectFrostedFrag []byte

//go:embed shaders/effects/tilt_shift.frag
var effectTiltShiftFrag []byte

//go:embed shaders/effects/skew.frag
var effectSkewFrag []byte

//go:embed shaders/effects/rotate.frag
var effectRotateFrag []byte

//go:embed shaders/effects/scale.frag
var effectScaleFrag []byte

//go:embed shaders/effects/perspective_transform.frag
var effectPerspectiveTransformFrag []byte

//go:embed shaders/effects/wave.frag
var effectWaveFrag []byte

//go:embed shaders/effects/page_curl.frag
var effectPageCurlFrag []byte

//go:embed shaders/effects/mirror.frag
var effectMirrorFrag []byte

//go:embed shaders/effects/kaleidoscope.frag
var effectKaleidoscopeFrag []byte

// --- Detail effects ---

//go:embed shaders/effects/sharpen.frag
var effectSharpenFrag []byte

//go:embed shaders/effects/emboss.frag
var effectEmbossFrag []byte

//go:embed shaders/effects/edge_detect.frag
var effectEdgeDetectFrag []byte

// --- Shadow effects ---

//go:embed shaders/effects/drop_shadow.frag
var effectDropShadowFrag []byte

//go:embed shaders/effects/box_shadow.frag
var effectBoxShadowFrag []byte

//go:embed shaders/effects/inner_shadow.frag
var effectInnerShadowFrag []byte

//go:embed shaders/effects/outer_glow.frag
var effectOuterGlowFrag []byte

// --- Stylization effects ---

//go:embed shaders/effects/vignette.frag
var effectVignetteFrag []byte

//go:embed shaders/effects/film_grain.frag
var effectFilmGrainFrag []byte

//go:embed shaders/effects/scanlines.frag
var effectScanlinesFrag []byte

//go:embed shaders/effects/crt.frag
var effectCRTFrag []byte

//go:embed shaders/effects/halftone.frag
var effectHalftoneFrag []byte

//go:embed shaders/effects/dot_matrix.frag
var effectDotMatrixFrag []byte

//go:embed shaders/effects/oil_painting.frag
var effectOilPaintingFrag []byte

// --- Lighting effects ---

//go:embed shaders/effects/diffuse_light.frag
var effectDiffuseLightFrag []byte

//go:embed shaders/effects/specular_light.frag
var effectSpecularLightFrag []byte

//go:embed shaders/effects/ambient_light.frag
var effectAmbientLightFrag []byte

// --- Gradient overlay effects ---

//go:embed shaders/effects/linear_gradient_overlay.frag
var effectLinearGradientOverlayFrag []byte

//go:embed shaders/effects/radial_gradient_overlay.frag
var effectRadialGradientOverlayFrag []byte

//go:embed shaders/effects/conic_gradient_overlay.frag
var effectConicGradientOverlayFrag []byte

// --- Mask effects ---

//go:embed shaders/effects/gradient_mask.frag
var effectGradientMaskFrag []byte

// --- Procedural effects ---

//go:embed shaders/effects/turbulence.frag
var effectTurbulenceFrag []byte

//go:embed shaders/effects/bloom.frag
var effectBloomFrag []byte

// --- Advanced effects ---

//go:embed shaders/effects/tone_mapping.frag
var effectToneMappingFrag []byte

//go:embed shaders/effects/fxaa.frag
var effectFXAAFrag []byte

// --- Blend effects ---

//go:embed shaders/effects/blend_multiply.frag
var effectBlendMultiplyFrag []byte

//go:embed shaders/effects/blend_screen.frag
var effectBlendScreenFrag []byte

// effectShaderSource returns the fragment shader source for a given effect type.
// Returns nil if the effect type has no built-in shader.
func effectShaderSource(kind effect.EffectType) []byte {
	switch kind {
	// Color
	case effect.Brightness:
		return effectBrightnessFrag
	case effect.Contrast:
		return effectContrastFrag
	case effect.Grayscale:
		return effectGrayscaleFrag
	case effect.Sepia:
		return effectSepiaFrag
	case effect.Invert:
		return effectInvertFrag
	case effect.Opacity:
		return effectOpacityFrag
	case effect.HueRotate:
		return effectHueRotateFrag
	case effect.Saturate:
		return effectSaturateFrag
	case effect.ColorMatrix:
		return effectColorMatrixFrag
	case effect.Posterize:
		return effectPosterizeFrag
	case effect.Threshold:
		return effectThresholdFrag
	case effect.Gamma:
		return effectGammaFrag
	case effect.Vibrance:
		return effectVibranceFrag
	case effect.Temperature:
		return effectTemperatureFrag

	// Blur
	case effect.GaussianBlur:
		return effectGaussianBlurFrag
	case effect.BoxBlur:
		return effectBoxBlurFrag
	case effect.DirectionalBlur:
		return effectDirectionalBlurFrag
	case effect.ZoomBlur:
		return effectZoomBlurFrag
	case effect.BilateralBlur:
		return effectBilateralBlurFrag

	// Distortion
	case effect.Pixelate:
		return effectPixelateFrag
	case effect.Ripple:
		return effectRippleFrag
	case effect.Swirl:
		return effectSwirlFrag
	case effect.Barrel:
		return effectBarrelFrag
	case effect.Spherize:
		return effectSpherizeFrag
	case effect.Fisheye:
		return effectFisheyeFrag
	case effect.ChromaticAberration:
		return effectChromaticAberrationFrag
	case effect.RGBShift:
		return effectRGBShiftFrag
	case effect.Frosted:
		return effectFrostedFrag
	case effect.TiltShift:
		return effectTiltShiftFrag
	case effect.Skew:
		return effectSkewFrag
	case effect.Rotate:
		return effectRotateFrag
	case effect.Scale:
		return effectScaleFrag
	case effect.PerspectiveTransform:
		return effectPerspectiveTransformFrag
	case effect.Wave:
		return effectWaveFrag
	case effect.PageCurl:
		return effectPageCurlFrag
	case effect.Mirror:
		return effectMirrorFrag
	case effect.Kaleidoscope:
		return effectKaleidoscopeFrag

	// Detail
	case effect.Sharpen:
		return effectSharpenFrag
	case effect.Emboss:
		return effectEmbossFrag
	case effect.EdgeDetect:
		return effectEdgeDetectFrag

	// Shadow
	case effect.DropShadow:
		return effectDropShadowFrag
	case effect.BoxShadow:
		return effectBoxShadowFrag
	case effect.InnerShadow:
		return effectInnerShadowFrag
	case effect.OuterGlow:
		return effectOuterGlowFrag

	// Stylization
	case effect.Vignette:
		return effectVignetteFrag
	case effect.FilmGrain:
		return effectFilmGrainFrag
	case effect.Scanlines:
		return effectScanlinesFrag
	case effect.CRT:
		return effectCRTFrag
	case effect.Halftone:
		return effectHalftoneFrag
	case effect.DotMatrix:
		return effectDotMatrixFrag
	case effect.OilPainting:
		return effectOilPaintingFrag

	// Lighting
	case effect.DiffuseLight:
		return effectDiffuseLightFrag
	case effect.SpecularLight:
		return effectSpecularLightFrag
	case effect.AmbientLight:
		return effectAmbientLightFrag

	// Gradient overlay
	case effect.LinearGradientOverlay:
		return effectLinearGradientOverlayFrag
	case effect.RadialGradientOverlay:
		return effectRadialGradientOverlayFrag
	case effect.ConicGradientOverlay:
		return effectConicGradientOverlayFrag

	// Mask
	case effect.GradientMask:
		return effectGradientMaskFrag

	// Procedural
	case effect.Turbulence:
		return effectTurbulenceFrag
	case effect.Bloom:
		return effectBloomFrag

	// Advanced
	case effect.ToneMapping:
		return effectToneMappingFrag
	case effect.FXAA:
		return effectFXAAFrag

	// Blend
	case effect.BlendMultiply:
		return effectBlendMultiplyFrag
	case effect.BlendScreen:
		return effectBlendScreenFrag
	}
	return nil
}
