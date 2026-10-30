package farm

import (
	"errors"

	"github.com/gabrielluizsf/sge/bond"
)

// errTaken is asking out one who lives as a couple already, or asking
// while the player is courting someone else or married.
var errTaken = errors.New("farm: one of them is taken")

// chat is talking with v, fonder a little, once a day.
func (f *farm) chat(v *villager) error {
	if err := v.bond.Chat(f.clock.Day); err != nil {
		f.float(f.spot(v), "Já conversamos hoje!", 'r')
		return err
	}
	f.float(f.spot(v), "Que bom conversar com você!", 'g')
	f.sound("ui_tick", nil)
	return nil
}

// moment is now, as the player asks v: the day, whether both are grown
// up, and whether they are kin.
func (f *farm) moment(v *villager) bond.Moment {
	day := f.clock.Day
	return bond.Moment{Day: day, Grown: grownOn(f.self.life, day) && grownOn(v.life, day), Kin: f.family.Kin(f.self.id, v.id)}
}

// courting is who the player is dating, engaged or married to, if anyone
// but v.
func (f *farm) courting(v *villager) bool {
	if f.self.spouse != nil {
		return true
	}
	for _, p := range f.people {
		if p != v && p.bond.Courted >= bond.Dating {
			return true
		}
	}
	return false
}

// askOut asks v to date the player: only one fond enough, grown up like
// the player, not kin, and free, as the player must be. It is why not,
// if not.
func (f *farm) askOut(v *villager) error {
	err := errTaken
	if v.partner == nil && !f.courting(v) {
		err = v.bond.AskOut(f.moment(v))
	}
	f.answer(v, err, "Sim! Vamos namorar!")
	return err
}

// propose asks v, dated and fond enough, to marry the player: they wed
// bond.WeddingAfter days on. It is why not, if not.
func (f *farm) propose(v *villager) error {
	err := v.bond.Propose(f.moment(v))
	f.answer(v, err, "Sim! O casamento é no dia "+itoa(v.bond.WeddingDay()+1)+"!")
	return err
}

// refusal is what one asked says when they will not, by why.
func refusal(err error) string {
	switch err {
	case bond.ErrYoung:
		return "Não... somos novos demais para isso."
	case bond.ErrKin:
		return "Não! Somos da mesma família!"
	case bond.ErrNotYet:
		return "Ainda não... vamos nos conhecer melhor."
	case bond.ErrAlready:
		return "Mas já estamos juntos!"
	}
	return "Desculpe, não posso."
}

// answer is v's answer, yes when err is nil, the reason they give when
// it is not.
func (f *farm) answer(v *villager, err error, yes string) {
	if err != nil {
		f.float(f.spot(v), refusal(err), 'r')
		f.sound("deny", nil)
		return
	}
	f.notify(v.name + ": " + yes)
	f.sound("sell", nil)
}
