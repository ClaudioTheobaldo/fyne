#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec3 color;
uniform float intensity;
varying vec2 fragTexCoord;
void main() {
    vec4 texColor = texture2D(tex, fragTexCoord);
    texColor.rgb += color * intensity;
    gl_FragColor = clamp(texColor, 0.0, 1.0);
}
