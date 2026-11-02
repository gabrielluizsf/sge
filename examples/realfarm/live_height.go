package main

import "github.com/gabrielluizsf/sge/gltf"

// heightOf is how tall m stands, in metres: the highest point of its
// meshes, glTF's Y up.
func heightOf(m *gltf.Model) float64 {
	top := 0.0
	for _, mesh := range m.Meshes {
		for _, p := range mesh.Primitives {
			for _, v := range p.Positions {
				top = max(top, float64(v[1]))
			}
		}
	}
	return top
}
