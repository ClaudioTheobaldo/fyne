#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
uniform float noiseScale;
varying vec2 fragTexCoord;
float hash(vec2 p) {
    return fract(sin(dot(p, vec2(12.9898, 78.233))) * 43758.5453);
}
void main() {
    vec4 color = vec4(0.0);
    float total = 0.0;
    int r = int(min(radius, 8.0));
    for (int x = -8; x <= 8; x++) {
        if (x < -r || x > r) continue;
        for (int y = -8; y <= 8; y++) {
            if (y < -r || y > r) continue;
            vec2 noise = vec2(hash(fragTexCoord + vec2(float(x), float(y))),
                              hash(fragTexCoord + vec2(float(y), float(x)))) * 2.0 - 1.0;
            vec2 offset = (vec2(float(x), float(y)) + noise * noiseScale) * texelSize;
            color += texture2D(tex, fragTexCoord + offset);
            total += 1.0;
        }
    }
    gl_FragColor = color / total;
}
