package farm

// spriteSVG is the SVG of the sprite name, one of the embedded assets.
func spriteSVG(name string) string {
	src, err := assets.ReadFile("assets/" + name + ".svg")
	if err != nil {
		panic(err)
	}
	return string(src)
}
