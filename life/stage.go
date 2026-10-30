package life

// Stage is a time of life.
type Stage int

const (
	Baby Stage = iota
	Child
	Teen
	Adult
	Elder
)

// Ages at which each stage begins: a child from 2, a teen from 12, an
// adult from 18, an elder from 60.
const (
	ChildFrom = 2
	TeenFrom  = 12
	AdultFrom = 18
	ElderFrom = 60
)

// StageOf is the stage of life at years.
func StageOf(years float64) Stage {
	switch {
	case years < ChildFrom:
		return Baby
	case years < TeenFrom:
		return Child
	case years < AdultFrom:
		return Teen
	case years < ElderFrom:
		return Adult
	}
	return Elder
}

// Grown is whether someone of years is an adult, or older.
func Grown(years float64) bool { return years >= AdultFrom }

// String is the stage's name: "baby", "child", "teen", "adult", "elder".
func (s Stage) String() string {
	if s < Baby || s > Elder {
		return "stage?"
	}
	return [...]string{"baby", "child", "teen", "adult", "elder"}[s]
}
