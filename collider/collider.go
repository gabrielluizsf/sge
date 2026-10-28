// Package collider decides whether two shapes overlap.
package collider

// Collider is a shape that can be tested against the others. [Box] and
// [Circle] are the shapes there are.
type Collider interface {
	// Bounds is the smallest box the shape fits in.
	Bounds() Box
}
