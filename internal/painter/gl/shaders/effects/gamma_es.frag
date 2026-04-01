#version 100
precision mediump float;
uniform sampler2D tex;
uniform float gamma;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float g = 1.0 / max(gamma, 0.01);
    color.rgb = pow(color.rgb, vec3(g));
    gl_FragColor = color;
}
