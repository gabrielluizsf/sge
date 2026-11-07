package main

// liveUpdate is what the frame takes from the preview: the person to
// show, made new (to be uploaded) or the one shown, its colours changed
// by applies, in order.
type liveUpdate struct {
	person  person
	fresh   bool
	applies []func() []int
}
