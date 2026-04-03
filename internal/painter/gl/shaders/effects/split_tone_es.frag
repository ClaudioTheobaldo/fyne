#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec3 shadowTint;
uniform vec3 highlightTint;
uniform float balance;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float lum = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    float shadowMask = 1.0 - smoothstep(balance - 0.1, balance + 0.1, lum);
    float highlightMask = smoothstep(balance - 0.1, balance + 0.1, lum);
    color.rgb += shadowTint * shadowMask * 0.3;
    color.rgb += highlightTint * highlightMask * 0.3;
    gl_FragColor = clamp(color, 0.0, 1.0);
}
