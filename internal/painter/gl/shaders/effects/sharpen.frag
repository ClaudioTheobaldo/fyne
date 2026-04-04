#version 110
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float strength;
varying vec2 fragTexCoord;
void main() {
    vec4 center = texture2D(tex, fragTexCoord);
    vec4 top    = texture2D(tex, fragTexCoord + vec2(0.0, texelSize.y));
    vec4 bottom = texture2D(tex, fragTexCoord - vec2(0.0, texelSize.y));
    vec4 left   = texture2D(tex, fragTexCoord - vec2(texelSize.x, 0.0));
    vec4 right  = texture2D(tex, fragTexCoord + vec2(texelSize.x, 0.0));
    vec4 sharpened = center * (1.0 + 4.0 * strength) - (top + bottom + left + right) * strength;
    gl_FragColor = clamp(sharpened, 0.0, 1.0);
}
