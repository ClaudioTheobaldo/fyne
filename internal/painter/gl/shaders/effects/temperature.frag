#version 110
uniform sampler2D tex;
uniform float temperature;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    // Warm: increase red, decrease blue. Cool: opposite.
    color.r += temperature * 0.1;
    color.b -= temperature * 0.1;
    gl_FragColor = clamp(color, 0.0, 1.0);
}
