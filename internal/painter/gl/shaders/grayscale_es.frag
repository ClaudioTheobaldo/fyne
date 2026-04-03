#version 100

#ifdef GL_ES
# ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
# else
precision mediump float;
#endif
precision mediump int;
precision lowp sampler2D;
#endif

uniform sampler2D texGray;
uniform float bitMax;
uniform float hasAlpha;
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
    float hibit = step(2.0, bitMax);
    float reconstructed = (raw.r * 255.0 + raw.a * 255.0 * 256.0) / bitMax;
    float gray = mix(raw.r, reconstructed, hibit);
    float pixAlpha = mix(1.0, raw.a, hasAlpha * (1.0 - hibit));

    vec4 texColor = vec4(gray, gray, gray, pixAlpha);
    texColor.a *= fragAlpha * alpha;
    texColor.rgb *= fragAlpha * alpha;

    if (texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
