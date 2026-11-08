package i18n

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// onThisPC gives the test a config folder of its own, and the default
// language back after it.
func onThisPC(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	Set(English)
	t.Cleanup(func() { Set(English) })
	return filepath.Join(dir, choiceFolder, choiceName+".json")
}

// The language chosen is still the one chosen the next time a game starts.
func TestTheChoiceIsKeptOnThePC(t *testing.T) {
	onThisPC(t)
	if err := Change(1); err != nil {
		t.Fatal(err)
	}
	if Get() != Portuguese {
		t.Fatalf("changed to %v", Get())
	}
	Set(English) // a new run of the game
	Load()
	if Get() != Portuguese {
		t.Errorf("the choice was lost: %v", Get())
	}
	if err := Change(1); err != nil {
		t.Fatal(err)
	}
	Set(Portuguese)
	Load()
	if Get() != English {
		t.Errorf("the second choice was lost: %v", Get())
	}
}

// With nothing chosen yet the language is the default.
func TestNoChoiceIsEnglish(t *testing.T) {
	onThisPC(t)
	Load()
	if Get() != English {
		t.Errorf("language %v", Get())
	}
}

// A file that is damaged or made to harm leaves the language as it is,
// and the game starts all the same.
func TestABadFileIsIgnored(t *testing.T) {
	bad := map[string]string{
		"not JSON":         "pt",
		"empty":            "",
		"unknown language": `{"version":1,"data":{"language":"xx"}}`,
		"a number":         `{"version":1,"data":{"language":1}}`,
		"unknown field":    `{"version":1,"data":{"language":"pt","run":"x"}}`,
		"trailing data":    `{"version":1,"data":{"language":"pt"}} {}`,
		"another version":  `{"version":2,"data":{"language":"pt"}}`,
		"deep":             strings.Repeat("[", 100000),
		"too large":        `{"version":1,"data":{"language":"pt"}}` + strings.Repeat(" ", choiceMost),
	}
	// A good file on the same path is read: the bad ones are refused for
	// what they hold, not for not being found.
	path := onThisPC(t)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"version":1,"data":{"language":"pt"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	Load()
	if Get() != Portuguese {
		t.Fatalf("a good file was not read: %v", Get())
	}
	for name, text := range bad {
		path := onThisPC(t)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
		Load()
		if Get() != English {
			t.Errorf("%s: language %v", name, Get())
		}
	}
}

// A choice that cannot be written says so, and the language still changes
// for this run.
func TestAChoiceThatCannotBeKeptSaysSo(t *testing.T) {
	path := onThisPC(t)
	if err := os.WriteFile(filepath.Dir(filepath.Dir(path))+"/"+choiceFolder, []byte("a file, not a folder"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Change(1); err == nil {
		t.Error("no error")
	}
	if Get() != Portuguese {
		t.Errorf("language %v", Get())
	}
}
