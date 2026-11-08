package installer

// The systems an installer knows how to install on.
const (
	Linux   = "linux"
	Windows = "windows"
	MacOS   = "darwin"
)

// Place is the system installed on and the player's own folders on it,
// as Here finds them. Data is where programs keep their files for the
// player alone: ~/.local/share on Linux, %LOCALAPPDATA% on Windows;
// Roaming (Windows) is %APPDATA%, which holds the Start menu; SystemRoot
// (Windows) is where the system's own programs are; Caches the
// player's cache folder, which a Setup's Cache must be inside.
type Place struct {
	System                                  string
	Home, Data, Roaming, SystemRoot, Caches string
}
