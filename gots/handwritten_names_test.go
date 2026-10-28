package gots

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

var exportedName = regexp.MustCompile(`(?m)^export (?:class|function|const|let|interface|type) ([A-Za-z_$][A-Za-z0-9_$]*)`)

// TestHandWrittenExportsAreMangled checks that hand-written standard
// packages export Go names the way the emitter refers to them: a Go name
// that is a JavaScript reserved word or global gets a $ suffix.
func TestHandWrittenExportsAreMangled(t *testing.T) {
	err := fs.WalkDir(runtimeFS, "runtime/std", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, "mod.ts") {
			return err
		}
		src, err := runtimeFS.ReadFile(p)
		if err != nil {
			return err
		}
		for _, m := range exportedName.FindAllStringSubmatch(string(src), -1) {
			if reserved[m[1]] {
				t.Errorf("%s exports %s, which emitted code calls %s$", p, m[1], m[1])
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
