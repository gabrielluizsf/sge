package farm

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// While the bag or the pause menu is open the day stands still, and it
// goes on when they close.
func TestTimeStandsStillInTheMenus(t *testing.T) {
	g, f := newGame(idle{})
	g.Step(frame)
	for _, hold := range []func(){
		func() { f.pack.Show() }, func() { f.ui.open = pausing },
	} {
		hold()
		g.Step(frame)
		h := f.clock.Hour
		for range 30 {
			g.Step(frame)
		}
		if f.clock.Hour != h || !g.Paused {
			t.Errorf("hour %v went to %v while open", h, f.clock.Hour)
		}
		f.pack.Hide()
		f.ui.open = playing
		g.Step(frame)
		g.Step(frame)
		if g.Paused || f.clock.Hour == h {
			t.Error("time did not go on")
		}
	}
}

// The pause menu is Continuar, Créditos and Sair, with room for a Salvar
// once a game can save; Credits show each component with a licence text.
func TestThePauseMenuAndTheCredits(t *testing.T) {
	_, f := newGame(idle{})
	var lines []string
	for _, l := range f.ui.menuLines() {
		lines = append(lines, l.Text)
	}
	if strings.Join(lines, ",") != "Continuar,Créditos,Sair" {
		t.Errorf("lines %v", lines)
	}
	f.ui.save = func() {}
	if len(f.ui.menuLines()) != 4 {
		t.Error("no Salvar once there is something to save")
	}
	if len(f.ui.list) == 0 {
		t.Fatal("no credits")
	}
	for _, c := range f.ui.list {
		if c.Name == "" || c.Author == "" || c.URL == "" || c.License == nil || !strings.Contains(c.License(), "License") {
			t.Errorf("credit %+v", c)
		}
	}
	cv, err := canvas.NewCanvas(1280, 720)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []screen{pausing, crediting, licensing} {
		f.ui.open = s
		f.ui.reader.Open("x", f.ui.list[0].License())
		f.ui.paint(cv)
	}
}
