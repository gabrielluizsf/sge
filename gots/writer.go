package gots

import (
	"fmt"
	"strings"
)

// writer collects indented lines of TypeScript.
type writer struct {
	b      strings.Builder
	indent int
}

// line writes one line at the current indentation.
func (w *writer) line(format string, args ...any) {
	if format == "" {
		w.b.WriteByte('\n')
		return
	}
	w.b.WriteString(strings.Repeat("  ", w.indent))
	fmt.Fprintf(&w.b, format, args...)
	w.b.WriteByte('\n')
}

// open writes a line ending a block opener and indents what follows.
func (w *writer) open(format string, args ...any) {
	w.line(format, args...)
	w.indent++
}

// close dedents and writes the line closing the block.
func (w *writer) close(format string, args ...any) {
	w.indent--
	w.line(format, args...)
}

func (w *writer) String() string { return w.b.String() }
