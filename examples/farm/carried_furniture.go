package farm

import (
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/examples/i18n"
)

// pieceLabels are the furniture's names on screen, as English keys, and
// pieceUses what each is for.
var (
	pieceLabels = map[string]string{"bed": "Bed", "table": "Table", "chair": "Chair", "shelf": "Shelf", "stove": "Stove", "plant": "Plant pot", "chest": "Chest"}
	pieceUses   = map[string]string{
		"bed": "Sleep in it to end the day", "table": "A place to set things down", "chair": "Sit on it with E",
		"shelf": "Tall and tidy", "stove": "Warms the house", "plant": "Greens up a corner", "chest": "Keeps 10 stacks of items",
	}
)

// furniture is the piece of furniture name carried, to be placed in the
// house, or the chest outside.
func (c carried) furniture(name string) backpack.Item {
	return backpack.Item{ID: name, Name: i18n.T(pieceLabels[name]), Icon: pieceIcon + name, Count: c.f.bag.Count(name), Hand: true,
		Info: []string{i18n.T(pieceUses[name]), i18n.T("B: arrange the house"), i18n.T("Pick it back up to carry it")}}
}

// stone is the stone from the rocks.
func (c carried) stone() backpack.Item {
	return backpack.Item{ID: stone, Name: i18n.T("Stone"), Icon: "icon_" + stone, Count: c.f.bag.Count(stone),
		Info: []string{i18n.T("Furniture is made of it")}}
}
