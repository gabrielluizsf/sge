package human

// hair is the proxy of the hair folder named name.
func hair(name string) string { return "hair/" + name + "/" + name + ".mhclo" }

// hairStyles are the hair of MakeHuman's CC0 packs: short for men and
// boys, long, bobbed, braided and tied for women and girls, and the afro
// for those of African ancestry. The pack's two first bobs and its braid are
// gone: their fringes lay over an eye.
var hairStyles = []garment{
	forMen(hair("short01"), 3, 90),
	forMen(hair("short02"), 3, 90),
	forMen(hair("short03"), 3, 90),
	forMen(hair("short04"), 3, 90),
	forWomen(hair("toigo_blunt_bob"), 3, 90),
	forWomen(hair("toigo_blunt_bob_with_bangs"), 3, 90),
	forWomen(hair("toigo_curled_under_bob"), 3, 90),
	forWomen(hair("toigo_curled_under_bob_with_bangs"), 3, 90),
	forWomen(hair("toigo_inverted_bob"), 3, 90),
	forWomen(hair("toigo_inverted_bob_with_bangs"), 3, 90),
	forBoth(hair("faydaen_hair_1"), 12, 90),
	forWomen(hair("long01"), 3, 90),
	forWomen(hair("ponytail01"), 3, 90),
	{item: hair("afro01"), women: true, men: true, from: 3, to: 90, african: 0.35},
}
