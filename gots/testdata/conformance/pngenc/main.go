package main

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
)

func main() {
	img := image.NewNRGBA(image.Rect(0, 0, 480, 320))
	for y := 0; y < 320; y++ {
		for x := 0; x < 480; x++ {
			img.Set(x, y, color.NRGBA{uint8(x), uint8(y), uint8(x ^ y), 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("png: %d bytes crc %08x\n", buf.Len(), crc32.ChecksumIEEE(buf.Bytes()))
	back, err := png.Decode(bytes.NewReader(buf.Bytes()))
	fmt.Println(err, back.Bounds(), back.At(100, 50))
}
