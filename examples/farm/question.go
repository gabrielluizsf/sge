package farm

import "github.com/gabrielluizsf/sge/examples/i18n"

// question is a question put to the player, in a panel in the middle,
// that waits for an answer: its words, and the answers they may choose,
// the chosen one on a gold bar.
type question struct {
	ask     []string
	answers []line
	choice  int
}

// answerQuestion answers the open question with the chosen answer: the
// question closes first, so the answer may ask another.
func (f *farm) answerQuestion() {
	q := f.question
	f.question = nil
	q.answers[q.choice].do(f)
}

// childQuestion is what the player is asked going to bed married, and
// nothing else: whether to have a child. Either way the night passes.
func childQuestion() *question {
	return &question{ask: []string{i18n.T("Should we have a child?")}, answers: []line{
		{text: i18n.T("Yes"), do: func(f *farm) { f.conceive(); f.sleep() }},
		{text: i18n.T("No"), do: func(f *farm) { f.sleep() }},
	}}
}
