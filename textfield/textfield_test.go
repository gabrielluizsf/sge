package textfield

import (
	"strings"
	"testing"
)

func TestInsertTypesLettersAccentsSpacesHyphensAndApostrophes(t *testing.T) {
	got := Insert("", "João d'Ávila-Souza", 20)
	if got != "João d'Ávila-Souza" {
		t.Errorf("got %q", got)
	}
}

func TestInsertRefusesWhatIsNotAPartOfAName(t *testing.T) {
	for _, bad := range []string{"\x00", "\n", "\t", "\x1b", "1", "<", "/", "\\", "..", "\u202e", "\u200b", "\u0301", "\xff"} {
		if got := Insert("Ana", bad, 20); got != "Ana" {
			t.Errorf("%q was taken: %q", bad, got)
		}
	}
	if got := Insert("", "A1b2", 20); got != "Ab" {
		t.Errorf("mixed typing gave %q", got)
	}
}

func TestInsertStopsAtTheLimitInCharactersNotBytes(t *testing.T) {
	got := Insert("", strings.Repeat("ã", 30), 20)
	if n := len([]rune(got)); n != 20 {
		t.Errorf("%d characters", n)
	}
	if Insert(got, "x", 20) != got {
		t.Error("typed past the limit")
	}
}

func TestBackspaceRemovesTheLastCharacterWhole(t *testing.T) {
	if got := Backspace("Fátima"); got != "Fátim" {
		t.Errorf("got %q", got)
	}
	if got := Backspace("Tupã"); got != "Tup" {
		t.Errorf("got %q", got)
	}
	if Backspace("") != "" {
		t.Error("backspace on nothing")
	}
}

func TestValidIsStrict(t *testing.T) {
	for _, ok := range []string{"Ana", "João", "Mary-Jane", "O'Neil", "Ana Maria", strings.Repeat("a", 20)} {
		if !Valid(ok, 20) {
			t.Errorf("%q refused", ok)
		}
	}
	for _, bad := range []string{"", " Ana", "Ana ", "  ", "Ana\n", "Ana\x00", "A\u0007", "<script>", "../x", strings.Repeat("a", 21), "\xff", "Ana1", "a‮b"} {
		if Valid(bad, 20) {
			t.Errorf("%q accepted", bad)
		}
	}
}
