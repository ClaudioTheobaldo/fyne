package effect

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEffect(t *testing.T) {
	e := NewEffect(GaussianBlur, []float32{4.0}, 0)

	assert.Equal(t, GaussianBlur, e.Type())
	assert.Equal(t, []float32{4.0}, e.Params())
	assert.Equal(t, 0, e.Order())
	assert.True(t, e.Enabled())
}

func TestNewEffect_AppliesDefaults(t *testing.T) {
	e := NewEffect(GaussianBlur, []float32{4.0}, 0)

	// Default mapping: params[0] -> "radius"
	v := e.Uniform("radius")
	require.NotNil(t, v)
	assert.Equal(t, float32(4.0), v.(float32))
}

func TestNewEffect_MultiParamDefaults(t *testing.T) {
	// DropShadow: offsetX, offsetY, blur, shadowColor(r,g,b,a)
	e := NewEffect(DropShadow, []float32{2, 4, 8, 0, 0, 0, 1}, 0)

	assert.Equal(t, float32(2.0), e.Uniform("offsetX"))
	assert.Equal(t, float32(4.0), e.Uniform("offsetY"))
	assert.Equal(t, float32(8.0), e.Uniform("blur"))
	assert.Equal(t, [4]float32{0, 0, 0, 1}, e.Uniform("shadowColor"))
}

func TestEffect_SetFloat(t *testing.T) {
	e := NewEffect(Brightness, []float32{1.0}, 0)

	e.SetFloat("brightness", 1.5)
	assert.Equal(t, float32(1.5), e.Uniform("brightness"))
}

func TestEffect_SetVec2(t *testing.T) {
	e := NewEffect(Custom, nil, 0)

	e.SetVec2("offset", 1.0, 2.0)
	assert.Equal(t, [2]float32{1.0, 2.0}, e.Uniform("offset"))
}

func TestEffect_SetVec3(t *testing.T) {
	e := NewEffect(Custom, nil, 0)

	e.SetVec3("color", 0.1, 0.2, 0.3)
	assert.Equal(t, [3]float32{0.1, 0.2, 0.3}, e.Uniform("color"))
}

func TestEffect_SetVec4(t *testing.T) {
	e := NewEffect(Custom, nil, 0)

	e.SetVec4("tint", 0.1, 0.2, 0.3, 0.4)
	assert.Equal(t, [4]float32{0.1, 0.2, 0.3, 0.4}, e.Uniform("tint"))
}

func TestEffect_SetInt(t *testing.T) {
	e := NewEffect(Custom, nil, 0)

	e.SetInt("mode", 2)
	assert.Equal(t, int32(2), e.Uniform("mode"))
}

func TestEffect_SetBool(t *testing.T) {
	e := NewEffect(Custom, nil, 0)

	e.SetBool("enabled", true)
	assert.Equal(t, true, e.Uniform("enabled"))
}

func TestEffect_SetMat3(t *testing.T) {
	e := NewEffect(Custom, nil, 0)

	m := [9]float32{1, 0, 0, 0, 1, 0, 0, 0, 1}
	e.SetMat3("transform", m)
	assert.Equal(t, m, e.Uniform("transform"))
}

func TestEffect_SetMat4(t *testing.T) {
	e := NewEffect(Custom, nil, 0)

	m := [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	e.SetMat4("projection", m)
	assert.Equal(t, m, e.Uniform("projection"))
}

func TestEffect_Enabled(t *testing.T) {
	e := NewEffect(GaussianBlur, []float32{4.0}, 0)

	assert.True(t, e.Enabled())

	e.SetEnabled(false)
	assert.False(t, e.Enabled())

	e.SetEnabled(true)
	assert.True(t, e.Enabled())
}

func TestEffect_Uniforms_Snapshot(t *testing.T) {
	e := NewEffect(Custom, nil, 0)
	e.SetFloat("a", 1.0)
	e.SetFloat("b", 2.0)

	snapshot := e.Uniforms()
	assert.Len(t, snapshot, 2)
	assert.Equal(t, float32(1.0), snapshot["a"])
	assert.Equal(t, float32(2.0), snapshot["b"])

	// Modifying snapshot doesn't affect the effect
	snapshot["c"] = float32(3.0)
	assert.Nil(t, e.Uniform("c"))
}

func TestEffect_Params_Copy(t *testing.T) {
	e := NewEffect(GaussianBlur, []float32{4.0}, 0)

	params := e.Params()
	params[0] = 99.0
	// Original should be unchanged
	assert.Equal(t, float32(4.0), e.Params()[0])
}

func TestNewCustomEffect(t *testing.T) {
	src := `void main() { gl_FragColor = vec4(1.0); }`
	uniforms := map[string][]float32{
		"intensity": {0.5},
		"offset":    {1.0, 2.0},
		"color":     {0.1, 0.2, 0.3},
		"tint":      {0.1, 0.2, 0.3, 0.4},
	}

	e := NewCustomEffect(src, uniforms)

	assert.Equal(t, Custom, e.Type())
	assert.Equal(t, src, e.CustomShaderSrc())
	assert.Equal(t, float32(0.5), e.Uniform("intensity"))
	assert.Equal(t, [2]float32{1.0, 2.0}, e.Uniform("offset"))
	assert.Equal(t, [3]float32{0.1, 0.2, 0.3}, e.Uniform("color"))
	assert.Equal(t, [4]float32{0.1, 0.2, 0.3, 0.4}, e.Uniform("tint"))
}

func TestEffect_ConcurrentAccess(t *testing.T) {
	e := NewEffect(Brightness, []float32{1.0}, 0)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(v float32) {
			defer wg.Done()
			e.SetFloat("brightness", v)
		}(float32(i))
		go func() {
			defer wg.Done()
			_ = e.Uniform("brightness")
		}()
	}
	wg.Wait()

	// Should not panic; final value is non-deterministic
	v := e.Uniform("brightness")
	assert.NotNil(t, v)
}

func TestPassCount(t *testing.T) {
	assert.Equal(t, 2, PassCount(GaussianBlur))
	assert.Equal(t, 3, PassCount(Bloom))
	assert.Equal(t, 1, PassCount(Brightness))
	assert.Equal(t, 1, PassCount(DropShadow))
	assert.Equal(t, 1, PassCount(Custom))
}

func TestDefaultUniforms_AllEffectTypesHandled(t *testing.T) {
	// Ensure DefaultUniforms doesn't panic for any effect type
	for kind := EffectType(0); kind < effectTypeCount; kind++ {
		_ = DefaultUniforms(kind)
	}
}

// Blur is not a defined constant, this tests using GaussianBlur alias behavior
func TestEffect_GaussianBlur_Alias(t *testing.T) {
	e := NewEffect(GaussianBlur, []float32{8.0}, 0)
	assert.Equal(t, float32(8.0), e.Uniform("radius"))
}
