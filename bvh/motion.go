package bvh

// Channel is a value a joint moves by: along or about an axis.
type Channel uint8

// The channels, as a file names them.
const (
	Xposition Channel = iota
	Yposition
	Zposition
	Xrotation
	Yrotation
	Zrotation
)

// channelNames are the channels by their names in a file.
var channelNames = map[string]Channel{
	"Xposition": Xposition, "Yposition": Yposition, "Zposition": Zposition,
	"Xrotation": Xrotation, "Yrotation": Yrotation, "Zrotation": Zrotation,
}

// Joint is a joint of the skeleton: its name, its parent's index, -1 for
// a root, its offset from its parent, and its channels, in the order the
// frames give their values.
type Joint struct {
	Name     string
	Parent   int
	Offset   [3]float64
	Channels []Channel
}

// Motion is a skeleton and its frames: each frame every joint's channels'
// values, the joints in order, rotations in degrees; FrameTime is the
// seconds between frames.
type Motion struct {
	Joints    []Joint
	FrameTime float64
	Frames    [][]float64
}

// Joint is the index of the joint called name, -1 when there is none.
func (m *Motion) Joint(name string) int {
	for i, j := range m.Joints {
		if j.Name == name {
			return i
		}
	}
	return -1
}

// values is the joint i's channels' values in frame f.
func (m *Motion) values(f, i int) []float64 {
	at := 0
	for _, j := range m.Joints[:i] {
		at += len(j.Channels)
	}
	return m.Frames[f][at : at+len(m.Joints[i].Channels)]
}
