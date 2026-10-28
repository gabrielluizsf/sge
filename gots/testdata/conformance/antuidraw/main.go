// Command antuidraw draws with antui's canvas and svg and prints checksums
// of the pixels and of their PNG encoding, which must not change between Go
// and the transpiled program.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	"image/png"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/svg"
)

const drawing = `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="80">
  <defs><linearGradient id="g"><stop offset="0" stop-color="#f40"/><stop offset="1" stop-color="#04f"/></linearGradient></defs>
  <rect x="5" y="5" width="110" height="70" rx="12" fill="url(#g)" stroke="#222" stroke-width="3"/>
  <circle cx="60" cy="40" r="22" fill="#fff" fill-opacity="0.6"/>
  <path d="M10 70 Q60 10 110 70" fill="none" stroke="#0a0" stroke-width="4" stroke-dasharray="6 3"/>
  <text x="18" y="46" font-size="16" fill="#000">sge</text>
</svg>`

func checksum(cv *canvas.Canvas) uint32 {
	buf := make([]byte, 4*len(cv.Pixels))
	for i, c := range cv.Pixels {
		binary.LittleEndian.PutUint32(buf[i*4:], uint32(c))
	}
	return crc32.ChecksumIEEE(buf)
}

func pngOf(cv *canvas.Canvas) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, cv.Width, cv.Height))
	for y := 0; y < cv.Height; y++ {
		for x := 0; x < cv.Width; x++ {
			c := cv.Pixels[y*cv.Stride+x]
			i := img.PixOffset(x, y)
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = c.R(), c.G(), c.B(), c.A()
		}
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	return b.Bytes()
}

func main() {
	cv, err := canvas.NewCanvas(160, 120)
	if err != nil {
		panic(err)
	}
	cv.Clear(canvas.RGB(30, 30, 40))
	cv.FillRect(canvas.Area{X: 10, Y: 10, Width: 60, Height: 30}, canvas.RGB(200, 60, 60))
	cv.FillCircle(canvas.Circle{X: 110, Y: 40, R: 25}, canvas.RGBA(canvas.RGBAParams{R: 60, G: 200, B: 90, A: 160}))
	cv.Circle(canvas.Circle{X: 50, Y: 85, R: 20}, canvas.RGB(250, 250, 250))
	cv.Text(canvas.Point{X: 12, Y: 100}, "Olá, gots!", canvas.RGB(255, 220, 0))
	cv.Blur(canvas.Area{X: 0, Y: 0, Width: 80, Height: 60}, 2)
	fmt.Printf("shapes: %08x\n", checksum(cv))

	doc, err := svg.Parse(drawing)
	if err != nil {
		panic(err)
	}
	pic := doc.Render(120, 80)
	fmt.Printf("svg: %08x\n", checksum(pic))
	fmt.Printf("png: %08x\n", crc32.ChecksumIEEE(pngOf(cv)))
}
