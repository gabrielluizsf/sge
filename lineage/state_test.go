package lineage

import (
	"reflect"
	"testing"
)

func TestEntriesRestoreTheTree(t *testing.T) {
	var tr Tree
	tr.Add(Person{ID: "a", Born: 1})
	tr.Add(Person{ID: "b", Born: 2})
	tr.Add(Person{ID: "c", Parents: [2]string{"a", "b"}, Born: 30})
	tr.Die("a")
	got, err := Restore(tr.Entries())
	if err != nil || !reflect.DeepEqual(got.Entries(), tr.Entries()) || !reflect.DeepEqual(got.Children("a"), []string{"c"}) {
		t.Fatalf("%v %+v", err, got.Entries())
	}
	if _, err := Restore([]Entry{{Person: Person{ID: "x", Parents: [2]string{"nobody", ""}}}}); err == nil {
		t.Error("restored a child without parents")
	}
	if _, err := Restore([]Entry{{Person: Person{ID: "x"}}, {Person: Person{ID: "x"}}}); err == nil {
		t.Error("restored a duplicate")
	}
}
