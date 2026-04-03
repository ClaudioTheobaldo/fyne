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

uniform sampler2D texY;
uniform sampler2D texUV;
uniform float bitMax;
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

    vec2 ys = texture2D(texY, fragTexCoord).ra;
    float y = (ys.x * 255.0 + ys.y * 255.0 * 256.0) / bitMax;

    vec4 uvs = texture2D(texUV, fragTexCoord);
    float uRaw = (uvs.r * 255.0 + uvs.g * 255.0 * 256.0) / bitMax;
    float vRaw = (uvs.b * 255.0 + uvs.a * 255.0 * 256.0) / bitMax;
    float u = mix(uRaw, vRaw, swapUV);
    float v = mix(vRaw, uRaw, swapUV);

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
