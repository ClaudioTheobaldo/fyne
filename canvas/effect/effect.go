// Package effect provides composable visual effects for Fyne canvas objects.
//
// Effects are applied to any canvas primitive via AddEffect and processed
// through an FBO-based shader pipeline in the GL painter. Multiple effects
// stack in insertion order and can be animated via the Animate method.
//
// Since: 2.8
package effect

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// EffectType identifies a built-in or custom shader effect.
type EffectType int

const (
	// Custom allows a user-supplied GLSL fragment shader.
	Custom EffectType = iota

	// Color effects

	Brightness  // Multiplies color values to brighten/darken
	Contrast    // Adjusts difference between light and dark
	Saturate    // Increases or decreases color saturation
	Grayscale   // Converts to grayscale via luminance weighting
	HueRotate   // Rotates hue on the color wheel
	Sepia       // Applies warm brownish tone via color matrix
	Invert      // Inverts color values (1.0 - color)
	Opacity     // Adjusts alpha transparency
	ColorMatrix // Full 5x4 color matrix transform
	Posterize   // Reduces color levels for poster-like appearance
	Threshold   // Binary black/white at a cutoff value
	Gamma       // Power curve adjustment
	Vibrance    // Selective saturation boost (less saturated colors boosted more)
	Temperature // Warm/cool color shift

	// Blur effects

	GaussianBlur   // Two-pass separable Gaussian blur
	BoxBlur        // Simple NxN average blur
	DirectionalBlur // Motion blur along a direction vector
	ZoomBlur       // Radial blur from center point
	BilateralBlur  // Edge-preserving blur

	// Distortion effects

	Ripple             // Sine-wave UV displacement
	Swirl              // Rotational UV distortion
	Barrel             // Barrel/pincushion lens distortion
	Pixelate           // Quantize UV for blocky pixel effect
	Spherize           // Sphere-mapped UV distortion
	Fisheye            // Fisheye lens distortion
	ChromaticAberration // Per-channel UV offset simulating lens defect
	RGBShift           // Random block displacement and channel shifts
	Frosted            // Blur with noise-displaced sampling
	TiltShift          // Position-dependent selective focus blur

	// Edge and detail effects

	Sharpen    // Increases edge contrast via convolution
	Emboss     // Raised surface simulation via convolution
	EdgeDetect // Sobel edge detection
	Convolve   // Arbitrary NxN convolution kernel

	// Blend modes

	BlendMultiply    // Multiplies src and dst; darkens
	BlendScreen      // Inverse multiply; lightens
	BlendOverlay     // Multiply or screen depending on dst brightness
	BlendDarken      // Per-channel minimum
	BlendLighten     // Per-channel maximum
	BlendColorDodge  // Brightens dst by dividing by inverse of src
	BlendColorBurn   // Darkens dst by dividing inverse by src
	BlendHardLight   // Overlay with src/dst swapped
	BlendSoftLight   // Gentler version of hard-light
	BlendDifference  // Absolute difference
	BlendExclusion   // Like difference but lower contrast
	BlendHue         // Hue from src, sat+lum from dst
	BlendSaturation  // Saturation from src, hue+lum from dst
	BlendColor       // Hue+sat from src, lum from dst
	BlendLuminosity  // Lum from src, hue+sat from dst
	BlendPlusLighter // Additive blending

	// Compositing

	ArithmeticComposite // k1*A*B + k2*A + k3*B + k4

	// Shadow effects

	DropShadow  // Blurred offset shadow behind element
	BoxShadow   // Shadow around element box
	InnerShadow // Shadow inside element
	OuterGlow   // Blurred glow outside element

	// Stylization effects

	Vignette    // Radial edge darkening
	FilmGrain   // Noise overlay
	Scanlines   // Horizontal line pattern
	CRT         // CRT monitor simulation
	Halftone    // Dot pattern by luminance
	DotMatrix   // Grid of circles
	OilPainting // Kuwahara filter

	// Backdrop effects (read behind the element)

	BackdropBlur       // Frosted glass on background content
	BackdropBrightness // Brighten/darken backdrop

	// Lighting effects

	DiffuseLight  // Lambertian diffuse using alpha as height map
	SpecularLight // Phong specular highlights on height map
	AmbientLight  // Constant ambient light contribution

	// Gradient overlay effects

	LinearGradientOverlay // Linear color gradient composited on top
	RadialGradientOverlay // Radial gradient composited on top
	ConicGradientOverlay  // Conic/angular gradient composited on top

	// Mask effects

	AlphaMask    // Mask via alpha channel of a texture
	LuminanceMask // Mask via luminance of a texture
	GradientMask // Mask via a generated gradient

	// Procedural effects

	Turbulence // Perlin noise / fractal noise generation
	Bloom      // Threshold bright pixels, blur, add back

	// Advanced effects

	ColorLUT   // 3D lookup table for cinematic color grading
	ToneMapping // HDR to SDR conversion
	FXAA       // Fast approximate anti-aliasing

	effectTypeCount // sentinel for iteration
)

// Effect represents a single shader effect instance applied to a canvas object.
// Each call to AddEffect returns a unique Effect handle that can be used to
// modify parameters, animate values, or remove the effect.
type Effect struct {
	kind    EffectType
	params  []float32
	order   int
	owner   fyne.CanvasObject
	enabled bool

	// uniforms holds user-set uniform values keyed by name.
	// Values can be: float32, [2]float32, [3]float32, [4]float32,
	// int32, bool, [9]float32 (mat3), [16]float32 (mat4).
	uniforms map[string]any

	// customShaderSrc holds the fragment shader source for Custom effects.
	customShaderSrc string

	mu sync.RWMutex
}

// Type returns the effect type.
func (e *Effect) Type() EffectType {
	return e.kind
}

// Params returns a copy of the constructor parameters.
func (e *Effect) Params() []float32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]float32, len(e.params))
	copy(out, e.params)
	return out
}

// Enabled returns whether this effect is active.
func (e *Effect) Enabled() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enabled
}

// SetEnabled toggles the effect on or off without removing it from the chain.
func (e *Effect) SetEnabled(enabled bool) {
	e.mu.Lock()
	e.enabled = enabled
	e.mu.Unlock()
	e.refresh()
}

// SetFloat sets a named float uniform.
func (e *Effect) SetFloat(name string, v float32) {
	e.mu.Lock()
	e.uniforms[name] = v
	e.mu.Unlock()
	e.refresh()
}

// SetVec2 sets a named vec2 uniform.
func (e *Effect) SetVec2(name string, v0, v1 float32) {
	e.mu.Lock()
	e.uniforms[name] = [2]float32{v0, v1}
	e.mu.Unlock()
	e.refresh()
}

// SetVec3 sets a named vec3 uniform.
func (e *Effect) SetVec3(name string, v0, v1, v2 float32) {
	e.mu.Lock()
	e.uniforms[name] = [3]float32{v0, v1, v2}
	e.mu.Unlock()
	e.refresh()
}

// SetVec4 sets a named vec4 uniform.
func (e *Effect) SetVec4(name string, v0, v1, v2, v3 float32) {
	e.mu.Lock()
	e.uniforms[name] = [4]float32{v0, v1, v2, v3}
	e.mu.Unlock()
	e.refresh()
}

// SetInt sets a named integer uniform.
func (e *Effect) SetInt(name string, v int32) {
	e.mu.Lock()
	e.uniforms[name] = v
	e.mu.Unlock()
	e.refresh()
}

// SetBool sets a named boolean uniform.
func (e *Effect) SetBool(name string, v bool) {
	e.mu.Lock()
	e.uniforms[name] = v
	e.mu.Unlock()
	e.refresh()
}

// SetMat3 sets a named 3x3 matrix uniform (9 floats, column-major).
func (e *Effect) SetMat3(name string, m [9]float32) {
	e.mu.Lock()
	e.uniforms[name] = m
	e.mu.Unlock()
	e.refresh()
}

// SetMat4 sets a named 4x4 matrix uniform (16 floats, column-major).
func (e *Effect) SetMat4(name string, m [16]float32) {
	e.mu.Lock()
	e.uniforms[name] = m
	e.mu.Unlock()
	e.refresh()
}

// Uniform returns the current value for a named uniform, or nil if not set.
func (e *Effect) Uniform(name string) any {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.uniforms[name]
}

// Uniforms returns a snapshot of all uniform values.
func (e *Effect) Uniforms() map[string]any {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make(map[string]any, len(e.uniforms))
	for k, v := range e.uniforms {
		out[k] = v
	}
	return out
}

// CustomShaderSrc returns the fragment shader source for Custom effects.
func (e *Effect) CustomShaderSrc() string {
	return e.customShaderSrc
}

// Order returns the insertion order of this effect in the chain.
func (e *Effect) Order() int {
	return e.order
}

// Animate creates and starts an animation that interpolates a float uniform
// from `from` to `to` over the given duration using the specified easing curve.
// Returns the started animation which can be stopped if needed.
func (e *Effect) Animate(name string, from, to float32, duration time.Duration, curve fyne.AnimationCurve) *fyne.Animation {
	anim := fyne.NewAnimation(duration, func(progress float32) {
		value := from + (to-from)*progress
		e.SetFloat(name, value)
	})
	anim.Curve = curve
	anim.Start()
	return anim
}

func (e *Effect) refresh() {
	if e.owner != nil {
		e.owner.Refresh()
	}
}

// NewEffect creates a new Effect with the given type, parameters, and insertion order.
// The owner is set separately via SetOwner after the effect is added to a canvas object.
func NewEffect(kind EffectType, params []float32, order int) *Effect {
	e := &Effect{
		kind:     kind,
		params:   params,
		order:    order,
		enabled:  true,
		uniforms: make(map[string]any),
	}
	ApplyDefaults(e)
	return e
}

// SetOwner sets the canvas object that owns this effect (used for refresh callbacks).
func SetOwner(e *Effect, owner fyne.CanvasObject) {
	e.owner = owner
}

// SetOrder sets the insertion order of this effect in the chain.
func SetOrder(e *Effect, order int) {
	e.order = order
}

// NewCustomEffect creates a new Effect with a user-supplied fragment shader source.
func NewCustomEffect(shaderSrc string, uniforms map[string][]float32) *Effect {
	e := &Effect{
		kind:            Custom,
		enabled:         true,
		uniforms:        make(map[string]any),
		customShaderSrc: shaderSrc,
	}
	for k, v := range uniforms {
		switch len(v) {
		case 1:
			e.uniforms[k] = v[0]
		case 2:
			e.uniforms[k] = [2]float32{v[0], v[1]}
		case 3:
			e.uniforms[k] = [3]float32{v[0], v[1], v[2]}
		case 4:
			e.uniforms[k] = [4]float32{v[0], v[1], v[2], v[3]}
		}
	}
	return e
}
