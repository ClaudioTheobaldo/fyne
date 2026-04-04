#version 100
precision mediump float;
uniform sampler2D tex;
uniform float amount;
uniform float direction;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord;
    // direction: 0=bottom narrows, 1=top narrows, 2=left narrows, 3=right narrows
    float t;
    if (direction < 0.5) {
        t = uv.y;
        float s = mix(1.0, 1.0 - amount, t);
        uv.x = 0.5 + (uv.x - 0.5) / max(s, 0.01);
    } else if (direction < 1.5) {
        t = 1.0 - uv.y;
        float s = mix(1.0, 1.0 - amount, t);
        uv.x = 0.5 + (uv.x - 0.5) / max(s, 0.01);
    } else if (direction < 2.5) {
        t = 1.0 - uv.x;
        float s = mix(1.0, 1.0 - amount, t);
        uv.y = 0.5 + (uv.y - 0.5) / max(s, 0.01);
    } else {
        t = uv.x;
        float s = mix(1.0, 1.0 - amount, t);
        uv.y = 0.5 + (uv.y - 0.5) / max(s, 0.01);
    }
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else
        gl_FragColor = texture2D(tex, uv);
}
