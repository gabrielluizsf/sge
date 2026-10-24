// encoding/binary reads and writes structs through reflect.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	Magic   [4]byte
	Version uint16
	Flags   uint8
	_       uint8
	Count   int32
	Scale   float32
	Offsets [3]uint64
}

func main() {
	h := Header{Magic: [4]byte{'G', 'O', 'T', 'S'}, Version: 2, Flags: 5, Count: -3, Scale: 1.5, Offsets: [3]uint64{1, 1 << 40, 1<<63 + 1}}
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, h); err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf.Len(), binary.Size(h))
	fmt.Printf("% x\n", buf.Bytes()[:16])
	var back Header
	err := binary.Read(bytes.NewReader(buf.Bytes()), binary.LittleEndian, &back)
	fmt.Println(err, back == h, string(back.Magic[:]), back.Offsets[2])

	nums := []int16{1, -2, 300}
	var nb bytes.Buffer
	binary.Write(&nb, binary.BigEndian, nums)
	fmt.Printf("% x\n", nb.Bytes())
	out := make([]int16, 3)
	binary.Read(&nb, binary.BigEndian, out)
	fmt.Println(out)
}
