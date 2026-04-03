#version 110

// GRAY8:  GL_LUMINANCE, bitMax=1.0,     hasAlpha=0.0
// GRAY16: GL_LUMINANCE_ALPHA, bitMax=65535.0, hasAlpha=0.0
// YA8:    GL_LUMINANCE_ALPHA, bitMax=1.0,     hasAlpha=1.0
//
// For 16-bit: gray = (sample.r * 255 + sample.a * 255 * 256) / bitMax
// For 8-bit:  gray = sample.r  (sample.a is 1.0 for LUMINANCE, or the alpha for YA8)

uniform sampler2D texGray;
uniform float bitMax;   // 1.0 for 8-bit (no reconstruction), 65535.0 for 16-bit
uniform float hasAlpha; // 1.0 for YA8, 0.0 otherwise
uniform float cornerRadius;
uniform vec2 size;
uniform vec4 inset;

varying vec2 fragTexCoord;
varying float fragAlpha;

void main() {
    float alpha = 1.0;
    if (cornerRadius > 0.5) {
        vec2 normalizedCoord = (fragTexCoord - inset.xy) / (inset.zw - inset.xy);
        vec2 pos = normalizedCoord * size;
        vec2 halfSize = size * 0.5;
        float dist = length(max(abs(pos - halfSize) - halfSize + cornerRadius, 0.0)) - cornerRadius;
        alpha = 1.0 - smoothstep(-1.0, 1.0, dist);
    }

    vec4 raw = texture2D(texGray, fragTexCoord);
    // hibit=1.0 when we need to reconstruct 16-bit value from two 8-bit halves
    float hibit = step(2.0, bitMax);
    float reconstructed = (raw.r * 255.0 + raw.a * 255.0 * 256.0) / bitMax;
    float gray = mix(raw.r, reconstructed, hibit);
    // Alpha: use .a only for YA8 (8-bit luminance+alpha). For GRAY16, .a holds MSB.
    float pixAlpha = mix(1.0, raw.a, hasAlpha * (1.0 - hibit));

    vec4 texColor = vec4(gray, gray, gray, pixAlpha);
    texColor.a *= fragAlpha * alpha;
    texColor.rgb *= fragAlpha * alpha;

    if (texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
