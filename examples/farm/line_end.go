package farm

// endLine is the player old dying with no heir: the line ends, and the
// player is asked only to start again.
func (f *farm) endLine(old *self) {
	f.talk, f.shop.open = nil, false
	f.question = &question{ask: []string{"Fim da linhagem", old.name + " morreu sem deixar herdeiros."},
		answers: []line{{text: "Recomeçar", do: func(f *farm) { f.restart() }}}}
	f.sound("deny", nil)
}

// restart starts the farm again after a line ends: a new farmer of
// founderYears comes to it with the money a game starts with, an empty
// bag and no axe, its field bare and the smith's errand to do anew; the
// towns live on as they were, and whoever lived with the player goes
// back to their own house.
func (f *farm) restart() {
	f.self = f.found()
	*f.kit = kit{money: 20, axe: noAxe}
	f.errand, f.errandDone = errandBook.Draw(f.luck), false
	for _, p := range f.plots {
		if f.onFarm(p.col, p.row) {
			p.crop, p.stage, p.watered, p.growing = -1, 0, false, 0
			f.show(p)
		}
	}
	for _, v := range f.people {
		if v.ours && v.home != nil {
			v.ours, v.inside = false, true
		}
	}
	f.forgetBonds()
	if f.room != nil {
		f.leave() // to lay the house out again without them
	}
	f.goHome()
	f.redressPlayer(f.self.who(f.clock.Day))
	f.notify("Uma nova família chega à fazenda.")
}
