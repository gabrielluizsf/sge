package gltf

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// The GLB container: a 12-byte header, then chunks, each its length, type
// and data, padded to four bytes.
const (
	glbMagic   = 0x46546C67 // "glTF"
	chunkJSON  = 0x4E4F534A // "JSON"
	chunkBIN   = 0x004E4942 // "BIN\0"
	glbVersion = 2
)

// errFormat is what every malformed file is, wrapped with why.
var errFormat = errors.New("gltf: malformed")

func malformed(format string, a ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{errFormat}, a...)...)
}

// splitGLB is a GLB's JSON chunk and its binary chunk, if it has one.
func splitGLB(data []byte) (js, bin []byte, err error) {
	if len(data) < 12 || binary.LittleEndian.Uint32(data) != glbMagic {
		return nil, nil, malformed("not a GLB")
	}
	if v := binary.LittleEndian.Uint32(data[4:]); v != glbVersion {
		return nil, nil, malformed("GLB version %d", v)
	}
	if total := binary.LittleEndian.Uint32(data[8:]); uint64(total) != uint64(len(data)) {
		return nil, nil, malformed("GLB says it is %d bytes, not %d", total, len(data))
	}
	for at := 12; at < len(data); {
		if len(data)-at < 8 {
			return nil, nil, malformed("truncated chunk header")
		}
		n := uint64(binary.LittleEndian.Uint32(data[at:]))
		kind := binary.LittleEndian.Uint32(data[at+4:])
		at += 8
		if n > uint64(len(data)-at) {
			return nil, nil, malformed("chunk of %d bytes past the end", n)
		}
		chunk := data[at : at+int(n)]
		at += int(n)
		switch {
		case kind == chunkJSON && js == nil:
			js = chunk
		case kind == chunkBIN && bin == nil && js != nil:
			bin = chunk
		}
	}
	if js == nil {
		return nil, nil, malformed("no JSON chunk")
	}
	return js, bin, nil
}
