// Package random is a small random sequence that is the same for the same
// seed every run, in Go and on Deno alike: what a game draws by chance, a
// villager's wish or a tree's yield, plays back the same.
package random

// Source is the sequence, a linear congruential generator in 32 bits,
// which JavaScript does fast.
type Source uint32

// New is the sequence that starts from seed.
func New(seed uint32) *Source {
	s := Source(seed)
	return &s
}
