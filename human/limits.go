package human

import (
	"fmt"
	"io"
	"io/fs"
	"path"
)

// Limits bound what a kit's files may make it allocate: they are read as
// untrusted, a pack downloaded from anyone.
type Limits struct {
	// Bytes is the most a file may be, as read, unpacked if it is packed.
	Bytes int64
	// Pixels is the most a texture may have, width times height.
	Pixels int64
	// Texture is the most pixels across a texture is kept, shrunk to it.
	Texture int
}

// DefaultLimits are generous for MakeHuman's own files: 64 MiB a file, 16
// megapixels a texture, kept at 1024 pixels across.
var DefaultLimits = Limits{Bytes: 64 << 20, Pixels: 16 << 20, Texture: 1024}

// read is the file name of fsys, refused if its name could climb out of
// fsys or it is larger than the limit.
func (l Limits) read(fsys fs.FS, name string) ([]byte, error) {
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("human: %q is not a name inside the kit", name)
	}
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, l.Bytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > l.Bytes {
		return nil, fmt.Errorf("human: %s is over %d bytes", name, l.Bytes)
	}
	return data, nil
}

// beside is the name of the file a file at from names as rel: next to it,
// or up and down from there, but never out of the kit.
func beside(from, rel string) (string, error) {
	if path.IsAbs(rel) {
		return "", fmt.Errorf("human: %q in %s is not relative", rel, from)
	}
	name := path.Join(path.Dir(from), rel)
	if !fs.ValidPath(name) {
		return "", fmt.Errorf("human: %q in %s leads out of the kit", rel, from)
	}
	return name, nil
}
