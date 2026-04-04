#version 100
precision mediump float;
uniform sampler2D tex;
uniform float angle;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float a = angle * 3.14159265 / 180.0;
    float cosA = cos(a);
    float sinA = sin(a);
    mat3 hueRotation = mat3(
        0.213 + cosA * 0.787 - sinA * 0.213,
        0.213 - cosA * 0.213 + sinA * 0.143,
        0.213 - cosA * 0.213 - sinA * 0.787,
        0.715 - cosA * 0.715 - sinA * 0.715,
        0.715 + cosA * 0.285 + sinA * 0.140,
        0.715 - cosA * 0.715 + sinA * 0.715,
        0.072 - cosA * 0.072 + sinA * 0.928,
        0.072 - cosA * 0.072 - sinA * 0.283,
        0.072 + cosA * 0.928 + sinA * 0.072
    );
    color.rgb = hueRotation * color.rgb;
    gl_FragColor = color;
}
