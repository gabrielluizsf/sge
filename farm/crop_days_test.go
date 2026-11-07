package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/tilemap"
)

// sown is a plot of f with crop c in it, sown as the player does.
func sown(f *farm, c int) *plot {
	p := f.plots[0]
	p.clear()
	f.slot, f.seed = slotSeeds, c
	f.bag.Add(seedsOf(c), 1)
	f.use(f.kit, p.col, p.row)
	return p
}

// water waters p with the can, as the player does.
func water(f *farm, p *plot) bool {
	f.slot = slotCan
	return f.use(f.kit, p.col, p.row)
}

// A crop shows a stage a day at most, from the seed to ripe, ripe on the
// day it takes, whatever it takes.
func TestDaysMapToStages(t *testing.T) {
	for _, c := range crops {
		last := 0
		for d := 0; d <= c.daysToRipe+2; d++ {
			s := stageOf(c, d)
			if s < last || s > ripe || (s == ripe) != (d >= c.daysToRipe) {
				t.Errorf("%s: day %d is stage %d after %d", c.name, d, s, last)
			}
			last = s
		}
	}
}

// Crops take days: watered, one grows a day each morning, never a
// second on the same one, and not at all in the seconds of the day.
func TestACropGrowsADayEachMorning(t *testing.T) {
	g, f := newGame(idle{})
	p := sown(f, 2) // potato, 6 days
	water(f, p)
	for range 600 {
		g.Step(frame)
	}
	if p.days != 0 || p.stage != 0 {
		t.Fatalf("grew in seconds: %+v", p)
	}
	f.passDays(1)
	if p.days != 1 || p.stage != 1 || p.dead {
		t.Fatalf("a morning on: %+v", p)
	}
	f.days()
	if p.days != 1 {
		t.Fatalf("a day counted twice: %+v", p)
	}
}

// A growing crop left dry for a whole day is dead the next morning,
// shown withered, picked by no one; the hoe clears it.
func TestADryDayKillsAGrowingCrop(t *testing.T) {
	_, f := newGame(idle{})
	p := sown(f, 0)
	f.passDays(1)
	if !p.dead || f.crops.At(p.col, p.row) != tileID[witheredTile] || p.thirsty() {
		t.Fatalf("not withered: %+v", p)
	}
	f.slot = slotAxe
	if f.use(f.kit, p.col, p.row) || p.crop < 0 {
		t.Fatal("the axe cleared it")
	}
	f.slot = slotSeeds
	f.bag.Add(seedsOf(1), 1)
	if f.use(f.kit, p.col, p.row) || p.crop < 0 {
		t.Fatal("sowed over a dead plant")
	}
	f.slot = slotCan
	if f.use(f.kit, p.col, p.row) {
		t.Fatal("watered a dead plant")
	}
	f.slot = slotHoe
	if !f.use(f.kit, p.col, p.row) || p.crop != -1 || f.crops.At(p.col, p.row) != tilemap.Empty || f.bag.Count(cropOf(0)) != 0 {
		t.Fatalf("not cleared: %+v", p)
	}
}

// A crop watered every day lives; one missed a single day dies.
func TestOneMissedDayIsEnough(t *testing.T) {
	_, f := newGame(idle{})
	p := sown(f, 0)
	for range 2 {
		water(f, p)
		f.passDays(1)
	}
	if p.dead || p.days != 2 {
		t.Fatalf("died watered: %+v", p)
	}
	f.passDays(1)
	if !p.dead {
		t.Fatalf("lived a dry day: %+v", p)
	}
}

// A ripe crop waiting to be picked needs no water, however long it waits.
func TestARipeCropNeedsNoWater(t *testing.T) {
	_, f := newGame(idle{})
	p := sown(f, 0)
	p.days, p.stage = crops[0].daysToRipe, ripe
	f.passDays(5)
	if p.dead || p.stage != ripe || p.thirsty() {
		t.Fatalf("a ripe crop withered: %+v", p)
	}
	f.slot = slotCan
	f.use(f.kit, p.col, p.row)
	if f.inCan() != canCapacity {
		t.Fatal("watered a ripe crop")
	}
}

// A crop that bears again stays planted when picked, and is ripe again
// after its regrow days of water.
func TestACropBearsAgain(t *testing.T) {
	_, f := newGame(idle{})
	c := 3 // tomato: 7 days, again every 3
	p := sown(f, c)
	for range crops[c].daysToRipe {
		water(f, p)
		f.passDays(1)
	}
	f.slot = slotAxe
	if !f.use(f.kit, p.col, p.row) || f.bag.Count(cropOf(c)) != 1 || p.crop != c || p.stage == ripe || p.dead {
		t.Fatalf("not left planted: %+v", p)
	}
	for d := range crops[c].regrowDays {
		if p.stage == ripe {
			t.Fatalf("ripe after %d days: %+v", d, p)
		}
		water(f, p)
		f.passDays(1)
	}
	if p.stage != ripe || !f.use(f.kit, p.col, p.row) || f.bag.Count(cropOf(c)) != 2 {
		t.Fatalf("did not bear again: %+v", p)
	}
}

// A crop that does not bear again is gone once picked.
func TestACropThatDoesNotRegrowIsGone(t *testing.T) {
	_, f := newGame(idle{})
	p := sown(f, 0)
	p.days, p.stage = crops[0].daysToRipe, ripe
	f.use(f.kit, p.col, p.row)
	if p.crop != -1 {
		t.Fatalf("still planted: %+v", p)
	}
}

// The can starts full, each watering takes a drop, and only water fills it.
func TestTheCanEmptiesAndRefillsAtWater(t *testing.T) {
	_, f := newGame(idle{})
	p := sown(f, 0)
	if f.inCan() != canCapacity {
		t.Fatalf("starts with %d", f.inCan())
	}
	for i := 1; i <= canCapacity; i++ {
		p.watered = false
		if !water(f, p) || f.inCan() != canCapacity-i {
			t.Fatalf("watering %d left %d", i, f.inCan())
		}
	}
	p.watered = false
	f.notices = nil
	if water(f, p) || p.watered || len(f.notices) == 0 || f.notices[0].Text != "The watering can is empty — fill it at the water" {
		t.Fatalf("watered with an empty can: notices %v", f.notices)
	}
	dc, dr := dryCell(f)
	f.slot = slotCan
	f.use(f.kit, dc, dr)
	if f.inCan() != 0 {
		t.Fatal("filled on dry land")
	}
	wc, wr := waterCell(f)
	if !f.use(f.kit, wc, wr) || f.inCan() != canCapacity {
		t.Fatalf("not filled at the water: %d", f.inCan())
	}
	if f.use(f.kit, wc, wr) {
		t.Fatal("filled a full can")
	}
	f.slot = slotSeeds
	f.kit.used = 5
	f.use(f.kit, wc, wr)
	if f.inCan() != canCapacity-5 {
		t.Fatal("filled without the can in hand")
	}
}

// waterCell is a cell of the farm's water, dryCell one of its grass.
func waterCell(f *farm) (col, row int) {
	for row := 0; row < f.water.Rows; row++ {
		for col := 0; col < f.water.Cols; col++ {
			if f.isWater(col, row) {
				return col, row
			}
		}
	}
	panic("no water on the farm")
}

func dryCell(f *farm) (col, row int) {
	for row := 0; row < f.water.Rows; row++ {
		for col := 0; col < f.water.Cols; col++ {
			if !f.isWater(col, row) && f.plotAt(col, row) == nil {
				return col, row
			}
		}
	}
	panic("no dry land on the farm")
}

// Every crop is sold as seeds, sells its harvest for more the longer it
// takes, has a name on screen, and is asked for by villagers.
func TestEveryCropIsSoldAndAsked(t *testing.T) {
	_, f := newGame(idle{})
	if len(crops) != 15 {
		t.Fatalf("%d crops", len(crops))
	}
	f.money = 100000
	for c := range crops {
		f.choose(f.kit, &shop{choice: c, at: farmStall})
		if f.bag.Count(seedsOf(c)) != 1 {
			t.Errorf("%s seeds not sold", crops[c].name)
		}
		f.bag.Add(cropOf(c), 1)
		if itemLabels[cropOf(c)] != crops[c].label || crops[c].label == "" {
			t.Errorf("%s has no label", crops[c].name)
		}
	}
	want := 0
	for _, c := range crops {
		want += c.sellPrice
	}
	f.money = 0
	f.choose(f.kit, &shop{choice: sellChoice, at: farmStall})
	if f.money != want || worth(f.kit) != 0 {
		t.Errorf("sold for %d, not %d", f.money, want)
	}
	for i, c := range crops {
		for _, o := range crops[i+1:] {
			if o.daysToRipe > c.daysToRipe && o.regrowDays == c.regrowDays && o.sellPrice <= c.sellPrice {
				t.Errorf("%s (%dd) pays no more than %s (%dd)", o.name, o.daysToRipe, c.name, c.daysToRipe)
			}
		}
	}
	for _, c := range crops {
		found := false
		for _, o := range orderBook.Offers {
			found = found || o.Item == c.name
		}
		if !found {
			t.Errorf("no villager asks for %s", c.name)
		}
	}
}

// Left to themselves, the farmers of a town keep their gardens alive over
// several days: they water each day and refill the can.
func TestFarmersKeepTheirCropsAlive(t *testing.T) {
	g, f := newGame(idle{})
	tn := eastTown(f)
	g.Find("farmer").Transform.Position = standIn(f, tn.Center.Col, tn.Center.Row+1)
	var gardens []*plot
	for _, v := range f.people {
		if v.town == tn && v.role == roleFarmer && v.home != nil {
			p := f.plotAt(v.home.Garden.Col, v.home.Garden.Row)
			p.crop = 1
			v.kit.used = canCapacity - 1 // nearly out of water
			gardens = append(gardens, p)
		}
	}
	if len(gardens) == 0 {
		t.Skip("no farmer in this town")
	}
	f.clock.Hour = 6
	f.clock.Speed *= 6 // a day in a minute
	for range 3600 * 4 {
		f.fainted = f.clock.Day
		g.Step(frame)
	}
	grown := 0
	for _, p := range gardens {
		if p.dead {
			t.Errorf("a garden withered: %+v", p)
		}
		grown += p.days
	}
	if grown < len(gardens) {
		t.Errorf("the gardens grew %d days in all", grown)
	}
}

// Both arts draw every tile, crops at every stage and withered, and every
// icon: the tools, the seeds and the harvest of each crop.
func TestTheArtDrawsEveryCrop(t *testing.T) {
	a := SVGArt()
	for _, name := range TileNames() {
		if cv := a.Tile(name); cv == nil || cv.Width == 0 {
			t.Errorf("no picture for tile %s", name)
		}
	}
	icons := []string{"icon_hoe", "icon_can", "icon_wood"}
	for _, c := range crops {
		icons = append(icons, "icon_"+c.seedIcon, "icon_"+c.cropIcon)
	}
	for _, name := range icons {
		if cv := a.Icon(name); cv == nil || cv.Width == 0 {
			t.Errorf("no picture for %s", name)
		}
	}
}
