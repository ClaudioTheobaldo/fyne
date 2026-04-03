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

    vec4 pack = texture2D(texPacked, fragTexCoord);
    float isRight = step(0.5, fract(fragTexCoord.x * texPackedWidth));

    float y, u, v;
    if (packingMode < 0.5) {
        y = mix(pack.r, pack.b, isRight);
        u = pack.g;
        v = pack.a;
    } else if (packingMode < 1.5) {
        y = mix(pack.g, pack.a, isRight);
        u = pack.r;
        v = pack.b;
    } else {
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
