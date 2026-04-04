#version 110
uniform sampler2D tex;
uniform float time;
uniform float speed;
uniform float brightMin;
uniform float brightMax;
uniform float scaleMin;
uniform float scaleMax;
varying vec2 fragTexCoord;
void main() {
    float t = sin(time * speed * 6.28318) * 0.5 + 0.5;
    float bright = mix(brightMin, brightMax, t);
    float sc = mix(scaleMin, scaleMax, t);
    vec2 uv = (fragTexCoord - 0.5) / max(sc, 0.01) + 0.5;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0)
        gl_FragColor = vec4(0.0);
    else {
        vec4 color = texture2D(tex, uv);
        color.rgb *= bright;
        gl_FragColor = color;
    }
}
