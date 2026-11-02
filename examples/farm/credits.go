package farm

import _ "embed"

// Credit is a third-party component the game uses: what it is, who made
// it, where it is from, and the text of the licence it is used under,
// which is read when the player asks to see it.
type Credit struct {
	Name, Author, URL string
	// License is the licence's text; nil when there is none to show.
	License func() string
}

// Credited is an Art that says what it is made of, for the credits.
type Credited interface {
	Credits() []Credit
}

//go:embed fonts/OFL.txt
var snigletLicense string

// Credits are the stylized art's third-party components: the font of its
// letters.
func (a *svgArt) Credits() []Credit {
	return []Credit{{
		Name: "Fonte Sniglet", Author: "Haley Fiege, Brenda Gallo e Pablo Impallari",
		URL: "https://www.fontsquirrel.com/fonts/sniglet", License: func() string { return snigletLicense },
	}}
}

// creditsOf are the credits of the art a, none if it says none.
func creditsOf(a Art) []Credit {
	if c, ok := a.(Credited); ok {
		return c.Credits()
	}
	return nil
}
