package main

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/save"
)

// focusName puts the cursor in the name field.
func focusName(t *testing.T, c *creator) {
	t.Helper()
	c.pick(0)
	c.row = 0
	if c.at().label != "Name" {
		t.Fatalf("the first row is %q, not the name", c.at().label)
	}
	c.handle(cin{Select: true})
	if !c.typing {
		t.Fatal("Enter did not focus the name")
	}
}

func TestTheNameStartsAsAFittingRandomOne(t *testing.T) {
	c, _ := newTestCreator(t)
	n := c.st.character().Name
	if !strings.Contains(strings.Join(farmNames(), " "), n) || kit.PlayerName() == "x" {
		t.Errorf("the default name %q is not from the lists", n)
	}
	if err := c.st.character().Check(); err != nil {
		t.Error(err)
	}
}

func farmNames() []string {
	var all []string
	for seed := range uint32(200) {
		all = append(all, farm.RandomName(true, seed), farm.RandomName(false, seed))
	}
	return all
}

func TestTypingReplacesNothingUntilTheNameIsFocused(t *testing.T) {
	c, _ := newTestCreator(t)
	c.pick(0)
	before := c.st.name
	c.handle(cin{Typed: "zzz"})
	if c.st.name != before {
		t.Errorf("typed into an unfocused field: %q", c.st.name)
	}
}

func TestTypingBackspaceAndTheLimit(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	for range 30 {
		c.handle(cin{Back: true})
	}
	if c.st.name != "" {
		t.Fatalf("backspace left %q", c.st.name)
	}
	c.handle(cin{Typed: "João "})
	c.handle(cin{Typed: "d'Ávila-Souza"})
	if c.st.name != "João d'Ávila-Souza" {
		t.Errorf("typed %q", c.st.name)
	}
	c.handle(cin{Typed: strings.Repeat("x", 40)})
	if n := len([]rune(c.st.name)); n != kit.MaxName {
		t.Errorf("%d characters", n)
	}
	c.handle(cin{Back: true})
	if n := len([]rune(c.st.name)); n != kit.MaxName-1 {
		t.Errorf("%d characters after a backspace", n)
	}
}

func TestInvalidCharactersAreRefusedWhileTyping(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	for range 30 {
		c.handle(cin{Back: true})
	}
	c.handle(cin{Typed: "A1\x00<b>\n/../‮Zé"})
	if c.st.name != "AbZé" {
		t.Errorf("name %q", c.st.name)
	}
}

// With the field focused, typing is typing: it does not move the row or
// the category; the arrows up and down leave it.
func TestTypingDoesNotNavigate(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	cat, row := c.cat, c.row
	c.handle(cin{Typed: "wasd qe"})
	c.handle(cin{Dx: 1})
	if c.cat != cat || c.row != row || !c.typing {
		t.Errorf("moved to %d/%d typing=%v", c.cat, c.row, c.typing)
	}
	c.handle(cin{Dy: 1})
	if c.typing || c.row == row && c.cat == cat {
		t.Errorf("the arrows did not leave the field: %d/%d typing=%v", c.cat, c.row, c.typing)
	}
	c.handle(cin{Typed: "zzz"})
	if strings.HasSuffix(c.st.name, "zzz") {
		t.Error("typed after leaving the field")
	}
}

func TestEnterLeavesTheFieldAndAClickFocusesIt(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	c.handle(cin{Select: true})
	if c.typing {
		t.Error("Enter did not leave the field")
	}
	c.hits = []hit{{at: rect{0, 0, 100, 40}, kind: hitRow, cat: 0, row: 0}}
	c.handle(cin{X: 10, Y: 10, Press: true, Held: true})
	if !c.typing {
		t.Error("a click did not focus the field")
	}
	c.hits = []hit{{at: rect{0, 0, 100, 40}, kind: hitRow, cat: 0, row: 1}}
	c.handle(cin{X: 10, Y: 10, Press: true, Held: true})
	if c.typing {
		t.Error("a click elsewhere left the field focused")
	}
}

func TestTheNameIsTrimmedAndEmptyFallsBackToTheDefault(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	for range 30 {
		c.handle(cin{Back: true})
	}
	c.handle(cin{Typed: "  Ana  "})
	if got := c.st.character().Name; got != "Ana" {
		t.Errorf("name %q", got)
	}
	for range 30 {
		c.handle(cin{Back: true})
	}
	got := c.st.character().Name
	if got == "" || got != c.st.defaultName() {
		t.Errorf("an empty field gave %q, default %q", got, c.st.defaultName())
	}
	if err := c.st.character().Check(); err != nil {
		t.Error(err)
	}
}

func TestRandomChangesTheNameAndASexChangeRepicksADefaultOne(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	c.handle(cin{Typed: "Z"})
	typed := c.st.name
	c.randomise()
	if c.st.name == typed {
		t.Errorf("Aleatório kept the name %q", typed)
	}
	c.st.named = false
	c.st.sex = human.Woman
	c.st.fit()
	if !strings.Contains(strings.Join(womenOf(), " "), c.st.name) {
		t.Errorf("a woman is called %q", c.st.name)
	}
	c.st.sex = human.Man
	c.st.fit()
	if strings.Contains(strings.Join(womenOf(), " "), c.st.name) && !strings.Contains(strings.Join(menOf(), " "), c.st.name) {
		t.Errorf("a man is called %q", c.st.name)
	}
	c.st.name, c.st.named = "Mine", true
	c.st.fit()
	if c.st.name != "Mine" {
		t.Error("the player's own name was replaced")
	}
}

func womenOf() (out []string) {
	for s := range uint32(200) {
		out = append(out, farm.RandomName(true, s))
	}
	return
}

func menOf() (out []string) {
	for s := range uint32(200) {
		out = append(out, farm.RandomName(false, s))
	}
	return
}

func TestFinishGivesTheName(t *testing.T) {
	c, got := newTestCreator(t)
	focusName(t, c)
	for range 30 {
		c.handle(cin{Back: true})
	}
	c.handle(cin{Typed: "Iara"})
	c.finish()
	if len(*got) != 1 || (*got)[0].Name != "Iara" {
		t.Errorf("done got %+v", *got)
	}
}

func TestTheNameIsPainted(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	c.handle(cin{Typed: "Ana"})
	if got := c.at().text(c.st); got != c.st.name {
		t.Errorf("row shows %q", got)
	}
}

func TestChosenKeepsTheNameForTheFarmAndTheSave(t *testing.T) {
	defer kit.ResetPlayer()
	a := testArt(t).(*realArt)
	a.store = saveIn(t)
	c := someone(12)
	c.Name = "Kwame"
	a.Chosen(c)
	if a.PlayerName() != "Kwame" {
		t.Errorf("the farm is told %q", a.PlayerName())
	}
	if got, ok := kit.LoadCharacter(a.store); !ok || got.Name != "Kwame" {
		t.Errorf("saved %q, ok=%v", got.Name, ok)
	}
	c.Name = "<bad>\n"
	a.Chosen(c)
	if a.PlayerName() != "Kwame" {
		t.Errorf("a bad name was taken: %q", a.PlayerName())
	}
}

func saveIn(t *testing.T) *save.Store { return save.In(t.TempDir(), 1) }

func TestTypingOverTheDefaultNameReplacesIt(t *testing.T) {
	c, _ := newTestCreator(t)
	focusName(t, c)
	c.handle(cin{Typed: "Ana"})
	if c.st.name != "Ana" {
		t.Errorf("name %q", c.st.name)
	}
}
