package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/menu"
)

// Start is how the main menu ended.
type Start int

// The ways the main menu ends: the player quit (or closed the window),
// asked to go on with the saved game, or to begin a new one.
const (
	StartQuit Start = iota
	StartContinue
	StartNew
)

// Pick is what the main menu came to: how, and what the art's NewGame
// step chose, if it has one.
type Pick struct {
	Start  Start
	Choice any
}

// The main menu's screens, besides the credits and licences of the pause
// menu's.
const (
	titling screen = iota + 100
	confirming
	creating
)

// The main menu's lines, by the key that names them whatever the
// language.
const (
	lineStart    = "Start"
	lineNew      = "New game"
	lineBegin    = "Begin"
	lineBack     = "Back"
	lineLanguage = "Language"
)

// mainMenu is the menu shown when the game starts: the title and its
// choices over a farm that plays by itself, the world going on behind it.
type mainMenu struct {
	f       *farm
	g       *game.Game
	hasSave bool
	open    screen
	main    menu.List
	confirm menu.List
	credits menu.List
	reader  menu.Reader
	list    []Credit
	// pick is what the menu ended in; done is set once it has.
	pick Pick
	done bool
}

// newMainMenu is the main menu of the game titled title, over f in g.
func newMainMenu(f *farm, g *game.Game, title string) *mainMenu {
	m := &mainMenu{f: f, g: g, open: titling, hasSave: HasSave(title), list: creditsOf(f.art)}
	m.main = menu.List{Title: title, Scale: scale, MinW: 130}
	m.confirm = menu.List{Scale: scale, MinW: 200}
	m.credits = menu.List{Scale: scale, MinW: 260}
	m.reader.Scale = scale
	m.refresh()
	return m
}

// refresh rewrites the menus in the language chosen, keeping the line
// chosen.
func (m *mainMenu) refresh() {
	m.main.Lines = m.main.Lines[:0]
	if m.hasSave {
		m.main.Lines = append(m.main.Lines, menu.Line{Key: lineStart, Text: i18n.T(lineStart)})
	}
	m.main.Lines = append(m.main.Lines,
		menu.Line{Key: lineNew, Text: i18n.T(lineNew)},
		menu.Line{Key: lineCredits, Text: i18n.T(lineCredits)},
		menu.Line{Key: lineLanguage, Text: i18n.T(lineLanguage) + ": " + i18n.Name()},
		menu.Line{Key: lineQuit, Text: i18n.T(lineQuit)},
	)
	m.confirm.Title = i18n.T("Start a new game?")
	m.confirm.Hint = i18n.T("The saved game will be replaced when you save.")
	m.confirm.Lines = []menu.Line{
		{Key: lineBegin, Text: i18n.T(lineBegin)},
		{Key: lineBack, Text: i18n.T(lineBack)},
	}
	m.credits.Title = i18n.T(lineCredits)
	m.credits.Lines = m.credits.Lines[:0]
	for _, c := range m.list {
		m.credits.Lines = append(m.credits.Lines, menu.Line{Text: i18n.T(c.Name) + " - " + i18n.T(c.Author), Sub: c.URL, Button: i18n.T("[View license]")})
	}
	m.f.pack.Title = i18n.T("Backpack")
	m.f.pack.Hint = i18n.T("Enter: pick up   1-6: shortcut   Drag: mouse   I/Tab/Esc: close")
}

// finish ends the menu in p and ends the game it runs in.
func (m *mainMenu) finish(p Pick) {
	m.pick, m.done = p, true
	m.g.Quit()
}
