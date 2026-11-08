package human

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gabrielluizsf/sge/bvh"
)

// The files of MakeHuman's face pose units: the poses, a frame each, and
// their names and licence.
const (
	faceUnitsBVH  = "poseunits/face-poseunits.bvh"
	faceUnitsJSON = "poseunits/face-poseunits.json"
)

// faceUnits are MakeHuman's face pose units, each a frame of the face's
// bones turned, by its name: "JawDrop", "LeftUpperLidClosed" and so on.
type faceUnits map[string][]boneTurn

// boneTurn is a bone, by name, and its turn, in the model's axes.
type boneTurn struct {
	bone string
	turn quat
}

// faceUnits are the kit's face pose units, read once: their poses read as
// untrusted, within bvh's limits, and only if their file says they are
// CC0.
func (k *Kit) faceUnits() (faceUnits, error) {
	if k.faces != nil {
		return k.faces, nil
	}
	data, err := k.lim.read(k.fsys, faceUnitsJSON)
	if err != nil {
		return nil, err
	}
	var about struct {
		License      string   `json:"license"`
		Framemapping []string `json:"framemapping"`
	}
	if err := json.Unmarshal(data, &about); err != nil {
		return nil, fmt.Errorf("human: %s: %w", faceUnitsJSON, err)
	}
	if normalLicence(about.License) != "CC0" {
		return nil, fmt.Errorf("human: the face pose units are %q, not CC0", about.License)
	}
	if data, err = k.lim.read(k.fsys, faceUnitsBVH); err != nil {
		return nil, err
	}
	m, err := bvh.Read(bytes.NewReader(data), bvh.DefaultLimits)
	if err != nil {
		return nil, err
	}
	if len(m.Frames) != len(about.Framemapping) {
		return nil, fmt.Errorf("human: %d face poses for %d names", len(m.Frames), len(about.Framemapping))
	}
	k.faces = unitsOf(m, about.Framemapping)
	return k.faces, nil
}

// unitsOf are the frames of m by names, each the turns of the bones of
// the head, those beyond it, that it turns at all.
func unitsOf(m *bvh.Motion, names []string) faceUnits {
	face := headBones(m)
	out := faceUnits{}
	for f, name := range names {
		var turns []boneTurn
		for _, j := range face {
			q := zUp(m.Rotation(f, j))
			if q[3] < 1-1e-9 {
				turns = append(turns, boneTurn{bone: m.Joints[j].Name, turn: q})
			}
		}
		out[name] = turns
	}
	return out
}

// headBones are the joints of m below its head.
func headBones(m *bvh.Motion) []int {
	head := m.Joint("head")
	var out []int
	for i := range m.Joints {
		for p := m.Joints[i].Parent; p >= 0; p = m.Joints[p].Parent {
			if p == head {
				out = append(out, i)
				break
			}
		}
	}
	return out
}

// zUp is q, a turn in the axes of MakeHuman's BVH files, Z up and the
// face toward -Y, in the model's, Y up and the face toward +Z.
func zUp(q [4]float64) quat {
	if q[3] < 0 {
		q = [4]float64{-q[0], -q[1], -q[2], -q[3]}
	}
	return quat{q[0], q[2], -q[1], q[3]}
}
