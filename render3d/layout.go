package render3d

// A Static or Unlit vertex is StaticFloats floats: its position (3), its
// normal (3) and its texture coordinates (2). A Skinned one is
// SkinnedFloats: those, then four joints, whole numbers, and their four
// weights, which add up to 1.
const (
	StaticFloats  = 8
	SkinnedFloats = 16
)

// A frame's uniforms are FrameFloats floats: matrices column-major, and
// every vector four floats, as std140 and WGSL lay them out.
const (
	FrameViewProj      = 0  // the camera's view and projection
	FrameLightViewProj = 16 // the sun's, for its shadow map
	FrameSun           = 32 // toward the sun, x, y, z, and the shadow's bias
	FrameSunColor      = 36 // the sun's light, linear red, green, blue
	FrameAmbient       = 40 // the ambient light
	FrameEye           = 44 // where the camera is
	FrameFloats        = 48
)

// A material's uniforms are MaterialFloats floats.
const (
	MaterialColor  = 0 // its base colour, linear, and alpha
	MaterialParams = 4 // metallic, roughness, alpha cutoff (0 for none)
	MaterialFloats = 8
)

// Objects are a matrix, ObjectFloats floats, an instance or a joint, at
// most MaxObjects a draw: 64 KiB, what a uniform block may hold. Every
// range of uniforms starts on a multiple of UniformAlign floats, 256
// bytes, as WebGPU asks.
const (
	ObjectFloats = 16
	MaxObjects   = 1024
	UniformAlign = 64
)

// aligned is n rounded up to a multiple of UniformAlign.
func aligned(n int) int { return (n + UniformAlign - 1) / UniformAlign * UniformAlign }
