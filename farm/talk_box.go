package farm

import (
	"github.com/gabrielluizsf/sge/dialogue"
	"github.com/gabrielluizsf/sge/i18n"
)

// typing is how many letters a second the people speak.
const typing = 60

// openTalk begins a conversation with v: the box opens with what they say
// first.
func (f *farm) openTalk(v *villager) {
	f.talk = &talk{v: v}
	f.talkBox().Say(f.greeting(v))
	f.sound("ui_open", nil)
}

// talkBox is the open conversation's box, made the first time it is
// asked for.
func (f *farm) talkBox() *dialogue.Box {
	t := f.talk
	if t.box == nil {
		t.box = &dialogue.Box{Scale: scale, Speed: typing}
		t.box.Say(f.greeting(t.v))
	}
	t.box.Speaker = i18n.T("%s - %s, %d years old", t.v.name, f.labelOf(t.v), t.v.who.Years)
	return t.box
}

// speak is v saying text: in the conversation's box, which opens if it
// was not.
func (f *farm) speak(v *villager, text string) {
	if f.talk == nil || f.talk.v != v {
		f.talk = &talk{v: v}
	}
	f.talkBox().Say(text)
}

// greeting is what v says as the player comes to them: who they are to
// the player, what they want and what it pays, or what they sell.
func (f *farm) greeting(v *villager) string {
	var say []string
	switch {
	case v.ours:
		say = append(say, i18n.T("It's good to have you home!"))
	case v.smith() && !f.errandDone:
		say = append(say, i18n.T("I need this for the forge:"))
		say = append(say, wants(f, f.errand.Want)...)
		say = append(say, i18n.T("In exchange: stone axe"))
	case v.smith():
		say = append(say, i18n.T("Better axes cut faster."))
	case len(v.order.Want) > 0:
		say = append(say, i18n.T("You can bring me:"))
		say = append(say, wants(f, v.order.Want)...)
		say = append(say, i18n.T("I pay $%d", v.order.Reward))
	default:
		say = append(say, i18n.T("Hi! I don't need anything today."))
	}
	if v.role == roleLumberjack {
		say = append(say, i18n.T("I have %d wood", v.kit.bag.Count(wood)))
	}
	text := ""
	for i, s := range say {
		if i > 0 {
			text += "\n"
		}
		text += s
	}
	return text
}

// talkUpdate types the open conversation's words, keeps its buttons the
// lines the player may say now, and does what the player did in it.
func (f *farm) talkUpdate(dt float64, in dialogue.Input) {
	if f.talk == nil {
		return
	}
	b := f.talkBox()
	b.Update(dt)
	f.faceUpdate(dt)
	ls := f.lines(f.talk.v)
	b.Choices = b.Choices[:0]
	for _, l := range ls {
		b.Choices = append(b.Choices, l.text)
	}
	b.Choice = min(b.Choice, len(ls)-1)
	switch b.Handle(in) {
	case dialogue.Left:
		f.talk = nil
		f.sound("ui_close", nil)
	case dialogue.Chosen:
		ls[b.Choice].do(f)
	case dialogue.Moved, dialogue.Advanced:
		f.sound("ui_tick", nil)
	}
}
