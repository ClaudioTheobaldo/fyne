#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
uniform vec2 direction;
varying vec2 fragTexCoord;
void main() {
    vec4 color = vec4(0.0);
    float total = 0.0;
    int samples = int(min(radius, 32.0));
    vec2 dir = normalize(direction) * texelSize;
    for (int i = -32; i <= 32; i++) {
        if (i < -samples || i > samples) continue;
        float fi = float(i);
        float weight = 1.0 - abs(fi) / float(samples + 1);
        color += texture2D(tex, fragTexCoord + dir * fi) * weight;
        total += weight;
    }
    gl_FragColor = color / total;
}
