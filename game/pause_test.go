package game

import "testing"

// While the game is paused, the scripts of objects stop, but for those
// that are Always; they go on when it is not.
func TestPausedScriptsStopButAlwaysOnesRun(t *testing.T) {
	var log []string
	g := New()
	g.Add(NewObject("world")).AddScript(recorder{&log})
	menu := g.Add(NewObject("menu"))
	menu.Always = true
	menu.AddScript(recorder{&log})
	g.Step(0.1)
	log = nil
	g.Paused = true
	g.Step(0.1)
	if len(log) != 1 || log[0] != "update menu" {
		t.Fatalf("paused log = %v", log)
	}
	log = nil
	g.Paused = false
	g.Step(0.1)
	if len(log) != 2 {
		t.Errorf("log = %v after resuming", log)
	}
}
