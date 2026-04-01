#version 110
uniform sampler2D tex;
uniform vec4 startColor;
uniform vec4 endColor;
uniform vec2 center;
uniform float angle;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec2 d = fragTexCoord - center;
    float a = atan(d.y, d.x) + 3.14159265;
    float offset = angle * 3.14159265 / 180.0;
    float t = mod(a + offset, 6.28318) / 6.28318;
    vec4 gradient = mix(startColor, endColor, t);
    color.rgb = mix(color.rgb, gradient.rgb, gradient.a);
    gl_FragColor = color;
}
