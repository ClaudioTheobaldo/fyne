#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
varying vec2 fragTexCoord;
void main() {
    vec4 color = vec4(0.0);
    float count = 0.0;
    int r = int(min(radius, 16.0));
    for (int x = -16; x <= 16; x++) {
        if (x < -r || x > r) continue;
        for (int y = -16; y <= 16; y++) {
            if (y < -r || y > r) continue;
            color += texture2D(tex, fragTexCoord + vec2(float(x), float(y)) * texelSize);
            count += 1.0;
        }
    }
    gl_FragColor = color / count;
}
