package bakes

import (
	"testing"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

// A model is drawn at its own scale when it has one, and otherwise as all
// are; and a change of scale is a task of another name, baked again.
func TestAModelIsBakedAtItsOwnScale(t *testing.T) {
	plain := kit.Asset{ID: "rock"}
	if modelOptions(plain).Unit != Options().Unit {
		t.Error("a model with no scale was scaled")
	}
	big := kit.Asset{ID: "rock", Scale: 2.5}
	if got, want := modelOptions(big).Unit, Options().Unit*2.5; got != want {
		t.Errorf("unit %v, want %v", got, want)
	}
	if modelTaskName(plain) == modelTaskName(big) {
		t.Error("a new scale would not be baked again")
	}
	other := kit.Asset{ID: "rock", Scale: 3}
	if modelTaskName(other) == modelTaskName(big) {
		t.Error("another scale would not be baked again")
	}
}
