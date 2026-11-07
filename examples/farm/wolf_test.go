package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
)

// wolfGame is a game at night with wolves on, stepping by seconds.
func wolfGame(t *testing.T, ctl Controller) (*farm, func(seconds float64)) {
	t.Helper()
	g, f := build(SVGArt(), ctl)
	f.wolves.on = true
	f.clock.Hour = duskAt + 1
	return f, func(seconds float64) {
		for range int(seconds / frame) {
			g.Step(frame)
		}
	}
}

// Wolves come after dusk and never by day, and go with the night.
func TestWolvesComeAtDuskAndGoAtDawn(t *testing.T) {
	f, step := wolfGame(t, still{})
	f.clock.Hour = 10
	step(5)
	if len(f.wolves.wolves) != 0 {
		t.Fatalf("%d wolves by day", len(f.wolves.wolves))
	}
	f.clock.Hour = duskAt + 0.01
	step(6)
	if len(f.wolves.wolves) == 0 {
		t.Fatal("no wolf at night")
	}
	f.clock.Hour = dawnAt + 1
	step(0.1)
	if len(f.wolves.wolves) != 0 {
		t.Error("wolves left at dawn")
	}
}

// Only to one playing: the demo's bot is left alone unless told.
func TestWolvesLeaveTheBotAlone(t *testing.T) {
	_, f := build(SVGArt(), Bot(2))
	if f.wolves.on {
		t.Error("wolves on for the bot")
	}
	_, kb := build(SVGArt(), Keyboard())
	if !kb.wolves.on {
		t.Error("no wolves for the one playing")
	}
}

// A wolf runs to the player and bites them, and it hurts.
func TestAWolfBitesThePlayer(t *testing.T) {
	f, step := wolfGame(t, still{})
	f.wolves.on = false
	at := feet(f.player)
	f.wolves.add(f.player.Game(), at.Add(vector.Vector3{X: 6 * tile}))
	before := f.vit.Health
	step(8)
	if f.vit.Health >= before {
		t.Errorf("health still %v", f.vit.Health)
	}
	w := f.wolves.wolves[0]
	if feet(w.obj).Sub(feet(f.player)).Length() > wolfReach+tile {
		t.Errorf("the wolf is %v away", feet(w.obj).Sub(feet(f.player)).Length())
	}
}

// The player kills a wolf with blows, which lies dead a while, has a look of
// its own, bites no more, and gives XP.
func TestThePlayerKillsAWolf(t *testing.T) {
	f, step := wolfGame(t, still{})
	f.wolves.on = false
	at := feet(f.player)
	w := f.wolves.add(f.player.Game(), at.Add(vector.Vector3{X: tile}))
	f.hands.facing = vector.Vector3{X: 1}
	xp := f.xpValue()
	f.axe = 0
	f.kit.slot = slotAxe
	for range 30 {
		if w.fallen {
			break
		}
		f.hands.swing(f.player)
	}
	if !w.fallen {
		t.Fatal("the wolf took every blow")
	}
	if f.xpValue() <= xp {
		t.Error("no XP for the wolf")
	}
	if f.wolves.ahead(at, f.hands.facing) != nil {
		t.Error("a dead wolf still struck")
	}
	health := f.vit.Health
	step(wolfDead + 1)
	if len(f.wolves.wolves) != 0 || f.vit.Health < health {
		t.Errorf("%d wolves, health %v from %v", len(f.wolves.wolves), f.vit.Health, health)
	}
}

// A grown villager in sight takes a wolf on with an axe, and kills it,
// without turning on the player; a bite hurts them.
func TestAVillagerKillsAWolf(t *testing.T) {
	f, tn, step := seenTown(t, SVGArt())
	f.wolves.on = false
	f.clock.Hour = duskAt + 0.2
	v := tradesman(f, tn, roleLumberjack)
	if v.obj == nil {
		step(10)
	}
	v.obj.Hidden, v.inside = false, false
	f.player.Transform.Position = feet(v.obj)
	f.player.Transform.Position.Y -= 6 * tile
	w := f.wolves.add(f.player.Game(), feet(v.obj).Add(vector.Vector3{X: 3 * tile}))
	for range 40 * 60 {
		step(1)
		if w.fallen {
			break
		}
		v.wounds = 0
	}
	if !w.fallen {
		t.Fatalf("the wolf at %v, the villager at %v", feet(w.obj), feet(v.obj))
	}
	if v.angry > 0 {
		t.Error("the villager is angry at the player")
	}
}

// A wolf's bite can kill a villager, with no anger at the player.
func TestAWolfCanKillAVillager(t *testing.T) {
	f, tn, step := seenTown(t, SVGArt())
	v := tradesman(f, tn, roleLumberjack)
	if v.obj == nil {
		step(10)
	}
	f.bitten(v, npcHealth, feet(v.obj).Add(vector.Vector3{X: tile}))
	if f.personByID(v.id) != nil {
		t.Error("alive")
	}
}

// A wolf finds its way round what is solid.
func TestAWolfGoesRoundAFence(t *testing.T) {
	f, step := wolfGame(t, still{})
	f.wolves.on = false
	r := penRect
	f.player.Transform.Position = f.center(pathfind.Cell{Col: r.col + r.cols/2, Row: r.row + 2})
	w := f.wolves.add(f.player.Game(), f.center(pathfind.Cell{Col: r.col + r.cols/2, Row: r.row - 2}))
	step(14)
	if feet(w.obj).Sub(feet(f.player)).Length() > 2*tile {
		t.Errorf("the wolf is stuck at %v", feet(w.obj))
	}
}

// A villager's swing that misses a wolf does not land on the player beside
// them, whom they have no quarrel with.
func TestASwingAtAWolfDoesNotHurtThePlayer(t *testing.T) {
	f, tn, step := seenTown(t, SVGArt())
	v := tradesman(f, tn, roleLumberjack)
	if v.obj == nil {
		step(10)
	}
	f.player.Transform.Position = feet(v.obj).Add(vector.Vector3{X: tile})
	before := f.vit.Health
	v.hands.swing(v.obj)
	if f.vit.Health != before {
		t.Errorf("health %v from %v", f.vit.Health, before)
	}
	v.angry = npcAnger
	v.hands.acting = 0
	v.hands.swing(v.obj)
	if f.vit.Health >= before {
		t.Error("an angry villager's swing missed")
	}
}

// press is a player pressing what intent says, for one frame.
type press struct{ in *intent }

func (p press) control(*farm, *game.Object, float64) intent {
	in := *p.in
	*p.in = intent{}
	return in
}

// A wolf at the player's back, or pushed off by its own bite, is still hit
// by a blow, which turns them to it; and with the axe in hand Space, or a
// click, swings at it as F does, and does not go off at a tree.
func TestThePlayerHitsAWolfBesideThem(t *testing.T) {
	for name, in := range map[string]intent{"F": {strike: true}, "Space or a click": {use: true}} {
		t.Run(name, func(t *testing.T) {
			var next intent
			f, step := wolfGame(t, press{&next})
			f.wolves.on = false
			f.axe = 0
			f.kit.slot = slotAxe
			f.hands.facing = vector.Vector3{X: 1}
			w := f.wolves.add(f.player.Game(), feet(f.player).Add(vector.Vector3{X: -1.6 * tile}))
			w.next = 100 // it does not bite, for this
			next = in
			step(3 * frame)
			if w.wounds == 0 {
				t.Fatalf("the blow missed the wolf %v away", feet(w.obj).Sub(feet(f.player)).Length()/tile)
			}
			if f.hands.facing.X >= 0 {
				t.Errorf("not turned to the wolf: %v", f.hands.facing)
			}
		})
	}
}

// With no wolf near, Space with the axe in hand is what it was: the blow is
// not swung at the air, and a wolf far off is not struck.
func TestSpaceWithNoWolfNearSwingsNothing(t *testing.T) {
	var next intent
	f, step := wolfGame(t, press{&next})
	f.wolves.on = false
	f.axe = 0
	f.kit.slot = slotAxe
	w := f.wolves.add(f.player.Game(), feet(f.player).Add(vector.Vector3{X: 9 * tile}))
	w.next = 100
	next = intent{use: true}
	step(3 * frame)
	if w.wounds != 0 {
		t.Error("struck a wolf far off")
	}
}
