// byteorderfuse: byte-order reads and writes on subslices, written fused
// with them under gots, and their bounds, as Go checks them.
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func main() {
	b := make([]byte, 6, 8)
	for i := range b {
		b[i] = byte(0xa0 + i)
	}
	be, le := binary.BigEndian, binary.LittleEndian
	fmt.Printf("%x %x %x %x\n", be.Uint32(b), be.Uint32(b[2:]), le.Uint16(b[4:]), le.Uint32(b[1:5]))
	// Up to the capacity, past the length, as a subslice may.
	fmt.Printf("%x\n", be.Uint32(b[4:8]))
	be.PutUint32(b[1:], 0xdeadbeef)
	le.PutUint16(b[4:], 0x1234)
	fmt.Println(hex.EncodeToString(b))
	fmt.Println(panics(func() { be.Uint32(b[3:]) }), panics(func() { be.Uint32(b[5:9]) }),
		panics(func() { be.Uint32(b[3:6]) }), panics(func() { le.PutUint32(b[3:], 1) }),
		panics(func() { be.Uint16(b[7:]) }), panics(func() { be.Uint16(b[-len(b)+7:]) }))
	// The standard library's own: sha256 reads blocks and writes its sum.
	sum := sha256.Sum256(b)
	fmt.Println(hex.EncodeToString(sum[:]))
}

func panics(f func()) (p bool) {
	defer func() { p = recover() != nil }()
	f()
	return
}
