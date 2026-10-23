package gots

import (
	"fmt"
	"strings"
)

// jsString is a JavaScript string literal holding the bytes of the Go
// string s, one code unit per byte, the bytes above ASCII moved to the
// private use area as the runtime keeps them.
func jsString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			b.WriteString(`\"`)
		case c == '\\':
			b.WriteString(`\\`)
		case c == '\n':
			b.WriteString(`\n`)
		case c == '\r':
			b.WriteString(`\r`)
		case c == '\t':
			b.WriteString(`\t`)
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(c)
		case c >= 0x80:
			// A byte above ASCII is held in the private use area; see
			// runtime/rt/string.ts.
			fmt.Fprintf(&b, `\u%04x`, 0xf700+int(c))
		default:
			fmt.Fprintf(&b, `\x%02x`, c)
		}
	}
	b.WriteByte('"')
	return b.String()
}
