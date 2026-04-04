#version 100
precision mediump float;
uniform sampler2D tex;
uniform vec2 texelSize;
varying vec2 fragTexCoord;
void main() {
    float FXAA_REDUCE_MIN = 1.0 / 128.0;
    float FXAA_REDUCE_MUL = 1.0 / 8.0;
    float FXAA_SPAN_MAX = 8.0;
    vec3 luma = vec3(0.299, 0.587, 0.114);
    float lumaTL = dot(texture2D(tex, fragTexCoord + vec2(-1.0, -1.0) * texelSize).rgb, luma);
    float lumaTR = dot(texture2D(tex, fragTexCoord + vec2(1.0, -1.0) * texelSize).rgb, luma);
    float lumaBL = dot(texture2D(tex, fragTexCoord + vec2(-1.0, 1.0) * texelSize).rgb, luma);
    float lumaBR = dot(texture2D(tex, fragTexCoord + vec2(1.0, 1.0) * texelSize).rgb, luma);
    float lumaM  = dot(texture2D(tex, fragTexCoord).rgb, luma);
    vec2 dir;
    dir.x = -((lumaTL + lumaTR) - (lumaBL + lumaBR));
    dir.y = ((lumaTL + lumaBL) - (lumaTR + lumaBR));
    float dirReduce = max((lumaTL + lumaTR + lumaBL + lumaBR) * 0.25 * FXAA_REDUCE_MUL, FXAA_REDUCE_MIN);
    float rcpDirMin = 1.0 / (min(abs(dir.x), abs(dir.y)) + dirReduce);
    dir = min(vec2(FXAA_SPAN_MAX), max(vec2(-FXAA_SPAN_MAX), dir * rcpDirMin)) * texelSize;
    vec4 rgbA = 0.5 * (texture2D(tex, fragTexCoord + dir * (1.0/3.0 - 0.5)) + texture2D(tex, fragTexCoord + dir * (2.0/3.0 - 0.5)));
    vec4 rgbB = rgbA * 0.5 + 0.25 * (texture2D(tex, fragTexCoord + dir * -0.5) + texture2D(tex, fragTexCoord + dir * 0.5));
    float lumaB = dot(rgbB.rgb, luma);
    float lumaMin = min(lumaM, min(min(lumaTL, lumaTR), min(lumaBL, lumaBR)));
    float lumaMax = max(lumaM, max(max(lumaTL, lumaTR), max(lumaBL, lumaBR)));
    if (lumaB < lumaMin || lumaB > lumaMax)
        gl_FragColor = rgbA;
    else
        gl_FragColor = rgbB;
}
