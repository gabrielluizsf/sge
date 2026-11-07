package art

// woolIcon is a ball of wool.
func woolIcon() string {
	d := newDoc(Size, Size)
	d.add(`<circle cx="8" cy="8.6" r="5.4" fill="%s" stroke="%s" stroke-width="0.8"/>`, "#f4f0e6", line)
	d.add(`<path d="M4 6.4 Q8 9 12 6.6 M3.4 9 Q8 12 12.6 9 M6 3.8 Q9 8 7 13.8" fill="none" stroke="%s" stroke-width="0.6"/>`, "#cfc7b6")
	return d.String()
}

// breadIcon is a golden loaf with slashes on its crust.
func breadIcon() string {
	d := newDoc(Size, Size)
	d.add(`<path d="M2.4 10 Q2 4.6 8 4.4 Q14 4.6 13.6 10 Q13.6 12.6 11 12.6 L5 12.6 Q2.4 12.6 2.4 10 Z" fill="%s" stroke="%s" stroke-width="0.8" stroke-linejoin="round"/>`, d.down("#efc06a", "#b9782e"), line)
	d.add(`<path d="M5.4 6.4 L6.6 9 M8 5.8 L9.2 8.6 M10.6 6.4 L11.6 8.8" stroke="%s" stroke-width="0.7" stroke-linecap="round"/>`, "#8a5a24")
	return d.String()
}

// meatIcon is a cut of meat of the colours flesh and fat, with a bone.
func meatIcon(flesh, fat string) string {
	d := newDoc(Size, Size)
	d.add(`<path d="M3 7 Q3.4 3 8 3.4 Q13 4 13 8.4 Q13 13 8.4 12.6 Q3 12.4 3 7 Z" fill="%s" stroke="%s" stroke-width="0.8" stroke-linejoin="round"/>`, fat, line)
	d.add(`<path d="M4.4 7.2 Q4.8 4.6 8 4.8 Q11.6 5.2 11.6 8.4 Q11.6 11.4 8.4 11.2 Q4.6 11 4.4 7.2 Z" fill="%s"/>`, flesh)
	d.add(`<circle cx="8" cy="8" r="1.5" fill="%s" stroke="%s" stroke-width="0.5"/>`, white, line)
	return d.String()
}
