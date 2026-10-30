package farm

// Who is a person as the art sees them: the player or a villager, a
// woman or a man, their trade ("farmer", "lumberjack", "smith" or
// "wanderer"; none for the player), and five numbers from 0 to 63 their
// town's seed drew for them, for the art to pick their skin, hair and
// clothes by. The same person is always drawn the same.
type Who struct {
	Player bool
	Woman  bool
	Role   string
	Genes  [5]int
}
