package main

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/kit"
	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/save"
)

// newTestCreator is a creation screen with nothing to bake with, and the
// character it was done with.
func newTestCreator(t *testing.T) (*creator, *[]kit.Character) {
	t.Helper()
	var got []kit.Character
	c := newCreator(testArt(t), nil, func(v any) { got = append(got, v.(kit.Character)) })
	return c, &got
}

// coherent reports why c is not a person that suits their sex and age.
func coherent(t *testing.T, s *charState) {
	t.Helper()
	ch := s.character()
	if err := ch.Check(); err != nil {
		t.Fatalf("not a valid character: %v", err)
	}
	for _, w := range ch.Look.Wear {
		p, ok := human.PartOf(w)
		if !ok {
			continue
		}
		if !slices.Contains(human.Options(p, s.sex, s.years, s.african()), w) {
			t.Errorf("%s does not suit a %s of %.0f", w, sexName(s.sex), s.years)
		}
	}
	if s.wantsBottom() != (s.bottom != "") {
		t.Errorf("main %s, bottom %q", s.main, s.bottom)
	}
}

func TestEveryRowChangesTheCharacter(t *testing.T) {
	c, _ := newTestCreator(t)
	seen := 0
	for cat := range c.cats {
		for _, i := range c.shown(cat) {
			r := &c.cats[cat].rows[i]
			if r.kind == buttonRow || r.kind == textRow {
				continue
			}
			if r.label == "Language" {
				continue // the game's, not the character's
			}
			c.st = stateOf(human.Random(human.Spec{Sex: human.Woman, Age: 30, Seed: 5}))
			c.st.main = human.Options(human.Top, human.Woman, 30, 0)[0] // a top, so the bottom row is shown
			c.st.fit()
			c.cat, c.row = cat, i
			before := c.st.character().Tag()
			c.change(1)
			if c.st.character().Tag() == before {
				c.change(-1) // at the end of a slider
			}
			if c.st.character().Tag() == before {
				t.Errorf("%s / %s changes nothing", c.cats[cat].name, r.label)
			}
			coherent(t, c.st)
			seen++
		}
	}
	if seen < 18 {
		t.Errorf("only %d rows", seen)
	}
}

func TestEveryCategoryHasRowsToShow(t *testing.T) {
	c, _ := newTestCreator(t)
	for cat := range c.cats {
		if len(c.shown(cat)) == 0 {
			t.Errorf("%s shows no row", c.cats[cat].name)
		}
	}
	want := []string{"Body", "Skin", "Face", "Hair", "Eyes", "Clothes", "Random", "Begin"}
	for i, w := range want {
		if c.cats[i].name != w {
			t.Errorf("category %d is %s", i, c.cats[i].name)
		}
	}
}

func TestChangingTheSexKeepsThePersonCoherent(t *testing.T) {
	c, _ := newTestCreator(t)
	for seed := uint32(1); seed < 30; seed++ {
		c.st = stateOf(human.Random(human.Spec{Seed: seed}))
		c.cat, c.row = 0, 1 // Sexo
		sex := c.st.sex
		c.change(1)
		if c.st.sex == sex || c.st.body().Gender == float64(sex-1) {
			t.Fatalf("seed %d: the sex did not change", seed)
		}
		coherent(t, c.st)
	}
}

func TestAgeAndAncestryKeepThePersonCoherent(t *testing.T) {
	c, _ := newTestCreator(t)
	for seed := uint32(1); seed < 20; seed++ {
		c.st = stateOf(human.Random(human.Spec{Seed: seed}))
		for _, cr := range [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, 2}} {
			c.cat, c.row = cr[0], cr[1]
			for range 25 {
				c.change(1)
			}
			coherent(t, c.st)
			for range 40 {
				c.change(-1)
			}
			coherent(t, c.st)
		}
	}
}

func TestAncestryMovesTheSkinToneAlong(t *testing.T) {
	c, _ := newTestCreator(t)
	c.st.anc = human.Ancestry{European: 1}
	c.cat, c.row = 1, 0
	c.cats[1].rows[0].slider.set(c.st, 1) // all African
	c.st.anc = human.Ancestry{African: 1}
	ancestryChanged(c.st)
	dark := c.st.tone
	c.st.anc = human.Ancestry{European: 1}
	ancestryChanged(c.st)
	if dark <= c.st.tone+0.3 {
		t.Errorf("African %.2f, European %.2f", dark, c.st.tone)
	}
}

func TestAleatorioMakesANewCoherentPerson(t *testing.T) {
	c, _ := newTestCreator(t)
	seen := map[string]bool{c.st.character().Tag(): true}
	c.pick(6)
	if c.at().label != "Random" {
		t.Fatalf("chose %s", c.at().label)
	}
	for range 25 {
		c.handle(cin{Select: true})
		coherent(t, c.st)
		seen[c.st.character().Tag()] = true
	}
	if len(seen) < 24 {
		t.Errorf("%d different people in 26", len(seen))
	}
}

func TestComecarCallsDoneOnceWithTheLook(t *testing.T) {
	c, got := newTestCreator(t)
	c.pick(0)
	c.change(1)
	want := c.st.character()
	c.handle(cin{Dy: -1}) // wraps to the last row: Começar
	if c.at().label != "Begin" {
		t.Fatalf("on %s", c.at().label)
	}
	c.handle(cin{Select: true})
	c.handle(cin{Select: true})
	if len(*got) != 1 || (*got)[0].Tag() != want.Tag() || (*got)[0].Check() != nil {
		t.Errorf("done got %d characters: %v", len(*got), *got)
	}
}

func TestTheKeysPickRowsAndChangeValues(t *testing.T) {
	c, _ := newTestCreator(t)
	c.pick(0)
	c.handle(cin{Dy: 1}) // past the name
	if c.at().label != "Sex" {
		t.Fatalf("starts on %s", c.at().label)
	}
	sex := c.st.sex
	c.handle(cin{Dx: 1})
	if c.st.sex == sex {
		t.Error("right did not change the sex")
	}
	c.handle(cin{Dy: 1})
	c.handle(cin{Dy: 1})
	if c.at().label != "Muscle" {
		t.Errorf("down twice is %s", c.at().label)
	}
	m := c.st.muscle
	c.handle(cin{Dx: -1})
	if c.st.muscle >= m {
		t.Error("left did not lower the slider")
	}
	for range 10 {
		c.handle(cin{Dy: 1})
	}
	if c.cat == 0 {
		t.Error("down did not go across categories")
	}
}

func TestTheMouseClicksArrowsSlidersAndCategories(t *testing.T) {
	c, got := newTestCreator(t)
	cv, _ := canvas.NewCanvas(1280, 720)
	c.paint(cv)
	click := func(kind hitKind, cat, row int) {
		t.Helper()
		c.paint(cv)
		for _, h := range c.hits {
			if h.kind == kind && h.cat == cat && (kind == hitCategory || h.row == row) {
				c.handle(cin{Press: true, Held: true, X: h.at.X + h.at.W/2, Y: h.at.Y + h.at.H/2})
				c.handle(cin{})
				return
			}
		}
		t.Fatalf("nothing to click: kind %v cat %d row %d", kind, cat, row)
	}
	click(hitCategory, 3, 0)
	if c.cat != 3 {
		t.Fatalf("category %d", c.cat)
	}
	hair := c.st.hair
	click(hitMore, 3, 0)
	if c.st.hair == hair {
		t.Error("the right arrow did not change the hair")
	}
	click(hitLess, 3, 0)
	if c.st.hair != hair {
		t.Error("the left arrow did not bring the hair back")
	}
	click(hitCategory, 0, 0)
	c.paint(cv)
	for _, h := range c.hits {
		if h.kind == hitTrack && h.row == 3 { // muscle: left end, then right end
			c.handle(cin{Press: true, Held: true, X: h.at.X, Y: h.at.Y + 5})
			c.handle(cin{})
			if c.st.muscle != 0 {
				t.Errorf("left end of the track is %v", c.st.muscle)
			}
			c.handle(cin{Press: true, Held: true, X: h.at.X + h.at.W/2, Y: h.at.Y + 5})
			c.handle(cin{Held: true, X: h.at.X + h.at.W, Y: h.at.Y + 5}) // dragged to the right end
			c.handle(cin{})
			if c.st.muscle < 0.99 {
				t.Errorf("dragged to the end it is %v", c.st.muscle)
			}
		}
	}
	click(hitCategory, 7, 0) // Começar
	if len(*got) != 1 {
		t.Errorf("done called %d times", len(*got))
	}
}

func TestThePaintedScreenHasText(t *testing.T) {
	c, _ := newTestCreator(t)
	cv, _ := canvas.NewCanvas(1280, 720)
	c.paint(cv)
	lit := 0
	for y := 0; y < cv.Height; y += 4 {
		for x := 0; x < cv.Width; x += 4 {
			if cv.At(x, y).A() > 0 {
				lit++
			}
		}
	}
	if lit == 0 || len(c.hits) < 12 {
		t.Errorf("lit %d, hits %d", lit, len(c.hits))
	}
}

func TestNewGameAddsTheScreenAndReleasesTheBakerWhenDone(t *testing.T) {
	a := testArt(t).(*realArt)
	a.world = &world{baker: &bakes.Baker{}}
	g := game.New()
	var chosen any
	c := a.screenFor(g, func(v any) { chosen = v })
	if !a.world.baker.Hold.Load() {
		t.Error("the background bakes were not held")
	}
	var o *game.Object
	for _, x := range g.Objects() {
		if x.Name == "creator" {
			o = x
		}
	}
	if o == nil || !o.Always || !o.Screen {
		t.Fatalf("object %+v", o)
	}
	c.pick(7)
	c.handle(cin{Select: true})
	if chosen == nil || a.world.baker.Hold.Load() {
		t.Errorf("chosen %v, hold %v", chosen, a.world.baker.Hold.Load())
	}
}

func TestChosenMakesThePlayerAndKeepsTheLook(t *testing.T) {
	defer kit.ResetPlayer()
	a := testArt(t).(*realArt)
	a.store = save.In(t.TempDir(), 1)
	b, l := human.Random(human.Spec{Seed: 99})
	c := kit.Character{Body: b, Look: l}
	a.Chosen(c)
	if kit.PlayerCharacter().Tag() != c.Tag() || kit.SpriteTag("player") == "" {
		t.Error("the player is not the character")
	}
	if got, ok := kit.LoadCharacter(a.store); !ok || got.Tag() != c.Tag() {
		t.Error("the look was not kept")
	}
	if _, ok := farmPerson(a); !ok {
		t.Error("the art draws no player")
	}
}

func TestChosenRefusesWhatIsNotACharacter(t *testing.T) {
	defer kit.ResetPlayer()
	a := testArt(t).(*realArt)
	a.store = save.In(t.TempDir(), 1)
	before := kit.PlayerCharacter().Tag()
	bad := kit.Character{Look: human.Look{Wear: []string{"../../etc/passwd"}}}
	for _, v := range []any{nil, "Ana", bad} {
		a.Chosen(v)
	}
	if kit.PlayerCharacter().Tag() != before || a.store.Has("look") {
		t.Error("something that is not a character was taken")
	}
}

// farmPerson is the art's player.
func farmPerson(a *realArt) (farm.Body, bool) {
	b := a.Person(farm.Who{Player: true})
	return b, b != nil
}
