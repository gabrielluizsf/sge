package main

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"io"
)

func main() {
	data := []byte("hello hello hello, compress me! " + string(bytes.Repeat([]byte("abcabcabc"), 40)))
	fmt.Printf("crc32=%08x adler32=%08x\n", crc32.ChecksumIEEE(data), adler32.Checksum(data))
	for _, level := range []int{flate.NoCompression, flate.BestSpeed, flate.DefaultCompression, flate.BestCompression, flate.HuffmanOnly} {
		var buf bytes.Buffer
		w, err := flate.NewWriter(&buf, level)
		if err != nil {
			fmt.Println(err)
			continue
		}
		w.Write(data)
		w.Close()
		fmt.Printf("level %d: %d bytes %x\n", level, buf.Len(), buf.Bytes())
		back, err := io.ReadAll(flate.NewReader(bytes.NewReader(buf.Bytes())))
		fmt.Println("round trip:", bytes.Equal(back, data), err)
	}
	var zb bytes.Buffer
	zw := zlib.NewWriter(&zb)
	zw.Write(data)
	zw.Close()
	fmt.Printf("zlib: %x\n", zb.Bytes())
}

func init() {
	// Data larger than a deflate window and block, as an image's is.
	big := make([]byte, 480*320*4+1)
	seed := uint32(1)
	for i := range big {
		seed = seed*1664525 + 1013904223
		if i%7 == 0 {
			big[i] = byte(seed >> 24)
		} else {
			big[i] = byte(i / 977)
		}
	}
	for _, level := range []int{flate.BestSpeed, flate.DefaultCompression} {
		var buf bytes.Buffer
		w, _ := flate.NewWriter(&buf, level)
		w.Write(big)
		w.Close()
		back, err := io.ReadAll(flate.NewReader(bytes.NewReader(buf.Bytes())))
		fmt.Printf("big level %d: %d bytes crc %08x round trip %v %v\n", level, buf.Len(), crc32.ChecksumIEEE(buf.Bytes()), bytes.Equal(back, big), err)
	}
}
