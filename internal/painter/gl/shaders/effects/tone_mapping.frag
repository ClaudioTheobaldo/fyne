#version 110
uniform sampler2D tex;
uniform float exposure;
uniform float gamma;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    // Reinhard tone mapping with exposure
    color.rgb *= exposure;
    color.rgb = color.rgb / (1.0 + color.rgb);
    // Gamma correction
    float g = 1.0 / max(gamma, 0.01);
    color.rgb = pow(color.rgb, vec3(g));
    gl_FragColor = color;
}
