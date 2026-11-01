//go:build linux && amd64

package gl

// The uniform blocks are render3d's layout.go in std140: the frame's at
// binding 0, the material's at 1, the objects' at 2. A kind's shaders are
// these, with SKINNED, SHADOW or UNLIT defined.
const frameBlock = `
layout(std140, binding = 0) uniform Frame {
	mat4 viewProj;
	mat4 lightViewProj;
	vec4 sun; // toward the sun, and the shadow's bias
	vec4 sunColor;
	vec4 ambient;
	vec4 eye;
};
`

// vertexShader places a vertex in the world by its instance's matrix,
// or its joints' mixed, and in GL's clip space, its depth moved from
// render3d's 0 to 1 into GL's -1 to 1.
const vertexShader = frameBlock + `
layout(std140, binding = 2) uniform Objects { mat4 objects[1024]; };
layout(location = 0) in vec3 position;
layout(location = 1) in vec3 normal;
layout(location = 2) in vec2 uv;
#ifdef SKINNED
layout(location = 3) in vec4 joints;
layout(location = 4) in vec4 weights;
#endif
out vec3 vWorld;
out vec3 vNormal;
out vec2 vUV;

void main() {
#ifdef SKINNED
	mat4 m = objects[int(joints.x)] * weights.x + objects[int(joints.y)] * weights.y +
		objects[int(joints.z)] * weights.z + objects[int(joints.w)] * weights.w;
#else
	mat4 m = objects[gl_InstanceID];
#endif
	vec4 world = m * vec4(position, 1.0);
	vWorld = world.xyz;
	vNormal = (m * vec4(normal, 0.0)).xyz;
	vUV = uv;
#ifdef SHADOW
	vec4 clip = lightViewProj * world;
#else
	vec4 clip = viewProj * world;
#endif
	clip.z = 2.0 * clip.z - clip.w;
	gl_Position = clip;
}
`

// fragmentShader is soft's shading: the base colour times the texture's,
// cut under the material's cutoff, lit by the ambient light and, by
// Lambert's law, the sun's, but where its shadow map is nearer the sun.
const fragmentShader = frameBlock + `
layout(std140, binding = 1) uniform Material {
	vec4 color;
	vec4 params; // metallic, roughness, alpha cutoff
};
layout(binding = 0) uniform sampler2D base;
layout(binding = 1) uniform sampler2D shadowMap;
in vec3 vWorld;
in vec3 vNormal;
in vec2 vUV;
#ifndef SHADOW
layout(location = 0) out vec4 outColor;
#endif

float lit(vec3 world) {
	vec4 c = lightViewProj * vec4(world, 1.0);
	vec3 p = c.xyz / c.w;
	ivec2 size = textureSize(shadowMap, 0);
	ivec2 at = ivec2(floor((p.xy + 1.0) * 0.5 * vec2(size)));
	if (at.x < 0 || at.y < 0 || at.x >= size.x || at.y >= size.y || p.z > 1.0) {
		return 1.0;
	}
	return p.z - sun.w > texelFetch(shadowMap, at, 0).r ? 0.0 : 1.0;
}

void main() {
	vec4 albedo = color * texture(base, vUV);
	if (params.z > 0.0 && albedo.a < params.z) {
		discard;
	}
#ifndef SHADOW
#ifdef UNLIT
	outColor = vec4(albedo.rgb, 1.0);
#else
	vec3 n = normalize(vNormal);
	if (!gl_FrontFacing) {
		n = -n;
	}
	float s = max(0.0, dot(n, sun.xyz)) * lit(vWorld);
	outColor = vec4(albedo.rgb * (ambient.rgb + sunColor.rgb * s), 1.0);
#endif
#endif
}
`
