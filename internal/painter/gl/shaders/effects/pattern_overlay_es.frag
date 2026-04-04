#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 resolution;
uniform float patternSize;
uniform float patternType;
uniform float opacity;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec2 px = fragTexCoord * resolution;
    float pattern = 0.0;
    if (patternType < 0.5) {
        float cx = step(0.5, fract(px.x / patternSize));
        float cy = step(0.5, fract(px.y / patternSize));
        pattern = abs(cx - cy);
    } else if (patternType < 1.5) {
        pattern = step(0.5, fract(px.y / patternSize));
    } else if (patternType < 2.5) {
        pattern = step(0.5, fract(px.x / patternSize));
    } else {
        vec2 cell = mod(px, patternSize) - patternSize * 0.5;
        pattern = 1.0 - step(patternSize * 0.3, length(cell));
    }
    color.rgb = mix(color.rgb, vec3(pattern), opacity);
    gl_FragColor = color;
}
