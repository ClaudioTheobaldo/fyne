#version 110
uniform sampler2D tex;
uniform float time;
uniform float width;
uniform float angle;
uniform float intensity;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float a = angle * 3.14159265 / 180.0;
    float pos = dot(fragTexCoord, vec2(cos(a), sin(a)));
    float t = fract(time);
    float dist = abs(pos - t);
    float w = max(width, 0.01);
    float shine = exp(-dist * dist / (w * w * 0.5)) * intensity;
    color.rgb += shine;
    gl_FragColor = clamp(color, 0.0, 1.0);
}
