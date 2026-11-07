package farm

import (
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/i18n"
)

// syncPackHint writes the inventory's hint for the tab or the chest shown.
func (f *farm) syncPackHint() {
	switch {
	case f.pack.Side != nil:
		f.pack.Hint = i18n.T("Enter: pick up   Shift+Enter: send across   Esc: close")
	case f.pack.Crafting():
		f.pack.Hint = i18n.T("Enter: craft   Q/E: tabs   Esc: close")
	default:
		f.pack.Hint = i18n.T("Enter: pick up   1-6: shortcut   Q/E: tabs   I/Tab/Esc: close")
	}
}

// afterPack answers what the inventory did: its sound, the hint of the tab
// it shows, and a notice for a piece made.
func (f *farm) afterPack(ev backpack.Event) {
	if name := packSound[ev]; name != "" {
		f.sound(name, nil)
	}
	f.syncPackHint()
	if ev == backpack.Crafted {
		f.notify(i18n.T("Made: %s", f.itemOf(f.pack.Made()).Name))
	}
}
