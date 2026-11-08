package installer

import (
	"io"
	"io/fs"
)

// maxIconBytes bounds the icon read from the files.
const maxIconBytes = 16 << 20

// iconData is the app's icon read from the files, or nil with none.
func iconData(s Setup) ([]byte, error) {
	if s.App.Icon == "" {
		return nil, nil
	}
	f, err := s.Files.Open(s.App.Icon)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if info, err := f.Stat(); err != nil || !info.Mode().IsRegular() {
		return nil, fs.ErrInvalid
	}
	return io.ReadAll(io.LimitReader(f, maxIconBytes))
}
