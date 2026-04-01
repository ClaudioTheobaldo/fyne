#version 100
precision mediump float;
uniform sampler2D tex;
uniform float distortion;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord * 2.0 - 1.0;
    float r2 = dot(uv, uv);
    uv *= 1.0 + distortion * r2;
    uv = uv * 0.5 + 0.5;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else
        gl_FragColor = texture2D(tex, uv);
}
