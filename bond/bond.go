// Package bond is how two people come to love each other in a game: gifts
// and talk, a day at a time, make strangers friends; a friend fond enough
// may be asked out, one dated long enough asked to marry, and a wedding
// some days on makes them a couple who may have a child. Only grown-ups
// court, never kin.
package bond

import "errors"

// Stage is how far a bond has come.
type Stage int

const (
	Strangers Stage = iota
	Acquainted
	Friends
	Dating
	Engaged
	Married
)

// Affection, from 0 to Most, makes strangers acquainted at Acquaint and
// friends at Befriend; a friend may be asked out at DateAt, one dated may
// be asked to marry at ProposeAt; a wedding is WeddingAfter days after.
const (
	Most         = 1000
	Acquaint     = 100
	Befriend     = 300
	DateAt       = 600
	ProposeAt    = 850
	WeddingAfter = 3
)

// Bond is one person's with another: how fond of them they are, how far
// it has come, the day after they were last given a gift and last talked
// to, 0 for never, and the day they are to wed, once engaged.
type Bond struct {
	Affection int
	Courted   Stage // Dating or later, once asked; Strangers before
	gifted    int
	talked    int
	wedding   int
}

// Why a bond will not go on.
var (
	ErrToday   = errors.New("bond: already today")
	ErrYoung   = errors.New("bond: not both grown up")
	ErrKin     = errors.New("bond: they are kin")
	ErrNotYet  = errors.New("bond: not fond enough yet")
	ErrAlready = errors.New("bond: that is behind them")
)

// Stage is how far the bond has come: courted, or as fond as they are.
func (b Bond) Stage() Stage {
	switch {
	case b.Courted >= Dating:
		return b.Courted
	case b.Affection >= Befriend:
		return Friends
	case b.Affection >= Acquaint:
		return Acquainted
	}
	return Strangers
}

// MayHaveChild is whether the two are married: only then.
func (b Bond) MayHaveChild() bool { return b.Courted == Married }

// fonder adds n to the affection, held from 0 to Most.
func (b *Bond) fonder(n int) { b.Affection = min(Most, max(0, b.Affection+n)) }
