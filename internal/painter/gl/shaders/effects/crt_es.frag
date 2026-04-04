#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 resolution;
uniform float curvature;
uniform float scanlineIntensity;
uniform float vignetteStrength;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord * 2.0 - 1.0;
    uv *= 1.0 + curvature * dot(uv, uv);
    uv = uv * 0.5 + 0.5;
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0) {
        gl_FragColor = vec4(0.0);
        return;
    }
    vec4 color = texture2D(tex, uv);
    float scanline = sin(uv.y * resolution.y * 3.14159) * 0.5 + 0.5;
    color.rgb *= 1.0 - (1.0 - scanline) * scanlineIntensity;
    float vig = 1.0 - dot(uv * 2.0 - 1.0, uv * 2.0 - 1.0) * vignetteStrength;
    color.rgb *= vig;
    gl_FragColor = color;
}
