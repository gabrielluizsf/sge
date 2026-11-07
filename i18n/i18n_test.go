package i18n

import "testing"

func TestEnglishIsTheDefaultAndNextChangesItBack(t *testing.T) {
	Set(English)
	if got := T("Save"); got != "Save" {
		t.Errorf("default: %q", got)
	}
	if got := TF(true, "farmer"); got != "farmer" {
		t.Errorf("default woman: %q", got)
	}
	Next(1)
	if got := T("Save"); got != "Salvar" {
		t.Errorf("Portuguese: %q", got)
	}
	if got := TF(true, "farmer"); got != "fazendeira" {
		t.Errorf("Portuguese woman: %q", got)
	}
	if got := TF(false, "farmer"); got != "fazendeiro" {
		t.Errorf("Portuguese man: %q", got)
	}
	if got := T("Day %d", 3); got != "Dia 3" {
		t.Errorf("formatted: %q", got)
	}
	Next(1)
	if got := T("Save"); got != "Save" {
		t.Errorf("back to English: %q", got)
	}
}
