#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec4 startColor;
uniform vec4 endColor;
uniform vec2 center;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float dist = length(fragTexCoord - center) * 1.414;
    float t = clamp(dist, 0.0, 1.0);
    vec4 gradient = mix(startColor, endColor, t);
    color.rgb = mix(color.rgb, gradient.rgb, gradient.a);
    gl_FragColor = color;
}
