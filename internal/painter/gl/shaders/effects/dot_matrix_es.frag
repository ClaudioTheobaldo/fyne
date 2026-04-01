#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 resolution;
uniform float dotSize;
uniform float spacing;
varying vec2 fragTexCoord;
void main() {
    vec2 pixelPos = fragTexCoord * resolution;
    float cellSize = dotSize + spacing;
    vec2 cell = mod(pixelPos, cellSize) - cellSize * 0.5;
    float dist = length(cell);
    vec4 color = texture2D(tex, fragTexCoord);
    float lum = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    float r = dotSize * 0.5 * lum;
    float alpha = smoothstep(r + 0.5, r - 0.5, dist);
    gl_FragColor = vec4(color.rgb * alpha, color.a * alpha);
}
