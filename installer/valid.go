package installer

import (
	"io/fs"
	"strings"
	"unicode"
)

// validID is whether id may name the app's folders on every system.
func validID(id string) bool {
	if id == "" || len(id) > 64 || !alnum(rune(id[0])) || strings.HasSuffix(id, ".") || reserved(id) {
		return false
	}
	for _, r := range id {
		if !alnum(r) && r != '.' && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

func alnum(r rune) bool {
	return r < 128 && (unicode.IsLetter(r) || unicode.IsDigit(r))
}

// validName is whether name may be shown and name a bundle or shortcut.
func validName(name string) bool {
	t := strings.TrimSpace(name)
	if t != name || name == "" || strings.HasSuffix(name, ".") || reserved(name) {
		return false
	}
	return validText(name, 64) && !strings.ContainsAny(name, `/\:*?"<>|%`)
}

// validText is whether s is at most max bytes of printable characters.
func validText(s string, max int) bool {
	if len(s) > max {
		return false
	}
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// validFile is whether name is a path inside the installed files, the
// same on every system: no "..", no root, no drive, no stream, no device.
func validFile(name string) bool {
	if !fs.ValidPath(name) || name == "." || strings.ContainsAny(name, `\:`) || !validText(name, 255) {
		return false
	}
	for _, part := range strings.Split(name, "/") {
		if reserved(part) || strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") {
			return false
		}
	}
	return true
}

// reserved is whether Windows takes name, with any extension, as a device.
func reserved(name string) bool {
	base, _, _ := strings.Cut(strings.ToUpper(name), ".")
	switch strings.TrimRight(base, " ") {
	case "CON", "PRN", "AUX", "NUL", "CONIN$", "CONOUT$":
		return true
	}
	r := []rune(strings.TrimRight(base, " "))
	if len(r) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) {
		return strings.ContainsRune("0123456789¹²³", r[3])
	}
	return false
}
