#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
uniform vec4 glowColor;
varying vec2 fragTexCoord;
void main() {
    vec4 original = texture2D(tex, fragTexCoord);
    float maxAlpha = 0.0;
    int r = int(min(radius, 16.0));
    for (int x = -16; x <= 16; x++) {
        if (x < -r || x > r) continue;
        for (int y = -16; y <= 16; y++) {
            if (y < -r || y > r) continue;
            float a = texture2D(tex, fragTexCoord + vec2(float(x), float(y)) * texelSize).a;
            float dist = length(vec2(float(x), float(y)));
            float weight = 1.0 - dist / float(r + 1);
            maxAlpha = max(maxAlpha, a * max(weight, 0.0));
        }
    }
    float glowAlpha = maxAlpha * (1.0 - original.a) * glowColor.a;
    vec4 glow = vec4(glowColor.rgb * glowAlpha, glowAlpha);
    gl_FragColor = original + glow * (1.0 - original.a);
}
