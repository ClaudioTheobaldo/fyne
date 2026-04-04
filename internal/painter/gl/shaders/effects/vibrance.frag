#version 110
uniform sampler2D tex;
uniform float vibrance;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float avg = (color.r + color.g + color.b) / 3.0;
    float mx = max(color.r, max(color.g, color.b));
    float sat = mx - avg;
    float amt = (mx - avg) * (-vibrance * 3.0);
    color.rgb = mix(vec3(mx), color.rgb, 1.0 + amt * (1.0 - sat));
    gl_FragColor = color;
}
