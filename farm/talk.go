package farm

import (
	"github.com/gabrielluizsf/sge/dialogue"
	"github.com/gabrielluizsf/sge/face"
	"github.com/gabrielluizsf/sge/i18n"
)

// talk is a conversation with a villager: the choice on its menu.
type talk struct {
	v      *villager
	choice int
	// box is the dialogue box the conversation is shown in.
	box *dialogue.Box
	// blink and mouth are the speaker's face, alive while they talk.
	blink *face.Blinker
	mouth face.Speaker
}

// line is a choice in a conversation: its words, and what choosing it
// does.
type line struct {
	text string
	do   func(f *farm)
}

// The prices a lumberjack trades wood at with the player: they sell dear
// and buy cheap.
const (
	lumberSells = woodPrice + 1
	lumberBuys  = woodPrice - 1
)

// lines are what the player may say to v: deliver what they ask for,
// trade wood with a lumberjack, do the smith's errand or buy their axes,
// chat, give a gift and court them, and say goodbye; to the player's own
// family, only chat and gifts.
func (f *farm) lines(v *villager) []line {
	if v.ours {
		return append(f.bondLines(v), byeLine())
	}
	var ls []line
	if !v.smith() && len(v.order.Want) > 0 {
		ls = append(ls, line{text: i18n.T("Deliver order"), do: func(f *farm) { f.deliver(v) }})
	}
	switch v.role {
	case roleLumberjack:
		ls = append(ls,
			line{text: i18n.T("Buy wood $%d", lumberSells), do: func(f *farm) { f.tradeWood(deal{v.kit, f.kit, lumberSells, v}) }},
			line{text: i18n.T("Sell wood $%d", lumberBuys), do: func(f *farm) { f.tradeWood(deal{f.kit, v.kit, lumberBuys, v}) }})
	case roleSmith:
		if !f.errandDone {
			ls = append(ls, line{text: i18n.T("Deliver errand"), do: func(f *farm) { f.doErrand(v) }})
			break
		}
		for i := 1; i < len(axes); i++ {
			ls = append(ls, line{text: i18n.T("%s axe $%d", i18n.T(axes[i].label), axes[i].price), do: func(f *farm) { f.buyAxe(i, v) }})
		}
		ls = append(ls, line{text: i18n.T("Sell firewood $%d", woodPrice), do: func(f *farm) { f.tradeWood(deal{f.kit, v.kit, woodPrice, v}) }})
	}
	ls = append(ls, f.bondLines(v)...)
	return append(ls, byeLine())
}

// byeLine ends the conversation.
func byeLine() line {
	return line{text: i18n.T("Bye"), do: func(f *farm) { f.talk = nil; f.sound("ui_close", nil) }}
}
