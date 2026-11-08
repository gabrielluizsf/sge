package render3d

// BufferKind is what a buffer holds.
type BufferKind uint8

const (
	// Vertices are floats, in the layout of the pipeline that draws them.
	Vertices BufferKind = iota
	// Indices are uint32s, three a triangle.
	Indices
	// Uniforms are floats bound by Range, each range starting on a
	// multiple of UniformAlign.
	Uniforms
)

// BufferDesc is a buffer to make: of Kind, holding Floats or, for
// indices, Indices; or, with no data, Size floats of zeros.
type BufferDesc struct {
	Kind    BufferKind
	Floats  []float32
	Indices []uint32
	Size    int
}

// Format is how a texture holds its pixels.
type Format uint8

const (
	// RGBA8 is four bytes a pixel, red, green, blue and alpha.
	RGBA8 Format = iota
	// Depth is a depth a pixel, 0 nearest to 1, compared less-than by
	// a shadow lookup.
	Depth
)

// TextureDesc is a texture to make, Width by Height of Format: Pixels its
// RGBA8 bytes, top row first, or nil for a target to draw into. SRGB says
// its colours are sRGB, decoded to linear light when sampled and encoded
// when drawn into; Mipmaps that the backend makes its smaller levels, for
// sampling far away.
type TextureDesc struct {
	Width, Height int
	Format        Format
	Pixels        []byte
	SRGB          bool
	Mipmaps       bool
}

// Kind is a pipeline's built-in shaders. Every kind places a vertex in
// the world by Objects' matrix for the instance drawn (rigid kinds) or by
// its joints' matrices mixed by its weights (skinned kinds), and on the
// screen by the frame's ViewProj, or for shadows its LightViewProj.
type Kind uint8

const (
	// Static is rigid PBR-lite: the base colour times the texture's,
	// cut where its alpha is below the material's cutoff, lit by the
	// ambient light and the sun's (Lambert at least; a backend may add
	// the specular of metal and roughness), the sun's shut out where the
	// shadow map is nearer the sun. Its vertices are StaticFloats each.
	Static Kind = iota
	// Skinned is Static's shading on skinned vertices, SkinnedFloats each.
	Skinned
	// Shadow and ShadowSkinned only write depth, seen from the sun.
	Shadow
	ShadowSkinned
	// Unlit is the base colour times the texture's, unshaded, on Static's
	// vertices.
	Unlit
	kinds
)

// PipelineDesc is a pipeline to make: its Kind, and whether it draws
// both faces of a triangle, or only those turning counterclockwise to the
// camera.
type PipelineDesc struct {
	Kind        Kind
	DoubleSided bool
}

// PassDesc is a pass to begin: drawing into Target (Screen or a texture)
// with Depth (0 for the screen's own); DepthOnly for a pass with no
// colour, Depth its target. Clear clears the colour, ClearColor in linear
// light, and ClearDepth the depth, to 1.
type PassDesc struct {
	Target     Texture
	Depth      Texture
	DepthOnly  bool
	Clear      bool
	ClearColor [4]float32
	ClearDepth bool
}

// Range is At and Count floats of a uniform buffer.
type Range struct {
	Buffer    Buffer
	At, Count int
}

// Bindings are what a draw reads besides its vertices: the Frame's
// uniforms, the Material's and the Objects' (a matrix an instance, or
// one a joint for a skinned draw), laid out as layout.go says; the base
// colour Texture (0 for white) and the sun's Shadow map (0 for none).
type Bindings struct {
	Frame, Material, Objects Range
	Texture, Shadow          Texture
}

// Draw is Count indices from First, three a triangle, drawn Instances
// times, the instance numbering Objects' matrices.
type Draw struct {
	First, Count, Instances int
}
