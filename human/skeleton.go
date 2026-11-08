package human

import (
	"encoding/json"
	"fmt"
	"sort"
)

// skeleton is MakeHuman's skeleton: its bones, every parent before its
// children, and the joints their ends are at, each the middle of some of
// the base mesh's vertices.
type skeleton struct {
	bones  []bone
	joints map[string][]int32
}

// bone is a bone: its name, its parent's index (-1 for a root), and the
// joints at its head and tail.
type bone struct {
	name       string
	parent     int
	head, tail string
}

// maxBones is the most bones a skeleton may have: a primitive's joints
// are 16-bit.
const maxBones = 1024

// parseSkeleton reads a .mhskel of a mesh of n vertices: every joint's
// vertices must be in the mesh, every bone's joints and parent must be
// there, and the bones must make a tree.
func parseSkeleton(data []byte, n int) (*skeleton, error) {
	var doc struct {
		Bones map[string]struct {
			Head, Tail, Parent string
		}
		Joints map[string][]int32
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("human: skeleton: %w", err)
	}
	if len(doc.Bones) == 0 || len(doc.Bones) > maxBones {
		return nil, fmt.Errorf("human: skeleton of %d bones", len(doc.Bones))
	}
	for name, vs := range doc.Joints {
		if len(vs) == 0 {
			return nil, fmt.Errorf("human: joint %s has no vertices", name)
		}
		for _, v := range vs {
			if v < 0 || int(v) >= n {
				return nil, fmt.Errorf("human: joint %s: vertex %d of %d", name, v, n)
			}
		}
	}
	s := &skeleton{joints: doc.Joints}
	at := map[string]int{}
	names := make([]string, 0, len(doc.Bones))
	for name := range doc.Bones {
		names = append(names, name)
	}
	sort.Strings(names)
	// Bones are placed once their parent is: a pass places at least one,
	// or the rest wait on a parent missing or on each other.
	for len(at) < len(names) {
		placed := len(at)
		for _, name := range names {
			b := doc.Bones[name]
			if _, done := at[name]; done {
				continue
			}
			parent, ok := -1, b.Parent == ""
			if !ok {
				parent, ok = at[b.Parent]
			}
			if !ok {
				continue
			}
			if doc.Joints[b.Head] == nil || doc.Joints[b.Tail] == nil {
				return nil, fmt.Errorf("human: bone %s's joints are not in the skeleton", name)
			}
			at[name] = len(s.bones)
			s.bones = append(s.bones, bone{name: name, parent: parent, head: b.Head, tail: b.Tail})
		}
		if len(at) == placed {
			return nil, fmt.Errorf("human: skeleton's bones do not make a tree")
		}
	}
	return s, nil
}

// joint is where joint name is among pos.
func (s *skeleton) joint(name string, pos [][3]float64) [3]float64 {
	var c [3]float64
	vs := s.joints[name]
	for _, v := range vs {
		c = add(c, pos[v])
	}
	return scale(c, 1/float64(len(vs)))
}

// index is the bone named name's index, -1 if there is none.
func (s *skeleton) index(name string) int {
	for i, b := range s.bones {
		if b.name == name {
			return i
		}
	}
	return -1
}
