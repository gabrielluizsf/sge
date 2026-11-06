package farm

import "github.com/gabrielluizsf/sge/backpack"

// openHouse opens the screens of the furniture on s, to look at them:
// "craft" is the inventory's Craft tab with materials to make some of it,
// "arrange" the player's house with a few pieces, the cursor on a chair
// about to be put down, and "chest" a chest with things in it beside the
// bag. It reports whether what was one of them.
func (s *Session) openHouse(what string) bool {
	f := s.f
	switch what {
	case "craft":
		f.bag.Add(wood, 9)
		f.bag.Add(stone, 3)
		f.pack.Show()
		f.pack.SetTab(1)
		f.pack.Handle(backpackDown)
		f.syncPackHint()
	case "arrange", "chest":
		f.hands.ctl = still{}
		f.enter(houseDoor{at: homeDoor, faceOut: "down", seed: seed, own: true})
		f.player.Transform.Position = f.ground.CellCenter(entry[0], entry[1])
		f.bag.Add("chair", 2)
		f.bag.Add(wood, 12)
		f.bag.Add(seedsOf(1), 4)
		f.home = append(f.home, newPiece("table", 4, 3, 0), newPiece("shelf", 6, 1, 0), newPiece("chest", 1, 4, 0))
		f.home[3].held[0] = chestSlot{item: stone, n: 5}
		f.relay()
		if what == "chest" {
			f.openChest(&f.home[3])
			break
		}
		f.startEdit("chair")
		f.edit.col, f.edit.row = 5, 3
	default:
		return false
	}
	return true
}

// backpackDown is one step down the craft list.
var backpackDown = backpack.Input{Dy: 1}
