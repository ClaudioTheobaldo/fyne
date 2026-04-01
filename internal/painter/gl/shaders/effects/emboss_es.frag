#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float strength;
varying vec2 fragTexCoord;
void main() {
    vec4 tl = texture2D(tex, fragTexCoord + vec2(-texelSize.x, texelSize.y));
    vec4 br = texture2D(tex, fragTexCoord + vec2(texelSize.x, -texelSize.y));
    vec4 emboss = (br - tl) * strength + vec4(0.5, 0.5, 0.5, 1.0);
    emboss.a = texture2D(tex, fragTexCoord).a;
    gl_FragColor = emboss;
}
