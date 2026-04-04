#version 100
precision mediump float;
uniform sampler2D tex;
uniform float brightness;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    color.rgb *= brightness;
    gl_FragColor = color;
}
