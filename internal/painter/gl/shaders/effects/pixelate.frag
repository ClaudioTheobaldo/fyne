#version 110
uniform sampler2D tex;
uniform vec2 resolution;
uniform float pixelSize;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord;
    vec2 pixels = resolution / pixelSize;
    uv = floor(uv * pixels) / pixels;
    gl_FragColor = texture2D(tex, uv);
}
