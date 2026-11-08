package installer

import (
	"strings"
	"testing"
)

// A path with the characters the spec reserves is escaped so the entry
// still starts that one program, and none adds a key of its own.
func TestADesktopEntryQuotesTheProgram(t *testing.T) {
	got := desktopEntry(App{Name: `Farm \ 2`, Game: true}, `/home/a "b"/$x/100%/f\g`, "/d", "/d/i.png")
	for _, want := range []string{
		"Name=Farm \\\\ 2\n",
		"Exec=\"/home/a \\\\\"b\\\\\"/\\\\$x/100%%/f\\\\\\\\g\"\n",
		"Categories=Game;\n",
		"Icon=/d/i.png\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the entry has no %q:\n%s", want, got)
		}
	}
	if strings.Count(got, "\n") != 9 {
		t.Errorf("the entry is not nine lines:\n%s", got)
	}
}
