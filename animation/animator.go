package animation

// Animator plays one of its animations at a time.
type Animator struct {
	anims   map[string]*Animation
	current string
	elapsed float64
	playing bool
}
