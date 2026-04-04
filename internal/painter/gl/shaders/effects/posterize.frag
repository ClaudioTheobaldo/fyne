#version 110
uniform sampler2D tex;
uniform float levels;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float n = max(levels, 2.0);
    color.rgb = floor(color.rgb * n) / (n - 1.0);
    gl_FragColor = color;
}
