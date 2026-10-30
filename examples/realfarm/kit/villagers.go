package kit

import (
	"fmt"

	"github.com/gabrielluizsf/sge/human"
)

// Villager is one of the farm's people: who they are, their work, their
// body and how they look.
type Villager struct {
	ID, Role string
	Body     human.Body
	Look     human.Look
}

// age is a body's Age for years.
func age(years float64) float64 {
	if years < 25 {
		return (years - 1) / 48
	}
	return 0.5 + (years-25)/130
}

// face is the eyes, brows and lashes everyone has.
func face(brows, lashes int) []string {
	return []string{"eyes/high-poly/high-poly.mhclo",
		fmt.Sprintf("eyebrows/eyebrow%03d/eyebrow%03d.mhclo", brows, brows),
		fmt.Sprintf("eyelashes/eyelashes%02d/eyelashes%02d.mhclo", lashes, lashes)}
}

// clothes is the .mhclo of a pack's item.
func clothes(item string) string { return "clothes/" + item + "/" + item + ".mhclo" }

func hair(style string) string { return "hair/" + style + "/" + style + ".mhclo" }

func eyes(colour string) string { return "eyes/materials/" + colour + ".mhmat" }

// Villagers are the farm's people: women and men of Africa, Asia and
// Europe and of all of them, young to old, each dressed for their work.
var Villagers = []Villager{
	{"joao", "farmer", human.Body{Gender: 1, Age: age(42), African: 0.4, European: 0.6, Muscle: 0.65, Weight: 0.55, Height: 0.55, Proportions: 0.5},
		human.Look{Eyes: eyes("brown"), Wear: append(face(1, 1), hair("short02"), clothes("male_worksuit01"), clothes("shoes02"), clothes("fedora_cocked"))}},
	{"ana", "farmer", human.Body{Gender: 0, Age: age(32), African: 0.8, European: 0.2, Muscle: 0.55, Weight: 0.5, Height: 0.5, Proportions: 0.6},
		human.Look{Eyes: eyes("brown"), Wear: append(face(2, 2), hair("braid01"), clothes("toigo_keyhole_tank_top"), clothes("toigo_long_full_skirt"), clothes("toigo_flats"))}},
	{"beto", "smith", human.Body{Gender: 1, Age: age(46), African: 1, Muscle: 0.85, Weight: 0.6, Height: 0.6, Proportions: 0.5},
		human.Look{Eyes: eyes("brown"), Wear: append(face(3, 1), hair("short01"), clothes("toigo_basic_tucked_t-shirt"), clothes("toigo_wool_pants"), clothes("culturalibre_male_boots"))}},
	{"rui", "lumberjack", human.Body{Gender: 1, Age: age(35), European: 1, Muscle: 0.75, Weight: 0.6, Height: 0.65, Proportions: 0.5},
		human.Look{Eyes: eyes("bluegreen"), Wear: append(face(4, 1), hair("short03"), clothes("toigo_fisherman_sweater"), clothes("toigo_wool_pants"), clothes("toigo_ankle_boots_male"))}},
	{"mei", "walker", human.Body{Gender: 0, Age: age(27), Asian: 1, Muscle: 0.45, Weight: 0.4, Height: 0.4, Proportions: 0.6},
		human.Look{Eyes: eyes("brown"), Wear: append(face(5, 3), hair("bob01"), clothes("toigo_shift_dress"), clothes("toigo_ballet_flats"))}},
	{"lucia", "walker", human.Body{Gender: 0, Age: age(66), African: 0.3, European: 0.7, Muscle: 0.4, Weight: 0.65, Height: 0.4, Proportions: 0.5},
		human.Look{Eyes: eyes("green"), Wear: append(face(6, 2), hair("bob02"), clothes("toigo_bodice-style_top"), clothes("toigo_tiered_skirt"), clothes("toigo_flats"))}},
	{"kenji", "walker", human.Body{Gender: 1, Age: age(71), Asian: 1, Muscle: 0.4, Weight: 0.45, Height: 0.35, Proportions: 0.5},
		human.Look{Eyes: eyes("brown"), Wear: append(face(7, 1), hair("short04"), clothes("male_casualsuit03"), clothes("shoes03"), clothes("fedora01"))}},
	{"iara", "walker", human.Body{Gender: 0, Age: age(22), African: 0.35, Asian: 0.3, European: 0.35, Muscle: 0.5, Weight: 0.45, Height: 0.55, Proportions: 0.6},
		human.Look{Eyes: eyes("brownlight"), Wear: append(face(8, 4), hair("long01"), clothes("female_casualsuit02"), clothes("shoes01"))}},
	{"player", "player", human.Body{Gender: 1, Age: age(26), African: 0.3, European: 0.5, Asian: 0.2, Muscle: 0.6, Weight: 0.5, Height: 0.55, Proportions: 0.55},
		human.Look{Eyes: eyes("brown"), Wear: append(face(9, 1), hair("short02"), clothes("male_casualsuit02"), clothes("shoes04"))}},
}
