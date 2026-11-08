package life

import "github.com/gabrielluizsf/sge/random"

// Spans are the years a life lasts, drawn between: none dies young of old
// age, few live past ninety.
const (
	SpanFrom = 68
	SpanTo   = 94
)

// Draw is a life drawn by luck, of someone years old on day, living a
// span drawn between SpanFrom and SpanTo and longer than they are old.
func (c Calendar) Draw(luck *random.Source, years float64, day int) Life {
	span := SpanFrom + luck.Float()*(SpanTo-SpanFrom)
	span = max(span, years+1+luck.Float()*5)
	return Life{Born: day - int(years*float64(max(c.DaysPerYear, 1))), Span: span}
}

// Newborn is a life born on day, its span drawn by luck.
func (c Calendar) Newborn(luck *random.Source, day int) Life { return c.Draw(luck, 0, day) }
