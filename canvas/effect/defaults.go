package effect

// UniformMapping describes how constructor params map to shader uniform names.
type UniformMapping struct {
	Name string // Uniform name in the shader
	Size int    // Number of float32 values (1, 2, 3, or 4)
}

// DefaultUniforms returns the uniform mappings for a given effect type.
// Each entry maps a constructor parameter (by position in the flat params slice)
// to a shader uniform name and its component count.
func DefaultUniforms(kind EffectType) []UniformMapping {
	switch kind {

	// Color effects
	case Brightness:
		return []UniformMapping{{"brightness", 1}}
	case Contrast:
		return []UniformMapping{{"contrast", 1}}
	case Saturate:
		return []UniformMapping{{"saturation", 1}}
	case Grayscale:
		return []UniformMapping{{"amount", 1}}
	case HueRotate:
		return []UniformMapping{{"angle", 1}}
	case Sepia:
		return []UniformMapping{{"amount", 1}}
	case Invert:
		return []UniformMapping{{"amount", 1}}
	case Opacity:
		return []UniformMapping{{"opacity", 1}}
	case Posterize:
		return []UniformMapping{{"levels", 1}}
	case Threshold:
		return []UniformMapping{{"threshold", 1}}
	case Gamma:
		return []UniformMapping{{"gamma", 1}}
	case Vibrance:
		return []UniformMapping{{"vibrance", 1}}
	case Temperature:
		return []UniformMapping{{"temperature", 1}}

	// Blur effects
	case GaussianBlur:
		return []UniformMapping{{"radius", 1}}
	case BoxBlur:
		return []UniformMapping{{"radius", 1}}
	case DirectionalBlur:
		return []UniformMapping{{"radius", 1}, {"direction", 2}}
	case ZoomBlur:
		return []UniformMapping{{"strength", 1}, {"center", 2}}
	case BilateralBlur:
		return []UniformMapping{{"radius", 1}, {"sigmaSpace", 1}, {"sigmaColor", 1}}

	// Distortion effects
	case Ripple:
		return []UniformMapping{{"amplitude", 1}, {"frequency", 1}, {"speed", 1}}
	case Swirl:
		return []UniformMapping{{"radius", 1}, {"angle", 1}, {"center", 2}}
	case Barrel:
		return []UniformMapping{{"distortion", 1}}
	case Pixelate:
		return []UniformMapping{{"pixelSize", 1}}
	case Spherize:
		return []UniformMapping{{"radius", 1}, {"center", 2}}
	case Fisheye:
		return []UniformMapping{{"strength", 1}}
	case ChromaticAberration:
		return []UniformMapping{{"offset", 1}}
	case RGBShift:
		return []UniformMapping{{"amount", 1}, {"angle", 1}}
	case Frosted:
		return []UniformMapping{{"radius", 1}, {"noiseScale", 1}}
	case TiltShift:
		return []UniformMapping{{"radius", 1}, {"focusY", 1}, {"focusWidth", 1}}
	case Skew:
		return []UniformMapping{{"skewX", 1}, {"skewY", 1}}
	case Rotate:
		return []UniformMapping{{"angle", 1}}
	case Scale:
		return []UniformMapping{{"scaleX", 1}, {"scaleY", 1}}
	case PerspectiveTransform:
		return []UniformMapping{{"amount", 1}, {"direction", 1}}
	case Wave:
		return []UniformMapping{{"amplitude", 1}, {"frequency", 1}, {"phase", 1}, {"direction", 1}}
	case PageCurl:
		return []UniformMapping{{"curl", 1}, {"radius", 1}}
	case Mirror:
		return []UniformMapping{{"flipX", 1}, {"flipY", 1}}
	case Kaleidoscope:
		return []UniformMapping{{"segments", 1}, {"rotation", 1}}

	// Edge/Detail effects
	case Sharpen:
		return []UniformMapping{{"strength", 1}}
	case Emboss:
		return []UniformMapping{{"strength", 1}}
	case EdgeDetect:
		return []UniformMapping{{"threshold", 1}}
	case Convolve:
		// 9 values for 3x3 kernel + divisor
		return []UniformMapping{{"kernel", 4}, {"kernel2", 4}, {"kernelExtra", 1}, {"divisor", 1}}

	// Shadow effects
	case DropShadow:
		return []UniformMapping{{"offsetX", 1}, {"offsetY", 1}, {"blur", 1}, {"shadowColor", 4}}
	case BoxShadow:
		return []UniformMapping{{"offsetX", 1}, {"offsetY", 1}, {"blur", 1}, {"spread", 1}, {"shadowColor", 4}}
	case InnerShadow:
		return []UniformMapping{{"offsetX", 1}, {"offsetY", 1}, {"blur", 1}, {"shadowColor", 4}}
	case OuterGlow:
		return []UniformMapping{{"radius", 1}, {"glowColor", 4}}

	// Stylization effects
	case Vignette:
		return []UniformMapping{{"intensity", 1}, {"smoothness", 1}}
	case FilmGrain:
		return []UniformMapping{{"intensity", 1}, {"speed", 1}}
	case Scanlines:
		return []UniformMapping{{"density", 1}, {"opacity", 1}}
	case CRT:
		return []UniformMapping{{"curvature", 1}, {"scanlineIntensity", 1}, {"vignetteStrength", 1}}
	case Halftone:
		return []UniformMapping{{"dotSize", 1}, {"angle", 1}}
	case DotMatrix:
		return []UniformMapping{{"dotSize", 1}, {"spacing", 1}}
	case OilPainting:
		return []UniformMapping{{"radius", 1}, {"levels", 1}}

	// Backdrop effects
	case BackdropBlur:
		return []UniformMapping{{"radius", 1}}
	case BackdropBrightness:
		return []UniformMapping{{"brightness", 1}}

	// Lighting effects
	case DiffuseLight:
		return []UniformMapping{{"lightPos", 3}, {"diffuseConstant", 1}, {"surfaceScale", 1}}
	case SpecularLight:
		return []UniformMapping{{"lightPos", 3}, {"specularConstant", 1}, {"specularExponent", 1}, {"surfaceScale", 1}}
	case AmbientLight:
		return []UniformMapping{{"color", 3}, {"intensity", 1}}

	// Gradient overlay effects
	case LinearGradientOverlay:
		return []UniformMapping{{"startColor", 4}, {"endColor", 4}, {"angle", 1}}
	case RadialGradientOverlay:
		return []UniformMapping{{"startColor", 4}, {"endColor", 4}, {"center", 2}}
	case ConicGradientOverlay:
		return []UniformMapping{{"startColor", 4}, {"endColor", 4}, {"center", 2}, {"angle", 1}}

	// Mask effects
	case AlphaMask:
		return nil // mask texture set via SetTexture
	case LuminanceMask:
		return nil // mask texture set via SetTexture
	case GradientMask:
		return []UniformMapping{{"startAlpha", 1}, {"endAlpha", 1}, {"angle", 1}}

	// Procedural effects
	case Turbulence:
		return []UniformMapping{{"baseFrequency", 2}, {"numOctaves", 1}, {"seed", 1}}
	case Bloom:
		return []UniformMapping{{"threshold", 1}, {"intensity", 1}, {"radius", 1}}

	// Advanced effects
	case ColorLUT:
		return nil // LUT texture set via SetTexture
	case ToneMapping:
		return []UniformMapping{{"exposure", 1}, {"gamma", 1}}
	case FXAA:
		return nil // no params, uses texelSize

	// Blend modes — all take a blend opacity parameter
	case BlendMultiply, BlendScreen, BlendOverlay, BlendDarken, BlendLighten,
		BlendColorDodge, BlendColorBurn, BlendHardLight, BlendSoftLight,
		BlendDifference, BlendExclusion, BlendHue, BlendSaturation,
		BlendColor, BlendLuminosity, BlendPlusLighter:
		return []UniformMapping{{"blendOpacity", 1}}

	case ArithmeticComposite:
		return []UniformMapping{{"k1", 1}, {"k2", 1}, {"k3", 1}, {"k4", 1}}

	case ColorMatrix:
		// 20 values for 5x4 matrix, passed as mat4 + extra vec4
		return []UniformMapping{{"matrix", 4}, {"matrix2", 4}, {"matrix3", 4}, {"matrix4", 4}, {"matrixOffset", 4}}
	}

	return nil
}

// PassCount returns the number of shader passes required for an effect type.
// Most effects are single-pass; Gaussian blur is two-pass (horizontal + vertical).
func PassCount(kind EffectType) int {
	switch kind {
	case GaussianBlur:
		return 2
	case Bloom:
		return 3 // threshold + horizontal blur + vertical blur
	}
	return 1
}

// ApplyDefaults sets the default uniform values on an Effect based on its
// constructor params and the DefaultUniforms mapping.
func ApplyDefaults(e *Effect) {
	mappings := DefaultUniforms(e.kind)
	if mappings == nil || len(e.params) == 0 {
		return
	}

	paramIdx := 0
	for _, m := range mappings {
		if paramIdx >= len(e.params) {
			break
		}
		switch m.Size {
		case 1:
			e.uniforms[m.Name] = e.params[paramIdx]
			paramIdx++
		case 2:
			if paramIdx+1 < len(e.params) {
				e.uniforms[m.Name] = [2]float32{e.params[paramIdx], e.params[paramIdx+1]}
				paramIdx += 2
			}
		case 3:
			if paramIdx+2 < len(e.params) {
				e.uniforms[m.Name] = [3]float32{e.params[paramIdx], e.params[paramIdx+1], e.params[paramIdx+2]}
				paramIdx += 3
			}
		case 4:
			if paramIdx+3 < len(e.params) {
				e.uniforms[m.Name] = [4]float32{e.params[paramIdx], e.params[paramIdx+1], e.params[paramIdx+2], e.params[paramIdx+3]}
				paramIdx += 4
			}
		}
	}
}
