package farm

import "github.com/gabrielluizsf/antui/canvas"

// questionPanel is the open question's panel: its words, and the answers
// under them, the chosen one on a gold bar.
func (h *hud) questionPanel(cv *canvas.Canvas) {
	q := h.f.question
	lines := append([]string{}, q.ask...)
	for _, a := range q.answers {
		lines = append(lines, a.text)
	}
	h.menu(cv, menuOf{lines: lines, chosen: len(q.ask) + q.choice})
}
