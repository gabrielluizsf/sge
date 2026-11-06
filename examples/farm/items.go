package farm

// wood is what a tree felled gives, and what it sells for.
const (
	wood      = "wood"
	woodPrice = 4
)

// seedsOf and cropOf are the items of crop c's seeds and harvest.
func seedsOf(c int) string { return crops[c].name + "_seeds" }
func cropOf(c int) string  { return crops[c].name }

// itemLabels are the items' names on screen, as English keys.
var itemLabels = labelled()

// labelled is the wood's name and every crop's, by item.
func labelled() map[string]string {
	m := map[string]string{wood: "Wood"}
	for _, c := range crops {
		m[c.name] = c.label
	}
	return m
}

// worth is what kit k's harvest and wood sell for.
func worth(k *kit) int {
	w := k.bag.Count(wood) * woodPrice
	for c := range crops {
		w += k.bag.Count(cropOf(c)) * crops[c].sellPrice
	}
	return w
}
