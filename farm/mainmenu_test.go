package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/menu"
)

// titleOf is the main menu of a game titled title, config in a temp dir,
// with a save if saved.
func titleOf(t *testing.T, a Art, saved bool) (*Session, *mainMenu) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	const title = "sge menu test"
	if saved {
		if err := New(a, idle{}).Save(title); err != nil {
			t.Fatal(err)
		}
	}
	s := New(a, Bot(10))
	return s, s.attract(title)
}

// choose picks the line text of the open list.
func choose(t *testing.T, m *mainMenu, l *menu.List, text string) {
	t.Helper()
	for i, ln := range l.Lines {
		if ln.Text == text {
			l.Choice = i
			m.handle(menu.Input{Select: true})
			return
		}
	}
	t.Fatalf("no line %q in %v", text, l.Lines)
}

func linesOf(m *mainMenu) []string {
	var out []string
	for _, l := range m.main.Lines {
		out = append(out, l.Text)
	}
	return out
}

func TestIniciarIsOfferedOnlyWithASave(t *testing.T) {
	_, m := titleOf(t, SVGArt(), false)
	if got := linesOf(m); len(got) != 4 || got[0] != lineNew {
		t.Errorf("without a save: %v", got)
	}
	_, m = titleOf(t, SVGArt(), true)
	if got := linesOf(m); len(got) != 5 || got[0] != lineStart {
		t.Errorf("with a save: %v", got)
	}
}

func TestEachLineDoesItsThing(t *testing.T) {
	_, m := titleOf(t, SVGArt(), true)
	choose(t, m, &m.main, lineStart)
	if !m.done || m.pick.Start != StartContinue {
		t.Errorf("Iniciar: %+v", m.pick)
	}
	_, m = titleOf(t, SVGArt(), false)
	choose(t, m, &m.main, lineNew)
	if !m.done || m.pick.Start != StartNew {
		t.Errorf("Novo jogo: %+v", m.pick)
	}
	_, m = titleOf(t, SVGArt(), false)
	choose(t, m, &m.main, lineQuit)
	if !m.done || m.pick.Start != StartQuit {
		t.Errorf("Sair: %+v", m.pick)
	}
	_, m = titleOf(t, SVGArt(), false)
	choose(t, m, &m.main, lineCredits)
	if m.open != crediting {
		t.Fatal("Créditos did not open the credits")
	}
	m.handle(menu.Input{Select: true})
	if m.open != licensing {
		t.Fatal("no licence")
	}
	m.handle(menu.Input{Cancel: true})
	m.handle(menu.Input{Cancel: true})
	if m.open != titling || m.done {
		t.Error("Escape did not go back to the menu")
	}
}

func TestANewGameOverASaveAsksFirst(t *testing.T) {
	_, m := titleOf(t, SVGArt(), true)
	choose(t, m, &m.main, lineNew)
	if m.open != confirming || m.done {
		t.Fatalf("no question: open %v done %v", m.open, m.done)
	}
	m.handle(menu.Input{Cancel: true})
	if m.open != titling || m.done {
		t.Fatal("Escape did not decline")
	}
	choose(t, m, &m.main, lineNew)
	choose(t, m, &m.confirm, lineBack)
	if m.open != titling || m.done {
		t.Fatal("Voltar did not decline")
	}
	choose(t, m, &m.main, lineNew)
	choose(t, m, &m.confirm, lineBegin)
	if !m.done || m.pick.Start != StartNew {
		t.Errorf("confirmed: %+v", m.pick)
	}
}

func TestEscapeOnTheMainMenuDoesNothing(t *testing.T) {
	_, m := titleOf(t, SVGArt(), true)
	m.handle(menu.Input{Cancel: true})
	if m.done || m.open != titling {
		t.Error("Escape did something")
	}
}

func TestTheWorldGoesOnBehindTheMenu(t *testing.T) {
	s, m := titleOf(t, SVGArt(), false)
	s.Game.Step(frame)
	h, at := s.f.clock.Hour, s.f.player.Transform.Position
	for range 120 {
		s.Game.Step(frame)
	}
	if s.Game.Paused || s.f.clock.Hour == h || s.f.player.Transform.Position == at {
		t.Errorf("the world stood still: paused %v hour %v->%v", s.Game.Paused, h, s.f.clock.Hour)
	}
	cv, err := canvas.NewCanvas(1280, 720)
	if err != nil {
		t.Fatal(err)
	}
	s.Game.Draw(cv)
	m.open = confirming
	s.Game.Draw(cv)
}

// creator is an art with a step before a new game.
type creator struct {
	Art
	done func(any)
}

func (c *creator) NewGame(_ *game.Game, done func(any)) { c.done = done }

func TestTheArtsStepRunsBeforeANewGame(t *testing.T) {
	c := &creator{Art: SVGArt()}
	_, m := titleOf(t, c, false)
	choose(t, m, &m.main, lineNew)
	if c.done == nil || m.done || m.open != creating {
		t.Fatal("the hook was not called, or the game started before it was done")
	}
	c.done("Ana")
	if !m.done || m.pick.Start != StartNew || m.pick.Choice != "Ana" {
		t.Errorf("pick %+v", m.pick)
	}
}

func TestIniciarNeverShowsTheArtsStep(t *testing.T) {
	c := &creator{Art: SVGArt()}
	_, m := titleOf(t, c, true)
	choose(t, m, &m.main, lineStart)
	if c.done != nil || !m.done || m.pick.Start != StartContinue || m.pick.Choice != nil {
		t.Errorf("Iniciar went through the step: %+v", m.pick)
	}
}

// Changing the language in the menu is kept on the computer: the next time
// the game starts it speaks the same language, and its lines change with
// it, meaning the same.
func TestTheLanguageChosenInTheMenuIsKept(t *testing.T) {
	_, m := titleOf(t, SVGArt(), false)
	t.Cleanup(func() { i18n.Set(i18n.English) })
	for i, l := range m.main.Lines {
		if l.Key == lineLanguage {
			m.main.Choice = i
		}
	}
	m.titleKeys(menu.Input{Select: true})
	if i18n.Get() != i18n.Portuguese {
		t.Fatalf("language %v", i18n.Get())
	}
	if got := linesOf(m)[0]; got != "Novo jogo" {
		t.Errorf("the menu did not change: %q", got)
	}
	if m.main.Lines[0].Key != lineNew {
		t.Errorf("the line lost its meaning: %q", m.main.Lines[0].Key)
	}
	i18n.Set(i18n.English) // the game starts again
	i18n.Load()
	if i18n.Get() != i18n.Portuguese {
		t.Errorf("the choice was not kept: %v", i18n.Get())
	}
}
