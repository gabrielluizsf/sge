package human

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// target is a MakeHuman target: how far it moves some of the base mesh's
// vertices, at full weight.
type target []shift

// shift is how far a target moves vertex v.
type shift struct {
	v int32
	d [3]float64
}

// parseTarget reads a target, a line "index dx dy dz" a vertex, of a mesh
// of n vertices: an index outside it is refused.
func parseTarget(data []byte, n int) (target, error) {
	var t target
	sc := bufio.NewScanner(bytes.NewReader(data))
	for line := 1; sc.Scan(); line++ {
		f := strings.Fields(sc.Text())
		if len(f) == 0 || strings.HasPrefix(f[0], "#") {
			continue
		}
		v, err := strconv.ParseInt(f[0], 10, 32)
		if err != nil || v < 0 || int(v) >= n {
			return nil, fmt.Errorf("human: target line %d: vertex %q of %d", line, f[0], n)
		}
		s := shift{v: int32(v)}
		if err := floats(f[1:], s.d[:]); err != nil {
			return nil, fmt.Errorf("human: target line %d: %w", line, err)
		}
		t = append(t, s)
	}
	return t, sc.Err()
}

// apply moves pos by the target at weight w.
func (t target) apply(pos [][3]float64, w float64) {
	for _, s := range t {
		p := &pos[s.v]
		p[0] += w * s.d[0]
		p[1] += w * s.d[1]
		p[2] += w * s.d[2]
	}
}
