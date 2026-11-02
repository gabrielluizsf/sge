package main

// framesHead reports whether the category called name changes the face,
// which the camera then comes close to; the others change the body or
// what it wears, which it shows whole.
func framesHead(name string) bool {
	switch name {
	case "Pele", "Rosto", "Cabelo", "Olhos":
		return true
	}
	return false
}
