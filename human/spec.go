package human

// Sex is the sex a random person is: AnySex leaves it to chance.
type Sex int

const (
	AnySex Sex = iota
	Woman
	Man
)

// Ancestry is how much of each ancestry a person has, in any units: all
// zero leaves it to chance.
type Ancestry struct{ African, Asian, European float64 }

// Spec is what is asked of a random person; what is left zero is chosen
// by chance. Age is in years: 0 for a grown person of any age.
type Spec struct {
	Sex      Sex
	Age      float64
	Ancestry Ancestry
	Seed     uint32
}
