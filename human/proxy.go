package human

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// proxy is a MakeHuman proxy, .mhclo: a mesh worn on the body, clothes,
// hair, eyes or brows, each of its vertices placed by three of the body's
// and an offset, which grows with the body.
type proxy struct {
	obj, material string
	refs          []ref
	scales        [3]scaleBy
	hides         []int32
}

// ref places a proxy's vertex: at the body's vertices v, by shares w,
// moved by d, scaled as the body is.
type ref struct {
	v [3]int32
	w [3]float64
	d [3]float64
}

// scaleBy is how much the body has grown along an axis: the distance
// between two of its vertices along it, over what it was when the proxy
// was fitted.
type scaleBy struct {
	a, b int32
	was  float64
}

// parseProxy reads a .mhclo for a body of n vertices: every vertex it
// names must be in the body.
func parseProxy(data []byte, n int) (*proxy, error) {
	p := &proxy{scales: [3]scaleBy{{was: 0}, {was: 0}, {was: 0}}}
	mode := ""
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64<<10), 1<<20)
	for line := 1; sc.Scan(); line++ {
		f := strings.Fields(sc.Text())
		if len(f) == 0 || strings.HasPrefix(f[0], "#") {
			continue
		}
		var err error
		switch {
		case f[0] == "obj_file" && len(f) > 1:
			p.obj = strings.Join(f[1:], " ")
		case f[0] == "material" && len(f) > 1:
			p.material = strings.Join(f[1:], " ")
		case strings.HasSuffix(f[0], "_scale") && len(f[0]) == 7:
			err = p.scale(f, n)
		case f[0] == "verts" || f[0] == "delete_verts" || f[0] == "weights":
			mode = f[0]
		case mode == "verts" && isNumber(f[0]):
			var r ref
			r, err = parseRef(f, n)
			p.refs = append(p.refs, r)
		case mode == "delete_verts" && isNumber(f[0]):
			p.hides, err = hidden(p.hides, f, n)
		}
		if err != nil {
			return nil, fmt.Errorf("human: proxy line %d: %w", line, err)
		}
	}
	if p.obj == "" || len(p.refs) == 0 {
		return nil, fmt.Errorf("human: a proxy with no mesh")
	}
	return p, sc.Err()
}

// scale reads "x_scale a b was" and its y and z.
func (p *proxy) scale(f []string, n int) error {
	axis := strings.IndexByte("xyz", f[0][0])
	if axis < 0 || len(f) < 4 {
		return fmt.Errorf("%s", f[0])
	}
	a, err1 := vertex(f[1], n)
	b, err2 := vertex(f[2], n)
	var was [1]float64
	if err := floats(f[3:], was[:]); err != nil || err1 != nil || err2 != nil || was[0] <= 0 {
		return fmt.Errorf("%s", strings.Join(f, " "))
	}
	p.scales[axis] = scaleBy{a, b, was[0]}
	return nil
}

// parseRef reads a vertex's place: one body vertex, or three, their
// shares and an offset.
func parseRef(f []string, n int) (ref, error) {
	switch len(f) {
	case 1:
		v, err := vertex(f[0], n)
		return ref{v: [3]int32{v, v, v}, w: [3]float64{1}}, err
	case 6, 9:
		var r ref
		for k := range 3 {
			v, err := vertex(f[k], n)
			if err != nil {
				return r, err
			}
			r.v[k] = v
		}
		if err := floats(f[3:6], r.w[:]); err != nil {
			return r, err
		}
		if len(f) == 9 {
			// Read before r is returned: Go leaves unsaid whether r would be
			// read before or after a call in the same return fills it.
			if err := floats(f[6:9], r.d[:]); err != nil {
				return r, err
			}
		}
		return r, nil
	}
	return ref{}, fmt.Errorf("a vertex of %d numbers", len(f))
}

// hidden adds the body's vertices a line of delete_verts names, alone or
// "a - b" runs, to hides.
func hidden(hides []int32, f []string, n int) ([]int32, error) {
	run := false
	var last int32
	for _, w := range f {
		if w == "-" {
			run = true
			continue
		}
		v, err := vertex(w, n)
		if err != nil {
			return nil, err
		}
		if run {
			if v < last {
				return nil, fmt.Errorf("a run from %d down to %d", last, v)
			}
			for u := last + 1; u <= v; u++ {
				hides = append(hides, u)
			}
		} else {
			hides = append(hides, v)
		}
		run, last = false, v
	}
	return hides, nil
}

// vertex is a body vertex's index, which must be in a body of n.
func vertex(s string, n int) (int32, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil || v < 0 || int(v) >= n {
		return 0, fmt.Errorf("vertex %q of %d", s, n)
	}
	return int32(v), nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseInt(s, 10, 32)
	return err == nil
}

// fit places the proxy's vertices on the body at pos.
func (p *proxy) fit(pos [][3]float64) [][3]float64 {
	var grown [3]float64
	for k, s := range p.scales {
		grown[k] = 1
		if s.was > 0 {
			grown[k] = abs(pos[s.a][k]-pos[s.b][k]) / s.was
		}
	}
	out := make([][3]float64, len(p.refs))
	for i, r := range p.refs {
		var q [3]float64
		for k := range 3 {
			q = add(q, scale(pos[r.v[k]], r.w[k]))
		}
		for k := range 3 {
			q[k] += r.d[k] * grown[k]
		}
		out[i] = q
	}
	return out
}
