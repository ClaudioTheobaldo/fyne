#version 110
uniform sampler2D tex;
uniform float threshold;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float lum = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    float v = step(threshold, lum);
    gl_FragColor = vec4(v, v, v, color.a);
}
