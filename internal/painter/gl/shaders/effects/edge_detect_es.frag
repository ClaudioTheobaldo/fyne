#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float threshold;
varying vec2 fragTexCoord;
void main() {
    float tl = dot(texture2D(tex, fragTexCoord + vec2(-texelSize.x, texelSize.y)).rgb, vec3(0.333));
    float t  = dot(texture2D(tex, fragTexCoord + vec2(0.0, texelSize.y)).rgb, vec3(0.333));
    float tr = dot(texture2D(tex, fragTexCoord + vec2(texelSize.x, texelSize.y)).rgb, vec3(0.333));
    float l  = dot(texture2D(tex, fragTexCoord + vec2(-texelSize.x, 0.0)).rgb, vec3(0.333));
    float r  = dot(texture2D(tex, fragTexCoord + vec2(texelSize.x, 0.0)).rgb, vec3(0.333));
    float bl = dot(texture2D(tex, fragTexCoord + vec2(-texelSize.x, -texelSize.y)).rgb, vec3(0.333));
    float b  = dot(texture2D(tex, fragTexCoord + vec2(0.0, -texelSize.y)).rgb, vec3(0.333));
    float br = dot(texture2D(tex, fragTexCoord + vec2(texelSize.x, -texelSize.y)).rgb, vec3(0.333));
    float gx = -tl - 2.0*l - bl + tr + 2.0*r + br;
    float gy = -tl - 2.0*t - tr + bl + 2.0*b + br;
    float edge = sqrt(gx*gx + gy*gy);
    edge = step(threshold, edge);
    gl_FragColor = vec4(vec3(edge), texture2D(tex, fragTexCoord).a);
}
