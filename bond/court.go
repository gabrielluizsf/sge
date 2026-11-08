package bond

// Moment is when one asks: the day, whether both are grown up, and
// whether they are kin.
type Moment struct {
	Day   int
	Grown bool
	Kin   bool
}

// may is why the two may not court at m, if they may not.
func may(m Moment) error {
	if !m.Grown {
		return ErrYoung
	}
	if m.Kin {
		return ErrKin
	}
	return nil
}

// AskOut asks a friend fond enough, DateAt, to date, at m.
func (b *Bond) AskOut(m Moment) error {
	if err := may(m); err != nil {
		return err
	}
	if b.Courted >= Dating {
		return ErrAlready
	}
	if b.Affection < DateAt {
		return ErrNotYet
	}
	b.Courted = Dating
	return nil
}

// Propose asks one dated and fond enough, ProposeAt, to marry, at m:
// they wed WeddingAfter days on.
func (b *Bond) Propose(m Moment) error {
	if err := may(m); err != nil {
		return err
	}
	switch {
	case b.Courted > Dating:
		return ErrAlready
	case b.Courted < Dating || b.Affection < ProposeAt:
		return ErrNotYet
	}
	b.Courted, b.wedding = Engaged, m.Day+WeddingAfter
	return nil
}

// Wed marries the engaged on their wedding day, or after: whether they
// are married now.
func (b *Bond) Wed(day int) bool {
	if b.Courted == Engaged && day >= b.wedding {
		b.Courted = Married
	}
	return b.Courted == Married
}

// WeddingDay is the day the engaged are to wed.
func (b Bond) WeddingDay() int { return b.wedding }
