package main

import (
	"errors"
	"testing"

	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/screenpick"
)

// fakeScreen is a screen that shows one colour under a mouse that clicks
// when it is told to.
type fakeScreen struct {
	at     screenpick.Reading
	err    error
	closed bool
}

func (f *fakeScreen) Read() (screenpick.Reading, error) { return f.at, f.err }
func (f *fakeScreen) Close()                            { f.closed = true }

// codeCreator is a creator with a top chosen, on its item code row.
func codeCreator(t *testing.T) *creator {
	t.Helper()
	c, _ := newTestCreator(t)
	c.st.main = human.Options(human.Top, human.Woman, 30, 0)[0]
	c.st.sex, c.st.years = human.Woman, 30
	c.st.fit()
	c.st.mainTint = human.Colour{}
	for cat := range c.cats {
		for _, i := range c.shown(cat) {
			if c.cats[cat].rows[i].label == "Item code" {
				c.cat, c.row = cat, i
			}
		}
	}
	if c.at().label != "Item code" {
		t.Fatal("no item code row")
	}
	return c
}

// typeCodeText presses the row and types s.
func typeCodeText(c *creator, s string) {
	c.handle(cin{Select: true})
	c.handle(cin{Typed: s})
}

func TestACodeTypedColoursTheClothesAtOnce(t *testing.T) {
	c := codeCreator(t)
	typeCodeText(c, "#FF8000")
	got := c.st.mainTint
	if !got.Set || got.R != 1 || got.B != 0 || got.G < 0.5 || got.G > 0.51 {
		t.Errorf("tint = %+v", got)
	}
	if err := c.st.character().Check(); err != nil {
		t.Errorf("the character is not valid: %v", err)
	}
	c.handle(cin{Select: true})
	if c.typing {
		t.Error("Enter on a good code did not leave the field")
	}
}

func TestAnRGBACodeWorksAndAlphaIsSaidToBeIgnored(t *testing.T) {
	c := codeCreator(t)
	typeCodeText(c, "rgba(0, 255, 0, 0.5)")
	if got := c.st.mainTint; !got.Set || got.G != 1 || got.R != 0 {
		t.Errorf("tint = %+v", got)
	}
	if c.note == "" {
		t.Error("no word that transparency is ignored")
	}
}

func TestABadCodeChangesNothingAndSaysSo(t *testing.T) {
	c := codeCreator(t)
	before := c.st.mainTint
	for _, bad := range []string{"#12", "rgba(999999999999, 0, 0, 1)", "red", "#GGGGGG", "rgba(1,2,3,9)"} {
		typeCodeText(c, bad)
		c.handle(cin{Select: true})
		if c.st.mainTint != before {
			t.Errorf("%q changed the tint to %+v", bad, c.st.mainTint)
		}
		if c.note == "" || !c.typing {
			t.Errorf("%q: note %q, typing %v", bad, c.note, c.typing)
		}
		c.handle(cin{Cancel: true})
		if c.typing || c.note != "" {
			t.Errorf("Escape left typing=%v note=%q", c.typing, c.note)
		}
	}
}

func TestACodeKeepsOnlyWhatIsAllowedAndLimited(t *testing.T) {
	c := codeCreator(t)
	typeCodeText(c, "<script>\x00#ab")
	for _, r := range c.code {
		if r == '<' || r == 0 {
			t.Errorf("code kept %q", c.code)
		}
	}
	c.handle(cin{Typed: string(make([]byte, 5000))})
	c.handle(cin{Typed: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	if len(c.code) > 40 {
		t.Errorf("code grew to %d", len(c.code))
	}
}

func TestACustomColourSavesAndLoads(t *testing.T) {
	c := codeCreator(t)
	typeCodeText(c, "#336699")
	ch := c.st.character()
	if err := ch.Check(); err != nil {
		t.Fatal(err)
	}
	if codeOf(c.st.mainTint) != "#336699" {
		t.Errorf("code = %s", codeOf(c.st.mainTint))
	}
}

func pressDropper(c *creator, screen *fakeScreen, err error) {
	c.openScreen = func() (screenpick.Source, error) {
		if err != nil {
			return nil, err
		}
		return screen, nil
	}
	c.handle(cin{Pick: true})
}

func TestTheEyedropperTakesTheColourOfTheFirstClickAfterTheButtonIsUp(t *testing.T) {
	c := codeCreator(t)
	s := &fakeScreen{at: screenpick.Reading{Down: true, R: 10, G: 20, B: 30}}
	pressDropper(c, s, nil)
	c.handle(cin{})
	if c.st.mainTint.Set {
		t.Fatal("the click that began it took the colour")
	}
	s.at = screenpick.Reading{Down: false, R: 200, G: 100, B: 50}
	c.handle(cin{})
	if c.st.mainTint.Set {
		t.Fatal("hovering changed the clothes")
	}
	s.at.Down = true
	c.handle(cin{})
	if got := codeOf(c.st.mainTint); got != "#C86432" {
		t.Errorf("took %s, want #C86432", got)
	}
	if c.drop != nil || !s.closed {
		t.Error("the screen was kept after the pick")
	}
}

func TestEscapeCancelsTheEyedropperAndChangesNothing(t *testing.T) {
	c := codeCreator(t)
	s := &fakeScreen{at: screenpick.Reading{R: 1}}
	pressDropper(c, s, nil)
	c.handle(cin{})
	c.handle(cin{Cancel: true})
	if c.drop != nil || !s.closed || c.st.mainTint.Set {
		t.Errorf("drop %v closed %v tint %+v", c.drop, s.closed, c.st.mainTint)
	}
}

func TestTheEyedropperSaysWhenItIsUnavailable(t *testing.T) {
	c := codeCreator(t)
	pressDropper(c, nil, errors.New("no screen"))
	if c.drop != nil || c.note == "" {
		t.Errorf("drop %v note %q", c.drop, c.note)
	}
	c.handle(cin{})
}

func TestAScreenThatFailsEndsTheEyedropper(t *testing.T) {
	c := codeCreator(t)
	s := &fakeScreen{}
	pressDropper(c, s, nil)
	s.err = errors.New("gone")
	c.handle(cin{})
	if c.drop != nil || !s.closed || c.note == "" {
		t.Errorf("drop %v closed %v note %q", c.drop, s.closed, c.note)
	}
}
