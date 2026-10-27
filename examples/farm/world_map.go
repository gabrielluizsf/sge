package main

// worldMap is the farm, a character a cell: g grass, G flowers, p path,
// w water, d a plot of soil, f fence, T a tree, C the stall's
// cell. t, where a tree's crown reaches over from the cell in front, is
// grass.
const worldMap = `
ffffffffffffffffffffffffffffff
gGgtggggggggggggggggggtgggggGg
gggTggggggggggggggggggTgggtggg
gggggGgggggCggggggGgggggggTggg
pppppppppppppppppppppppppppppp
ggggggggggggggpgggggggggggggGg
gddddddddddgggpggggtgggwwwwwgg
gddddddddddgggpggggTggwwwwwwwg
gddddddddddgGgpgggggggwwwwwwwg
gddddddddddgggpgggggggwwwwwwgg
gddddddddddgggpGgggggggwwwwggg
ggggggggggggggpggggggggggggggg
gGggtggggGggggpgggggtgggGggggg
ggggTgggggggggpgggggTggggggtgg
ggggggggggGgggpggGgggggggggTgg
gggggggtggggggpggggggggGgggggg
ggGggggTggggggpggggtgggggggGgg
ggggggggggggGgpggggTgggggggggg
ggggGgggggggggpggggggggGgggggg
ffffffffffffffffffffffffffffff`
