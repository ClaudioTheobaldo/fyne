package canvas

import (
	"sync"

	"fyne.io/fyne/v2/canvas/effect"
)

// effectList provides AddEffect/RemoveEffect/ClearEffects methods for canvas objects.
// It is embedded in baseObject (and separately in Circle and Line which don't use baseObject).
type effectList struct {
	effects []*effect.Effect
	mu      sync.RWMutex
}

// AddEffect applies a new shader effect to this canvas object and returns a handle
// to the created effect instance. Effects stack in insertion order.
// The params are mapped to shader uniforms based on the effect type's default mappings.
//
// Note: To enable auto-refresh when effect parameters change (e.g. via SetFloat or
// Animate), call effect.SetOwner(e, canvasObject) after adding the effect.
// The painter also sets owners automatically when it first encounters effects.
//
// Example:
//
//	blur := rect.AddEffect(effect.GaussianBlur, 4.0)
//	effect.SetOwner(blur, rect) // enables auto-refresh on parameter changes
//
// Since: 2.8
func (el *effectList) AddEffect(kind effect.EffectType, params ...float32) *effect.Effect {
	el.mu.Lock()
	e := effect.NewEffect(kind, params, len(el.effects))
	el.effects = append(el.effects, e)
	el.mu.Unlock()
	return e
}

// AddCustomEffect applies a custom shader effect with user-supplied GLSL fragment source.
// The effect participates in the same pipeline as built-in effects and can be stacked.
//
// Since: 2.8
func (el *effectList) AddCustomEffect(shaderSrc string, uniforms map[string][]float32) *effect.Effect {
	el.mu.Lock()
	e := effect.NewCustomEffect(shaderSrc, uniforms)
	effect.SetOrder(e, len(el.effects))
	el.effects = append(el.effects, e)
	el.mu.Unlock()
	return e
}

// RemoveEffect removes a specific effect from this canvas object's effect chain.
//
// Since: 2.8
func (el *effectList) RemoveEffect(e *effect.Effect) {
	el.mu.Lock()
	for i, existing := range el.effects {
		if existing == e {
			el.effects = append(el.effects[:i], el.effects[i+1:]...)
			break
		}
	}
	el.mu.Unlock()
}

// ClearEffects removes all effects from this canvas object.
//
// Since: 2.8
func (el *effectList) ClearEffects() {
	el.mu.Lock()
	el.effects = nil
	el.mu.Unlock()
}

// Effects returns a copy of the current effect chain in stacking order.
//
// Since: 2.8
func (el *effectList) Effects() []*effect.Effect {
	el.mu.RLock()
	defer el.mu.RUnlock()
	out := make([]*effect.Effect, len(el.effects))
	copy(out, el.effects)
	return out
}

// HasEffects returns true if this canvas object has any effects applied.
//
// Since: 2.8
func (el *effectList) HasEffects() bool {
	el.mu.RLock()
	defer el.mu.RUnlock()
	return len(el.effects) > 0
}
