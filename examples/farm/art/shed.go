package art

// IsShed reports whether name is one of the towns' shops that are drawn
// like the stall, on a cell of their own.
func IsShed(name string) bool { return name == "stable" }

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
