#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform vec2 direction;
uniform float radius;
varying vec2 fragTexCoord;

// Separable Gaussian, one axis per pass (see PassCount).
//
// GLSL 110 requires a compile-time constant loop bound, so the trip count is
// fixed at the maximum and lanes outside the current radius are masked. Two
// things keep that from being as expensive as it looks: the Gaussian divisor
// depends only on `radius` (a uniform) and is hoisted out of the loop, and the
// masked iterations skip the texture fetch rather than sampling and discarding.
void main() {
    vec4 color = vec4(0.0);
    float total = 0.0;

    int samples = int(min(radius, 32.0));
    float invDenom = 1.0 / (radius * radius * 0.25 + 0.001);

    for (int i = -32; i <= 32; i++) {
        if (i < -samples || i > samples) continue;
        float fi = float(i);
        float weight = exp(-0.5 * fi * fi * invDenom);
        color += texture2D(tex, fragTexCoord + direction * texelSize * fi) * weight;
        total += weight;
    }
    gl_FragColor = color / total;
}
