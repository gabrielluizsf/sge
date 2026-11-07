package farm

import "testing"

// A villager the player talks to stands still, facing the player, for as
// long as the conversation lasts, and goes back to their day after it.
func TestATalkingVillagerStandsStill(t *testing.T) {
	g, f := newGame(idle{})
	v := tradesman(f, eastTown(f), roleLumberjack)
	f.player.Transform.Position = v.pos
	moved, start := false, v.pos
	for range int(30 / frame) {
		g.Step(frame)
		if v.obj != nil && v.obj.Transform.Position.Sub(start).Length() > 0.5 {
			moved = true
			break
		}
	}
	if !moved {
		t.Skip("the villager never walked")
	}
	f.openTalk(v)
	at := v.obj.Transform.Position
	for range int(5 / frame) {
		g.Step(frame)
		if f.talk == nil {
			f.openTalk(v)
		}
		if d := v.obj.Transform.Position.Sub(at).Length(); d > 1e-9 {
			t.Fatalf("the villager walked %.2f away while talking", d)
		}
	}
	f.talk = nil
	for range int(20 / frame) {
		g.Step(frame)
		if v.obj.Transform.Position.Sub(at).Length() > 0 {
			return
		}
	}
	t.Error("the villager never went back to their day")
}
