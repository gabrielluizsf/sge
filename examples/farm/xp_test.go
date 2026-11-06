package farm

import "testing"

func xpOf(f *farm) int { return f.prog.XP + levelsSum(f.prog.Level) }

// levelsSum is the XP spent on the levels below l.
func levelsSum(l int) int {
	n := 0
	for i := 1; i < l; i++ {
		n += xpNeed(i)
	}
	return n
}

func TestDeedsGiveXPAndRefusedUsesGiveNone(t *testing.T) {
	_, f := newGame(idle{})
	if f.prog.Level != 1 || f.prog.XP != 0 {
		t.Fatalf("starts at %+v", f.prog)
	}
	p := sown(f, 0)
	if xpOf(f) != xpSow {
		t.Errorf("sowing gave %d", xpOf(f))
	}
	water(f, p)
	if xpOf(f) != xpSow+xpWater {
		t.Errorf("watering gave %d", xpOf(f))
	}
	before := xpOf(f)
	p.watered = true
	water(f, p) // not thirsty: refused
	f.slot = slotSeeds
	f.use(f.kit, p.col, p.row) // occupied: nothing
	if xpOf(f) != before {
		t.Errorf("refused uses gave %d XP", xpOf(f)-before)
	}
	p.days, p.stage = crops[0].daysToRipe, ripe
	f.slot = slotAxe
	f.use(f.kit, p.col, p.row)
	if got := xpOf(f) - before; got != xpPick(0) {
		t.Errorf("harvest gave %d, want %d", got, xpPick(0))
	}
	if xpPick(8) <= xpPick(0) {
		t.Errorf("a slower crop gives no more")
	}
}

func TestAVillagersWorkGivesThePlayerNoXP(t *testing.T) {
	_, f := newGame(idle{})
	v := newHand(f, "Ana")
	f.clock.Day = 0
	p := f.plots[1]
	p.clear()
	v.kit.bag.Add(seedsOf(0), 1)
	sowIn(f, v, p)
	if xpOf(f) != 0 {
		t.Errorf("villager's sowing gave the player %d XP", xpOf(f))
	}
}

func TestSellingGivesXPOnlyWhenSomethingSold(t *testing.T) {
	_, f := newGame(idle{})
	at := f.player.Transform.Position
	f.sell(f.kit, at)
	if xpOf(f) != 0 {
		t.Errorf("an empty sale gave %d", xpOf(f))
	}
	f.bag.Add(cropOf(0), 3)
	f.sell(f.kit, at)
	if xpOf(f) < 1 {
		t.Errorf("a sale gave nothing")
	}
}

func TestALevelUpTellsRewardsAndAnnounces(t *testing.T) {
	_, f := newGame(idle{})
	money := f.money
	f.notices = nil
	f.earn(xpNeed(1) + 3)
	if f.prog.Level != 2 || f.prog.XP != 3 {
		t.Fatalf("%+v", f.prog)
	}
	if f.money <= money {
		t.Errorf("no reward")
	}
	if len(f.notices) == 0 || f.notices[len(f.notices)-1].Text != "Level 2!" {
		t.Errorf("notices %v", f.notices)
	}
	if f.glow <= 0 {
		t.Errorf("no glow")
	}
}

func TestTheHeirKeepsTheFamilysLevel(t *testing.T) {
	_, f := newGame(idle{})
	f.earn(xpNeed(1))
	level := f.prog.Level
	f.self = &self{id: "heir", name: "Heir", genes: f.self.genes, life: f.self.life} // who plays changes
	if f.prog.Level != level {
		t.Errorf("level %d, was %d", f.prog.Level, level)
	}
}
