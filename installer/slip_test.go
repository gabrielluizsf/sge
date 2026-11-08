package installer

import (
	"archive/zip"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

// packed is an installer file of a program and files named as given,
// written by hand as a hostile one would be.
func packed(t *testing.T, names ...string) string {
	var b bytes.Buffer
	b.WriteString("PROGRAM")
	z := zip.NewWriter(&b)
	z.SetOffset(7)
	for _, n := range names {
		w, err := z.Create(n)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte("evil"))
	}
	z.SetComment(payloadTag + "7")
	z.Close()
	out := filepath.Join(t.TempDir(), "setup")
	os.WriteFile(out, b.Bytes(), 0o755)
	return out
}

// The zip slip: names that climb out of the game's folder, or start at
// a root, never write a file outside it.
func TestAFileClimbingOutOfTheFolderIsRefused(t *testing.T) {
	for _, bad := range []string{"../evil", "../../evil", "/tmp/evil", `..\evil`, "a/../../evil", "C:/evil"} {
		setup := packed(t, "farm", bad)
		pl, err := Open(setup)
		if err == nil {
			t.Errorf("opened a payload holding %q", bad)
		} else {
			continue
		}
		pl.Close()
	}
}

// Links in the files are never made, whether to a file outside or in.
func TestALinkInTheFilesIsRefused(t *testing.T) {
	files := fstest.MapFS{
		"farm":  {Data: []byte("x")},
		"evil":  {Data: []byte("/etc/passwd"), Mode: fs.ModeSymlink},
		"inner": {Data: []byte("farm"), Mode: fs.ModeSymlink},
	}
	s := hostSetup(t, App{ID: "farm", Name: "Farm", Exe: "farm"}, files)
	if err := Install(s); err == nil {
		t.Error("installed files holding links")
	}
	if l := where(t, s); exists(l.Dir) || exists(l.stage()) {
		t.Error("a folder was made before the files were checked")
	}
}

// Files that would fill the disk are refused, whether they say so or lie.
func TestFilesPastTheLimitsAreRefused(t *testing.T) {
	app := App{ID: "farm", Name: "Farm", Exe: "farm"}
	s := hostSetup(t, app, fstest.MapFS{"farm": {Data: []byte("x")}, "a": {Data: []byte("a")}, "b": {Data: []byte("b")}})
	s.Limits = Limits{Files: 2}
	if err := Install(s); err == nil {
		t.Error("installed more files than the limit")
	}
	big := fstest.MapFS{"farm": {Data: bytes.Repeat([]byte("x"), 100)}}
	s.Files, s.Limits = big, Limits{Bytes: 99}
	if err := Install(s); err == nil {
		t.Error("installed more bytes than the limit")
	}
	s.Files, s.Limits = liar{big}, Limits{Bytes: 50}
	if err := Install(s); err == nil {
		t.Error("installed a file larger than it said")
	}
	if l := where(t, s); exists(l.Dir) || exists(l.stage()) {
		t.Error("files refused were left installed")
	}
}

// liar is files whose sizes say 1 byte, whatever they hold.
type liar struct{ fstest.MapFS }

func (l liar) Open(name string) (fs.File, error) {
	f, err := l.MapFS.Open(name)
	return liarFile{f}, err
}

type liarFile struct{ fs.File }

func (f liarFile) Stat() (fs.FileInfo, error) {
	i, err := f.File.Stat()
	return liarInfo{i}, err
}

type liarInfo struct{ fs.FileInfo }

func (liarInfo) Size() int64 { return 1 }

func (l liar) ReadDir(name string) ([]fs.DirEntry, error) {
	es, err := l.MapFS.ReadDir(name)
	for i, e := range es {
		es[i] = liarEntry{e}
	}
	return es, err
}

type liarEntry struct{ fs.DirEntry }

func (e liarEntry) Info() (fs.FileInfo, error) {
	i, err := e.DirEntry.Info()
	return liarInfo{i}, err
}
