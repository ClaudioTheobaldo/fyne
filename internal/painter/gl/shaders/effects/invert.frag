#version 110
uniform sampler2D tex;
uniform float amount;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    color.rgb = mix(color.rgb, 1.0 - color.rgb, amount);
    gl_FragColor = color;
}
