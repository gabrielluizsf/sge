//go:build windows

package installer

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf16"
)

func windowsPlace() Place {
	return Place{System: Windows, Home: `C:\Users\Ana`, Data: `C:\Users\Ana\AppData\Local`, Roaming: `C:\Users\Ana\AppData\Roaming`, SystemRoot: `C:\Windows`, Caches: `C:\Users\Ana\AppData\Local`}
}

func TestOnWindowsTheGameGoesInThePlayersPrograms(t *testing.T) {
	l, err := layoutOf(farmApp, windowsPlace(), testPNG(t, 64, 64))
	if err != nil {
		t.Fatal(err)
	}
	if l.Dir != `C:\Users\Ana\AppData\Local\Programs\com.example.farm` || l.Exe != l.Dir+`\farm` {
		t.Errorf("the game goes to %s, %s", l.Dir, l.Exe)
	}
	if l.Shortcut != `C:\Users\Ana\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Farm.lnk` {
		t.Errorf("the shortcut goes to %s", l.Shortcut)
	}
	if l.IconLoc != l.Dir+`\icon.ico` {
		t.Errorf("the icon is %s", l.IconLoc)
	}
	ico := l.Inside[1].Data
	if ico[2] != 1 || ico[4] != 1 || ico[6] != 64 || string(ico[22:26]) != "\x89PNG" {
		t.Errorf("icon.ico is not one 64-pixel PNG: % x", ico[:26])
	}
}

// Names reach PowerShell only through its environment: a name that is
// code in PowerShell is never in the script's text.
func TestTheWindowsScriptTakesNoNameAsCode(t *testing.T) {
	app := farmApp
	app.Name = "Farm'; Remove-Item -Recurse ~; '"
	var program string
	var args, env []string
	s := Setup{App: app, Place: windowsPlace(), Run: func(p string, a, e []string) error {
		program, args, env = p, a, e
		return nil
	}}
	l, _ := layoutOf(app, s.Place, nil)
	l.Uninstaller = "uninstall.exe"
	if err := runScript(s, windowsInstall(app, l, 1234)); err != nil {
		t.Fatal(err)
	}
	if program != `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe` {
		t.Errorf("ran %s, not the system's own PowerShell", program)
	}
	text := decoded(t, args[len(args)-1])
	if strings.Contains(text, "Remove-Item -Recurse ~") || !strings.Contains(text, "CreateShortcut($env:SGE_SHORTCUT)") {
		t.Errorf("the script holds the name:\n%s", text)
	}
	want := map[string]bool{"SGE_NAME=" + app.Name: false, `SGE_UninstallString="` + l.Dir + `\uninstall.exe" --uninstall`: false, "SGE_SIZE=1234": false}
	for _, e := range env {
		if _, ok := want[e]; ok {
			want[e] = true
		}
	}
	for e, found := range want {
		if !found {
			t.Errorf("the environment has no %q", e)
		}
	}
}

func decoded(t *testing.T, b64 string) string {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Fatal(err)
	}
	units := make([]uint16, len(b)/2)
	for i := range units {
		units[i] = uint16(b[2*i]) | uint16(b[2*i+1])<<8
	}
	return string(utf16.Decode(units))
}

// Installing for real puts the shortcut in the Start menu and the entry
// among the installed apps, and uninstalling takes both away.
func TestInstallingOnWindowsRegistersTheApp(t *testing.T) {
	app := farmApp
	app.ID = fmt.Sprintf("com.example.sgetest%d", time.Now().UnixNano())
	s := Setup{App: app, Files: gameFiles(t), Place: testPlace(t), Uninstaller: []byte("u")}
	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\` + app.ID
	t.Cleanup(func() { exec.Command("reg", "delete", key, "/f").Run() })
	if err := Install(s); err != nil {
		t.Fatal(err)
	}
	l := where(t, s)
	if read(t, filepath.Join(l.Dir, "assets", "tree.png")) != "tree" || !exists(filepath.Join(l.Dir, "icon.ico")) {
		t.Error("the files are not in the app's folder")
	}
	if !exists(l.Shortcut) {
		t.Errorf("no shortcut at %s", l.Shortcut)
	}
	out, err := exec.Command("reg", "query", key, "/v", "UninstallString").CombinedOutput()
	if err != nil || !strings.Contains(string(out), `uninstall.exe" --uninstall`) {
		t.Errorf("the installed apps' entry is wrong: %v %s", err, out)
	}
	if err := Uninstall(s); err != nil {
		t.Fatal(err)
	}
	if exists(l.Dir) || exists(l.Shortcut) {
		t.Error("the folder or the shortcut is still there")
	}
	if exec.Command("reg", "query", key).Run() == nil {
		t.Error("the installed apps' entry is still there")
	}
}
