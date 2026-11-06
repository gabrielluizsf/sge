package farm

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/dialogue"
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/menu"
)

// screen is which of the menus is open.
type screen int

// The screens: none, the pause menu, the credits, and a licence on them.
const (
	playing screen = iota
	pausing
	crediting
	licensing
	commanding
)

// The pause menu's lines, by the key that names them whatever the
// language.
const (
	lineContinue = "Continue"
	lineCommands = "Commands"
	lineCredits  = "Credits"
	lineQuit     = "Quit"
	lineSave     = "Save"
)

// ui is the farm's menus: the inventory, the pause menu, the credits and
// the licences. While one is open the world's time is held: cosy farm
// games stop the day, the crops and the people while the bag is open, so
// nothing is lost to sorting it, and so does a pause menu. Animations, the
// music and the ambience go on. It runs while the game is paused.
type ui struct {
	f       *farm
	open    screen
	pause   menu.List
	credits menu.List
	reader  menu.Reader
	list    []Credit
	// save, when set, is what the pause menu's Salvar button does; it has
	// no such button while it is nil.
	save func()
	// attract is set when the farm plays by itself behind the main menu:
	// the player has none of its keys.
	attract bool
}

// newUI is the farm's menus.
func newUI(f *farm) *ui {
	u := &ui{f: f, list: creditsOf(f.art)}
	u.pause = menu.List{Scale: scale, MinW: 130}
	u.credits = menu.List{Scale: scale, MinW: 260}
	u.reader.Scale = scale
	u.refresh()
	return u
}

// refresh rewrites the menus in the language chosen, keeping the line
// chosen.
func (u *ui) refresh() {
	u.pause.Title = i18n.T("Pause")
	u.pause.Lines = u.menuLines()
	u.credits.Title = i18n.T(lineCredits)
	u.credits.Lines = u.credits.Lines[:0]
	for _, c := range u.list {
		u.credits.Lines = append(u.credits.Lines, menu.Line{Text: i18n.T(c.Name) + " - " + i18n.T(c.Author), Sub: c.URL, Button: i18n.T("[View license]")})
	}
	u.f.pack.Title = i18n.T("Backpack")
	u.f.pack.Tabs = craftTabs()
	u.f.syncPackHint()
}

// menuLines are the pause menu's lines.
func (u *ui) menuLines() []menu.Line {
	ls := []menu.Line{{Key: lineContinue, Text: i18n.T(lineContinue)}, {Key: lineCommands, Text: i18n.T(lineCommands)}}
	if u.save != nil {
		ls = append(ls, menu.Line{Key: lineSave, Text: i18n.T(lineSave)})
	}
	return append(ls,
		menu.Line{Key: lineCredits, Text: i18n.T(lineCredits)},
		menu.Line{Key: lineLanguage, Text: i18n.T(lineLanguage) + ": " + i18n.Name()},
		menu.Line{Key: lineQuit, Text: i18n.T(lineQuit)},
	)
}

// Update runs the open menu, or opens one on its key, and holds the
// world's time while one is open.
func (u *ui) Update(o *game.Object, dt float64) {
	f := u.f
	g := o.Game()
	win := g.Window()
	var in dialogue.Input
	if win != nil && u.open == playing && !f.pack.Open && !u.attract {
		in = talkInput(win)
	}
	f.talkUpdate(dt, in)
	if win != nil && !u.attract {
		u.keys(win, g)
	}
	g.Paused = f.pack.Open || u.open != playing
}

// busy reports whether something else has the keys: a question, the shop
// or a conversation.
func (u *ui) busy() bool {
	f := u.f
	return f.question != nil || f.shop.open || f.talk != nil
}

// keys does what the player pressed.
func (u *ui) keys(win *antui.Window, g *game.Game) {
	f := u.f
	switch {
	case u.open == licensing:
		if u.reader.Handle(menuInput(win)) == menu.Cancelled {
			u.open = crediting
			f.sound("ui_close", nil)
		}
	case u.open == commanding:
		u.commandsInput(menuInput(win))
	case u.open == crediting:
		u.creditsKeys(win)
	case u.open == pausing:
		u.pauseKeys(win, g)
	case f.pack.Open:
		u.packKeys(win)
	case u.busy():
	case win.KeyPressed(antui.KeyEscape):
		u.pause.Lines, u.pause.Choice = u.menuLines(), 0
		u.open = pausing
		f.sound("ui_open", nil)
	case win.KeyPressed(antui.KeyI) || win.KeyPressed(antui.KeyTab):
		f.sound("ui_open", nil)
		f.pack.Show()
		f.syncPackHint()
	}
}

// creditsKeys runs the credits: Enter or a click on a line shows its
// licence, Escape goes back to the pause menu.
func (u *ui) creditsKeys(win *antui.Window) {
	f := u.f
	switch u.credits.Handle(menuInput(win)) {
	case menu.Cancelled:
		u.open = pausing
		f.sound("ui_close", nil)
	case menu.Chosen:
		c := u.list[u.credits.Choice]
		text := i18n.T("No license text.")
		if c.License != nil {
			text = c.License()
		}
		u.reader.Open(i18n.T("License - %s", c.Name), text)
		u.open = licensing
		f.sound("ui_open", nil)
	case menu.Moved:
		f.sound("ui_tick", nil)
	}
}

// pauseKeys runs the pause menu.
func (u *ui) pauseKeys(win *antui.Window, g *game.Game) {
	f := u.f
	switch u.pause.Handle(menuInput(win)) {
	case menu.Cancelled:
		u.open = playing
		f.sound("ui_close", nil)
	case menu.Moved:
		f.sound("ui_tick", nil)
	case menu.Chosen:
		switch u.pause.Lines[u.pause.Choice].Key {
		case lineContinue:
			u.open = playing
			f.sound("ui_close", nil)
		case lineCommands:
			u.showCommands()
		case lineCredits:
			u.open = crediting
			f.sound("ui_open", nil)
		case lineSave:
			u.save()
		case lineLanguage:
			f.changeLanguage()
			u.refresh()
		case lineQuit:
			g.Quit()
		}
	}
}

// packKeys runs the inventory: its keys and mouse, with I and Tab to
// close it too.
func (u *ui) packKeys(win *antui.Window) {
	f := u.f
	if win.KeyPressed(antui.KeyI) || win.KeyPressed(antui.KeyTab) {
		f.sound(packSound[f.pack.Hide()], nil)
		f.pack.Side = nil
		return
	}
	ev := f.pack.Handle(packInput(win))
	f.keepHand()
	f.afterPack(ev)
	if !f.pack.Open {
		f.pack.Side = nil
	}
}

// packSound is the sound each thing the inventory does makes.
var packSound = map[backpack.Event]string{
	backpack.Opened: "ui_open", backpack.Closed: "ui_close", backpack.Moved: "ui_tick", backpack.Picked: "ui_tick",
	backpack.Bound: "ui_open", backpack.Unbound: "ui_close", backpack.Dropped: "ui_close", backpack.Refused: "deny",
	backpack.Crafted: "ui_open", backpack.Stored: "ui_tick", backpack.Tabbed: "ui_tick",
}
