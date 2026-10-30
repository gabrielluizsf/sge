package render3d

// Level is one model of a chain of detail, drawn while what shows it is
// at least Pixels across the screen; the last level's Pixels is not
// asked, as it is drawn when no other is.
type Level struct {
	Model  *Model
	Pixels float64
}

// LOD is a chain of detail: its Levels, most detailed first, and its
// Hysteresis, the fraction of a level's Pixels past its threshold an
// instance must go to change level, so that one hovering at a threshold
// does not flicker between two. Its levels' models should be posed by the
// same animations, by name.
type LOD struct {
	Levels     []Level
	Hysteresis float64
}

// One is the chain of just m.
func One(m *Model) *LOD { return &LOD{Levels: []Level{{Model: m}}} }

// choose is the level an instance at level cur takes when it is pixels
// across: a more detailed one only when over that level's threshold by
// the hysteresis, a less detailed one only when under its own by as much.
func (l *LOD) choose(cur int, pixels float64) int {
	last := len(l.Levels) - 1
	for i := range last {
		f := 1 - l.Hysteresis
		if i < cur {
			f = 1 + l.Hysteresis
		}
		if pixels >= l.Levels[i].Pixels*f {
			return i
		}
	}
	return last
}
