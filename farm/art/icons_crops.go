package art

// cropIcons are the colours that tell each new crop's seed packet and
// harvest apart: its fruit and its shade.
var cropIcons = map[string][2]string{
	"lettuce":    {"#c8ec8c", "#7fc65a"},
	"potato":     {"#e0b878", "#b48a4e"},
	"beans":      {"#9ad060", "#5a9a3a"},
	"corn":       {yellow, "#d6a428"},
	"wheat":      {"#f4dc88", "#d4aa48"},
	"pumpkin":    {orange, orangeShade},
	"cassava":    {"#c8a070", "#8a6540"},
	"strawberry": {red, redShade},
	"watermelon": {"#5fb14a", "#2f7a3a"},
	"pepper":     {"#ff5a40", "#b02a20"},
	"sugarcane":  {"#a8d070", "#6a9a40"},
	"sunflower":  {yellow, "#7a4b2a"},
	"coffee":     {red, "#4a2c1a"},
}

// harvestIcon is a crop picked: a round fruit in its colours with a
// little stalk and leaf.
func harvestIcon(c [2]string) string {
	d := newDoc(Size, Size)
	blob{8, 9.4, 5.4, 4.8, c[0], c[1]}.draw(d)
	d.add(`<path d="M8 5 q-0.4 -2 1 -3 M8 5.2 q-2.4 -2 -3.6 -0.8" stroke="%s" stroke-width="1" fill="none" stroke-linecap="round"/>`, leafBase)
	return d.String()
}

// hoeIcon is a hoe: a wooden handle and a blade of metal across its end.
func hoeIcon() string {
	d := newDoc(Size, Size)
	rod{3, 14, 8, -9, 1.2, woodBase}.draw(d)
	d.add(`<path d="M10.4 5 L14.6 6.4 L13.8 9 L9.6 7.4 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, d.down(metal, metalShade), line)
	return d.String()
}
