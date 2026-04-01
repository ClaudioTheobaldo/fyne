#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform vec2 direction;
uniform float radius;
varying vec2 fragTexCoord;
void main() {
    vec4 color = vec4(0.0);
    float total = 0.0;
    int samples = int(min(radius, 32.0));
    for (int i = -32; i <= 32; i++) {
        if (i < -samples || i > samples) continue;
        float fi = float(i);
        float weight = exp(-0.5 * (fi * fi) / (radius * radius * 0.25 + 0.001));
        vec2 offset = direction * texelSize * fi;
        color += texture2D(tex, fragTexCoord + offset) * weight;
        total += weight;
    }
    gl_FragColor = color / total;
}
