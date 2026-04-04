#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform vec3 lightPos;
uniform float specularConstant;
uniform float specularExponent;
uniform float surfaceScale;
varying vec2 fragTexCoord;
void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    float h = color.a * surfaceScale;
    float hx = texture2D(tex, fragTexCoord + vec2(texelSize.x, 0.0)).a * surfaceScale;
    float hy = texture2D(tex, fragTexCoord + vec2(0.0, texelSize.y)).a * surfaceScale;
    vec3 normal = normalize(vec3(h - hx, h - hy, 1.0));
    vec3 lightDir = normalize(lightPos - vec3(fragTexCoord, h));
    vec3 viewDir = vec3(0.0, 0.0, 1.0);
    vec3 halfVec = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfVec), 0.0), specularExponent) * specularConstant;
    gl_FragColor = vec4(color.rgb + vec3(spec), color.a);
}
