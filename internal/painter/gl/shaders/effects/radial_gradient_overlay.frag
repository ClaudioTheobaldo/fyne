#version 110
uniform sampler2D tex;
uniform vec4 startColor;
uniform vec4 endColor;
uniform vec2 center;
uniform vec2 radius;       // (radiusX, radiusY); (1,1)=circular
uniform float stopCount;   // 0 = legacy 2-stop mode
uniform vec4 stop0Color; uniform float stop0Pos;
uniform vec4 stop1Color; uniform float stop1Pos;
uniform vec4 stop2Color; uniform float stop2Pos;
uniform vec4 stop3Color; uniform float stop3Pos;
uniform vec4 stop4Color; uniform float stop4Pos;
varying vec2 fragTexCoord;

vec4 sampleStops(float t) {
    int n = int(stopCount);
    // Collect into arrays via if-chain (GLSL 110 has no dynamic indexing)
    vec4 colors[5]; float positions[5];
    colors[0] = stop0Color; positions[0] = stop0Pos;
    colors[1] = stop1Color; positions[1] = stop1Pos;
    colors[2] = stop2Color; positions[2] = stop2Pos;
    colors[3] = stop3Color; positions[3] = stop3Pos;
    colors[4] = stop4Color; positions[4] = stop4Pos;
    if (t <= positions[0]) return colors[0];
    for (int i = 1; i < 5; i++) {
        if (i >= n) break;
        if (t <= positions[i]) {
            float span = positions[i] - positions[i-1];
            if (span <= 0.0) return colors[i];
            float d = (t - positions[i-1]) / span;
            return mix(colors[i-1], colors[i], d);
        }
    }
    return colors[n - 1];
}

void main() {
    vec4 color = texture2D(tex, fragTexCoord);
    vec2 r = radius;
    if (r.x <= 0.0) r.x = 1.0;
    if (r.y <= 0.0) r.y = 1.0;
    vec2 delta = (fragTexCoord - center) / r;
    float dist = length(delta) * 1.414;
    float t = clamp(dist, 0.0, 1.0);
    vec4 gradient;
    if (stopCount > 0.5) {
        gradient = sampleStops(t);
    } else {
        gradient = mix(startColor, endColor, t);
    }
    color.rgb = mix(color.rgb, gradient.rgb, gradient.a);
    gl_FragColor = color;
}
