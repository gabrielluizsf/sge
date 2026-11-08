package save

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type data struct {
	Name string `json:"name"`
	N    int    `json:"n"`
}

func TestRoundTrip(t *testing.T) {
	s := In(filepath.Join(t.TempDir(), "g"), 1)
	if s.Has("game") {
		t.Fatal("has a save before any")
	}
	if err := s.Write("game", data{"a", 3}); err != nil {
		t.Fatal(err)
	}
	var got data
	if err := s.Read("game", &got, 1<<10); err != nil || got != (data{"a", 3}) || !s.Has("game") {
		t.Fatalf("%v %+v", err, got)
	}
}

func TestModes(t *testing.T) {
	s := In(filepath.Join(t.TempDir(), "g"), 1)
	s.Write("game", data{})
	d, _ := os.Stat(s.Dir())
	f, _ := os.Stat(filepath.Join(s.Dir(), "game.json"))
	if d.Mode().Perm() != 0o700 || f.Mode().Perm() != 0o600 {
		t.Errorf("modes %v %v", d.Mode(), f.Mode())
	}
	es, _ := os.ReadDir(s.Dir())
	if len(es) != 1 {
		t.Errorf("files left: %v", es)
	}
}

func TestBadNames(t *testing.T) {
	s := In(t.TempDir(), 1)
	for _, n := range []string{"", "../x", "a/b", "..", ".", "a.b", "A", `a\b`, strings.Repeat("a", 65), "a\x00"} {
		if s.Write(n, data{}) == nil || s.Read(n, &data{}, 100) == nil || s.Has(n) {
			t.Errorf("accepted %q", n)
		}
	}
}

func TestHostileFiles(t *testing.T) {
	dir := t.TempDir()
	s := In(dir, 1)
	cases := map[string]string{
		"unknown":   `{"version":1,"data":{"name":"a","n":1,"x":2}}`,
		"envelope":  `{"version":1,"data":{},"x":1}`,
		"trailing":  `{"version":1,"data":{"name":"a","n":1}} {}`,
		"junk":      `{"version":1,"data":{"name":"a","n":1}}x`,
		"version":   `{"version":2,"data":{"name":"a","n":1}}`,
		"noversion": `{"data":{"name":"a","n":1}}`,
		"type":      `{"version":1,"data":{"name":5}}`,
		"deep":      strings.Repeat("[", 100000),
		"deepdata":  `{"version":1,"data":` + strings.Repeat("[", 100000) + `}`,
		"empty":     ``,
	}
	for name, body := range cases {
		os.WriteFile(filepath.Join(dir, "game.json"), []byte(body), 0o600)
		if err := s.Read("game", &data{}, 1<<20); err == nil {
			t.Errorf("%s accepted", name)
		}
	}
}

func TestTooLarge(t *testing.T) {
	s := In(t.TempDir(), 1)
	s.Write("game", data{Name: strings.Repeat("a", 500)})
	if s.Read("game", &data{}, 100) == nil {
		t.Error("read past the limit")
	}
	if s.Read("game", &data{}, 1000) != nil {
		t.Error("refused within the limit")
	}
}

func TestFailedWriteKeepsTheOldSave(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "g")
	s := In(dir, 1)
	s.Write("game", data{"old", 1})
	// An unencodable value fails before anything is touched.
	if s.Write("game", make(chan int)) == nil {
		t.Fatal("wrote a channel")
	}
	// A folder that cannot be written to fails at the temporary file.
	os.Chmod(dir, 0o500)
	defer os.Chmod(dir, 0o700)
	err := s.Write("game", data{"new", 2})
	os.Chmod(dir, 0o700)
	var got data
	if s.Read("game", &got, 1<<10) != nil || got.Name != "old" {
		t.Errorf("old save lost: %+v (write error %v)", got, err)
	}
	if es, _ := os.ReadDir(dir); len(es) != 1 {
		t.Errorf("leftovers: %v", es)
	}
}

func TestId(t *testing.T) {
	for in, want := range map[string]string{"Real Farm": "real-farm", "../../etc": "etc", "": "game", "///": "game", "Fazenda!!": "fazenda"} {
		if got := Id(in); got != want {
			t.Errorf("Id(%q) = %q, want %q", in, got, want)
		}
	}
	if len(Id(strings.Repeat("a", 500))) > maxName {
		t.Error("long id")
	}
}
