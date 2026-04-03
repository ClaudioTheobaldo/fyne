#version 110

// High bit-depth planar YUV (9–16 bit, little-endian).
// Each plane is uploaded as GL_LUMINANCE_ALPHA: .r = low byte, .a = high byte.
// Reconstruction: val = (sample.r * 255 + sample.a * 255 * 256) / bitMax
// bitMax = (1 << bitDepth) - 1  (e.g. 1023 for 10-bit, 65535 for 16-bit)
//
// For YUVA variants (hasAlpha=1.0), a 4th sampler texA holds the alpha plane.

uniform sampler2D texY;
uniform sampler2D texU;
uniform sampler2D texV;
uniform sampler2D texA;   // only used when hasAlpha=1.0
uniform float bitMax;
uniform float hasAlpha;   // 0.0 = no alpha plane, 1.0 = read texA
uniform vec3 colorRow0;
uniform vec3 colorRow1;
uniform vec3 colorRow2;
uniform vec3 colorOffset;
uniform float cornerRadius;
uniform vec2 size;
uniform vec4 inset;

varying vec2 fragTexCoord;
varying float fragAlpha;

float unpack16(sampler2D tex, vec2 coord) {
    vec2 s = texture2D(tex, coord).ra;
    return (s.x * 255.0 + s.y * 255.0 * 256.0) / bitMax;
}

void main() {
    float alpha = 1.0;
    if (cornerRadius > 0.5) {
        vec2 normalizedCoord = (fragTexCoord - inset.xy) / (inset.zw - inset.xy);
        vec2 pos = normalizedCoord * size;
        vec2 halfSize = size * 0.5;
        float dist = length(max(abs(pos - halfSize) - halfSize + cornerRadius, 0.0)) - cornerRadius;
        alpha = 1.0 - smoothstep(-1.0, 1.0, dist);
    }

    float y = unpack16(texY, fragTexCoord);
    float u = unpack16(texU, fragTexCoord);
    float v = unpack16(texV, fragTexCoord);
    float a = mix(1.0, unpack16(texA, fragTexCoord), hasAlpha);

    vec3 yuv = vec3(y, u, v) - colorOffset;
    float r = dot(colorRow0, yuv);
    float g = dot(colorRow1, yuv);
    float b = dot(colorRow2, yuv);

    vec4 texColor = vec4(r, g, b, a);
    texColor.a *= fragAlpha * alpha;
    texColor.rgb *= fragAlpha * alpha;

    if (texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
