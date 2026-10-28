package inventory

import (
	"slices"
	"testing"
)

func TestABagCountsWhatIsPutInAndTakenOut(t *testing.T) {
	var b Bag
	b.Add("wood", 3)
	b.Add("carrot", 1)
	b.Add("wood", 2)
	b.Add("stone", 0)
	if b.Count("wood") != 5 || b.Count("carrot") != 1 || b.Count("stone") != 0 {
		t.Fatalf("content %v", b.Content())
	}
	if b.Take("wood", 6) || b.Count("wood") != 5 {
		t.Error("took more wood than there was")
	}
	if !b.Take("wood", 5) || b.Count("wood") != 0 {
		t.Error("did not take the wood")
	}
	if b.Take("carrot", -1) {
		t.Error("took a negative number")
	}
	if want := (List{{"carrot", 1}}); !slices.Equal(b.Content(), want) {
		t.Errorf("content %v, not %v", b.Content(), want)
	}
}

// A list is taken whole or not at all.
func TestTakeAllTakesAllOrNothing(t *testing.T) {
	var b Bag
	b.Add("wood", 2)
	b.Add("tomato", 1)
	if b.TakeAll(List{{"wood", 2}, {"tomato", 2}}) || b.Count("wood") != 2 {
		t.Fatal("took part of a list the bag cannot fill")
	}
	if b.TakeAll(List{{"wood", 1}, {"wood", 2}}) {
		t.Fatal("counted the same wood twice")
	}
	if !b.TakeAll(List{{"wood", 1}, {"tomato", 1}}) || b.Count("wood") != 1 || b.Count("tomato") != 0 {
		t.Fatalf("content %v", b.Content())
	}
}

// What is in the bag is listed in the order it first came in.
func TestContentKeepsItsOrder(t *testing.T) {
	var b Bag
	for _, item := range []string{"c", "a", "b", "a"} {
		b.Add(item, 1)
	}
	if want := (List{{"c", 1}, {"a", 2}, {"b", 1}}); !slices.Equal(b.Content(), want) {
		t.Errorf("content %v, not %v", b.Content(), want)
	}
}
