#version 100
precision mediump float;
uniform sampler2D tex;
uniform float segments;
uniform float rotation;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord - 0.5;
    float a = atan(uv.y, uv.x) + rotation * 3.14159265 / 180.0;
    float r = length(uv);
    float seg = max(segments, 1.0);
    float sliceAngle = 6.28318 / seg;
    a = mod(a, sliceAngle);
    // Mirror alternate slices
    if (a > sliceAngle * 0.5)
        a = sliceAngle - a;
    // Map back to UV
    uv = vec2(cos(a), sin(a)) * r + 0.5;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = texture2D(tex, clamp(uv, 0.0, 1.0));
    else
        gl_FragColor = texture2D(tex, uv);
}
