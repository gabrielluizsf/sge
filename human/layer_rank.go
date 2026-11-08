package human

// layerOf is how far out from the body the proxy item is worn: 1 for
// bottoms, 2 for tops, dresses and suits, and 0
// for what is not worn over the body, shoes and hats, which wrap what they
// are on and are left as they are.
func layerOf(item string) int {
	p, ok := PartOf(item)
	if !ok {
		return 0
	}
	switch p {
	case Bottom:
		return 1
	case Top, Dress, Suit:
		return 2
	}
	return 0
}
