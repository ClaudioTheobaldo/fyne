#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec3 darkColor;
uniform vec3 lightColor;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float lum = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    color.rgb = mix(darkColor, lightColor, lum);
    gl_FragColor = color;
}
