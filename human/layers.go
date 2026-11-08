package human

import (
	"errors"
	"io/fs"
)

// layers is several file systems as one, a file of the first that has it.
type layers []fs.FS

// Layers is fss as one file system: a name is opened in the first that
// has it, so MakeHuman's data folder and its asset packs make one kit.
func Layers(fss ...fs.FS) fs.FS { return layers(fss) }

func (l layers) Open(name string) (fs.File, error) {
	for _, f := range l {
		file, err := f.Open(name)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			return file, err
		}
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}
