package art

// Skins are skin tones from the lightest to the deepest, each a light and
// a shade, so people of every part of the world live in the towns.
var Skins = [][2]string{
	{"#ffe3cc", "#f0bd9d"},
	{"#f6cfae", "#dba684"},
	{"#e5b088", "#c68b62"},
	{"#c98e5f", "#a46b40"},
	{"#9c6644", "#7a4a2e"},
	{"#6e4532", "#4f2f20"},
}

// Hairs are hair colours, light, base and shade: black, dark brown,
// brown, auburn, blond and grey.
var Hairs = [][3]string{
	{"#4a4a58", "#26262f", "#121218"},
	{"#6e4a33", "#47301f", "#2b1b10"},
	{"#9a6238", "#6e4126", "#4a2a16"},
	{"#e0804c", "#b4532c", "#7f3519"},
	{"#ffe6a3", "#e8c060", "#b8903a"},
	{"#e4e4ea", "#b4b6be", "#80838d"},
}

// Shirts and Trousers are clothes' colours, light, base and shade.
var (
	Shirts = [][3]string{
		{"#ff8a80", "#e5544b", "#b23a33"}, // red
		{"#9be08a", "#5fb453", "#3f8740"}, // green
		{"#8fe3e0", "#3fb5b0", "#2a8584"}, // teal
		{"#d3a6f5", "#a46ad6", "#7748a8"}, // purple
		{"#ffe98a", "#f2c53d", "#c7962a"}, // yellow
		{"#ffc2d8", "#f2839f", "#c75a7a"}, // pink
		{"#fffaf0", "#efe3c8", "#cdbb98"}, // cream
		{"#9cc4ff", "#5b8fe0", "#3c66ad"}, // blue
	}
	Trousers = [][3]string{
		{denimLight, denimBase, denimShade},
		{"#b98a5e", "#8d613b", "#654326"}, // brown
		{"#6c6c78", "#45454f", "#2c2c33"}, // charcoal
		{"#a3ad6a", "#7a8445", "#58602f"}, // olive
		{"#e07a6a", "#b5503f", "#86362a"}, // brick
	}
)

// Villager is the look of a villager, a woman or not, from five numbers
// picking their skin, hair colour, hair style, shirt and trousers; each
// number may be any size, as it is taken round the choices.
func Villager(woman bool, picks [5]int) Look {
	styles := []string{Short, Buzz, Curly, Ponytail}
	if woman {
		styles = []string{Long, Bun, Ponytail, Curly, Short}
	}
	at := func(i, n int) int { return ((i % n) + n) % n }
	l := Look{
		Skin:  Skins[at(picks[0], len(Skins))],
		Hair:  Hairs[at(picks[1], len(Hairs))],
		Style: styles[at(picks[2], len(styles))],
		Shirt: Shirts[at(picks[3], len(Shirts))],
		Pants: Trousers[at(picks[4], len(Trousers))],
	}
	// Some women wear skirts, as the picks fall.
	l.Skirt = woman && at(picks[3]+picks[4], 2) == 0
	return l
}
