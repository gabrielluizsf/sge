package farm

import (
	"strings"

	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/menu"
)

// commandsText is the controls, a line each, in the language chosen.
func commandsText() string {
	return strings.Join([]string{
		i18n.T("Arrows or WASD: walk"),
		i18n.T("Shift: run"),
		i18n.T("1-6: take an item in hand"),
		i18n.T("I or Tab: backpack"),
		i18n.T("Space or click: use the item"),
		i18n.T("F: attack (the axe hits harder)"),
		i18n.T("E: interact (talk, shop, doors, sit, chests)"),
		i18n.T("B: arrange furniture (in your house; chests outside)"),
		i18n.T("Backpack: Q/E change tab, Craft makes furniture"),
		i18n.T("Esc: menu"),
	}, "\n")
}

// showCommands opens the controls over the pause menu.
func (u *ui) showCommands() {
	u.reader.Open(i18n.T("Commands"), commandsText())
	u.open = commanding
	u.f.sound("ui_open", nil)
}

// commandsInput runs the controls screen: Escape goes back to the pause menu.
func (u *ui) commandsInput(in menu.Input) {
	if u.reader.Handle(in) == menu.Cancelled {
		u.open = pausing
		u.f.sound("ui_close", nil)
	}
}
