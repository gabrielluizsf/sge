package life

import (
	"testing"

	"github.com/gabrielluizsf/sge/random"
)

// A year of 112 days: someone born on day -2240 is 20 on day 0, 21 a
// year on, and gone the day their span ends, not before.
func TestAgeGrowsWithTheYears(t *testing.T) {
	c := Calendar{DaysPerYear: 112}
	l := Life{Born: -2240, Span: 70}
	if a := c.Age(l, 0); a != 20 {
		t.Errorf("age %v", a)
	}
	if a := c.Age(l, 112); a != 21 {
		t.Errorf("a year on, age %v", a)
	}
	end := l.Born + 70*112
	if c.Gone(l, end-1) || !c.Gone(l, end) {
		t.Errorf("gone the day before %v, the day %v", c.Gone(l, end-1), c.Gone(l, end))
	}
}

func TestStages(t *testing.T) {
	for years, want := range map[float64]Stage{0: Baby, 1.9: Baby, 2: Child, 11: Child, 12: Teen, 17.9: Teen, 18: Adult, 59: Adult, 60: Elder, 95: Elder} {
		if got := StageOf(years); got != want {
			t.Errorf("%v years: %v, not %v", years, got, want)
		}
	}
	if Grown(17.99) || !Grown(18) {
		t.Error("grown before 18 or not at it")
	}
}

// A drawn life is as old as asked on the day, lives within the spans, and
// always outlives the age it starts at, even an old one's.
func TestDrawnLives(t *testing.T) {
	c := Calendar{DaysPerYear: 112}
	luck := random.New(7)
	for i := range 200 {
		years := float64(i % 99)
		l := c.Draw(luck, years, 50)
		if a := c.Age(l, 50); a < years-0.01 || a > years+0.01 {
			t.Errorf("drawn %v years old, is %v", years, a)
		}
		if l.Span < SpanFrom || l.Span <= years || c.Gone(l, 50) {
			t.Errorf("%v years old, span %v", years, l.Span)
		}
	}
	if b := c.Newborn(luck, 9); c.Age(b, 9) != 0 || b.Born != 9 {
		t.Errorf("newborn %v", b)
	}
}

func TestStageNames(t *testing.T) {
	if Baby.String() != "baby" || Elder.String() != "elder" || Stage(9).String() != "stage?" {
		t.Error(Baby, Elder, Stage(9))
	}
}
