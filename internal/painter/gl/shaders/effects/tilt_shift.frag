#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
uniform float focusY;
uniform float focusWidth;
varying vec2 fragTexCoord;
void main() {
    float dist = abs(fragTexCoord.y - focusY);
    float blurAmount = smoothstep(0.0, max(focusWidth, 0.001), dist) * radius;
    vec4 color = vec4(0.0);
    float total = 0.0;
    int samples = int(min(blurAmount, 16.0));
    if (samples == 0) {
        gl_FragColor = texture2D(tex, fragTexCoord);
        return;
    }
    for (int i = -16; i <= 16; i++) {
        if (i < -samples || i > samples) continue;
        float fi = float(i);
        float weight = 1.0 - abs(fi) / float(samples + 1);
        vec2 offset = vec2(0.0, fi * texelSize.y);
        color += texture2D(tex, fragTexCoord + offset) * weight;
        total += weight;
    }
    gl_FragColor = color / total;
}
