package main

// worldMap is the farm, a character a cell: g grass, G flowers, p path,
// w water, d a plot of soil, f fence, T a tree, C c the stall's two
// cells. t, R and r, where a tree's crown and the stall's awning reach
// over from the cell in front, are grass.
const worldMap = `
ffffffffffffffffffffffffffffff
gGgtggggggggggggggggggtgggggGg
gggTgggggggRrgggggggggTgggtggg
gggggGgggggCcgggggGgggggggTggg
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
