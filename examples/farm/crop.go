package main

// crop is a kind of plant: what its seeds cost, what it sells for, and how
// long each stage takes to grow once watered.
type crop struct {
	name      string // its sprites' name
	label     string // in the shop
	seedPrice int
	sellPrice int
	stageTime float64
	seedIcon  string
	cropIcon  string
}

// crops are what the farm grows. A crop is planted at stage 0 and ripe at
// stage ripe.
var crops = []crop{
	{name: "carrot", label: "cenoura", seedPrice: 5, sellPrice: 12, stageTime: 4, seedIcon: "carrotSeeds", cropIcon: "carrot"},
	{name: "tomato", label: "tomate", seedPrice: 10, sellPrice: 25, stageTime: 6, seedIcon: "tomatoSeeds", cropIcon: "tomato"},
}

const ripe = 3
