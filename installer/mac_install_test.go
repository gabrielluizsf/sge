//go:build darwin

package installer

import (
	"bytes"
	"image/png"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallingOnMacOSMakesAnAppBundle(t *testing.T) {
	p := homeOn(t, MacOS)
	app := farmApp
	app.Name = "Farm & Co's"
	if err := Install(Setup{App: app, Files: gameFiles(t), Place: p, Uninstaller: []byte("u")}); err != nil {
		t.Fatal(err)
	}
	bundle := filepath.Join(p.Home, "Applications", "Farm & Co's.app")
	if read(t, filepath.Join(bundle, "Contents/MacOS/assets/tree.png")) != "tree" {
		t.Error("the files are not in Contents/MacOS")
	}
	plist := read(t, filepath.Join(bundle, "Contents/Info.plist"))
	for _, want := range []string{"<string>Farm &amp; Co&#39;s</string>", "<key>CFBundleExecutable</key>\n\t<string>farm</string>", "<string>com.example.farm</string>", "public.app-category.games", "<string>icon.icns</string>"} {
		if !strings.Contains(plist, want) {
			t.Errorf("Info.plist has no %q:\n%s", want, plist)
		}
	}
	if out, err := exec.Command("/usr/bin/plutil", "-lint", filepath.Join(bundle, "Contents/Info.plist")).CombinedOutput(); err != nil {
		t.Errorf("plutil finds Info.plist wrong: %v %s", err, out)
	}
	icns := []byte(read(t, filepath.Join(bundle, "Contents/Resources/icon.icns")))
	if string(icns[:4]) != "icns" || string(icns[8:12]) != "ic08" {
		t.Fatalf("the icon is not a 256-pixel icns: %q", icns[:12])
	}
	img, err := png.Decode(bytes.NewReader(icns[16:]))
	if err != nil || img.Bounds().Dx() != 256 || img.Bounds().Dy() != 256 {
		t.Errorf("the icon's picture is not 256 square: %v", err)
	}
	if exists(filepath.Join(bundle, "uninstall")) || exists(filepath.Join(bundle, "Contents/MacOS/uninstall")) {
		t.Error("a bundle has an uninstaller: on macOS it goes in the bin")
	}
	if err := Uninstall(Setup{App: app, Place: p}); err != nil || exists(bundle) {
		t.Errorf("the bundle was not taken away: %v", err)
	}
}
