package main

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
)

const frame = 1.0 / 60

func TestWorldMapRowsAreFull(t *testing.T) {
	rows := strings.Split(strings.Trim(worldMap, "\n"), "\n")
	for i, r := range rows {
		if len(r) != len(rows[0]) {
			t.Errorf("row %d is %d cells wide, not %d", i, len(r), len(rows[0]))
		}
	}
}

// idle is a controller that does nothing: tests drive the farm.
type idle struct{}

func (idle) control(*farm, *game.Object) intent { return intent{} }

func TestACropIsSownWateredGrownAndPicked(t *testing.T) {
	g, f := newGame(idle{})
	p := f.plots[0]
	f.slot = slotCarrot
	if f.use(p.col, p.row) {
		t.Fatal("sowed without seeds")
	}
	f.seeds[slotCarrot] = 1
	if !f.use(p.col, p.row) || p.crop != slotCarrot || f.seeds[slotCarrot] != 0 {
		t.Fatalf("not sown: %+v", p)
	}
	if f.crops.At(p.col, p.row) != tileID["carrot0"] {
		t.Error("the crop layer does not show the seeds")
	}
	for stage := 1; stage <= ripe; stage++ {
		// Dry, it does not grow.
		for range 600 {
			g.Step(frame)
		}
		if p.stage != stage-1 {
			t.Fatalf("grew to %d without water", p.stage)
		}
		f.slot = slotCan
		if !f.use(p.col, p.row) || f.ground.At(p.col, p.row) != tileID["soil_wet"] {
			t.Fatal("not watered")
		}
		for range int(crops[slotCarrot].stageTime/frame) + 2 {
			g.Step(frame)
		}
		if p.stage != stage || p.watered || f.ground.At(p.col, p.row) != tileID["soil_dry"] {
			t.Fatalf("stage %d: %+v", stage, p)
		}
	}
	if f.use(p.col, p.row) != true || f.harvest[slotCarrot] != 1 || p.crop != -1 {
		t.Fatalf("not picked: %+v, harvest %v", p, f.harvest)
	}
}

func TestTheShopSellsSeedsAndBuysTheHarvest(t *testing.T) {
	_, f := newGame(idle{})
	f.money = 20
	f.choose(&shop{choice: 1})
	f.choose(&shop{choice: 0})
	f.choose(&shop{choice: 0})
	f.choose(&shop{choice: 0}) // with no money left
	if f.money != 0 || f.seeds[0] != 2 || f.seeds[1] != 1 {
		t.Fatalf("money %d, seeds %v", f.money, f.seeds)
	}
	f.harvest[0], f.harvest[1] = 2, 1
	f.choose(&shop{choice: sellChoice})
	if f.money != 2*12+25 || f.harvest[0]+f.harvest[1] != 0 {
		t.Fatalf("money %d, harvest %v", f.money, f.harvest)
	}
}

func TestTheFarmerStopsAtTheFenceAndTheWater(t *testing.T) {
	g, f := newGame(idle{})
	o := g.Find("farmer")
	o.Transform.Position = standIn(f, 5, 1)
	tilemap.Move(o, 0, -200, f.water, f.things, f.wild)
	if _, row, _ := f.ground.CellAt(feet(o).X, feet(o).Y-feetBox.Height/2); row != 1 {
		t.Errorf("walked into the fence, to row %d", row)
	}
	o.Transform.Position = standIn(f, 18, 7)
	tilemap.Move(o, 400, 0, f.water, f.things, f.wild)
	if col, _, _ := f.ground.CellAt(feet(o).X+feetBox.Width/2, feet(o).Y); col >= 22 {
		t.Errorf("walked into the water, to column %d", col)
	}
}

// The demo plays the whole loop: it buys seeds, grows them, sells the
// crop and ends up richer.
func TestTheDemoMakesMoney(t *testing.T) {
	g, f := newGame(&bot{tends: 10})
	for range 7200 {
		g.Step(frame)
	}
	if f.sold == 0 || f.money+f.seeds[0]*5+f.seeds[1]*10 <= 20 {
		t.Errorf("sold %d, money %d, seeds %v", f.sold, f.money, f.seeds)
	}
}
