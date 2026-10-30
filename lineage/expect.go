package lineage

// Gestation is how many days after a couple decide to have a child it is
// born.
const Gestation = 28

// Expecting is a child on the way: their parents, and the day they are
// due.
type Expecting struct {
	Parents [2]string
	Due     int
}

// Expect is the child a and b decide on day to have.
func Expect(parents [2]string, day int) Expecting {
	return Expecting{Parents: parents, Due: day + Gestation}
}

// Born is whether the child is born by day.
func (e Expecting) Born(day int) bool { return day >= e.Due }
