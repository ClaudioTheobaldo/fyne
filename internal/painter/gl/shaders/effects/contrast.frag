#version 110
uniform sampler2D tex;
uniform float contrast;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    color.rgb = (color.rgb - 0.5) * contrast + 0.5;
    gl_FragColor = color;
}
