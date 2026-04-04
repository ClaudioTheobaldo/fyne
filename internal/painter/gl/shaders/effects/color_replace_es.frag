#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec3 sourceColor;
uniform vec3 targetColor;
uniform float tolerance;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float dist = length(color.rgb - sourceColor);
    float blend = 1.0 - smoothstep(tolerance * 0.8, tolerance, dist);
    color.rgb = mix(color.rgb, targetColor, blend);
    gl_FragColor = color;
}
