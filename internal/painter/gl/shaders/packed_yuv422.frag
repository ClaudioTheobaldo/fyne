#version 110

// Packed YUV 4:2:2 formats: YUYV422, UYVY422, YVYU422.
// The packed data is uploaded as GL_RGBA at half-width:
//   YUYV: texel = (Y0, U, Y1, V)  — .r=Y0 .g=U .b=Y1 .a=V
//   UYVY: texel = (U, Y0, V, Y1)  — .r=U  .g=Y0 .b=V  .a=Y1
//   YVYU: texel = (Y0, V, Y1, U)  — .r=Y0 .g=V  .b=Y1 .a=U
// packingMode: 0.0=YUYV, 1.0=UYVY, 2.0=YVYU
// texPackedWidth: width of the packed texture in texels (= frame_width / 2)

uniform sampler2D texPacked;
uniform float packingMode;
uniform float texPackedWidth;
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

    // The packed texture is at half-width. Sample it directly — GL's texture
    // coordinate mapping already accounts for the half-width upload.
    vec4 pack = texture2D(texPacked, fragTexCoord);

    // Determine if we're rendering the right pixel of a pair (sub-pixel = [0.5, 1.0))
    float isRight = step(0.5, fract(fragTexCoord.x * texPackedWidth));

    float y, u, v;
    if (packingMode < 0.5) {
        // YUYV: (Y0=.r, U=.g, Y1=.b, V=.a)
        y = mix(pack.r, pack.b, isRight);
        u = pack.g;
        v = pack.a;
    } else if (packingMode < 1.5) {
        // UYVY: (U=.r, Y0=.g, V=.b, Y1=.a)
        y = mix(pack.g, pack.a, isRight);
        u = pack.r;
        v = pack.b;
    } else {
        // YVYU: (Y0=.r, V=.g, Y1=.b, U=.a)
        y = mix(pack.r, pack.b, isRight);
        u = pack.a;
        v = pack.g;
    }

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
