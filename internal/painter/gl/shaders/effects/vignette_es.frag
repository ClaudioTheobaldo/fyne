#version 100
precision mediump float;
uniform sampler2D tex;
uniform float intensity;
uniform float smoothness;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec2 uv = fragTexCoord * 2.0 - 1.0;
    float dist = length(uv);
    float vig = smoothstep(1.0 - smoothness, 1.0, dist);
    color.rgb *= 1.0 - vig * intensity;
    gl_FragColor = color;
}
