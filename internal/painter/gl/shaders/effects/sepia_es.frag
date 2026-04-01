#version 100
precision mediump float;
uniform sampler2D tex;
uniform float amount;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float r = dot(color.rgb, vec3(0.393, 0.769, 0.189));
    float g = dot(color.rgb, vec3(0.349, 0.686, 0.168));
    float b = dot(color.rgb, vec3(0.272, 0.534, 0.131));
    vec3 sepia = vec3(r, g, b);
    color.rgb = mix(color.rgb, sepia, amount);
    gl_FragColor = color;
}
