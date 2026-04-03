#version 110
uniform sampler2D tex;
uniform float scaleX;
uniform float scaleY;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord - 0.5;
    uv.x /= max(scaleX, 0.01);
    uv.y /= max(scaleY, 0.01);
    uv += 0.5;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else
        gl_FragColor = texture2D(tex, uv);
}
