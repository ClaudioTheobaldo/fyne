#version 110
uniform sampler2D tex;
uniform float time;
uniform float amount;
uniform float blockSize;
varying vec2 fragTexCoord;
float hash(vec2 p) {
    return fract(sin(dot(p, vec2(12.9898, 78.233))) * 43758.5453);
}
void main() {
    vec2 uv = fragTexCoord;
    float block = floor(uv.y * blockSize);
    float rand = hash(vec2(block, floor(time * 10.0)));
    if (rand > 1.0 - amount * 0.3) {
        float shift = (hash(vec2(block + 1.0, floor(time * 10.0))) - 0.5) * amount * 0.1;
        uv.x += shift;
    }
    float r = texture2D(tex, uv + vec2(amount * 0.003, 0.0)).r;
    float g = texture2D(tex, uv).g;
    float b = texture2D(tex, uv - vec2(amount * 0.003, 0.0)).b;
    float a = texture2D(tex, uv).a;
    gl_FragColor = vec4(r, g, b, a);
}
