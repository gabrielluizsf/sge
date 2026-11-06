package progress

import "testing"

func TestTheCurveRisesAndIsFixed(t *testing.T) {
	want := map[int]int{1: 100, 2: 283, 3: 520, 4: 800, 5: 1118, 10: 3162}
	for l, n := range want {
		if got := Need(l); got != n {
			t.Errorf("Need(%d) = %d, want %d", l, got, n)
		}
	}
	for l := 1; l < MaxLevel; l++ {
		if Need(l+1) <= Need(l) {
			t.Fatalf("level %d needs no more than %d", l+1, l)
		}
	}
}

func TestAddReportsLevelsGainedAndCarriesTheRest(t *testing.T) {
	p := New()
	if p.Level != 1 || p.XP != 0 {
		t.Fatalf("new: %+v", p)
	}
	if n := p.Add(99); n != 0 || p.XP != 99 {
		t.Errorf("99: %d %+v", n, p)
	}
	if n := p.Add(1); n != 1 || p.Level != 2 || p.XP != 0 {
		t.Errorf("100: %d %+v", n, p)
	}
	if n := p.Add(283 + 520 + 5); n != 2 || p.Level != 4 || p.XP != 5 {
		t.Errorf("two at once: %d %+v", n, p)
	}
	if p.Add(-5) != 0 || p.XP != 5 {
		t.Errorf("negative changed XP: %+v", p)
	}
}

func TestTheLastLevelTakesNoMore(t *testing.T) {
	p := Progress{Level: MaxLevel}
	if n := p.Add(1 << 30); n != 0 || p.Level != MaxLevel || p.XP != 0 {
		t.Errorf("%d %+v", n, p)
	}
}

func TestValid(t *testing.T) {
	bad := []Progress{{0, 0}, {1, -1}, {1, 100}, {MaxLevel + 1, 0}, {-3, 5}, {MaxLevel, 1}}
	for _, p := range bad {
		if p.Valid() == nil {
			t.Errorf("%+v accepted", p)
		}
	}
	for _, p := range []Progress{{1, 0}, {1, 99}, {7, Need(7) - 1}, {MaxLevel, 0}} {
		if err := p.Valid(); err != nil {
			t.Errorf("%+v: %v", p, err)
		}
	}
}
