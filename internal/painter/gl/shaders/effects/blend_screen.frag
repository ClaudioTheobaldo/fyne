#version 110
uniform sampler2D tex;
uniform float blendOpacity;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec3 blended = 1.0 - (1.0 - color.rgb) * (1.0 - color.rgb);
    color.rgb = mix(color.rgb, blended, blendOpacity);
    gl_FragColor = color;
}
