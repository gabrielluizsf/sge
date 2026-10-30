package farm

import "github.com/gabrielluizsf/sge/bond"

// weddings marries the player to the villager they are engaged to, on
// their wedding day.
func (f *farm) weddings(day int) {
	for _, v := range f.people {
		if v.bond.Courted == bond.Engaged && v.bond.Wed(day) {
			f.marry(v)
			return
		}
	}
}

// marry makes v the player's spouse: they leave their town's life and
// move into the player's house, wanting nothing more of them.
func (f *farm) marry(v *villager) {
	f.self.spouse = v
	v.ours, v.inside, v.order.Want = true, true, nil
	if v.obj != nil {
		v.obj.Destroy()
		v.obj = nil
	}
	f.notify("Você se casou com " + v.name + "! Agora vivem juntos.")
	f.sound("sell", nil)
}
