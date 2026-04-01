#version 110
uniform sampler2D tex;
uniform float opacity;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    color *= opacity;
    gl_FragColor = color;
}
