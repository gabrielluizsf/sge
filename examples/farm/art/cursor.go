package art

// Cursor marks the cell the farmer works: its diamond, IsoW by Size, in a
// glowing outline open in the middle so the tile shows.
func Cursor() string {
	d := newDoc(IsoW, Size)
	const pts = "16,1 31,8 16,15 1,8"
	d.add(`<polygon points="%s" fill="%s" fill-opacity="0.18" stroke="%s" stroke-width="2.2" stroke-opacity="0.45" stroke-linejoin="round"/>`, pts, yellow, yellow)
	d.add(`<polygon points="%s" fill="none" stroke="%s" stroke-width="0.8" stroke-linejoin="round"/>`, pts, white)
	return d.String()
}
