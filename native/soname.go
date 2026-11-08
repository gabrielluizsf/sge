package native

// validSoname accepts only bare sonames like "libEGL.so.1": "lib", a
// name of [A-Za-z0-9_+-], ".so", then optional ".N" version parts. No
// '/', so dlopen never treats it as a path (relative or absolute).
func validSoname(s string) bool {
	const pre, mid = "lib", ".so"
	if len(s) > 64 || len(s) < len(pre)+1+len(mid) || s[:3] != pre {
		return false
	}
	i := 3
	for i < len(s) && nameByte(s[i]) {
		i++
	}
	if i == 3 || len(s)-i < len(mid) || s[i:i+3] != mid {
		return false
	}
	return versionTail(s[i+3:])
}

func nameByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '+' || c == '-'
}

// versionTail accepts "" or ".1", ".0.0", ...
func versionTail(t string) bool {
	for len(t) > 0 {
		if t[0] != '.' || len(t) < 2 || t[1] < '0' || t[1] > '9' {
			return false
		}
		t = t[2:]
		for len(t) > 0 && t[0] >= '0' && t[0] <= '9' {
			t = t[1:]
		}
	}
	return true
}
