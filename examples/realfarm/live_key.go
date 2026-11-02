package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

// shapeKey tells characters whose people have the same mesh from those
// whose have not: the body and what is worn decide the mesh; the skin
// tone, the hair, eye and clothes colours, and the material of the eyes,
// only the textures.
func shapeKey(c kit.Character) string {
	return fmt.Sprintf("%v|%q|%q|%t", c.Body, c.Look.Wear, c.Look.Skin, c.Look.Eyes == "")
}
