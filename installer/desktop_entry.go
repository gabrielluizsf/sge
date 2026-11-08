package installer

import (
	"fmt"
	"strings"
)

// desktopEntry is the Desktop Entry file starting exe in dir, by the
// freedesktop.org spec: values escaped as strings, the program quoted as
// an argument of Exec, and no character that could start a new line, so
// nothing in a name or path can add a key of its own.
func desktopEntry(app App, exe, dir, icon string) string {
	category := "Utility;"
	if app.Game {
		category = "Game;"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[Desktop Entry]\nType=Application\nVersion=1.5\nName=%s\n", entryString(app.Name))
	fmt.Fprintf(&b, "Exec=%s\nPath=%s\n", entryString(execArg(exe)), entryString(dir))
	if icon != "" {
		fmt.Fprintf(&b, "Icon=%s\n", entryString(icon))
	}
	fmt.Fprintf(&b, "Terminal=false\nCategories=%s\n", category)
	return b.String()
}

// execArg is path quoted as one argument of Exec: in double quotes, with
// ", `, $ and \ escaped by a backslash, and % doubled, as no field code.
func execArg(path string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "`", "\\`", `$`, `\$`, `%`, `%%`)
	return `"` + r.Replace(path) + `"`
}

// entryString is s escaped as a string value: backslashes doubled, and
// control characters as their escapes. A leading space is kept by \s.
func entryString(s string) string {
	s = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "\t", `\t`, "\r", `\r`).Replace(s)
	if strings.HasPrefix(s, " ") {
		s = `\s` + s[1:]
	}
	return s
}
