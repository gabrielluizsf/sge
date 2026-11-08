package bvh

import (
	"math"
	"strings"
	"testing"
)

// arm is a root that moves and turns, and a joint that turns, at its end
// a site; two frames, the second turning the joint a quarter about Z.
const arm = `HIERARCHY
ROOT root
{
	OFFSET 0 1 0
	CHANNELS 6 Xposition Yposition Zposition Xrotation Yrotation Zrotation
	JOINT elbow
	{
		OFFSET 0 0 2
		CHANNELS 3 Zrotation Xrotation Yrotation
		End Site
		{
			OFFSET 0 0 1
		}
	}
}
MOTION
Frames: 2
Frame Time: 0.5
0 0 0 0 0 0 0 0 0
1 2 3 0 0 0 90 0 0
`

// A file is read whole: its joints, their parents, offsets and channels,
// and its frames, each with a value for every channel.
func TestReadsAMotion(t *testing.T) {
	m, err := Read(strings.NewReader(arm), DefaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Joints) != 2 || m.Joints[1].Parent != 0 || m.Joints[1].Offset != [3]float64{0, 0, 2} || m.Joint("elbow") != 1 {
		t.Fatalf("joints %+v", m.Joints)
	}
	if m.FrameTime != 0.5 || len(m.Frames) != 2 || len(m.Frames[1]) != 9 || m.Frames[1][2] != 3 {
		t.Fatalf("frames %v at %v", m.Frames, m.FrameTime)
	}
	q := m.Rotation(1, 1)
	if want := math.Sqrt(0.5); math.Abs(q[2]-want) > 1e-12 || math.Abs(q[3]-want) > 1e-12 {
		t.Errorf("a quarter about Z is %v", q)
	}
	if q := m.Rotation(1, 0); q != [4]float64{0, 0, 0, 1} {
		t.Errorf("the root, unturned, is %v", q)
	}
}

// Rotations compose in the order of their channels, the first outermost:
// Z then X is a turn about Z of what X has turned.
func TestRotationsComposeInTheirOrder(t *testing.T) {
	file := strings.Replace(arm, "1 2 3 0 0 0 90 0 0", "1 2 3 0 0 0 90 90 0", 1)
	m, err := Read(strings.NewReader(file), DefaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	q := m.Rotation(1, 1)
	// Y, turned a quarter about X, is Z; that turned a quarter about Z
	// stays Z.
	v := turn(q, [3]float64{0, 1, 0})
	if math.Abs(v[2]-1) > 1e-9 {
		t.Errorf("Y went to %v, not Z", v)
	}
}

// turn is v turned by q.
func turn(q [4]float64, v [3]float64) [3]float64 {
	p := mul(mul(q, [4]float64{v[0], v[1], v[2], 0}), [4]float64{-q[0], -q[1], -q[2], q[3]})
	return [3]float64{p[0], p[1], p[2]}
}

// A file shaped to take all memory or time, or malformed, is refused:
// the shapes of the attacks on other BVH readers, after CVE-2025-6119.
func TestRefusesHostileFiles(t *testing.T) {
	deep := "HIERARCHY\nROOT r\n{\n" + strings.Repeat("JOINT j\n{\n", 100)
	for i := range 100 {
		deep = strings.Replace(deep, "JOINT j\n", "JOINT j"+itoa(i)+"\n", 1)
	}
	many := "HIERARCHY\nROOT r\n{\n"
	for i := range 2000 {
		many += "JOINT j" + itoa(i) + "\n{\n}\n"
	}
	cases := map[string]string{
		"nested past the depth":       deep,
		"more joints than the limit":  many,
		"frames past what fits":       strings.Replace(arm, "Frames: 2", "Frames: 16000", 1),
		"frames past the limit":       strings.Replace(arm, "Frames: 2", "Frames: 99999999", 1),
		"a NaN":                       strings.Replace(arm, "1 2 3", "NaN 2 3", 1),
		"an infinity":                 strings.Replace(arm, "OFFSET 0 0 2", "OFFSET 0 0 +Inf", 1),
		"too few values":              strings.Replace(arm, "1 2 3 0 0 0 90 0 0\n", "1 2 3\n", 1),
		"values after the last frame": arm + "7\n",
		"seven channels":              strings.Replace(arm, "CHANNELS 3 Zrotation", "CHANNELS 7 Zrotation", 1),
		"a channel of no name":        strings.Replace(arm, "Zrotation Xrotation", "Wrotation Xrotation", 1),
		"channels twice":              strings.Replace(arm, "CHANNELS 3", "CHANNELS 0\n\t\tCHANNELS 3", 1),
		"a joint in an End Site":      strings.Replace(arm, "OFFSET 0 0 1\n", "OFFSET 0 0 1\nJOINT x\n{\n}\n", 1),
		"a root inside a joint":       strings.Replace(arm, "JOINT elbow", "ROOT elbow", 1),
		"a name twice":                strings.Replace(arm, "JOINT elbow", "JOINT root", 1),
		"a brace too many":            strings.Replace(arm, "MOTION", "}\nMOTION", 1),
		"MOTION with a joint open":    strings.Replace(arm, "}\nMOTION", "MOTION", 1),
		"no joints":                   "HIERARCHY\nMOTION\nFrames: 0\nFrame Time: 1\n",
		"cut short":                   arm[:len(arm)/2],
		"a negative count":            strings.Replace(arm, "Frames: 2", "Frames: -1", 1),
		"a huge word":                 "HIERARCHY ROOT " + strings.Repeat("a", 1<<17),
	}
	for name, file := range cases {
		if _, err := Read(strings.NewReader(file), DefaultLimits); err == nil {
			t.Errorf("%s: read", name)
		}
	}
	small := DefaultLimits
	small.Bytes = int64(len(arm) - 1)
	if _, err := Read(strings.NewReader(arm), small); err == nil {
		t.Error("a file over its size limit was read")
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for ; i > 0; i /= 10 {
		s = string(rune('0'+i%10)) + s
	}
	return s
}
