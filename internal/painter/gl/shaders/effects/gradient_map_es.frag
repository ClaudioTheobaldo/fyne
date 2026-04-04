#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec3 color0;
uniform vec3 color1;
uniform vec3 color2;
uniform vec3 color3;
uniform vec3 color4;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float lum = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    vec3 mapped;
    if (lum < 0.25)
        mapped = mix(color0, color1, lum * 4.0);
    else if (lum < 0.5)
        mapped = mix(color1, color2, (lum - 0.25) * 4.0);
    else if (lum < 0.75)
        mapped = mix(color2, color3, (lum - 0.5) * 4.0);
    else
        mapped = mix(color3, color4, (lum - 0.75) * 4.0);
    gl_FragColor = vec4(mapped, color.a);
}
