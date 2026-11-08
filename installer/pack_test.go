package installer

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// An installer packed carries the files after its program, and gives
// both back: the files to install, the program to uninstall with.
func TestAPackedInstallerGivesBackItsFilesAndProgram(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "game/assets"), 0o755)
	os.WriteFile(filepath.Join(dir, "game/farm"), []byte("game"), 0o755)
	os.WriteFile(filepath.Join(dir, "game/assets/tree.png"), []byte("tree"), 0o644)
	os.WriteFile(filepath.Join(dir, "stub"), []byte("installer program"), 0o755)
	out := filepath.Join(dir, "Farm-Setup")
	if err := Pack(filepath.Join(dir, "stub"), filepath.Join(dir, "game"), out); err != nil {
		t.Fatal(err)
	}
	p, err := Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if string(p.Program) != "installer program" {
		t.Errorf("the program is %q", p.Program)
	}
	s := hostSetup(t, App{ID: "farm", Name: "Farm", Exe: "farm"}, nil)
	s.Files, s.Uninstaller = p.Files, p.Program
	if err := Install(s); err != nil {
		t.Fatal(err)
	}
	l := where(t, s)
	if read(t, filepath.Join(l.Dir, l.Files, "assets", "tree.png")) != "tree" {
		t.Error("the packed files were not installed")
	}
	if l.Uninstaller != "" && read(t, filepath.Join(l.Dir, l.Uninstaller)) != "installer program" {
		t.Error("the program was not installed as the uninstaller")
	}
}

func TestAProgramCarryingNoFilesIsNoInstaller(t *testing.T) {
	name := filepath.Join(t.TempDir(), "plain")
	os.WriteFile(name, []byte("just a program"), 0o755)
	if _, err := Open(name); err == nil {
		t.Error("a plain program opened as an installer")
	}
}

// A PNG claiming a huge picture is refused from its header, before its
// pixels would be decoded into memory.
func TestAHugeIconIsRefusedBeforeDecoding(t *testing.T) {
	b := testPNG(t, 1, 1)
	// Width and height are the IHDR's first 8 bytes, at 16; its CRC is
	// not checked by DecodeConfig's caller once refused by size.
	copy(b[16:24], []byte{0, 1, 0, 0, 0, 1, 0, 0})
	if _, err := iconImage(b); err == nil {
		t.Error("a 65536-pixel icon was taken")
	}
}

func TestAnIconIsCentredAndShrunk(t *testing.T) {
	pic, side, err := iconPNG(testPNG(t, 300, 200), []int{256, 128})
	if err != nil || side != 256 {
		t.Fatal(side, err)
	}
	img, _ := png.Decode(bytes.NewReader(pic))
	if _, _, _, a := img.At(128, 2).RGBA(); a != 0 {
		t.Error("above the picture is not clear")
	}
	if r, _, _, a := img.At(128, 128).RGBA(); a != 0xffff || r>>8 != 200 {
		t.Errorf("the middle is not the picture's red: %d %d", r>>8, a)
	}
}
