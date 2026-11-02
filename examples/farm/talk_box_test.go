package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/dialogue"
)

// said opens a conversation with v, does what it is asked, and is what
// the box says after; the screen's notices and floating words are empty,
// as everything a villager says is in the box.
func said(t *testing.T, f *farm, do func()) string {
	t.Helper()
	notices, floats := len(f.notices), len(f.floats)
	do()
	if len(f.notices) != notices || len(f.floats) != floats {
		t.Errorf("a villager spoke outside the box: %d notices, %d floats", len(f.notices), len(f.floats))
	}
	return f.talk.box.Text()
}

// Every kind of line a villager has goes through the box: the greeting
// with what they want, a gift's reaction, a chat, an order not filled and
// filled, wood and axe trades, the smith's errand, and the answers to
// being asked out.
func TestEveryVillagerLineIsSaidInTheBox(t *testing.T) {
	_, f := newGame(idle{})
	tn := eastTown(f)
	smith := tradesman(f, tn, roleSmith)
	lj := tradesman(f, tn, roleLumberjack)
	f.openTalk(smith)
	if g := f.talk.box.Text(); g == "" || f.talk.box.Speaker == "" {
		t.Fatalf("greeting %q by %q", g, f.talk.box.Speaker)
	}
	lines := map[string]func(){
		"errand not filled": func() { f.doErrand(smith) },
		"axe, no money":     func() { f.errandDone = true; f.money = 0; f.buyAxe(1, smith) },
		"axe bought":        func() { f.money = 500; f.buyAxe(1, smith) },
		"axe owned":         func() { f.buyAxe(1, smith) },
		"gift":              func() { f.bag.Add("carrot", 1); f.gift(smith, "carrot") },
		"second gift":       func() { f.bag.Add("carrot", 1); f.gift(smith, "carrot") },
		"chat":              func() { f.chat(smith) },
		"chat again":        func() { f.chat(smith) },
		"asked out, no":     func() { f.askOut(smith) },
		"wood, none":        func() { f.tradeWood(deal{lj.kit, f.kit, lumberSells, lj}) },
		"wood bought":       func() { lj.kit.bag.Add(wood, 1); f.money = 50; f.tradeWood(deal{lj.kit, f.kit, lumberSells, lj}) },
		"order not filled":  func() { lj.order.Want, lj.order.Reward = f.errand.Want, 9; f.deliver(lj) },
	}
	for name, do := range lines {
		f.talk.box.Say("")
		if got := said(t, f, do); got == "" {
			t.Errorf("%s: nothing said", name)
		}
	}
	f.talk = nil
	said(t, f, func() { f.speak(smith, "Sim! Vamos namorar!") })
}

// A long speech is in pages; a key shows a page whole and the next goes on,
// and the buttons are the lines the player may say.
func TestTheBoxPagesAndOffersTheLines(t *testing.T) {
	_, f := newGame(idle{})
	v := tradesman(f, eastTown(f), roleLumberjack)
	f.openTalk(v)
	f.speak(v, "Era uma vez uma floresta muito grande, cheia de árvores altas, onde ninguém ia cortar nada sem pedir licença à velha corujinha que morava no tronco mais antigo, e todos os dias os lenhadores subiam a colina para contar histórias e comparar machados e trocar madeira por moedas.")
	f.talkUpdate(0, dialogue.Input{})
	if len(f.talk.box.Choices) != len(f.lines(v)) || len(f.talk.box.Choices) < 3 {
		t.Fatalf("choices %v", f.talk.box.Choices)
	}
	if f.talk.box.Ready() {
		t.Error("ready before the words are said")
	}
	f.talkUpdate(0, dialogue.Input{Cancel: true})
	if f.talk != nil {
		t.Error("escape did not end the conversation")
	}
}
