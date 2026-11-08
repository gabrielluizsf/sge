package human

import "math"

// Colour is a colour a person's hair, eyes or clothes may be, in sRGB,
// each of R, G and B from 0 to 1. The zero Colour, Set false, asks for
// nothing: the pack's own colour stays.
type Colour struct {
	R, G, B float64
	Set     bool
}

// RGB is the Colour of r, g and b, sRGB from 0 to 1.
func RGB(r, g, b float64) Colour { return Colour{r, g, b, true} }

// linear is c in linear light, as lighting adds and multiplies it.
func (c Colour) linear() [3]float64 {
	return [3]float64{toLinear(c.R), toLinear(c.G), toLinear(c.B)}
}

// toLinear is the sRGB value v, 0 to 1, in linear light.
func toLinear(v float64) float64 {
	v = min(max(v, 0), 1)
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// fromLinear is the linear light v, 0 to 1, as an sRGB value.
func fromLinear(v float64) float64 {
	v = min(max(v, 0), 1)
	if v <= 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1/2.4) - 0.055
}

// encodeSteps is how finely linear light is looked up on its way back to
// 8 bits: finer than the 8 bits themselves are, even in the darks.
const encodeSteps = 1 << 15

// decode8 is each 8-bit sRGB value in linear light; encode8 is linear
// light, in encodeSteps, as an 8-bit sRGB value.
var decode8, encode8 = lightTables()

func lightTables() (dec [256]float64, enc [encodeSteps + 1]uint8) {
	for i := range dec {
		dec[i] = toLinear(float64(i) / 255)
	}
	for i := range enc {
		enc[i] = uint8(fromLinear(float64(i)/encodeSteps)*255 + 0.5)
	}
	return
}

// luma is the light of linear colour c.
func luma(c [3]float64) float64 { return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2] }
