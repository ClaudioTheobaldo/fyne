#version 110
uniform sampler2D tex;
uniform float startAlpha;
uniform float endAlpha;
uniform float angle;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float a = angle * 3.14159265 / 180.0;
    float t = dot(fragTexCoord - 0.5, vec2(cos(a), sin(a))) + 0.5;
    t = clamp(t, 0.0, 1.0);
    float maskAlpha = mix(startAlpha, endAlpha, t);
    color *= maskAlpha;
    gl_FragColor = color;
}
