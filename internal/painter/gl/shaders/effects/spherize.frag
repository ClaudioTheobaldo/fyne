#version 110
uniform sampler2D tex;
uniform float radius;
uniform vec2 center;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord - center;
    float dist = length(uv);
    float r = max(radius, 0.001);
    if (dist < r) {
        float percent = dist / r;
        float theta = asin(percent) / (3.14159 * 0.5);
        uv = normalize(uv) * theta * r;
    }
    gl_FragColor = texture2D(tex, uv + center);
}
