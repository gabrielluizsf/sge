package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/i18n"
)

// category is a group of rows, a line of the list on the left.
type category struct {
	name string
	rows []row
}

// The ancestry sliders move the skin tone to what suits the ancestry.
func ancestryChanged(s *charState) {
	s.tone = human.SkinToneOf(s.anc)
	s.fit()
}

func slide(label string, lo, hi float64, field func(*charState) *float64, after func(*charState)) row {
	return row{label: label, kind: sliderRow, slider: &slider{lo, hi, field, after}}
}

func partRow(label string, list func(*charState) []string, field func(*charState) *string) row {
	return row{label: label, kind: choiceRow,
		text: func(s *charState) string { return ordinal(list(s), *field(s)) },
		step: func(s *charState, d int) { *field(s) = cycle(list(s), *field(s), d) }}
}

func colourRow(label string, list []human.Colour, field func(*charState) *human.Colour, original bool) row {
	return row{label: label, kind: choiceRow,
		text: func(s *charState) string {
			if !field(s).Set {
				return i18n.T("Original")
			}
			if i := nearest(list, *field(s)); list[i] == *field(s) {
				return fmt.Sprintf("%d / %d", i+1, len(list))
			}
			return i18n.T("Custom")
		},
		swatch: func(s *charState) human.Colour { return *field(s) },
		step:   func(s *charState, d int) { *field(s) = cycleColour(list, *field(s), d, original) }}
}

// codeRow is a row where a colour's code is typed, or taken from the
// screen with the eyedropper.
func codeRow(label string, field func(*charState) *human.Colour) row {
	return row{label: label, kind: textRow, field: field,
		press: (*creator).startCode}
}

// categories are the creation screen's, in the order of its list: the
// body, skin, face, hair, eyes and clothes; Aleatório and Começar after.
func categories() []category {
	hasMain := func(s *charState) bool { return human.Colourable(s.main) }
	return []category{
		{"Body", []row{
			{label: "Name", kind: textRow, text: func(s *charState) string { return s.name },
				press: func(c *creator) { c.typing = true }},
			{label: "Sex", kind: choiceRow, text: func(s *charState) string { return sexName(s.sex) },
				step: func(s *charState, _ int) {
					s.sex = human.Woman + human.Man - s.sex
					s.fit()
				}},
			slide("Age", youngest, oldest, func(s *charState) *float64 { return &s.years }, (*charState).fit),
			slide("Muscle", 0, 1, func(s *charState) *float64 { return &s.muscle }, nil),
			slide("Weight", 0, 1, func(s *charState) *float64 { return &s.weight }, nil),
			slide("Height", 0, 1, func(s *charState) *float64 { return &s.height }, nil),
			{label: "Language", kind: choiceRow, text: func(*charState) string { return i18n.Name() },
				step: func(_ *charState, d int) { _ = i18n.Change(d) }},
		}},
		{"Skin", []row{
			slide("African ancestry", 0, 1, func(s *charState) *float64 { return &s.anc.African }, ancestryChanged),
			slide("Asian ancestry", 0, 1, func(s *charState) *float64 { return &s.anc.Asian }, ancestryChanged),
			slide("European ancestry", 0, 1, func(s *charState) *float64 { return &s.anc.European }, ancestryChanged),
			slide("Skin tone", 0.02, 1, func(s *charState) *float64 { return &s.tone }, nil),
		}},
		{"Face", []row{
			partRow("Eyebrows", func(*charState) []string { return human.Eyebrows() }, func(s *charState) *string { return &s.brows }),
			func() row {
				r := partRow("Eyelashes", func(s *charState) []string { return human.Eyelashes(s.sex) }, func(s *charState) *string { return &s.lashes })
				r.show = func(s *charState) bool { return len(human.Eyelashes(s.sex)) > 1 }
				return r
			}(),
		}},
		{"Hair", []row{
			partRow("Style", func(s *charState) []string { return s.options(human.Hair) }, func(s *charState) *string { return &s.hair }),
			colourRow("Color", human.HairColours(), func(s *charState) *human.Colour { return &s.hairColour }, false),
		}},
		{"Eyes", []row{
			{label: "Color", kind: choiceRow, text: func(s *charState) string { return eyeName(s.eyes) },
				step: func(s *charState, d int) { s.eyes = cycle(human.EyeMaterials(), s.eyes, d) }},
		}},
		{"Clothes", []row{
			partRow("Piece", (*charState).mains, func(s *charState) *string { return &s.main }),
			func() row {
				r := partRow("Bottom", func(s *charState) []string { return s.options(human.Bottom) }, func(s *charState) *string { return &s.bottom })
				r.show = (*charState).wantsBottom
				return r
			}(),
			partRow("Shoes", func(s *charState) []string { return s.options(human.Shoes) }, func(s *charState) *string { return &s.shoes }),
			{label: "Hat", kind: choiceRow,
				text: func(s *charState) string {
					if s.hat == "" {
						return i18n.T("None")
					}
					return ordinal(s.options(human.Hat), s.hat)
				},
				step: func(s *charState, d int) { s.hat = cycleHat(s.options(human.Hat), s.hat, d) }},
			func() row {
				r := colourRow("Item color", human.ClothColours(), func(s *charState) *human.Colour { return &s.mainTint }, true)
				r.show = hasMain
				return r
			}(),
			func() row {
				r := codeRow("Item code", func(s *charState) *human.Colour { return &s.mainTint })
				r.show = hasMain
				return r
			}(),
			func() row {
				r := colourRow("Bottom color", human.ClothColours(), func(s *charState) *human.Colour { return &s.bottomTint }, true)
				r.show = func(s *charState) bool { return s.wantsBottom() }
				return r
			}(),
			func() row {
				r := codeRow("Bottom code", func(s *charState) *human.Colour { return &s.bottomTint })
				r.show = func(s *charState) bool { return s.wantsBottom() }
				return r
			}(),
		}},
		{"Random", []row{{label: "Random", kind: buttonRow, press: (*creator).randomise}}},
		{"Begin", []row{{label: "Begin", kind: buttonRow, press: (*creator).finish}}},
	}
}
