package farm

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
)

// TreeBiomes are the kinds of land the wild's trees stand in, as a tree's
// tile names them: "tree_" biome, a variant, and a frame of its sway, as
// "tree_cerrado_4_0". Each biome has its own trees: the art says which
// species, there being many a look draws.
var TreeBiomes = []string{"campo", "pinhal", "cerrado", "tropical", "caatinga"}

// TreeVariants is how many trees a biome has, every one of them a tile of
// TreeFrames frames of sway. Variants 0 to 2 are saplings, 3 to 5 young
// trees and 6 to 8 trees grown, three of each size, each a tree of its
// own shape: a variant's size is its number over three, and its seed what
// is left.
const (
	TreeVariants = 9
	TreeFrames   = 4
)

// treeName is the tile of the tree of biome and variant, frame f of its
// sway.
func treeName(biome string, variant, f int) string {
	return fmt.Sprint("tree_", biome, "_", variant, "_", f)
}

// TreeFrame is the sway frame of the tile name when it is a tree's, and
// whether it is.
func TreeFrame(name string) (frame int, ok bool) {
	if !strings.HasPrefix(name, "tree_") {
		return 0, false
	}
	i := strings.LastIndex(name, "_")
	_, err := fmt.Sscanf(name[i+1:], "%d", &frame)
	return frame, err == nil
}

// treeIDs are the ids of the tile of every tree standing, the first of
// its frames, by id.
var treeIDs = map[int]bool{}

// isTree reports whether tile id is a tree standing.
func isTree(id int) bool { return treeIDs[id] }

// The wild's climate, the biomes it lays out in large patches, and the
// shares of its trees that are saplings and young trees; the rest are
// grown.
var (
	climate = terrain.Climate{Seed: seed + 5, Scale: 70, Ragged: 0.08, Biomes: []terrain.Biome{
		{Name: "campo", Warmth: 0.5, Wetness: 0.5},
		{Name: "pinhal", Warmth: 0.15, Wetness: 0.75},
		{Name: "cerrado", Warmth: 0.7, Wetness: 0.35},
		{Name: "tropical", Warmth: 0.9, Wetness: 0.9},
		{Name: "caatinga", Warmth: 0.9, Wetness: 0.05},
	}}
	saplings, young = 0.15, 0.35
)

// treeAt is the tile of the tree growing on cell col, row: its biome's, of
// a size and a shape its cell chooses, always the same.
func treeAt(col, row int) int {
	biome := TreeBiomes[max(climate.At(col, row), 0)%len(TreeBiomes)]
	size := 2
	if c := terrain.Chance(seed+9, col, row); c < saplings {
		size = 0
	} else if c < saplings+young {
		size = 1
	}
	shape := int(terrain.Chance(seed+10, col, row) * 3)
	return tileID[treeName(biome, size*3+min(shape, 2), 0)]
}

// plantTrees grows a tree of its biome on every cell of m a map marked
// with a stump, the farm's trees.
func plantTrees(m *tilemap.Map) {
	for row := range m.Rows {
		for col := range m.Cols {
			if m.At(col, row) == tileID["stump"] {
				m.Set(col, row, treeAt(col, row))
			}
		}
	}
}
