package farm

import "github.com/gabrielluizsf/sge/town"

// takeAway are the times the towns that sleep were last seen, in the order
// of the towns' own.
func takeAway(f *farm) []awayState {
	var out []awayState
	seen := map[*town.Town]bool{}
	for _, p := range f.plots {
		if at, ok := f.away[p.town]; ok && p.town != nil && !seen[p.town] {
			seen[p.town] = true
			r := takeRegion(p.town)
			out = append(out, awayState{Col: r.Col, Row: r.Row, At: at})
		}
	}
	return out
}
