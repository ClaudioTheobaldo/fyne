#version 100
precision mediump float;
uniform sampler2D tex;
uniform float amplitude;
uniform float frequency;
uniform float phase;
uniform float direction;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord;
    float amp = amplitude * 0.01;
    if (direction < 0.5) {
        // Horizontal wave: x displaced by sin(y)
        uv.x += sin(uv.y * frequency * 6.28318 + phase) * amp;
    } else {
        // Vertical wave: y displaced by sin(x)
        uv.y += sin(uv.x * frequency * 6.28318 + phase) * amp;
    }
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else
        gl_FragColor = texture2D(tex, uv);
}
