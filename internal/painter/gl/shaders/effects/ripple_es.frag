#version 100
precision mediump float;
uniform sampler2D tex;
uniform float amplitude;
uniform float frequency;
uniform float speed;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord;
    uv.x += sin(uv.y * frequency + speed) * amplitude * 0.01;
    uv.y += cos(uv.x * frequency + speed) * amplitude * 0.01;
    gl_FragColor = texture2D(tex, uv);
}
