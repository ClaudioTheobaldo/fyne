#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec3 redOut;
uniform vec3 greenOut;
uniform vec3 blueOut;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float r = dot(color.rgb, redOut);
    float g = dot(color.rgb, greenOut);
    float b = dot(color.rgb, blueOut);
    gl_FragColor = vec4(r, g, b, color.a);
}
