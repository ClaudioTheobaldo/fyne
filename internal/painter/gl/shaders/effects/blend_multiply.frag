#version 110
uniform sampler2D tex;
uniform float blendOpacity;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec3 blended = color.rgb * color.rgb;
    color.rgb = mix(color.rgb, blended, blendOpacity);
    gl_FragColor = color;
}
