package game

import "github.com/gabrielluizsf/sge/transform"

// NewObject is an object at the origin, at its natural size.
func NewObject(name string) *Object {
	return &Object{Name: name, Transform: transform.New()}
}
