package main

import (
	"slices"
	"strings"

	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/textfield"
)

// charState is what the character creation screen edits: every choice a
// value, from which the character is made. A choice that stops suiting
// another (a sex that has no dress) is fitted by fit.
type charState struct {
	sex                          human.Sex
	years                        float64
	anc                          human.Ancestry // each 0 to 1; the body has their shares
	tone, muscle, weight, height float64
	brows, lashes                string
	hair                         string // "" for none
	hairColour                   human.Colour
	eyes                         string
	main, bottom, shoes, hat     string // main is a top, a dress or a suit
	mainTint, bottomTint         human.Colour
	name                         string // what the player calls them, as typed
	named                        bool   // the player typed it: a change of sex leaves it
	nameSeed                     uint32 // picks the default name
}

// The ages the creator offers, in years.
const (
	youngest = 18
	oldest   = 80
)

// african is the share of African ancestry, for the hair that suits it.
func (s *charState) african() float64 {
	sum := s.anc.African + s.anc.Asian + s.anc.European
	if sum <= 0 {
		return 1.0 / 3
	}
	return s.anc.African / sum
}

// options are the proxies of part p that suit the person now.
func (s *charState) options(p human.Part) []string {
	return human.Options(p, s.sex, s.years, s.african())
}

// mains are the pieces worn on the body: tops, dresses and suits.
func (s *charState) mains() []string {
	return slices.Concat(s.options(human.Top), s.options(human.Dress), s.options(human.Suit))
}

// wantsBottom reports whether the main piece is a top, which a bottom
// goes with.
func (s *charState) wantsBottom() bool {
	p, _ := human.PartOf(s.main)
	return p == human.Top
}

// fit changes what no longer suits the person, the sex, age or ancestry
// having changed, to the nearest that does, and gives a top its bottom.
func (s *charState) fit() {
	s.hair = keepOr(s.hair, s.options(human.Hair), true)
	s.main = keepOr(s.main, s.mains(), false)
	s.shoes = keepOr(s.shoes, s.options(human.Shoes), false)
	if s.hat != "" {
		s.hat = keepOr(s.hat, s.options(human.Hat), true)
	}
	s.lashes = keepOr(s.lashes, human.Eyelashes(s.sex), false)
	if s.wantsBottom() {
		s.bottom = keepOr(s.bottom, s.options(human.Bottom), false)
	} else {
		s.bottom = ""
	}
	if !s.named {
		s.name = s.defaultName()
	}
}

// keepOr is cur if list has it, else the first of list, or "" for an
// empty list; it also keeps "" where it may be none.
func keepOr(cur string, list []string, noneOK bool) string {
	switch {
	case cur == "" && noneOK, slices.Contains(list, cur):
		return cur
	case len(list) == 0:
		return ""
	}
	return list[0]
}

// body is the body the state asks for.
func (s *charState) body() human.Body {
	a := s.anc
	if sum := a.African + a.Asian + a.European; sum <= 0 {
		a = human.Ancestry{African: 1, Asian: 1, European: 1}
	}
	sum := a.African + a.Asian + a.European
	b := human.Body{Age: human.AgeFor(s.years), African: a.African / sum, Asian: a.Asian / sum, European: a.European / sum,
		Muscle: s.muscle, Weight: s.weight, Height: s.height, Proportions: 0.5}
	if s.sex == human.Man {
		b.Gender = 1
	}
	return b
}

// character is the character the state makes.
func (s *charState) character() kit.Character {
	wear := []string{human.EyesProxy(), s.brows, s.lashes}
	var tints []human.Colour
	add := func(item string, tint human.Colour) {
		if item == "" {
			return
		}
		wear = append(wear, item)
		tints = append(tints, human.Colour{})
		if human.Colourable(item) {
			tints[len(tints)-1] = tint
		}
	}
	add(s.hair, human.Colour{})
	add(s.main, s.mainTint)
	add(s.bottom, s.bottomTint)
	add(s.shoes, human.Colour{})
	add(s.hat, human.Colour{})
	l := human.Look{Eyes: s.eyes, Wear: wear, SkinTone: min(max(s.tone, 0.02), 1), HairColour: s.hairColour}
	if slices.ContainsFunc(tints, func(c human.Colour) bool { return c.Set }) {
		l.ClothesColours = append(make([]human.Colour, 3), tints...)
	}
	name := textfield.Trim(s.name)
	if name == "" {
		name = s.defaultName()
	}
	return kit.Character{Body: s.body(), Look: l, Name: name}
}

// stateOf is the state of a body and look made by human.Random: the
// choices read back out of what they were made into.
func stateOf(b human.Body, l human.Look) *charState {
	s := &charState{sex: human.Woman, years: b.Years(), anc: human.Ancestry{African: b.African, Asian: b.Asian, European: b.European},
		tone: l.SkinTone, muscle: b.Muscle, weight: b.Weight, height: b.Height, eyes: l.Eyes, hairColour: l.HairColour}
	if b.Gender >= 0.5 {
		s.sex = human.Man
	}
	s.years = min(max(s.years, youngest), oldest)
	for i, w := range l.Wear {
		var tint human.Colour
		if i < len(l.ClothesColours) {
			tint = l.ClothesColours[i]
		}
		s.wear(w, tint)
	}
	if s.lashes == "" {
		s.lashes = human.Eyelashes(s.sex)[0]
	}
	s.fit()
	return s
}

// wear puts the proxy w, with tint, where it goes in the state.
func (s *charState) wear(w string, tint human.Colour) {
	switch p, ok := human.PartOf(w); {
	case strings.HasPrefix(w, "eyebrows/"):
		s.brows = w
	case strings.HasPrefix(w, "eyelashes/"):
		s.lashes = w
	case !ok:
	case p == human.Hair:
		s.hair = w
	case p == human.Top, p == human.Dress, p == human.Suit:
		s.main, s.mainTint = w, tint
	case p == human.Bottom:
		s.bottom, s.bottomTint = w, tint
	case p == human.Shoes:
		s.shoes = w
	case p == human.Hat:
		s.hat = w
	}
}

// defaultName is the name the person has until the player types one: a
// woman's or a man's from the game's lists.
func (s *charState) defaultName() string { return farm.RandomName(s.sex == human.Woman, s.nameSeed) }
