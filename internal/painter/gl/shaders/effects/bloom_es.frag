#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float threshold;
uniform float intensity;
uniform float radius;
varying vec2 fragTexCoord;
void main() {
    vec4 original = texture2D(tex, fragTexCoord);
    vec4 bright = vec4(0.0);
    float total = 0.0;
    int r = int(min(radius, 12.0));
    for (int x = -12; x <= 12; x++) {
        if (x < -r || x > r) continue;
        for (int y = -12; y <= 12; y++) {
            if (y < -r || y > r) continue;
            vec4 s = texture2D(tex, fragTexCoord + vec2(float(x), float(y)) * texelSize);
            float lum = dot(s.rgb, vec3(0.2126, 0.7152, 0.0722));
            if (lum > threshold) {
                float fi = float(x * x + y * y);
                float weight = exp(-fi / (2.0 * radius * radius + 0.001));
                bright += (s - vec4(vec3(threshold), 0.0)) * weight;
                total += weight;
            }
        }
    }
    if (total > 0.0) bright /= total;
    gl_FragColor = original + bright * intensity;
}
