package installer

import (
	"bytes"
	"encoding/binary"
)

// icnsTypes are the .icns types that hold a PNG, by side.
var icnsTypes = map[int]string{1024: "ic10", 512: "ic09", 256: "ic08", 128: "ic07", 64: "icp6", 32: "icp5", 16: "icp4"}

// icnsOf is the PNG icon as a macOS .icns holding one PNG picture.
func icnsOf(data []byte) ([]byte, error) {
	pic, side, err := iconPNG(data, []int{1024, 512, 256, 128, 64, 32, 16})
	if err != nil {
		return nil, err
	}
	be := binary.BigEndian
	var b bytes.Buffer
	b.WriteString("icns")
	b.Write(be.AppendUint32(nil, uint32(8+8+len(pic))))
	b.WriteString(icnsTypes[side])
	b.Write(be.AppendUint32(nil, uint32(8+len(pic))))
	b.Write(pic)
	return b.Bytes(), nil
}
