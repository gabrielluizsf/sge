package human

import (
	"bufio"
	"bytes"
	"strings"
)

// material is what of a MakeHuman material, .mhmat, a bake can use: its
// textures, named from the material's own folder, its colour, and whether
// it is cut out of its texture by its alpha, as hair and lashes are.
type material struct {
	diffuse, normal, occlusion string
	color                      [3]float64
	transparent, twoSided      bool
}

// parseMaterial reads a .mhmat.
func parseMaterial(data []byte) material {
	m := material{color: [3]float64{1, 1, 1}}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) < 2 || strings.HasPrefix(f[0], "#") {
			continue
		}
		rest := strings.Join(f[1:], " ")
		switch f[0] {
		case "diffuseTexture":
			m.diffuse = rest
		case "normalmapTexture":
			m.normal = rest
		case "aomapTexture":
			m.occlusion = rest
		case "diffuseColor":
			var c [3]float64
			if floats(f[1:], c[:]) == nil {
				m.color = c
			}
		case "transparent":
			m.transparent = strings.EqualFold(rest, "true")
		case "backfaceCull":
			m.twoSided = strings.EqualFold(rest, "false")
		}
	}
	return m
}
