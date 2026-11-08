package installer

import (
	"bytes"
	"testing"
	"testing/fstest"
)

// An installing that fails lets go of its new folder before taking it
// away: Windows cannot remove a folder still open, and left it there.
func TestAFailedInstallingLetsGoOfItsFolder(t *testing.T) {
	s := hostSetup(t, App{ID: "farm", Name: "Farm", Exe: "farm"}, nil)
	s.Files, s.Limits = liar{fstest.MapFS{"farm": {Data: bytes.Repeat([]byte("x"), 100)}}}, Limits{Bytes: 50}
	w, err := plan(s)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.run(nil); err == nil {
		t.Fatal("a file larger than it said was installed")
	}
	if w.root != nil {
		t.Error("the new folder is still open")
	}
	if l := where(t, s); exists(l.stage()) {
		t.Error("the new folder was left behind")
	}
}
