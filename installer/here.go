package installer

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// Here is the place of the player running this program.
func Here() (Place, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Place{}, err
	}
	caches, err := os.UserCacheDir()
	if err != nil {
		return Place{}, err
	}
	p := Place{System: system(), Home: home, Caches: caches}
	switch p.System {
	case Windows:
		p.Data, p.Roaming = os.Getenv("LOCALAPPDATA"), os.Getenv("APPDATA")
		p.SystemRoot = os.Getenv("SystemRoot")
		if p.Data == "" || p.Roaming == "" {
			return Place{}, errors.New("installer: %LOCALAPPDATA% or %APPDATA% is not set")
		}
		if p.SystemRoot == "" {
			p.SystemRoot = `C:\Windows`
		}
	case Linux:
		p.Data = os.Getenv("XDG_DATA_HOME")
		if !filepath.IsAbs(p.Data) { // the spec says to ignore a relative one
			p.Data = filepath.Join(home, ".local", "share")
		}
	}
	return p, nil
}

// system is the system this runs on. Transpiled to TypeScript, Go's is
// "js", and Deno's own system is told by what only it has.
func system() string {
	if runtime.GOOS != "js" {
		return runtime.GOOS
	}
	if os.Getenv("OS") == "Windows_NT" {
		return Windows
	}
	if _, err := os.Stat("/System/Library/CoreServices/SystemVersion.plist"); err == nil {
		return MacOS
	}
	return Linux
}
