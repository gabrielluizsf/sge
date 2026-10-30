package farm

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
	return &question{ask: []string{"Devemos ter um filho?"}, answers: []line{
		{text: "Sim", do: func(f *farm) { f.conceive(); f.sleep() }},
		{text: "Não", do: func(f *farm) { f.sleep() }},
	}}
}
