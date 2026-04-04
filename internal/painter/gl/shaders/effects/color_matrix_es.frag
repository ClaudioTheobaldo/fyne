#version 100
precision mediump float;
uniform sampler2D tex;
uniform mat4 colorMatrix;
uniform vec4 matrixOffset;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec4 result = colorMatrix * color + matrixOffset;
    result.a = color.a;
    gl_FragColor = clamp(result, 0.0, 1.0);
}
