package farm

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/sge/terrain"
)

// RockVariants is how many rocks a biome has, each a tile of its own, as
// "rock_cerrado_2": the first a small rock, the second a boulder and the
// third a few rocks together. The biomes are the trees'.
const RockVariants = 3

// rockName is the tile of the rock of biome and variant.
func rockName(biome string, variant int) string {
	return fmt.Sprint("rock_", biome, "_", variant)
}

// RockVariant is the variant the tile name says when it is a rock's, and
// whether it is.
func RockVariant(name string) (variant int, ok bool) {
	if !strings.HasPrefix(name, "rock_") {
		return 0, false
	}
	i := strings.LastIndex(name, "_")
	_, err := fmt.Sscanf(name[i+1:], "%d", &variant)
	return variant, err == nil && variant >= 0 && variant < RockVariants
}

// RockBiome is the biome the tile name of a rock says.
func RockBiome(name string) string {
	return strings.TrimPrefix(name[:strings.LastIndex(name, "_")], "rock_")
}

// rockIDs are the ids of the tile of every rock standing.
var rockIDs = map[int]bool{}

// isRock reports whether tile id is a rock standing.
func isRock(id int) bool { return rockIDs[id] }

// rockShare is how thick the rocks lie in each biome, as the chance a cell
// is a rock's before they are spaced apart: dry land has most.
var rockShare = []float64{0.012, 0.02, 0.025, 0.008, 0.035}

// rockAt is the tile of the rock on cell col, row: of the variants its
// biome has, a small rock the most often, always the same.
func rockAt(col, row int) int {
	biome := TreeBiomes[max(climate.At(col, row), 0)%len(TreeBiomes)]
	v := 0
	switch c := terrain.Chance(seed+12, col, row); {
	case c > 0.9:
		v = 2
	case c > 0.55:
		v = 1
	}
	return tileID[rockName(biome, v)]
}

// rockCeiling is the thickest any biome lays its rocks, which the others
// thin out from.
const rockCeiling = 0.035

// rockGap is how far apart rocks lie: a look with its pictures larger
// than a cell spaces them wider.
func rockGap(a Art) int {
	if a.Reach() > 0 {
		return 3
	}
	return 2
}

// rockSpace is the share of the biome at col, row that rocks lie in.
func rockSpace(col, row int) float64 {
	return rockShare[max(climate.At(col, row), 0)%len(rockShare)]
}
