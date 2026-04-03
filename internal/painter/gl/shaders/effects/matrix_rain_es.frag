#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 resolution;
uniform float time;
uniform float density;
uniform float speed;
uniform float opacity;
varying vec2 fragTexCoord;
float hash(vec2 p) {
    return fract(sin(dot(p, vec2(127.1, 311.7))) * 43758.5453);
}
float char_shape(vec2 p) {
    vec2 cell = fract(p) * vec2(5.0, 7.0);
    vec2 id = floor(cell);
    float h = hash(floor(p) + id * 0.1);
    return step(0.5, h) * step(0.2, fract(cell.x)) * step(0.2, fract(cell.y));
}
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float d = max(density, 1.0);
    vec2 px = fragTexCoord * resolution / d;
    float col = floor(px.x);
    float colSpeed = hash(vec2(col, 0.0)) * 0.5 + 0.5;
    float fall = fract(-time * speed * colSpeed + hash(vec2(col, 1.0)));
    float row = px.y - fall * (resolution.y / d);
    float trail = 1.0 - fract(row * 0.05);
    trail = pow(max(trail, 0.0), 3.0);
    float c = char_shape(vec2(col, floor(row)));
    vec3 green = vec3(0.0, 1.0, 0.3) * c * trail;
    color.rgb = mix(color.rgb, green, opacity * step(0.01, c * trail));
    gl_FragColor = color;
}
