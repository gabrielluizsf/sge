package human

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// mesh is a Wavefront OBJ mesh as MakeHuman writes them: positions,
// texture places, and faces, each in a group.
type mesh struct {
	pos    [][3]float64
	uv     [][2]float64
	faces  []face
	groups []string
}

// face is a polygon: its corners' positions and texture places, indices
// of the mesh's, and its group.
type face struct {
	v, t  []int32
	group int
}

// maxCorners is the most corners a face may have.
const maxCorners = 64

// parseOBJ reads an OBJ mesh's positions, texture places, faces and
// groups, and nothing else. After CVE-2020-28589, every index a face
// gives is checked against what the file has: none is negative, as OBJ's
// relative ones are, or past the end.
func parseOBJ(data []byte) (*mesh, error) {
	m := &mesh{groups: []string{""}}
	group := 0
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64<<10), 1<<20)
	for line := 1; sc.Scan(); line++ {
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		var err error
		switch f[0] {
		case "v":
			var p [3]float64
			err = floats(f[1:], p[:])
			m.pos = append(m.pos, p)
		case "vt":
			var t [2]float64
			err = floats(f[1:], t[:])
			m.uv = append(m.uv, t)
		case "g":
			m.groups = append(m.groups, strings.Join(f[1:], " "))
			group = len(m.groups) - 1
		case "f":
			var fc face
			fc, err = parseFace(f[1:])
			fc.group = group
			m.faces = append(m.faces, fc)
		}
		if err != nil {
			return nil, fmt.Errorf("human: OBJ line %d: %w", line, err)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return m, m.check()
}

// parseFace reads a face's corners, "v", "v/t", "v/t/n" or "v//n" each,
// counted from 1 as OBJ counts them, and keeps them from 0.
func parseFace(corners []string) (face, error) {
	if len(corners) < 3 || len(corners) > maxCorners {
		return face{}, fmt.Errorf("a face of %d corners", len(corners))
	}
	var fc face
	for _, c := range corners {
		parts := strings.Split(c, "/")
		v, err := index(parts[0])
		if err != nil {
			return face{}, err
		}
		fc.v = append(fc.v, v)
		if len(parts) > 1 && parts[1] != "" {
			t, err := index(parts[1])
			if err != nil {
				return face{}, err
			}
			fc.t = append(fc.t, t)
		}
	}
	if len(fc.t) != 0 && len(fc.t) != len(fc.v) {
		return face{}, fmt.Errorf("a face with texture places at only some corners")
	}
	return fc, nil
}

// index is an OBJ index, from 1, as an index from 0; a relative one, or
// one too large for the meshes MakeHuman makes, is refused.
func index(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("index %q", s)
	}
	return int32(n - 1), nil
}

// check refuses a face whose corners are not in the mesh.
func (m *mesh) check() error {
	for i, fc := range m.faces {
		for _, v := range fc.v {
			if int(v) >= len(m.pos) {
				return fmt.Errorf("human: OBJ face %d: position %d of %d", i, v+1, len(m.pos))
			}
		}
		for _, t := range fc.t {
			if int(t) >= len(m.uv) {
				return fmt.Errorf("human: OBJ face %d: texture place %d of %d", i, t+1, len(m.uv))
			}
		}
	}
	return nil
}

// floats reads len(out) numbers, finite each, from the first fields of f.
func floats(f []string, out []float64) error {
	if len(f) < len(out) {
		return fmt.Errorf("%d numbers, not %d", len(f), len(out))
	}
	for i := range out {
		v, err := strconv.ParseFloat(f[i], 64)
		if err != nil || v != v || v > 1e9 || v < -1e9 {
			return fmt.Errorf("number %q", f[i])
		}
		out[i] = v
	}
	return nil
}
