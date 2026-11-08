package installer

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// Payload is what an installer carries: its Files, and Program, the
// installer itself without them, which is what uninstalls.
type Payload struct {
	Files   fs.FS
	Program []byte
	file    *os.File
}

// maxProgram bounds the program read back from an installer.
const maxProgram = 1 << 30

// Open reads the payload of the installer at path (Pack's out).
func Open(path string) (*Payload, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	p, err := readPayload(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return p, nil
}

func readPayload(f *os.File) (*Payload, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	z, err := zip.NewReader(f, info.Size())
	if err != nil {
		return nil, errors.New("installer: this program carries no files")
	}
	n, err := strconv.ParseInt(strings.TrimPrefix(z.Comment, payloadTag), 10, 64)
	if !strings.HasPrefix(z.Comment, payloadTag) || err != nil || n <= 0 || n > info.Size() || n > maxProgram {
		return nil, errors.New("installer: this program's files are not an installer's")
	}
	for _, zf := range z.File {
		if !validFile(strings.TrimSuffix(zf.Name, "/")) {
			return nil, errors.New("installer: this program's files hold the name " + strconv.Quote(zf.Name))
		}
	}
	program := make([]byte, n)
	if _, err := io.ReadFull(io.NewSectionReader(f, 0, n), program); err != nil {
		return nil, err
	}
	return &Payload{Files: z, Program: program, file: f}, nil
}

// Close lets the installer's file go.
func (p *Payload) Close() error {
	return p.file.Close()
}
