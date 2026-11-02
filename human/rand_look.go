package human

// Random is a whole person, the same for the same seed: a body natural
// for the sex, age and ancestry asked, and a look to match: skin, eyes
// and hair of the ancestry, a hairstyle, clothes and shoes for sex and
// age, now and then a hat. What spec leaves zero is chosen by chance.
func Random(s Spec) (Body, Look) {
	r := newDraw(s.Seed)
	sex := s.Sex
	if sex == AnySex {
		sex = Woman + Sex(r.Intn(2))
	}
	a := ancestryOf(s.Ancestry, r)
	years := yearsOf(s.Age, r)
	b := bodyOf(sex, years, a, r)
	years = b.Years()
	l := Look{SkinTone: skinToneOf(a, r), HairColour: hairColourOf(a, r)}
	l.Eyes = eyesOf(a, r)
	l.Wear = append(faceOf(sex, r), hairOf(sex, years, a.African, r)...)
	l.Wear = append(l.Wear, outfitOf(sex, years, r)...)
	l.ClothesColours = colourClothes(l.Wear, r)
	return b, l
}
