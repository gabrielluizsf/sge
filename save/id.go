package save

import "strings"

// Id is a game's title as a folder name: lowercase letters and digits, runs
// of anything else one dash, never empty, never longer than maxName, so it
// is never a path whatever the title.
func Id(title string) string {
	var b strings.Builder
	dash := true
	for _, r := range strings.ToLower(title) {
		if b.Len() >= maxName {
			break
		}
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			dash = false
		} else if !dash {
			b.WriteByte('-')
			dash = true
		}
	}
	id := strings.Trim(b.String(), "-")
	if id == "" {
		return "game"
	}
	return id
}
