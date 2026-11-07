package farm

import (
	"errors"

	"github.com/gabrielluizsf/sge/bond"
	"github.com/gabrielluizsf/sge/i18n"
)

// errTaken is asking out one who lives as a couple already, or asking
// while the player is courting someone else or married.
var errTaken = errors.New("farm: one of them is taken")

// chat is talking with v, fonder a little, once a day.
func (f *farm) chat(v *villager) error {
	if err := v.bond.Chat(f.clock.Day); err != nil {
		f.speak(v, i18n.T("We already talked today!"))
		return err
	}
	f.speak(v, i18n.T("It's good to talk to you!"))
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
	f.answer(v, err, i18n.T("Yes! Let's date!"))
	return err
}

// propose asks v, dated and fond enough, to marry the player: they wed
// bond.WeddingAfter days on. It is why not, if not.
func (f *farm) propose(v *villager) error {
	err := v.bond.Propose(f.moment(v))
	f.answer(v, err, i18n.T("Yes! The wedding is on day %d!", v.bond.WeddingDay()+1))
	return err
}

// refusal is what one asked says when they will not, by why.
func refusal(err error) string {
	switch err {
	case bond.ErrYoung:
		return i18n.T("No... we're too young for that.")
	case bond.ErrKin:
		return i18n.T("No! We're family!")
	case bond.ErrNotYet:
		return i18n.T("Not yet... let's get to know each other better.")
	case bond.ErrAlready:
		return i18n.T("But we're already together!")
	}
	return i18n.T("Sorry, I can't.")
}

// answer is v's answer, yes when err is nil, the reason they give when
// it is not.
func (f *farm) answer(v *villager, err error, yes string) {
	if err != nil {
		f.speak(v, refusal(err))
		f.sound("deny", nil)
		return
	}
	f.speak(v, yes)
	f.sound("sell", nil)
}
