#version 110
uniform sampler2D tex;
uniform float strength;
uniform vec2 center;
varying vec2 fragTexCoord;
void main() {
    vec2 dir = fragTexCoord - center;
    vec4 color = vec4(0.0);
    float total = 0.0;
    for (int i = 0; i < 16; i++) {
        float t = float(i) / 16.0;
        float weight = 1.0 - t;
        vec2 uv = fragTexCoord - dir * t * strength * 0.01;
        color += texture2D(tex, uv) * weight;
        total += weight;
    }
    gl_FragColor = color / total;
}
