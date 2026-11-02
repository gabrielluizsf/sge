package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/menu"
)

// Update runs the menu on what the player does. The world behind goes on:
// the menu never pauses it.
func (m *mainMenu) Update(o *game.Object, _ float64) {
	if win := o.Game().Window(); win != nil && !m.done {
		m.handle(menuInput(win))
	}
}

// handle does what in says to the screen that is open.
func (m *mainMenu) handle(in menu.Input) {
	switch m.open {
	case titling:
		m.titleKeys(in)
	case confirming:
		m.confirmKeys(in)
	case crediting:
		m.creditsKeys(in)
	case licensing:
		if m.reader.Handle(in) == menu.Cancelled {
			m.open = crediting
			m.f.sound("ui_close", nil)
		}
	}
}

// titleKeys runs the main list. Escape does nothing: there is nowhere
// further back to go, and leaving is a line of its own.
func (m *mainMenu) titleKeys(in menu.Input) {
	switch m.main.Handle(in) {
	case menu.Moved:
		m.f.sound("ui_tick", nil)
	case menu.Chosen:
		switch m.main.Lines[m.main.Choice].Text {
		case lineStart:
			m.finish(Pick{Start: StartContinue})
		case lineNew:
			m.newGame()
		case lineCredits:
			m.open = crediting
			m.f.sound("ui_open", nil)
		case lineQuit:
			m.finish(Pick{})
		}
	}
}

// newGame asks first when a saved game would be left behind, then goes on.
func (m *mainMenu) newGame() {
	if m.hasSave {
		m.confirm.Choice = 1
		m.open = confirming
		m.f.sound("ui_open", nil)
		return
	}
	m.begin()
}

// begin starts a new game, through the art's step before it when it has
// one.
func (m *mainMenu) begin() {
	h, ok := m.f.art.(NewGameArt)
	if !ok {
		m.finish(Pick{Start: StartNew})
		return
	}
	m.open = creating
	h.NewGame(m.g, func(choice any) { m.finish(Pick{Start: StartNew, Choice: choice}) })
}

// confirmKeys runs the question asked before a new game replaces a save;
// Escape is "no".
func (m *mainMenu) confirmKeys(in menu.Input) {
	switch m.confirm.Handle(in) {
	case menu.Moved:
		m.f.sound("ui_tick", nil)
	case menu.Cancelled:
		m.open = titling
		m.f.sound("ui_close", nil)
	case menu.Chosen:
		if m.confirm.Lines[m.confirm.Choice].Text == lineBegin {
			m.begin()
			return
		}
		m.open = titling
		m.f.sound("ui_close", nil)
	}
}

// creditsKeys runs the credits: Enter or a click shows a licence, Escape
// goes back.
func (m *mainMenu) creditsKeys(in menu.Input) {
	switch m.credits.Handle(in) {
	case menu.Cancelled:
		m.open = titling
		m.f.sound("ui_close", nil)
	case menu.Moved:
		m.f.sound("ui_tick", nil)
	case menu.Chosen:
		c := m.list[m.credits.Choice]
		text := "Sem texto de licença."
		if c.License != nil {
			text = c.License()
		}
		m.reader.Open("Licença - "+c.Name, text)
		m.open = licensing
		m.f.sound("ui_open", nil)
	}
}
