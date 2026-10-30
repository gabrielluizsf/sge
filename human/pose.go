package human

// Pose is how a person holds themself at a moment, angles in radians,
// each side's Left then Right: Pitch tips the whole body forward, as a
// swimmer lies; Lean bends the back forward and Twist turns it, left
// positive; Nod bows the head. A Hip swings a thigh forward, a Knee bends
// back, an Ankle points the foot down. A Shoulder swings an arm forward,
// a Raise lifts it out to the side, 0 hanging, and an Elbow bends. Bob
// lifts the body, in metres.
type Pose struct {
	Pitch, Lean, Twist, Nod float64
	Hip, Knee, Ankle        [2]float64
	Shoulder, Raise, Elbow  [2]float64
	Bob                     float64
}

// side is a side's suffix in MakeHuman's bone names, and which way out
// from the body is along X.
var sides = [2]struct {
	suffix string
	out    float64
}{{".L", 1}, {".R", -1}}

var (
	xAxis = [3]float64{1, 0, 0}
	yAxis = [3]float64{0, 1, 0}
	zAxis = [3]float64{0, 0, 1}
)

// arms is how a body's arms lie at rest, each side's: the turn that lets
// one down to hang at the side, the turn that straightens its forearm,
// which MakeHuman's rest bends forward, and the axis its elbow bends
// forward about, in its own rest frame.
type arms struct {
	hang, straighten [2]quat
	elbow            [2][3]float64
}

// turns are the rotations of a pose's bones, by name, on a body whose
// arms lie at rest as a says: an arm is first let down to the side, then
// lifted and swung.
func (p Pose) turns(a arms) map[string]quat {
	t := map[string]quat{
		"root":    around(xAxis, p.Pitch),
		"spine04": around(xAxis, p.Lean/2),
		"spine02": around(xAxis, p.Lean/2).then(around(yAxis, p.Twist)),
		"neck01":  around(xAxis, p.Nod),
	}
	for s, side := range sides {
		t["upperleg01"+side.suffix] = around(xAxis, -p.Hip[s])
		t["lowerleg01"+side.suffix] = around(xAxis, p.Knee[s])
		t["foot"+side.suffix] = around(xAxis, p.Ankle[s])
		t["upperarm01"+side.suffix] = a.hang[s].then(around(zAxis, side.out*p.Raise[s])).then(around(xAxis, -p.Shoulder[s]))
		t["lowerarm01"+side.suffix] = a.straighten[s].then(around(a.elbow[s], p.Elbow[s]))
	}
	return t
}
