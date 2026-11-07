package farm

import "github.com/gabrielluizsf/sge/inventory"

// stone is what rocks give when mined; the recipes ask for it by name.
const stone = "stone"

// recipe is a piece of furniture the Craft tab makes: the item it makes,
// named as the piece is, and the materials it takes.
type recipe struct {
	item  string
	needs inventory.List
}

// recipes are the Craft tab's: the pickaxe first, the way to stone, which
// takes wood alone as stone itself takes the pickaxe; then the furniture,
// cheapest first, the second bed last, the player's house having only the
// first bed to begin with.
var recipes = []recipe{
	{itemPick, inventory.List{{Item: wood, N: 6}}},
	{"chair", inventory.List{{Item: wood, N: 3}}},
	{"table", inventory.List{{Item: wood, N: 6}}},
	{"chest", inventory.List{{Item: wood, N: 10}}},
	{"shelf", inventory.List{{Item: wood, N: 8}}},
	{"plant", inventory.List{{Item: wood, N: 2}, {Item: stone, N: 1}}},
	{"stove", inventory.List{{Item: wood, N: 4}, {Item: stone, N: 10}}},
	{"bed", inventory.List{{Item: wood, N: 10}, {Item: stone, N: 2}, {Item: wool, N: 2}}},
}

// recipeOf is the recipe that makes item.
func recipeOf(item string) (recipe, bool) {
	for _, r := range recipes {
		if r.item == item {
			return r, true
		}
	}
	return recipe{}, false
}
