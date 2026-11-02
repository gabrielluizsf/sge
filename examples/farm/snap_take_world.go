package farm

import (
	"slices"

	"github.com/gabrielluizsf/sge/town"
)

// takePlots are the plots, in order.
func takePlots(f *farm) []plotState {
	out := make([]plotState, 0, len(f.plots))
	for _, p := range f.plots {
		out = append(out, plotState{Col: p.col, Row: p.row, Town: takeRegion(p.town), Crop: p.crop, Stage: p.stage, Days: p.days, Watered: p.watered, Dead: p.dead})
	}
	return out
}

// takeRegion is t's region, nil for no town.
func takeRegion(t *town.Town) *cellState {
	if t == nil {
		return nil
	}
	return &cellState{Col: t.Region.Col, Row: t.Region.Row}
}

// takeStumps are the stumps.
func takeStumps(f *farm) []stumpState {
	out := make([]stumpState, 0, len(f.stumps))
	for _, s := range f.stumps {
		out = append(out, stumpState{Col: s.col, Row: s.row, Tree: s.tree, Left: s.left})
	}
	return out
}

// takeHits are the strokes of the trees being cut, in the order of their
// cells so a save is the same every time.
func takeHits(f *farm) []hitState {
	keys := make([]int, 0, len(f.hits))
	for k := range f.hits {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	out := make([]hitState, 0, len(keys))
	for _, k := range keys {
		col, row := cellOfKey(k)
		out = append(out, hitState{Col: col, Row: row, N: f.hits[k]})
	}
	return out
}

// takeTowns are the regions of the towns woken, in the order of their
// gardens' plots, which is the order they were woken in.
func takeTowns(f *farm) []cellState {
	out := []cellState{}
	seen := map[*town.Town]bool{}
	for _, p := range f.plots {
		if p.town != nil && !seen[p.town] {
			seen[p.town] = true
			out = append(out, *takeRegion(p.town))
		}
	}
	return out
}
