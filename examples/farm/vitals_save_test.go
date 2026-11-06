package farm

import (
	"encoding/json"
	"testing"
)

// Health and food come back as they were saved.
func TestVitalsRoundTrip(t *testing.T) {
	_, f := newGame(idle{})
	f.vit.Health, f.vit.Food = 37, 12
	b, _ := json.Marshal(takeSnapshot(f))
	var s snapshot
	json.Unmarshal(b, &s)
	_, g := newGame(idle{})
	if err := restore(g, &s); err != nil || g.vit != f.vit {
		t.Errorf("%v %+v", err, g.vit)
	}
}

// A save from before health and food has the player at their best.
func TestAnOldSaveWithoutVitalsIsFull(t *testing.T) {
	var old map[string]any
	json.Unmarshal(playedJSON(t), &old)
	delete(old, "vitals")
	b, _ := json.Marshal(old)
	var s snapshot
	json.Unmarshal(b, &s)
	_, g := newGame(idle{})
	g.vit.Health, g.vit.Food = 3, 3
	if err := restore(g, &s); err != nil || g.vit.Health != 100 || g.vit.Food != 100 {
		t.Errorf("%v %+v", err, g.vit)
	}
}

func TestHostileVitalsAreRefused(t *testing.T) {
	for name, edit := range map[string]func(*snapshot){
		"negative health": func(s *snapshot) { s.Vitals = &vitalsState{Health: -5, Food: 50} },
		"negative food":   func(s *snapshot) { s.Vitals = &vitalsState{Health: 50, Food: -1} },
		"huge health":     func(s *snapshot) { s.Vitals = &vitalsState{Health: 1e9, Food: 50} },
		"huge food":       func(s *snapshot) { s.Vitals = &vitalsState{Health: 50, Food: 101} },
		"dead player":     func(s *snapshot) { s.Vitals = &vitalsState{Health: 0, Food: 50} },
	} {
		hostile(t, name, edit)
	}
}
