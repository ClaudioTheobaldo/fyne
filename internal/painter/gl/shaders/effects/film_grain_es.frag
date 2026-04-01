#version 100
precision mediump float;
uniform sampler2D tex;
uniform float intensity;
uniform float speed;
varying vec2 fragTexCoord;
float hash(vec2 p) {
    return fract(sin(dot(p, vec2(12.9898, 78.233))) * 43758.5453);
}
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float noise = hash(fragTexCoord * 1000.0 + speed) * 2.0 - 1.0;
    color.rgb += noise * intensity;
    gl_FragColor = clamp(color, 0.0, 1.0);
}
