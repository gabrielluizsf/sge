package installer

import "testing"

func TestAnIDOrNameThatCouldReachTheSystemIsRefused(t *testing.T) {
	for _, id := range []string{"", ".farm", "farm.", "a/b", `a\b`, "../x", "a b", "CON", "nul.txt", "com1", "LPT¹", "fazenda\n", "fäzenda"} {
		if validID(id) {
			t.Errorf("the ID %q is taken", id)
		}
	}
	for _, id := range []string{"farm", "com.example.farm", "my_game-2"} {
		if !validID(id) {
			t.Errorf("the ID %q is refused", id)
		}
	}
	for _, name := range []string{"", " Farm", "Farm ", "Farm\nExec=sh", "Farm\x00", "a/b", `a\b`, "a:b", "100%", "Farm.", "PRN"} {
		if validName(name) {
			t.Errorf("the name %q is taken", name)
		}
	}
	for _, name := range []string{"Farm", "Fazenda Feliz", "Café d'Olé"} {
		if !validName(name) {
			t.Errorf("the name %q is refused", name)
		}
	}
}

func TestAFileNameThatLeavesTheFolderIsRefused(t *testing.T) {
	for _, name := range []string{"", ".", "..", "../x", "a/../../x", "/etc/passwd", `..\x`, `C:\x`, "C:x", "a:stream", "aux.png", "dir./x", "x ", "a\nb"} {
		if validFile(name) {
			t.Errorf("the file name %q is taken", name)
		}
	}
	for _, name := range []string{"farm", "assets/tree.png", "a/b/c.txt", ".hidden"} {
		if !validFile(name) {
			t.Errorf("the file name %q is refused", name)
		}
	}
}
