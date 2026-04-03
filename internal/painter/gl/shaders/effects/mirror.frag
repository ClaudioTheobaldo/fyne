#version 110
uniform sampler2D tex;
uniform float flipX;
uniform float flipY;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord;
    if (flipX > 0.5) uv.x = 1.0 - uv.x;
    if (flipY > 0.5) uv.y = 1.0 - uv.y;
    gl_FragColor = texture2D(tex, uv);
}
