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
}
