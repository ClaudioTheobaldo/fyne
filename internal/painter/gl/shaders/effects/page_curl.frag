#version 110
uniform sampler2D tex;
uniform float curl;
uniform float radius;
varying vec2 fragTexCoord;
void main() {
    vec2 uv = fragTexCoord;
    float curlPos = 1.0 - curl;
    float r = max(radius, 0.01);

    // Distance from the curl line (diagonal from bottom-right)
    float d = (uv.x + uv.y) * 0.5;

    if (d > curlPos + r) {
        // Behind the curl — show curled-back content (darkened, flipped)
        vec2 flipUV = vec2(1.0 - uv.x, 1.0 - uv.y);
        float dist = d - curlPos;
        float shadow = 1.0 - smoothstep(0.0, r * 2.0, dist) * 0.5;
        vec4 color = texture2D(tex, flipUV) * shadow * 0.6;
        color.a = smoothstep(curlPos + r * 3.0, curlPos + r, d);
        gl_FragColor = color;
    } else if (d > curlPos) {
        // On the curl cylinder — highlight
        float t = (d - curlPos) / r;
        float highlight = 1.0 + sin(t * 3.14159) * 0.3;
        vec4 color = texture2D(tex, uv) * highlight;
        gl_FragColor = color;
    } else {
        // Flat area — normal content with subtle shadow near curl
        float shadow = 1.0 - smoothstep(curlPos - r * 0.5, curlPos, d) * 0.15;
        gl_FragColor = texture2D(tex, uv) * shadow;
    }
}
