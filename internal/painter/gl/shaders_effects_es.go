//go:build ((gles || arm || arm64) && !android && !ios && !mobile && !darwin && !wasm && !test_web_driver) || ((android || ios || mobile) && (!wasm || !test_web_driver)) || wasm || test_web_driver

package gl

import (
	_ "embed"

	"fyne.io/fyne/v2/canvas/effect"
)

// --- Color effects ---

//go:embed shaders/effects/brightness_es.frag
var effectBrightnessFrag []byte

//go:embed shaders/effects/contrast_es.frag
var effectContrastFrag []byte

//go:embed shaders/effects/grayscale_es.frag
var effectGrayscaleFrag []byte

//go:embed shaders/effects/sepia_es.frag
var effectSepiaFrag []byte

//go:embed shaders/effects/invert_es.frag
var effectInvertFrag []byte

//go:embed shaders/effects/opacity_es.frag
var effectOpacityFrag []byte

//go:embed shaders/effects/hue_rotate_es.frag
var effectHueRotateFrag []byte

//go:embed shaders/effects/saturate_es.frag
var effectSaturateFrag []byte

//go:embed shaders/effects/color_matrix_es.frag
var effectColorMatrixFrag []byte

//go:embed shaders/effects/posterize_es.frag
var effectPosterizeFrag []byte

//go:embed shaders/effects/threshold_es.frag
var effectThresholdFrag []byte

//go:embed shaders/effects/gamma_es.frag
var effectGammaFrag []byte

//go:embed shaders/effects/vibrance_es.frag
var effectVibranceFrag []byte

//go:embed shaders/effects/temperature_es.frag
var effectTemperatureFrag []byte

// --- Blur effects ---

//go:embed shaders/effects/gaussian_blur_es.frag
var effectGaussianBlurFrag []byte

//go:embed shaders/effects/box_blur_es.frag
var effectBoxBlurFrag []byte

//go:embed shaders/effects/directional_blur_es.frag
var effectDirectionalBlurFrag []byte

//go:embed shaders/effects/zoom_blur_es.frag
var effectZoomBlurFrag []byte

//go:embed shaders/effects/bilateral_blur_es.frag
var effectBilateralBlurFrag []byte

// --- Distortion effects ---

//go:embed shaders/effects/pixelate_es.frag
var effectPixelateFrag []byte

//go:embed shaders/effects/ripple_es.frag
var effectRippleFrag []byte

//go:embed shaders/effects/swirl_es.frag
var effectSwirlFrag []byte

//go:embed shaders/effects/barrel_es.frag
var effectBarrelFrag []byte

//go:embed shaders/effects/spherize_es.frag
var effectSpherizeFrag []byte

//go:embed shaders/effects/fisheye_es.frag
var effectFisheyeFrag []byte

//go:embed shaders/effects/chromatic_aberration_es.frag
var effectChromaticAberrationFrag []byte

//go:embed shaders/effects/rgb_shift_es.frag
var effectRGBShiftFrag []byte

//go:embed shaders/effects/frosted_es.frag
var effectFrostedFrag []byte

//go:embed shaders/effects/tilt_shift_es.frag
var effectTiltShiftFrag []byte

//go:embed shaders/effects/skew_es.frag
var effectSkewFrag []byte

//go:embed shaders/effects/rotate_es.frag
var effectRotateFrag []byte

//go:embed shaders/effects/scale_es.frag
var effectScaleFrag []byte

//go:embed shaders/effects/perspective_transform_es.frag
var effectPerspectiveTransformFrag []byte

//go:embed shaders/effects/wave_es.frag
var effectWaveFrag []byte

//go:embed shaders/effects/page_curl_es.frag
var effectPageCurlFrag []byte

//go:embed shaders/effects/mirror_es.frag
var effectMirrorFrag []byte

//go:embed shaders/effects/kaleidoscope_es.frag
var effectKaleidoscopeFrag []byte

// --- Detail effects ---

//go:embed shaders/effects/sharpen_es.frag
var effectSharpenFrag []byte

//go:embed shaders/effects/emboss_es.frag
var effectEmbossFrag []byte

//go:embed shaders/effects/edge_detect_es.frag
var effectEdgeDetectFrag []byte

// --- Shadow effects ---

//go:embed shaders/effects/drop_shadow_es.frag
var effectDropShadowFrag []byte

//go:embed shaders/effects/box_shadow_es.frag
var effectBoxShadowFrag []byte

//go:embed shaders/effects/inner_shadow_es.frag
var effectInnerShadowFrag []byte

//go:embed shaders/effects/outer_glow_es.frag
var effectOuterGlowFrag []byte

// --- Stylization effects ---

//go:embed shaders/effects/vignette_es.frag
var effectVignetteFrag []byte

//go:embed shaders/effects/film_grain_es.frag
var effectFilmGrainFrag []byte

//go:embed shaders/effects/scanlines_es.frag
var effectScanlinesFrag []byte

//go:embed shaders/effects/crt_es.frag
var effectCRTFrag []byte

//go:embed shaders/effects/halftone_es.frag
var effectHalftoneFrag []byte

//go:embed shaders/effects/dot_matrix_es.frag
var effectDotMatrixFrag []byte

//go:embed shaders/effects/oil_painting_es.frag
var effectOilPaintingFrag []byte

// --- Lighting effects ---

//go:embed shaders/effects/diffuse_light_es.frag
var effectDiffuseLightFrag []byte

//go:embed shaders/effects/specular_light_es.frag
var effectSpecularLightFrag []byte

//go:embed shaders/effects/ambient_light_es.frag
var effectAmbientLightFrag []byte

// --- Gradient overlay effects ---

//go:embed shaders/effects/linear_gradient_overlay_es.frag
var effectLinearGradientOverlayFrag []byte

//go:embed shaders/effects/radial_gradient_overlay_es.frag
var effectRadialGradientOverlayFrag []byte

//go:embed shaders/effects/conic_gradient_overlay_es.frag
var effectConicGradientOverlayFrag []byte

// --- Mask effects ---

//go:embed shaders/effects/gradient_mask_es.frag
var effectGradientMaskFrag []byte

// --- Procedural effects ---

//go:embed shaders/effects/turbulence_es.frag
var effectTurbulenceFrag []byte

//go:embed shaders/effects/bloom_es.frag
var effectBloomFrag []byte

// --- Advanced effects ---

//go:embed shaders/effects/tone_mapping_es.frag
var effectToneMappingFrag []byte

//go:embed shaders/effects/fxaa_es.frag
var effectFXAAFrag []byte

// --- Blend effects ---

//go:embed shaders/effects/blend_multiply_es.frag
var effectBlendMultiplyFrag []byte

//go:embed shaders/effects/blend_screen_es.frag
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
