#version 100
precision mediump float;
uniform sampler2D tex;
uniform float angle;
varying vec2 fragTexCoord;

// Luminance-preserving hue rotation. The rotation depends only on `angle`,
// which is a uniform, so the sin/cos pair and the nine coefficients are the
// same for every pixel in the draw. Expressed as three dot products against
// coefficient vectors rather than by constructing a mat3 per fragment, which
// costs a full matrix build before the multiply.
void main() {
    vec4 color = texture2D(tex, fragTexCoord);

    float a = angle * 0.01745329252; // pi/180
    float c = cos(a);
    float s = sin(a);

    vec3 wr = vec3(0.213, 0.715, 0.072);
    vec3 r = wr + c * vec3( 0.787, -0.715, -0.072) + s * vec3(-0.213, -0.715,  0.928);
    vec3 g = wr + c * vec3(-0.213,  0.285, -0.072) + s * vec3( 0.143,  0.140, -0.283);
    vec3 b = wr + c * vec3(-0.213, -0.715,  0.928) + s * vec3(-0.787,  0.715,  0.072);

    color.rgb = vec3(dot(color.rgb, r), dot(color.rgb, g), dot(color.rgb, b));
    gl_FragColor = color;
}
