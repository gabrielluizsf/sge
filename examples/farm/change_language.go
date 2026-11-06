package farm

import "github.com/gabrielluizsf/sge/examples/i18n"

// changeLanguage moves to the next language and keeps the choice on this
// computer, saying so if it could not be kept.
func (f *farm) changeLanguage() {
	if err := i18n.Change(1); err != nil {
		f.notify(i18n.T("Could not save"))
	}
	f.sound("ui_tick", nil)
}
