#version 110
uniform sampler2D tex;
uniform float gamma;
varying vec2 fragTexCoord;
// `gamma` is a uniform, so the reciprocal below is the same for every pixel;
// expressed as a single vec3 so the pow() takes it directly.
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec3 g = vec3(1.0 / max(gamma, 0.01));
    color.rgb = pow(color.rgb, g);
    gl_FragColor = color;
}
