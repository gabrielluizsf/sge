package installer

import "errors"

// App is what is installed: ID names its folders and entries, letters,
// digits and . _ - only, such as "com.example.farm"; Name is shown to the
// player; Exe and Icon (a PNG, optional) are paths inside the files
// installed; Game files it among the system's games.
type App struct {
	ID, Name, Version, Publisher string
	Exe, Icon                    string
	Game                         bool
}

// check is an error when the app's fields could reach the system as
// something other than what they say: a path, a new line in an entry.
func (a App) check() error {
	switch {
	case !validID(a.ID):
		return errors.New("installer: ID must be 1 to 64 letters, digits, '.', '_' or '-', starting with a letter or digit")
	case !validName(a.Name):
		return errors.New("installer: Name must be 1 to 64 printable characters, none of / \\ : * ? \" < > |")
	case !validText(a.Version, 32) || !validText(a.Publisher, 64):
		return errors.New("installer: Version or Publisher holds a control character or is too long")
	case !validFile(a.Exe):
		return errors.New("installer: Exe must be a path inside the files")
	case a.Icon != "" && !validFile(a.Icon):
		return errors.New("installer: Icon must be a path inside the files")
	}
	return nil
}
