package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/vitals"
)

// vitalsOf is the health and food s holds; a save from before there were
// any holds none, and its player is at their best.
func vitalsOf(s *snapshot) vitals.Vitals {
	if s.Vitals == nil {
		return vitals.Full()
	}
	return vitals.Vitals{Health: s.Vitals.Health, Food: s.Vitals.Food}
}

// checkVitals is nil if s's health and food are ones a living player could
// have: from 0 up to the most, health above 0.
func checkVitals(s *snapshot) error {
	v := vitalsOf(s)
	if !v.Valid() || v.Dead() {
		return fmt.Errorf("save: health %v and food %v are out of range", v.Health, v.Food)
	}
	return nil
}
