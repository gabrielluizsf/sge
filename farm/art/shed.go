package art

// IsShed reports whether name is one of the towns' shops that are drawn
// like the stall, on a cell of their own.
func IsShed(name string) bool { return name == "stable" || name == "bakery" || name == "butcher" }

// Stable is a town's stable: the stall's counter and poles, under a green
// awning, and bales of hay in place of its baskets.
func Stable() string {
	return stallOf(stallLook{stripe: leafBase, hem: leafShade, goods: func(d *doc) {
		for i := range 3 {
			u := 6 + float64(i)*8
			d.box([4]float64{u, 7, u + 6, 12}, [2]float64{10, 14}, [3]string{"#f0d98a", "#d9b95e", "#b4933f"})
		}
	}})
}

// Bakery is a town's bakery: the stall under a yellow awning, loaves on its
// counter.
func Bakery() string {
	return stallOf(stallLook{stripe: yellow, hem: sandShade, goods: func(d *doc) {
		for i := range 3 {
			u := 6 + float64(i)*8
			d.box([4]float64{u, 7, u + 5, 11}, [2]float64{10, 12.4}, [3]string{"#efc06a", "#d29a48", "#b9782e"})
		}
	}})
}

// Butcher is a town's butcher's: the stall under a dark red awning, cuts of
// meat on its counter.
func Butcher() string {
	return stallOf(stallLook{stripe: redShade, hem: "#5a2420", goods: func(d *doc) {
		for i := range 3 {
			u := 6 + float64(i)*8
			d.box([4]float64{u, 7, u + 5, 11}, [2]float64{10, 11.6}, [3]string{"#c4605a", "#a8322e", "#7a2420"})
		}
	}})
}

// Whites is a baker's clothes, as light, middle and shade.
var Whites = [3]string{"#ffffff", "#f2efe8", "#d4d0c6"}
