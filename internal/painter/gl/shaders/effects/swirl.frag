#version 110
uniform sampler2D tex;
uniform float radius;
uniform float angle;
uniform vec2 center;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord - center;
    float dist = length(uv);
    float r = max(radius, 0.001);
    if (dist < r) {
        float percent = (r - dist) / r;
        float theta = percent * percent * angle;
        float s = sin(theta);
        float c = cos(theta);
        uv = vec2(uv.x * c - uv.y * s, uv.x * s + uv.y * c);
    }
    gl_FragColor = texture2D(tex, uv + center);
}
