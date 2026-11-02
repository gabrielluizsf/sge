package farm

import (
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

// The main menu's lines.
const (
	lineStart = "Iniciar"
	lineNew   = "Novo jogo"
	lineBegin = "Começar"
	lineBack  = "Voltar"
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
	if m.hasSave {
		m.main.Lines = append(m.main.Lines, menu.Line{Text: lineStart})
	}
	m.main.Lines = append(m.main.Lines, menu.Line{Text: lineNew}, menu.Line{Text: lineCredits}, menu.Line{Text: lineQuit})
	m.confirm = menu.List{
		Title: "Começar um novo jogo?", Scale: scale, MinW: 200,
		Hint:  "O jogo salvo será substituído quando você salvar.",
		Lines: []menu.Line{{Text: lineBegin}, {Text: lineBack}},
	}
	m.credits = menu.List{Title: "Créditos", Scale: scale, MinW: 260}
	m.reader.Scale = scale
	for _, c := range m.list {
		m.credits.Lines = append(m.credits.Lines, menu.Line{Text: c.Name + " - " + c.Author, Sub: c.URL, Button: "[Visualizar licença]"})
	}
	return m
}

// finish ends the menu in p and ends the game it runs in.
func (m *mainMenu) finish(p Pick) {
	m.pick, m.done = p, true
	m.g.Quit()
}
