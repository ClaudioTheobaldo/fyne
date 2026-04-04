#version 110
uniform sampler2D tex;
uniform float amount;
uniform float scale;
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
    vec2 uv = fragTexCoord;
    float n1 = noise(uv * scale) * 2.0 - 1.0;
    float n2 = noise(uv * scale + 100.0) * 2.0 - 1.0;
    uv += vec2(n1, n2) * amount * 0.01;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else
        gl_FragColor = texture2D(tex, uv);
}
