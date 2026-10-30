// Package gltf reads 3D models in glTF 2.0, binary (.glb) or not: their
// meshes, materials and textures, nodes, skins and animations. It reads
// files from anywhere as untrusted: every offset, count and index is
// checked against what it points into, sizes are capped, and the files a
// model names are only opened inside its own directory.
package gltf

import "image"

// Model is a model read: its scene's root nodes, and everything they use.
type Model struct {
	Roots      []int
	Nodes      []Node
	Meshes     []Mesh
	Materials  []Material
	Images     []image.Image
	Skins      []Skin
	Animations []Animation
}

// Node is a point of the scene graph: its children, the mesh and skin it
// shows, if any (-1 otherwise), and where it is in its parent, as a
// translation, rotation (a quaternion x, y, z, w) and scale.
type Node struct {
	Name        string
	Children    []int
	Mesh, Skin  int
	Translation [3]float64
	Rotation    [4]float64
	Scale       [3]float64
}

// Mesh is a mesh's primitives, triangles each.
type Mesh struct {
	Name       string
	Primitives []Primitive
}

// Primitive is triangles of a material: its vertices' positions, normals,
// texture coordinates, joints and weights (those it has), and their
// indices, three a triangle.
type Primitive struct {
	Positions [][3]float32
	Normals   [][3]float32
	UVs       [][2]float32
	Joints    [][4]uint16
	Weights   [][4]float32
	Indices   []uint32
	Material  int // -1 for the default
}

// Material is how a surface looks: its base colour, times its texture's
// when it has one (an index of Images, or -1), how metallic and rough it
// is, and how its alpha is taken: "OPAQUE", "MASK" below AlphaCutoff, or
// "BLEND".
type Material struct {
	Name                string
	BaseColor           [4]float64
	Texture             int
	Metallic, Roughness float64
	AlphaMode           string
	AlphaCutoff         float64
	DoubleSided         bool
}

// Skin is the joints a skinned mesh bends with, nodes each, and the
// inverse of each one's bind matrix, column-major.
type Skin struct {
	Joints      []int
	InverseBind [][16]float64
}

// Animation is channels moving nodes over time.
type Animation struct {
	Name     string
	Channels []Channel
}

// Channel moves a node's "translation", "rotation" or "scale" through
// keyframes: their times, and their values, three or four a key, by
// "LINEAR" or "STEP" interpolation.
type Channel struct {
	Node          int
	Path          string
	Interpolation string
	Times         []float32
	Values        [][4]float32
}
