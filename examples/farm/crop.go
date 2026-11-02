package farm

// crop is a kind of plant: what its seeds cost, what it sells for, and how
// many days of watering it takes to ripen.
type crop struct {
	name  string // its sprites' name
	label string // in the shop
	// daysToRipe is how many days, each watered, it takes to ripen;
	// regrowDays how many more, once picked, until it bears again, or 0
	// for a crop that is gone once picked.
	daysToRipe int
	regrowDays int
	seedPrice  int
	sellPrice  int
	seedIcon   string
	cropIcon   string
}

// crops are what the farm grows, in the order the shop lists them. A
// crop is planted at stage 0 and ripe at stage ripe.
var crops = priced([]crop{
	{name: "carrot", label: "cenoura", daysToRipe: 4},
	{name: "lettuce", label: "alface", daysToRipe: 3},
	{name: "potato", label: "batata", daysToRipe: 6},
	{name: "tomato", label: "tomate", daysToRipe: 7, regrowDays: 3},
	{name: "beans", label: "feijão", daysToRipe: 6, regrowDays: 3},
	{name: "corn", label: "milho", daysToRipe: 9},
	{name: "wheat", label: "trigo", daysToRipe: 5},
	{name: "pumpkin", label: "abóbora", daysToRipe: 12},
	{name: "cassava", label: "mandioca", daysToRipe: 14},
	{name: "strawberry", label: "morango", daysToRipe: 8, regrowDays: 4},
	{name: "watermelon", label: "melancia", daysToRipe: 11},
	{name: "pepper", label: "pimenta", daysToRipe: 7, regrowDays: 3},
	{name: "sugarcane", label: "cana-de-açúcar", daysToRipe: 13},
	{name: "sunflower", label: "girassol", daysToRipe: 8},
	{name: "coffee", label: "café", daysToRipe: 18, regrowDays: 5},
})

const ripe = 3

// priced is cs with their prices and icons filled in: a harvest sells for
// more the longer the crop takes, a crop that bears again for less each
// time, and the seeds for a third of the harvest.
func priced(cs []crop) []crop {
	for i := range cs {
		c := &cs[i]
		per := 3
		if c.regrowDays > 0 {
			per = 2
		}
		c.sellPrice = per*c.daysToRipe + c.daysToRipe*c.daysToRipe/8
		c.seedPrice = max(3, c.sellPrice/3)
		c.seedIcon, c.cropIcon = c.name+"Seeds", c.name
	}
	return cs
}
