#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float offsetX;
uniform float offsetY;
uniform float blur;
uniform float spread;
uniform vec4 shadowColor;
varying vec2 fragTexCoord;
void main() {
    vec4 original = texture2D(tex, fragTexCoord);
    vec2 shadowUV = fragTexCoord - vec2(offsetX, offsetY) * texelSize;
    float shadowAlpha = 0.0;
    float total = 0.0;
    float effectiveBlur = max(blur, 1.0);
    int r = int(min(effectiveBlur, 12.0));
    for (int x = -12; x <= 12; x++) {
        if (x < -r || x > r) continue;
        for (int y = -12; y <= 12; y++) {
            if (y < -r || y > r) continue;
            float fi = float(x * x + y * y);
            float weight = exp(-fi / (2.0 * effectiveBlur * effectiveBlur));
            shadowAlpha += texture2D(tex, shadowUV + vec2(float(x), float(y)) * texelSize).a * weight;
            total += weight;
        }
    }
    shadowAlpha = (shadowAlpha / total) * shadowColor.a;
    shadowAlpha = min(shadowAlpha + spread * 0.1, 1.0);
    vec4 shadow = vec4(shadowColor.rgb * shadowAlpha, shadowAlpha);
    gl_FragColor = original + shadow * (1.0 - original.a);
}
