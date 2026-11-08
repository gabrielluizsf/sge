package human

import (
	"encoding/json"
	"fmt"
	"sort"
)

// influence is how much bone moves a vertex.
type influence struct {
	bone   int
	amount float64
}

// parseWeights reads a .mhw, how much each bone moves each vertex of a
// mesh of n vertices, as every vertex's influences: a bone not in the
// skeleton, a vertex outside the mesh or an amount not a number from 0
// to 1 is refused.
func parseWeights(data []byte, s *skeleton, n int) ([][]influence, error) {
	var doc struct {
		Weights map[string][][2]float64
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("human: weights: %w", err)
	}
	out := make([][]influence, n)
	for name, pairs := range doc.Weights {
		b := s.index(name)
		if b < 0 {
			return nil, fmt.Errorf("human: weights of bone %s, not in the skeleton", name)
		}
		for _, p := range pairs {
			v, w := p[0], p[1]
			if v != float64(int32(v)) || v < 0 || int(v) >= n {
				return nil, fmt.Errorf("human: weights of %s: vertex %v of %d", name, v, n)
			}
			if !(w >= 0 && w <= 1) {
				return nil, fmt.Errorf("human: weights of %s: amount %v", name, w)
			}
			out[int(v)] = append(out[int(v)], influence{b, w})
		}
	}
	return out, nil
}

// strongest are the four influences that move a vertex most, as a
// primitive's joints and weights, the weights summing to one.
func strongest(in []influence) ([4]uint16, [4]float32) {
	in = append([]influence(nil), in...)
	sort.Slice(in, func(i, j int) bool {
		if in[i].amount != in[j].amount {
			return in[i].amount > in[j].amount
		}
		return in[i].bone < in[j].bone
	})
	var j [4]uint16
	var w [4]float32
	var sum float64
	for k := 0; k < 4 && k < len(in); k++ {
		sum += in[k].amount
	}
	if sum == 0 {
		return j, [4]float32{1}
	}
	for k := 0; k < 4 && k < len(in); k++ {
		j[k], w[k] = uint16(in[k].bone), float32(in[k].amount/sum)
	}
	return j, w
}

// mixed are the influences of several vertices, each at a share: a
// clothing vertex is moved as the body under it is.
func mixed(all [][]influence, vs [3]int32, shares [3]float64) []influence {
	by := map[int]float64{}
	for k, v := range vs {
		for _, in := range all[v] {
			by[in.bone] += in.amount * shares[k]
		}
	}
	out := make([]influence, 0, len(by))
	for b, a := range by {
		out = append(out, influence{b, a})
	}
	return out
}
