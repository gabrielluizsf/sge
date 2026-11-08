package human

// baseline is a person made before looks could be recoloured.
type baseline struct {
	name string
	body Body
	look Look
}

var baselines = []baseline{
	{"man", Body{Gender: 1, Age: 0.52, African: 0.4, European: 0.6, Muscle: 0.6, Weight: 0.5, Height: 0.55, Proportions: 0.5},
		Look{Eyes: "eyes/materials/brownlight.mhmat", Wear: []string{"eyes/high-poly/high-poly.mhclo", "eyebrows/eyebrow001/eyebrow001.mhclo", "eyelashes/eyelashes01/eyelashes01.mhclo", "hair/short02/short02.mhclo", "clothes/male_worksuit01/male_worksuit01.mhclo", "clothes/shoes02/shoes02.mhclo"}}},
	{"old woman", Body{Gender: 0, Age: 0.9, European: 1, Muscle: 0.4, Weight: 0.6, Height: 0.4, Proportions: 0.5},
		Look{Skin: "skins/old_caucasian_female/old_caucasian_female.mhmat", Eyes: "eyes/materials/grey.mhmat", Wear: []string{"eyes/high-poly/high-poly.mhclo", "eyebrows/eyebrow006/eyebrow006.mhclo", "hair/toigo_curled_under_bob/toigo_curled_under_bob.mhclo", "clothes/toigo_shift_dress/toigo_shift_dress.mhclo"}}},
}
