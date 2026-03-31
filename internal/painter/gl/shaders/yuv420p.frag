#version 110

uniform sampler2D texY;
uniform sampler2D texU;
uniform sampler2D texV;
uniform float cornerRadius;  // in pixels
uniform vec2 size;           // in pixels: size of the rendered image quad
uniform vec4 inset;          // texture coordinate insets (minX, minY, maxX, maxY)

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
    float u = texture2D(texU, fragTexCoord).r;
    float v = texture2D(texV, fragTexCoord).r;

    // BT.601 full-range (YUVJ420P, standard for CCTV/surveillance H.264)
    float r = y + 1.402 * (v - 0.5);
    float g = y - 0.344136 * (u - 0.5) - 0.714136 * (v - 0.5);
    float b = y + 1.772 * (u - 0.5);

    vec4 texColor = vec4(r, g, b, 1.0);
    texColor.a *= fragAlpha * alpha;
    texColor.rgb *= fragAlpha * alpha;

    if (texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
