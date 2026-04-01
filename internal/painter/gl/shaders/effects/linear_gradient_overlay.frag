#version 110
uniform sampler2D tex;
uniform vec4 startColor;
uniform vec4 endColor;
uniform float angle;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float a = angle * 3.14159265 / 180.0;
    float t = dot(fragTexCoord - 0.5, vec2(cos(a), sin(a))) + 0.5;
    t = clamp(t, 0.0, 1.0);
    vec4 gradient = mix(startColor, endColor, t);
    color.rgb = mix(color.rgb, gradient.rgb, gradient.a);
    gl_FragColor = color;
}
