package farm

import "github.com/gabrielluizsf/sge/life"

// Who is a person as the art sees them: the player or a villager, a
// woman or a man, their trade ("farmer", "lumberjack", "smith" or
// "wanderer"; none for the player), five numbers from 0 to 63 their
// town's seed drew for them, or their parents' gave them, for the art to
// pick their skin, hair and clothes by, and their time of life and whole
// years, as the calendar ages them. The founding farmer's numbers are all
// 0. The same person at the same age is always drawn the same; the farm
// asks for a new Body as each one's Stage changes.
type Who struct {
	Player bool
	Woman  bool
	Role   string
	Genes  [5]int
	Stage  life.Stage
	Years  int
}
