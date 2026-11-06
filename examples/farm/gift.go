package farm

import (
	"errors"

	"github.com/gabrielluizsf/sge/examples/i18n"
)

// gifts are what the player may give from the hotbar: the harvest and
// wood.
var gifts = []string{"carrot", "tomato", wood}

// How fond a gift makes one: of what they love, of what they do not mind,
// and of what they dislike.
const (
	loved    = 80
	liked    = 30
	disliked = -40
)

// favourite is the gift v loves: wood for the smith and the lumberjacks,
// whose trades burn and sell it, and for the rest a crop their numbers
// pick.
func favourite(v *villager) int {
	if v.role == roleSmith || v.role == roleLumberjack {
		return 2
	}
	return v.who.Genes[3] % 2
}

// liking is how fond item makes v: loved, disliked, or liked, as their
// trade and their numbers say.
func liking(v *villager, item string) int {
	fav := favourite(v)
	hate := (fav + 1 + v.who.Genes[4]%2) % len(gifts)
	switch item {
	case gifts[fav]:
		return loved
	case gifts[hate]:
		return disliked
	}
	return liked
}

// gift is the player giving v one of item, at most one gift a day: v is
// as much fonder as they like it, and says so. It is why not, if not.
func (f *farm) gift(v *villager, item string) error {
	if f.bag.Count(item) == 0 {
		return errNone
	}
	n := liking(v, item)
	if err := v.bond.Give(n, f.clock.Day); err != nil {
		f.speak(v, i18n.T("I already got a gift today."))
		f.sound("deny", nil)
		return err
	}
	f.bag.Take(item, 1)
	f.earn(xpGift)
	switch {
	case n == loved:
		f.speak(v, i18n.TF(v.woman, "I love it! Thank you!"))
	case n < 0:
		f.speak(v, i18n.T("Hmm... I don't like that."))
	default:
		f.speak(v, i18n.TF(v.woman, "Thank you!"))
	}
	f.sound("coin", nil)
	return nil
}

// errNone is having none of what was to be given.
var errNone = errors.New("farm: none to give")
