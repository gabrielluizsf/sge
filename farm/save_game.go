package farm

import (
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/save"
)

// saveVersion is the version of the save file's format; saveName its file.
const (
	saveVersion = 1
	saveName    = "game"
)

// HasSave reports whether the game titled title has a saved game, for a
// main menu to offer "Start" (continue) or not. It does not check the
// save is good: Load does.
func HasSave(title string) bool {
	st, err := save.Open(title, saveVersion)
	return err == nil && st.Has(saveName)
}

// EnableSave gives the pause menu a "Save" button, which saves the game
// as the one titled title, and says so on screen.
func (s *Session) EnableSave(title string) {
	s.f.ui.save = func() {
		text := i18n.T("Game saved")
		if err := s.Save(title); err != nil {
			text = i18n.T("Could not save")
		}
		s.f.notify(text)
		s.f.sound("ui_open", nil)
	}
}

// Save writes the game as the saved game of the one titled title, whole or
// not at all: a failure leaves the earlier save as it was.
func (s *Session) Save(title string) error {
	st, err := save.Open(title, saveVersion)
	if err != nil {
		return err
	}
	return st.Write(saveName, takeSnapshot(s.f))
}

// Load puts the saved game of the one titled title in s, which has to be a
// new session, before it is played. A save that is not good gives an
// error and leaves s as it was, a new game. A main menu's "Iniciar" is
// New(art, ctl), then Load(title).
func (s *Session) Load(title string) error {
	st, err := save.Open(title, saveVersion)
	if err != nil {
		return err
	}
	var snap snapshot
	if err := st.Read(saveName, &snap, maxSaveBytes); err != nil {
		return err
	}
	return restore(s.f, &snap)
}
