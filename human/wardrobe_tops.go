package human

// tops are shirts and sweaters, worn with a bottom.
var tops = []garment{
	forBoth(cloth("toigo_basic_tucked_t-shirt"), 1, 90),
	forBoth(cloth("toigo_fisherman_sweater"), 1, 90),
	forMen(cloth("namuhekam_male_polo_shirt"), 8, 90),
	forMen(cloth("elvs_crude_t-shirt_male"), 10, 90),
	forWomen(cloth("joepal_crude_t-shirt_female"), 6, 90),
	forWomen(cloth("toigo_keyhole_tank_top"), 12, 90),
	forWomen(cloth("toigo_bodice-style_top"), 12, 90),
	forWomen(cloth("toigo_camisole_top"), 16, 70),
	forWomen(cloth("toigo_turtleneck_halter_top"), 16, 70),
}
