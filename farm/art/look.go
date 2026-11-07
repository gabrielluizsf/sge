package art

// Look is how a person is drawn: their skin, light and shade; their hair,
// light, base and shade, and its Style; whether they wear the straw hat;
// their shirt and their trousers, light, base and shade; and whether the
// trousers are overalls, with a bib and straps, or a Skirt to the knee.
type Look struct {
	Skin     [2]string
	Hair     [3]string
	Style    string
	Hat      bool
	Shirt    [3]string
	Pants    [3]string
	Overalls bool
	Skirt    bool
}

// Hair styles.
const (
	Short    = "short"
	Long     = "long"
	Bun      = "bun"
	Ponytail = "ponytail"
	Curly    = "curly"
	Buzz     = "buzz"
)

// FarmerLook is the farmer's: the character sheet's straw hat, brown hair,
// cream shirt and denim overalls.
var FarmerLook = Look{
	Skin:     [2]string{skin, skinShade},
	Hair:     [3]string{hairLight, hairBase, hairShade},
	Style:    Short,
	Hat:      true,
	Shirt:    [3]string{shirtLight, shirtBase, shirtShade},
	Pants:    [3]string{denimLight, denimBase, denimShade},
	Overalls: true,
}
