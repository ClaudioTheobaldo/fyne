#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
uniform float sigmaSpace;
uniform float sigmaColor;
varying vec2 fragTexCoord;

// Edge-preserving blur: weight each tap by both its distance and its colour
// difference from the centre, so flat areas smooth while edges survive.
//
// MAX_RADIUS is rewritten by the painter to ceil(radius) before compiling, so
// each whole-number window size gets its own program and the loop runs only as
// far as the radius needs. GLSL 110 / ES 100 require a constant loop bound, so
// this cannot be a uniform. Keep the define below as the only place that
// spells it out, with the value 8, or the painter cannot find it.
//
// The outermost ring of the window fades in with the fractional part of the
// radius. Without that the window only grew in whole steps and any radius
// below 1.0 did nothing at all. At whole-number radii every tap gets a factor
// of exactly 1.0, so those render as before.
//
// Both Gaussian denominators are reciprocals of uniforms and are computed once
// up front, and the two exp() calls are folded into one: exp(a)*exp(b) is
// exp(a+b), which halves the transcendental count per tap.
#define MAX_RADIUS 8

void main() {
    vec4 center = texture2D(tex, fragTexCoord);
    vec4 color = vec4(0.0);
    float total = 0.0;

    float invSS2 = -1.0 / (2.0 * sigmaSpace * sigmaSpace + 0.001);
    float invSC2 = -1.0 / (2.0 * sigmaColor * sigmaColor + 0.001);

    for (int x = -MAX_RADIUS; x <= MAX_RADIUS; x++) {
        for (int y = -MAX_RADIUS; y <= MAX_RADIUS; y++) {
            vec2 offset = vec2(float(x), float(y)) * texelSize;
            vec4 sample_color = texture2D(tex, fragTexCoord + offset);

            float spatialDist = float(x * x + y * y);
            vec3 dc = sample_color.rgb - center.rgb;
            float colorDist = dot(dc, dc);

            // Ring index: 0 for the centre, 1 for the 3x3 border, and so on.
            float ring = max(abs(float(x)), abs(float(y)));
            float fade = clamp(radius - ring + 1.0, 0.0, 1.0);

            float weight = exp(spatialDist * invSS2 + colorDist * invSC2) * fade;
            color += sample_color * weight;
            total += weight;
        }
    }
    gl_FragColor = color / total;
}
