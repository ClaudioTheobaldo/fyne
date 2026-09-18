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
// As with the Gaussian, GLSL 110 forces a constant loop bound, so the window is
// fixed at the maximum and lanes outside `radius` are masked before sampling.
// Both Gaussian denominators are reciprocals of uniforms and are computed once
// up front, and the two exp() calls are folded into one: exp(a)*exp(b) is
// exp(a+b), which halves the transcendental count per tap.
void main() {
    vec4 center = texture2D(tex, fragTexCoord);
    vec4 color = vec4(0.0);
    float total = 0.0;

    int r = int(min(radius, 8.0));
    float invSS2 = -1.0 / (2.0 * sigmaSpace * sigmaSpace + 0.001);
    float invSC2 = -1.0 / (2.0 * sigmaColor * sigmaColor + 0.001);

    for (int x = -8; x <= 8; x++) {
        if (x < -r || x > r) continue;
        for (int y = -8; y <= 8; y++) {
            if (y < -r || y > r) continue;

            vec2 offset = vec2(float(x), float(y)) * texelSize;
            vec4 sample_color = texture2D(tex, fragTexCoord + offset);

            float spatialDist = float(x * x + y * y);
            vec3 dc = sample_color.rgb - center.rgb;
            float colorDist = dot(dc, dc);

            float weight = exp(spatialDist * invSS2 + colorDist * invSC2);
            color += sample_color * weight;
            total += weight;
        }
    }
    gl_FragColor = color / total;
}
