#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform vec3 lightPos;
uniform float diffuseConstant;
uniform float surfaceScale;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float h = color.a * surfaceScale;
    float hx = texture2D(tex, fragTexCoord + vec2(texelSize.x, 0.0)).a * surfaceScale;
    float hy = texture2D(tex, fragTexCoord + vec2(0.0, texelSize.y)).a * surfaceScale;
    vec3 normal = normalize(vec3(h - hx, h - hy, 1.0));
    vec3 lightDir = normalize(lightPos - vec3(fragTexCoord, h));
    float diff = max(dot(normal, lightDir), 0.0) * diffuseConstant;
    gl_FragColor = vec4(color.rgb * diff, color.a);
}
