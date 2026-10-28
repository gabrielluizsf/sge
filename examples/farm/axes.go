package main

// axe is a kind of axe: its name, as its icon's; its label; how many
// strokes fell a tree with it; what the smith sells it for; and at most
// how much wood a tree felled with it gives.
type axe struct {
	name, label string
	hits, price int
	most        int
}

// axes are from the worst to the best. The stone axe is the smith's
// reward for their errand; the others the smith sells.
var axes = []axe{
	{name: "stone", label: "PEDRA", hits: 4, most: 2},
	{name: "iron", label: "FERRO", hits: 3, price: 60, most: 3},
	{name: "gold", label: "OURO", hits: 2, price: 150, most: 3},
	{name: "diamond", label: "DIAMANTE", hits: 1, price: 400, most: 4},
}

// noAxe is the axe of one who has none.
const noAxe = -1
