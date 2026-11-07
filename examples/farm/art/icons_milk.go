package art

// milkIcon is a corked bottle of milk.
func milkIcon() string {
	d := newDoc(Size, Size)
	d.add(`<path d="M6.6 3.6 L9.4 3.6 L9.4 5.8 Q11.4 7 11.4 9 L11.4 13.2 Q11.4 14.4 10.2 14.4 L5.8 14.4 Q4.6 14.4 4.6 13.2 L4.6 9 Q4.6 7 6.6 5.8 Z" fill="%s" stroke="%s" stroke-width="0.8" stroke-linejoin="round"/>`, white, line)
	d.add(`<rect x="6.4" y="1.8" width="3.2" height="2.2" rx="0.5" fill="%s" stroke="%s" stroke-width="0.6"/>`, wood3[1], line)
	d.add(`<path d="M6 9 L6 12.6" stroke="%s" stroke-width="0.8" stroke-linecap="round" opacity="0.5"/>`, shirtShade)
	return d.String()
}
