#version 110

// NV12/NV21/NV16/NV24/NV42 semi-planar YUV.
// texY:  GL_LUMINANCE — one luma sample per texel.
// texUV: GL_LUMINANCE_ALPHA — interleaved UV: .r=first component, .a=second component.
// swapUV: 0.0 = NV12 order (U in .r, V in .a), 1.0 = NV21 order (V in .r, U in .a).

uniform sampler2D texY;
uniform sampler2D texUV;
uniform float swapUV;
uniform vec3 colorRow0;
uniform vec3 colorRow1;
uniform vec3 colorRow2;
uniform vec3 colorOffset;
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

    float y = texture2D(texY, fragTexCoord).r;
    vec2 uv = texture2D(texUV, fragTexCoord).ra;
    // swapUV=0: u=uv.x, v=uv.y;  swapUV=1: u=uv.y, v=uv.x
    float u = mix(uv.x, uv.y, swapUV);
    float v = mix(uv.y, uv.x, swapUV);

    vec3 yuv = vec3(y, u, v) - colorOffset;
    float r = dot(colorRow0, yuv);
    float g = dot(colorRow1, yuv);
    float b = dot(colorRow2, yuv);

    vec4 texColor = vec4(r, g, b, 1.0);
    texColor.a *= fragAlpha * alpha;
    texColor.rgb *= fragAlpha * alpha;

    if (texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
