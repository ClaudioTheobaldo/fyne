#version 110
uniform sampler2D tex;
uniform vec2 baseFrequency;
uniform float numOctaves;
uniform float seed;
varying vec2 fragTexCoord;
float hash(vec2 p) {
    return fract(sin(dot(p + seed, vec2(127.1, 311.7))) * 43758.5453);
}
float noise(vec2 p) {
    vec2 i = floor(p);
    vec2 f = fract(p);
    f = f * f * (3.0 - 2.0 * f);
    float a = hash(i);
    float b = hash(i + vec2(1.0, 0.0));
    float c = hash(i + vec2(0.0, 1.0));
    float d = hash(i + vec2(1.0, 1.0));
    return mix(mix(a, b, f.x), mix(c, d, f.x), f.y);
}
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float n = 0.0;
    float amp = 1.0;
    float freq = 1.0;
    float total = 0.0;
    int octaves = int(min(numOctaves, 8.0));
    for (int i = 0; i < 8; i++) {
        if (i >= octaves) break;
        n += noise(fragTexCoord * baseFrequency * freq) * amp;
        total += amp;
        amp *= 0.5;
        freq *= 2.0;
    }
    n /= total;
    color.rgb = mix(color.rgb, vec3(n), 0.5);
    gl_FragColor = color;
}
