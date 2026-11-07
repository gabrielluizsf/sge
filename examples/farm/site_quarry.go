package farm

import "github.com/gabrielluizsf/sge/pathfind"

// chopQuarry is the trees about v's town, that they fell with their axe.
func (f *farm) chopQuarry(v *villager) quarry {
	return quarry{
		area: area{center: pathfind.Cell(v.place().Center), cols: 15, rows: 9},
		is:   isTree, slot: slotAxe, uses: axes[v.kit.axe].hits,
		where: func(int, int) bool { return true },
	}
}

// mineQuarry is the rocks about v's town, that they break with a pickaxe;
// never the player's own land.
func (f *farm) mineQuarry(v *villager) quarry {
	return quarry{
		area: area{center: pathfind.Cell(v.place().Center), cols: 15, rows: 9},
		is:   isRock, slot: slotPick, uses: rockHits,
		where: func(col, row int) bool { return !f.onFarm(col, row) },
	}
}
