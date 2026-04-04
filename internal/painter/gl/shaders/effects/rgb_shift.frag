#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float amount;
uniform float angle;
varying vec2 fragTexCoord;
void main() {
    float a = angle * 3.14159265 / 180.0;
    vec2 dir = vec2(cos(a), sin(a)) * texelSize * amount;
    float r = texture2D(tex, fragTexCoord + dir).r;
    float g = texture2D(tex, fragTexCoord).g;
    float b = texture2D(tex, fragTexCoord - dir).b;
    float alpha = texture2D(tex, fragTexCoord).a;
    gl_FragColor = vec4(r, g, b, alpha);
}
