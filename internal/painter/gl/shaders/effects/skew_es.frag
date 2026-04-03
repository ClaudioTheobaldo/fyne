#version 100
precision mediump float;
uniform sampler2D tex;
uniform float skewX;
uniform float skewY;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord;
    // Skew by shifting UV based on the other axis, centered at 0.5
    uv.x += (uv.y - 0.5) * skewX;
    uv.y += (uv.x - 0.5) * skewY;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else
        gl_FragColor = texture2D(tex, uv);
}
