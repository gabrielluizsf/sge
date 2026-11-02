package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/human"
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
				return "Original"
			}
			return fmt.Sprintf("%d / %d", nearest(list, *field(s))+1, len(list))
		},
		swatch: func(s *charState) human.Colour { return *field(s) },
		step:   func(s *charState, d int) { *field(s) = cycleColour(list, *field(s), d, original) }}
}

// categories are the creation screen's, in the order of its list: the
// body, skin, face, hair, eyes and clothes; Aleatório and Começar after.
func categories() []category {
	hasMain := func(s *charState) bool { return human.Colourable(s.main) }
	return []category{
		{"Corpo", []row{
			{label: "Nome", kind: textRow, text: func(s *charState) string { return s.name },
				press: func(c *creator) { c.typing = true }},
			{label: "Sexo", kind: choiceRow, text: func(s *charState) string { return sexName(s.sex) },
				step: func(s *charState, _ int) {
					s.sex = human.Woman + human.Man - s.sex
					s.fit()
				}},
			slide("Idade", youngest, oldest, func(s *charState) *float64 { return &s.years }, (*charState).fit),
			slide("Músculos", 0, 1, func(s *charState) *float64 { return &s.muscle }, nil),
			slide("Peso", 0, 1, func(s *charState) *float64 { return &s.weight }, nil),
			slide("Altura", 0, 1, func(s *charState) *float64 { return &s.height }, nil),
		}},
		{"Pele", []row{
			slide("Ancestral africana", 0, 1, func(s *charState) *float64 { return &s.anc.African }, ancestryChanged),
			slide("Ancestral asiática", 0, 1, func(s *charState) *float64 { return &s.anc.Asian }, ancestryChanged),
			slide("Ancestral europeia", 0, 1, func(s *charState) *float64 { return &s.anc.European }, ancestryChanged),
			slide("Tom de pele", 0.02, 1, func(s *charState) *float64 { return &s.tone }, nil),
		}},
		{"Rosto", []row{
			partRow("Sobrancelhas", func(*charState) []string { return human.Eyebrows() }, func(s *charState) *string { return &s.brows }),
			func() row {
				r := partRow("Cílios", func(s *charState) []string { return human.Eyelashes(s.sex) }, func(s *charState) *string { return &s.lashes })
				r.show = func(s *charState) bool { return len(human.Eyelashes(s.sex)) > 1 }
				return r
			}(),
		}},
		{"Cabelo", []row{
			partRow("Estilo", func(s *charState) []string { return s.options(human.Hair) }, func(s *charState) *string { return &s.hair }),
			colourRow("Cor", human.HairColours(), func(s *charState) *human.Colour { return &s.hairColour }, false),
		}},
		{"Olhos", []row{
			{label: "Cor", kind: choiceRow, text: func(s *charState) string { return eyeName(s.eyes) },
				step: func(s *charState, d int) { s.eyes = cycle(human.EyeMaterials(), s.eyes, d) }},
		}},
		{"Roupa", []row{
			partRow("Peça", (*charState).mains, func(s *charState) *string { return &s.main }),
			func() row {
				r := partRow("Parte de baixo", func(s *charState) []string { return s.options(human.Bottom) }, func(s *charState) *string { return &s.bottom })
				r.show = (*charState).wantsBottom
				return r
			}(),
			partRow("Sapatos", func(s *charState) []string { return s.options(human.Shoes) }, func(s *charState) *string { return &s.shoes }),
			{label: "Chapéu", kind: choiceRow,
				text: func(s *charState) string {
					if s.hat == "" {
						return "Nenhum"
					}
					return ordinal(s.options(human.Hat), s.hat)
				},
				step: func(s *charState, d int) { s.hat = cycleHat(s.options(human.Hat), s.hat, d) }},
			func() row {
				r := colourRow("Cor da peça", human.ClothColours(), func(s *charState) *human.Colour { return &s.mainTint }, true)
				r.show = hasMain
				return r
			}(),
			func() row {
				r := colourRow("Cor de baixo", human.ClothColours(), func(s *charState) *human.Colour { return &s.bottomTint }, true)
				r.show = func(s *charState) bool { return s.wantsBottom() }
				return r
			}(),
		}},
		{"Aleatório", []row{{label: "Aleatório", kind: buttonRow, press: (*creator).randomise}}},
		{"Começar", []row{{label: "Começar", kind: buttonRow, press: (*creator).finish}}},
	}
}
