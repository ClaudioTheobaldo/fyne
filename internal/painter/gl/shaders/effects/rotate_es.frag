#version 100
precision mediump float;
uniform sampler2D tex;
uniform float angle;
varying vec2 fragTexCoord;
void main() {
    float a = angle * 3.14159265 / 180.0;
    float ca = cos(a);
    float sa = sin(a);
    vec2 uv = fragTexCoord - 0.5;
    uv = vec2(uv.x * ca - uv.y * sa, uv.x * sa + uv.y * ca);
    uv += 0.5;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else
        gl_FragColor = texture2D(tex, uv);
}
