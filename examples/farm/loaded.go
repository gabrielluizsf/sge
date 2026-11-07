package farm

import (
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/town"
)

// loaded is what a checked snapshot made, ready to put in a farm: nothing
// of it is the farm's until restore has checked all of the snapshot.
type loaded struct {
	tree   *lineage.Tree
	towns  []*town.Town
	away   map[*town.Town]float64
	plots  []*plot
	people []*villager
	stumps []*stump
	hits   map[int]int
	home   []piece
	yard   []piece
	herd   []*creature
	pens   []*pen
}
