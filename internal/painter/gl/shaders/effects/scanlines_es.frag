#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 resolution;
uniform float density;
uniform float opacity;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float line = mod(fragTexCoord.y * resolution.y * density, 2.0);
    float scanline = smoothstep(0.0, 1.0, line);
    color.rgb *= 1.0 - (1.0 - scanline) * opacity;
    gl_FragColor = color;
}
