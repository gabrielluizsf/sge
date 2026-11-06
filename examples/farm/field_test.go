package farm

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/random"
)

// cornCrop is the index of corn, whose note the tests read.
const cornCrop = 5

// newHand is a grown farmer of no town, in f's people, with seeds.
func newHand(f *farm, name string) *villager {
	v := &villager{name: name, role: roleFarmer, kit: &kit{axe: noAxe, exact: true}, luck: random.New(7)}
	v.who.Stage = life.Adult
	v.mind = &mind{v: v}
	v.kit.bag.Add(seedsOf(cornCrop), 10)
	v.kit.seed = cornCrop
	f.people = append(f.people, v)
	return v
}

// targetOf is the cell job j works, the one before where it stands.
func targetOf(j *job) (col, row int) {
	for _, b := range banks {
		if b.face == j.face {
			return j.to.Col + b.dc, j.to.Row + b.dr
		}
	}
	return j.to.Col, j.to.Row
}

// sowIn has v sow the plot p, as the field job would.
func sowIn(f *farm, v *villager, p *plot) {
	v.kit.slot = slotSeeds
	f.useAs(v.kit, p.col, p.row)
}

// A villager sows free soil only, never the player's plot, and keeps
// what they sowed.
func TestAVillagerSowsFreeSoilOnly(t *testing.T) {
	_, f := newGame(idle{})
	mine := f.plots[0]
	mine.crop = 0
	v := newHand(f, "Ana")
	f.clock.Day = 0
	for range 40 {
		j := v.mind.fieldSow(f)
		if j == nil {
			t.Fatal("no job on a farm of free soil")
		}
		col, row := targetOf(j)
		p := f.plotAt(col, row)
		if p == nil || p == mine || p.crop >= 0 || p.owner != nil {
			t.Fatalf("aims at %d, %d: %+v", col, row, p)
		}
	}
	v.kit.slot = slotSeeds
	if f.useAs(v.kit, mine.col, mine.row) != actNone || mine.owner != nil {
		t.Error("a villager tended the player's plot")
	}
	p := f.plots[1]
	sowIn(f, v, p)
	if p.owner != v || p.crop != cornCrop {
		t.Errorf("not kept by the sower: %+v", p)
	}
}

// The player cannot pick, water or hoe a villager's plot, and is told
// whose it is.
func TestThePlayerCannotTendAVillagersPlot(t *testing.T) {
	_, f := newGame(idle{})
	v := newHand(f, "Ana")
	p := f.plots[3]
	sowIn(f, v, p)
	f.bag.Add(seedsOf(0), 1)
	for _, c := range []struct {
		slot  int
		setup func()
	}{
		{slotCan, func() { p.watered = false }},
		{slotHoe, func() { p.dead = true }},
		{slotAxe, func() { p.dead, p.stage = false, ripe }},
		{slotSeeds, func() {}},
	} {
		p.dead, p.stage, p.watered = false, 0, false
		c.setup()
		before := *p
		f.notices = nil
		f.slot = c.slot
		if f.useAs(f.kit, p.col, p.row) != actNone || *p != before {
			t.Errorf("slot %d changed an owned plot: %+v", c.slot, p)
		}
		if len(f.notices) == 0 || f.notices[0].Text != "This plot belongs to Ana" {
			t.Errorf("slot %d: no notice", c.slot)
		}
	}
	if f.bag.Count(crops[cornCrop].name) != 0 {
		t.Error("the player took the crop")
	}
}

// The note over a villager's crop names its owner.
func TestTheNoteNamesTheOwner(t *testing.T) {
	_, f := newGame(idle{})
	p := f.plots[2]
	p.crop = cornCrop
	if strings.HasPrefix(cropNote(p), "Corn of") {
		t.Errorf("the player's crop: %q", cropNote(p))
	}
	p.owner = newHand(f, "Ana")
	if n := cropNote(p); !strings.HasPrefix(n, "Corn of Ana - sprouts in ") {
		t.Errorf("note %q", n)
	}
}

// A villager picks their ripe crop into their kit, and the plot is free
// again; they also water it when it is dry.
func TestAVillagerHarvestsAndFreesThePlot(t *testing.T) {
	_, f := newGame(idle{})
	v := newHand(f, "Ana")
	v.kit.fill()
	p := f.plots[4]
	sowIn(f, v, p)
	j := v.mind.fieldTend(f)
	if j == nil || j.slot != slotCan {
		t.Fatalf("a dry crop is to be watered: %+v", j)
	}
	v.kit.slot = slotCan
	if f.useAs(v.kit, p.col, p.row) != actWater || !p.watered {
		t.Fatal("not watered")
	}
	p.stage = ripe
	if j = v.mind.fieldTend(f); j == nil {
		t.Fatal("a ripe crop is to be picked")
	}
	if col, row := targetOf(j); col != p.col || row != p.row {
		t.Errorf("aims at %d, %d", col, row)
	}
	v.kit.slot = j.slot
	if f.useAs(v.kit, p.col, p.row) != actPick || v.kit.bag.Count(cropOf(cornCrop)) != 1 {
		t.Fatal("not picked into the kit")
	}
	if p.crop >= 0 || p.owner != nil {
		t.Errorf("not free: %+v", p)
	}
}

// Nobody takes more than fieldEach plots, and together they leave the
// player most of the soil.
func TestTheFieldCapsHold(t *testing.T) {
	_, f := newGame(idle{})
	a := newHand(f, "Ana")
	for _, p := range f.plots[:fieldEach] {
		sowIn(f, a, p)
	}
	if a.mind.fieldSow(f) != nil || f.roomInField(a) {
		t.Error("one villager took more than the cap")
	}
	farm := f.farmPlots()
	cap := int(fieldShare * float64(len(farm)))
	f.release(a)
	for i := 0; i < len(farm); i++ {
		b := newHand(f, "B")
		if f.roomInField(b) {
			sowIn(f, b, farm[i])
		}
	}
	if f.held() > cap {
		t.Errorf("villagers hold %d of %d plots, the cap is %d", f.held(), len(farm), cap)
	}
}

// Only a pair of each trade is on duty a day, and they take turns.
func TestOnlyAPairIsOnDutyADay(t *testing.T) {
	_, f := newGame(idle{})
	var all []*villager
	for range 6 {
		all = append(all, newHand(f, "x"))
	}
	seen := map[*villager]bool{}
	for d := 0; d < 3; d++ {
		f.clock.Day = d
		n := 0
		for _, v := range all {
			if f.onDuty(v) {
				n++
				seen[v] = true
			}
		}
		if n != fieldHands {
			t.Errorf("day %d: %d on duty", d, n)
		}
	}
	if len(seen) != 6 {
		t.Errorf("%d took a turn in 3 days", len(seen))
	}
}

// A lumberjack on duty fells trees around the farm, never inside it.
func TestALumberjackFellsOutsideTheFence(t *testing.T) {
	_, f := newGame(idle{})
	v := newHand(f, "Beto")
	v.role = roleLumberjack
	v.kit.axe = 0
	for range 30 {
		j := v.mind.chopIn(f, area{center: pathfind.Cell{Col: 15, Row: 10}, cols: 27, rows: 19}, func(c, r int) bool { return !f.onFarm(c, r) })
		if j == nil {
			t.Skip("no tree around the farm")
		}
		col, row := targetOf(j)
		if f.onFarm(col, row) || f.onFarm(j.to.Col, j.to.Row) || !isTree(f.thingAt(col, row)) {
			t.Fatalf("aims at %d, %d from %v", col, row, j.to)
		}
	}
}
