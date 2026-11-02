package dialogue

import "testing"

// While a page is said, Saying is the letter last shown; once it is all
// shown, nothing is being said.
func TestSayingIsTheLetterLastShown(t *testing.T) {
	b, _ := said("Oi, Ana")
	if _, ok := b.Saying(); !ok {
		t.Fatal("nothing said as the page begins")
	}
	b.Update(2.0 / b.Speed) // two letters
	if r, ok := b.Saying(); !ok || r != 'i' {
		t.Errorf("saying %q, %v, not the i", r, ok)
	}
	b.Update(1)
	if _, ok := b.Saying(); ok {
		t.Error("still saying with the page all shown")
	}
}
