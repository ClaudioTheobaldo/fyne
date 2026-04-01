#version 110
uniform sampler2D tex;
uniform vec2 resolution;
uniform float dotSize;
uniform float angle;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float lum = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    float a = angle * 3.14159265 / 180.0;
    float ca = cos(a);
    float sa = sin(a);
    vec2 pixelPos = fragTexCoord * resolution;
    vec2 rotated = vec2(pixelPos.x * ca - pixelPos.y * sa, pixelPos.x * sa + pixelPos.y * ca);
    vec2 grid = mod(rotated, dotSize) - dotSize * 0.5;
    float dist = length(grid) / (dotSize * 0.5);
    float dot_val = step(dist, lum);
    gl_FragColor = vec4(vec3(dot_val), color.a);
}
