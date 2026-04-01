#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float offset;
varying vec2 fragTexCoord;
void main() {
    vec2 dir = (fragTexCoord - 0.5) * texelSize * offset;
    float r = texture2D(tex, fragTexCoord + dir).r;
    float g = texture2D(tex, fragTexCoord).g;
    float b = texture2D(tex, fragTexCoord - dir).b;
    float a = texture2D(tex, fragTexCoord).a;
    gl_FragColor = vec4(r, g, b, a);
}
