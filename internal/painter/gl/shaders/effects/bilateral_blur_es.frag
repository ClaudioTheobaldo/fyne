#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform float radius;
uniform float sigmaSpace;
uniform float sigmaColor;
varying vec2 fragTexCoord;
void main() {
    vec4 center = texture2D(tex, fragTexCoord);
    vec4 color = vec4(0.0);
    float total = 0.0;
    float ss2 = 2.0 * sigmaSpace * sigmaSpace + 0.001;
    float sc2 = 2.0 * sigmaColor * sigmaColor + 0.001;
    int r = int(min(radius, 8.0));
    for (int x = -8; x <= 8; x++) {
        if (x < -r || x > r) continue;
        for (int y = -8; y <= 8; y++) {
            if (y < -r || y > r) continue;
            vec2 offset = vec2(float(x), float(y)) * texelSize;
            vec4 sample_color = texture2D(tex, fragTexCoord + offset);
            float spatialDist = float(x * x + y * y);
            float colorDist = dot(sample_color.rgb - center.rgb, sample_color.rgb - center.rgb);
            float weight = exp(-spatialDist / ss2) * exp(-colorDist / sc2);
            color += sample_color * weight;
            total += weight;
        }
    }
    gl_FragColor = color / total;
}
