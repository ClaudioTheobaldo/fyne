#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
uniform float levels;
varying vec2 fragTexCoord;
void main() {
    int r = int(min(radius, 4.0));
    int n = int(max(levels, 1.0));
    float bucketSize = 1.0 / float(n);
    vec3 bestColor = vec3(0.0);
    float bestCount = 0.0;
    for (int bucket = 0; bucket < 20; bucket++) {
        if (bucket >= n) break;
        float lo = float(bucket) * bucketSize;
        float hi = lo + bucketSize;
        vec3 sum = vec3(0.0);
        float count = 0.0;
        for (int x = -4; x <= 4; x++) {
            if (x < -r || x > r) continue;
            for (int y = -4; y <= 4; y++) {
                if (y < -r || y > r) continue;
                vec3 c = texture2D(tex, fragTexCoord + vec2(float(x), float(y)) * texelSize).rgb;
                float lum = dot(c, vec3(0.2126, 0.7152, 0.0722));
                if (lum >= lo && lum < hi) {
                    sum += c;
                    count += 1.0;
                }
            }
        }
        if (count > bestCount) {
            bestCount = count;
            bestColor = sum / count;
        }
    }
    gl_FragColor = vec4(bestColor, texture2D(tex, fragTexCoord).a);
}
