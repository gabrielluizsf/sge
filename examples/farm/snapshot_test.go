package farm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/gabrielluizsf/sge/save"
)

// played is a game the bot has played for a while, with a town woken.
func played(t *testing.T) *farm {
	t.Helper()
	g, f := newGame(Bot(6))
	for range 60 * 90 {
		g.Step(frame)
	}
	for c := -1; c <= 1; c++ {
		if tn := f.towns.Of(c, 0); tn != nil {
			f.settle(tn)
		}
	}
	f.bag.Add(wood, 3)
	f.errandDone = true
	if len(f.people) > 0 {
		f.people[0].bond.Affection = 400
		f.self.spouse = f.people[0]
		f.people[0].ours = true
	}
	return f
}

var playedOnce struct {
	sync.Once
	json []byte
}

// playedJSON is the JSON of a played game's snapshot, made once.
func playedJSON(t *testing.T) []byte {
	playedOnce.Do(func() { playedOnce.json, _ = json.Marshal(takeSnapshot(played(t))) })
	return playedOnce.json
}

func TestSnapshotRoundTrip(t *testing.T) {
	f := played(t)
	if len(f.people) == 0 {
		t.Skip("no town at region 0, 0")
	}
	a := takeSnapshot(f)
	_, g2 := newGame(idle{})
	if err := restore(g2, &a); err != nil {
		t.Fatal(err)
	}
	b := takeSnapshot(g2)
	if !reflect.DeepEqual(a, b) {
		va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
		for i := range va.NumField() {
			if !reflect.DeepEqual(va.Field(i).Interface(), vb.Field(i).Interface()) {
				t.Errorf("%s differs:\n%+v\n%+v", va.Type().Field(i).Name, va.Field(i).Interface(), vb.Field(i).Interface())
			}
		}
		t.FailNow()
	}
	if g2.money != f.money || g2.self.spouse == nil || g2.bag.Count(wood) != f.bag.Count(wood) {
		t.Error("state differs")
	}
}

func TestSaveAndLoadAFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	s := New(SVGArt(), Bot(4))
	for range 600 {
		s.Game.Step(frame)
	}
	if HasSave("Test Farm") {
		t.Fatal("a save before saving")
	}
	s.EnableSave("Test Farm")
	s.f.ui.save()
	if !HasSave("Test Farm") {
		t.Fatal("no save after Salvar")
	}
	n := New(SVGArt(), idle{})
	if err := n.Load("Test Farm"); err != nil || n.f.money != s.f.money || n.f.clock.Hour != s.f.clock.Hour {
		t.Fatalf("load: %v", err)
	}
}

// hostile applies edit to a good snapshot's JSON, as a hostile save would,
// and restores it: it has to be an error, and the game still whole.
func hostile(t *testing.T, name string, edit func(*snapshot)) {
	t.Helper()
	b := playedJSON(t)
	var c snapshot
	json.Unmarshal(b, &c)
	edit(&c)
	_, g := newGame(idle{})
	before := takeSnapshot(g)
	if err := restore(g, &c); err == nil {
		t.Errorf("%s: accepted", name)
		return
	}
	if !reflect.DeepEqual(before, takeSnapshot(g)) {
		t.Errorf("%s: the game changed though the save was refused", name)
	}
}

func TestHostileSnapshotsAreRefused(t *testing.T) {
	cases := map[string]func(*snapshot){
		"negative money":     func(s *snapshot) { s.Player.Money = -1 },
		"huge money":         func(s *snapshot) { s.Player.Money = 1 << 60 },
		"negative count":     func(s *snapshot) { s.Player.Bag = []stackState{{Item: wood, N: -4}} },
		"huge count":         func(s *snapshot) { s.Player.Bag = []stackState{{Item: wood, N: 1 << 40}} },
		"unknown item":       func(s *snapshot) { s.Player.Bag = []stackState{{Item: "gold", N: 1}} },
		"bad axe":            func(s *snapshot) { s.Player.Axe = 99 },
		"bad seed":           func(s *snapshot) { s.Player.Seed = -3 },
		"bad slot":           func(s *snapshot) { s.Player.Slot = 9 },
		"day negative":       func(s *snapshot) { s.Day = -1 },
		"day huge":           func(s *snapshot) { s.Day = 1 << 50 },
		"hour":               func(s *snapshot) { s.Hour = 99 },
		"seen ahead":         func(s *snapshot) { s.Seen = s.Day + 5 },
		"hotbar short":       func(s *snapshot) { s.Hotbar = nil },
		"hotbar item":        func(s *snapshot) { s.Hotbar[2] = "nope" },
		"active":             func(s *snapshot) { s.Active = 40 },
		"off the map":        func(s *snapshot) { s.At = point{X: 1e30, Y: 0} },
		"player name":        func(s *snapshot) { s.Self.Name = strings.Repeat("x", 5000) },
		"player id":          func(s *snapshot) { s.Self.ID = "nobody" },
		"spouse index":       func(s *snapshot) { s.Self.Spouse = 9999 },
		"genes":              func(s *snapshot) { s.Self.Genes[0] = 99 },
		"plot crop":          func(s *snapshot) { s.Plots[0].Crop = 77 },
		"plot stage":         func(s *snapshot) { s.Plots[0].Stage = -2 },
		"plot cell":          func(s *snapshot) { s.Plots[0].Col = 1 << 40 },
		"plot elsewhere":     func(s *snapshot) { s.Plots[0].Col, s.Plots[0].Row = 0, 0 },
		"plot twice":         func(s *snapshot) { s.Plots = append(s.Plots, s.Plots[0]) },
		"plot unwoken town":  func(s *snapshot) { s.Plots[0].Town = &cellState{Col: 5000, Row: 5000} },
		"stump no tree":      func(s *snapshot) { s.Stumps = []stumpState{{Col: 0, Row: 0, Tree: 0, Left: 1}} },
		"hit no tree":        func(s *snapshot) { s.Hits = []hitState{{Col: 0, Row: 0, N: 1}} },
		"town none":          func(s *snapshot) { s.Towns = append(s.Towns, cellState{Col: 1 << 30}) },
		"errand item":        func(s *snapshot) { s.Errand.Want = []stackState{{Item: "x", N: 1}} },
		"errand reward":      func(s *snapshot) { s.Errand.Reward = -1 },
		"family parent":      func(s *snapshot) { s.Family[len(s.Family)-1].Parents = [2]string{"ghost", ""} },
		"family twice":       func(s *snapshot) { s.Family = append(s.Family, s.Family[0]) },
		"villager role":      func(s *snapshot) { s.People[0].Role = "king" },
		"villager home":      func(s *snapshot) { s.People[0].Home = 500 },
		"villager town":      func(s *snapshot) { s.People[0].Town = &cellState{Col: 4000, Row: 4000} },
		"villager partner":   func(s *snapshot) { s.People[0].Partner = 9999 },
		"villager one-sided": func(s *snapshot) { s.People[0].Partner = 1; s.People[1].Partner = -1 },
		"villager affection": func(s *snapshot) { s.People[0].Bond.Affection = 1 << 30 },
		"villager courted":   func(s *snapshot) { s.People[0].Bond.Courted = 40 },
		"villager name":      func(s *snapshot) { s.People[0].Name = "a\x00b" },
		"villager place":     func(s *snapshot) { s.People[0].At.X = -1e300 },
		"villager id":        func(s *snapshot) { s.People[0].ID = "unlisted" },
		"villager order":     func(s *snapshot) { s.People[0].Order.Reward = 1 << 50 },
	}
	for name, edit := range cases {
		hostile(t, name, edit)
	}
}

func TestHostileFilesAreRefused(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	st, _ := save.Open("Evil", saveVersion)
	file := func(body string) {
		os.MkdirAll(st.Dir(), 0o700)
		os.WriteFile(filepath.Join(st.Dir(), saveName+".json"), []byte(body), 0o600)
	}
	for name, body := range map[string]string{
		"deep":     strings.Repeat("[", 100000),
		"unknown":  `{"version":1,"data":{"evil":1}}`,
		"trailing": `{"version":1,"data":{}} {}`,
		"version":  `{"version":9,"data":{}}`,
		"big":      `{"version":1,"data":{"hotbar":["` + strings.Repeat("a", maxSaveBytes) + `"]}}`,
	} {
		file(body)
		if err := New(SVGArt(), idle{}).Load("Evil"); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestAnEmptySnapshotIsRefused(t *testing.T) {
	_, g := newGame(idle{})
	if restore(g, &snapshot{}) == nil {
		t.Error("an empty save was taken")
	}
}

// A villager's ownership of a plot survives a save, and a hostile owner
// is refused.
func TestPlotOwnersRoundTripAndHostileOnesAreRefused(t *testing.T) {
	f := played(t)
	if len(f.people) < 2 {
		t.Skip("no town at region 0, 0")
	}
	own := f.people[len(f.people)-1]
	p := f.plots[0]
	p.crop, p.owner = cornCrop, own
	a := takeSnapshot(f)
	_, g := newGame(idle{})
	if err := restore(g, &a); err != nil {
		t.Fatal(err)
	}
	if g.plots[0].owner == nil || g.plots[0].owner.id != own.id {
		t.Fatalf("owner lost: %+v", g.plots[0].owner)
	}
	for name, owner := range map[string]int{"past the people": len(a.People) + 1, "negative": -1} {
		bad := takeSnapshot(f)
		bad.Plots[0].Owner = owner
		_, g2 := newGame(idle{})
		if err := restore(g2, &bad); err == nil {
			t.Errorf("%s owner accepted", name)
		}
	}
	bare := takeSnapshot(f)
	last := &bare.Plots[len(bare.Plots)-1]
	last.Crop, last.Owner = -1, 1
	_, g3 := newGame(idle{})
	if err := restore(g3, &bare); err == nil {
		t.Error("an owner of bare soil accepted")
	}
}
