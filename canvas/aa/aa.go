// Package aa defines anti-aliasing mode constants and helpers used by the
// Fyne renderer to select the active AA technique.
package aa

// Mode describes the anti-aliasing technique applied by the GL painter.
type Mode int

const (
	// Off disables all anti-aliasing (no SDF edge softness, no FXAA, no MSAA).
	Off Mode = iota
	// SDF enables signed-distance-field edge smoothing on all canvas primitives.
	// This is the default mode and provides smooth edges via smoothstep in the
	// fragment shaders (controlled by the edgeSoftness / feather uniforms).
	SDF
	// FXAA enables SDF edge smoothing plus a post-process Fast Approximate
	// Anti-Aliasing pass on canvas objects that opt in via the effect system.
	FXAA
	// MSAA2X requests 2x hardware multisample anti-aliasing from the window
	// system (GLFW Samples hint) and enables GL_MULTISAMPLE.  Desktop GL only;
	// falls back to SDF on GLES / WASM / mobile.
	MSAA2X
	// MSAA4X requests 4x hardware multisample anti-aliasing.  Same platform
	// constraints as MSAA2X.
	MSAA4X
)

// String returns the human-readable name of the AA mode.
func (m Mode) String() string {
	switch m {
	case Off:
		return "OFF"
	case SDF:
		return "SDF"
	case FXAA:
		return "FXAA"
	case MSAA2X:
		return "MSAA_2X"
	case MSAA4X:
		return "MSAA_4X"
	default:
		return "SDF"
	}
}

// Samples returns the MSAA sample count for modes that use hardware
// multisampling, or 0 for non-MSAA modes.
func (m Mode) Samples() int {
	switch m {
	case MSAA2X:
		return 2
	case MSAA4X:
		return 4
	default:
		return 0
	}
}

// NeedsMSAA reports whether the mode requires hardware multisampling.
func (m Mode) NeedsMSAA() bool {
	return m == MSAA2X || m == MSAA4X
}

// NeedsFXAA reports whether the mode requires a post-process FXAA pass.
func (m Mode) NeedsFXAA() bool {
	return m == FXAA
}

// NeedsSDF reports whether the mode uses SDF edge smoothing on primitives.
// Every mode except Off uses SDF.
func (m Mode) NeedsSDF() bool {
	return m != Off
}

// FromString parses a string (as persisted in config) into a Mode.
// Unrecognised values default to SDF.
func FromString(s string) Mode {
	switch s {
	case "OFF":
		return Off
	case "SDF":
		return SDF
	case "FXAA":
		return FXAA
	case "MSAA_2X":
		return MSAA2X
	case "MSAA_4X":
		return MSAA4X
	default:
		return SDF
	}
}
