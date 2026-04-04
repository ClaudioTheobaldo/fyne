#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float offsetX;
uniform float offsetY;
uniform float blur;
uniform vec4 shadowColor;
varying vec2 fragTexCoord;
void main() {
    vec4 original = texture2D(tex, fragTexCoord);
    vec2 shadowUV = fragTexCoord - vec2(offsetX, offsetY) * texelSize;
    float blurredAlpha = 0.0;
    float total = 0.0;
    int r = int(min(blur, 12.0));
    for (int x = -12; x <= 12; x++) {
        if (x < -r || x > r) continue;
        for (int y = -12; y <= 12; y++) {
            if (y < -r || y > r) continue;
            float fi = float(x * x + y * y);
            float weight = exp(-fi / (2.0 * blur * blur + 0.001));
            blurredAlpha += texture2D(tex, shadowUV + vec2(float(x), float(y)) * texelSize).a * weight;
            total += weight;
        }
    }
    blurredAlpha /= total;
    float inShadow = original.a * (1.0 - blurredAlpha) * shadowColor.a;
    original.rgb = mix(original.rgb, shadowColor.rgb, inShadow);
    gl_FragColor = original;
}
