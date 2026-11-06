package i18n

import "github.com/gabrielluizsf/sge/save"

// The player's choice of language, kept on their computer for every
// example game: in the folder choiceFolder of the user's config directory,
// through package save, which writes whole files and reads them as
// untrusted. choiceMost is the most of the file that is read.
const (
	choiceFolder  = "sge-examples"
	choiceName    = "language"
	choiceVersion = 1
	choiceMost    = 256
)

// choice is the file's content: the language by its code, not its number,
// so adding languages changes no file.
type choice struct {
	Language string `json:"language"`
}

// codeOf is the code language l is kept under.
func codeOf(l Language) string {
	if l == Portuguese {
		return "pt"
	}
	return "en"
}

// languageOf is the language kept under code, and whether there is one.
func languageOf(code string) (Language, bool) {
	switch code {
	case "en":
		return English, true
	case "pt":
		return Portuguese, true
	}
	return English, false
}

// Load puts the language the player last chose on this computer. With no
// choice yet, or a file that is damaged, the language stays as it is: no
// file stops a game from starting.
func Load() {
	st, err := save.Open(choiceFolder, choiceVersion)
	if err != nil || !st.Has(choiceName) {
		return
	}
	var c choice
	if err := st.Read(choiceName, &c, choiceMost); err != nil {
		return
	}
	if l, ok := languageOf(c.Language); ok {
		Set(l)
	}
}

// Change moves step languages along the list and keeps the choice on this
// computer. The language changes even if it could not be kept, which the
// error says.
func Change(step int) error {
	l := Next(step)
	st, err := save.Open(choiceFolder, choiceVersion)
	if err != nil {
		return err
	}
	return st.Write(choiceName, choice{Language: codeOf(l)})
}
