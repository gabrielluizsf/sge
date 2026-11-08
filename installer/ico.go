package installer

import (
	"bytes"
	"encoding/binary"
)

// icoOf is the PNG icon as a Windows .ico holding one PNG picture, which
// Windows reads from Vista on: 256 pixels a side or less.
func icoOf(data []byte) ([]byte, error) {
	pic, side, err := iconPNG(data, []int{256, 128, 64, 48, 32, 16})
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	le := binary.LittleEndian
	b.Write(le.AppendUint16(nil, 0))
	b.Write(le.AppendUint16(nil, 1)) // an icon
	b.Write(le.AppendUint16(nil, 1)) // one picture
	dim := byte(side % 256)          // 0 means 256
	b.Write([]byte{dim, dim, 0, 0})
	b.Write(le.AppendUint16(nil, 1))  // planes
	b.Write(le.AppendUint16(nil, 32)) // bits per pixel
	b.Write(le.AppendUint32(nil, uint32(len(pic))))
	b.Write(le.AppendUint32(nil, 6+16)) // where the picture starts
	b.Write(pic)
	return b.Bytes(), nil
}
