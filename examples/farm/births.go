package farm

// births has the player's child born on day, if one is due: they live in
// the player's house, are there to see if the player is home, and the
// player hears the news.
func (f *farm) births(day int) {
	e := f.self.expecting
	if e == nil || !e.Born(day) {
		return
	}
	f.self.expecting = nil
	other := f.self.genes
	if s := f.self.spouse; s != nil {
		other = s.who.Genes
	}
	v := f.baby(birth{parents: e.Parents, genes: [2][5]int{f.self.genes, other}, day: day})
	v.ours = true
	f.people = append(f.people, v)
	if f.room != nil && f.room.own && f.player.Game() != nil {
		f.player.Game().Add(f.resident(v, homeSpot(len(f.people))))
	}
	f.notify("Nasceu " + pick(v, "sua filha", "seu filho") + ", " + v.name + "!")
	f.sound("sell", nil)
}
