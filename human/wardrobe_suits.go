package human

// suits are whole outfits, top and bottom in one.
var suits = []garment{
	forMen(cloth("male_casualsuit01"), 6, 90),
	forMen(cloth("male_casualsuit02"), 6, 90),
	forMen(cloth("male_casualsuit03"), 6, 90),
	forMen(cloth("male_casualsuit04"), 6, 90),
	forMen(cloth("male_casualsuit05"), 6, 90),
	forMen(cloth("male_casualsuit06"), 6, 90),
	forMen(cloth("male_worksuit01"), 14, 75),
	forMen(cloth("male_elegantsuit01"), 14, 90),
	forWomen(cloth("female_casualsuit01"), 6, 90),
	forWomen(cloth("female_casualsuit02"), 6, 90),
	forWomen(cloth("female_sportsuit01"), 12, 60),
	forWomen(cloth("female_elegantsuit01"), 14, 90),
	forBoth(cloth("toigo_suit_with_dinner_jacket"), 16, 90),
	forBoth(cloth("toigo_suit_with_jacket_and_bowtie"), 16, 90),
	forMen(cloth("toigo_male_double-breasted_suit"), 16, 90),
	forMen(cloth("toigo_male_suit_tie_and_jacket"), 16, 90),
	forMen(cloth("toigo_male_suit_3"), 16, 90),
	forWomen(cloth("toigo_female_suit"), 16, 90),
	forWomen(cloth("toigo_female_double-breasted_suit"), 16, 90),
}
