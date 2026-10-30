package farm

// wood is what a tree felled gives, and what it sells for.
const (
	wood      = "wood"
	woodPrice = 4
)

// seedsOf and cropOf are the items of crop c's seeds and harvest.
func seedsOf(c int) string { return crops[c].name + "_seeds" }
func cropOf(c int) string  { return crops[c].name }

// itemLabels are the items' names on screen.
var itemLabels = map[string]string{wood: "madeira", "carrot": "cenoura", "tomato": "tomate"}

// worth is what kit k's harvest and wood sell for.
func worth(k *kit) int {
	w := k.bag.Count(wood) * woodPrice
	for c := range crops {
		w += k.bag.Count(cropOf(c)) * crops[c].sellPrice
	}
	return w
}
