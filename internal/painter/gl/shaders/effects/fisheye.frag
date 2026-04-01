#version 110
uniform sampler2D tex;
uniform float strength;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord * 2.0 - 1.0;
    float r = length(uv);
    float bind = sqrt(2.0);
    if (strength > 0.0) {
        uv = normalize(uv) * tan(r * strength) * bind / tan(bind * strength);
    } else {
        uv = normalize(uv) * atan(r * -strength * 10.0) * bind / atan(-strength * bind * 10.0);
    }
    gl_FragColor = texture2D(tex, uv * 0.5 + 0.5);
}
