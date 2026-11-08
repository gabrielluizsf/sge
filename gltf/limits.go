package gltf

// Limits bound what a model may make the reader allocate.
type Limits struct {
	// Bytes is the most a file, GLB, glTF JSON, buffer or image, may be.
	Bytes int64
	// Pixels is the most an image may have, width times height.
	Pixels int64
	// Depth is how deep the scene graph may nest.
	Depth int
}

// DefaultLimits are generous for any model a game would bake: 256 MiB a
// file, 64 megapixels an image, 64 levels of nodes.
var DefaultLimits = Limits{Bytes: 256 << 20, Pixels: 64 << 20, Depth: 64}
