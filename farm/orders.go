package farm

import "github.com/gabrielluizsf/sge/request"

// orderBook is what villagers ask for: wood and every crop, one or two
// of them, paying half as much again as the stall. The longer a crop
// takes to grow the fewer they ask for.
var orderBook = request.Book{
	Offers: append([]request.Offer{{Item: wood, Most: 6, Price: woodPrice}}, cropOffers(4)...),
	Kinds:  2,
	Bonus:  1.5,
}

// cropOffers are every crop's offer, at the stall's price, up to most of
// the quick ones and one of the slowest.
func cropOffers(most int) []request.Offer {
	offers := make([]request.Offer, 0, len(crops))
	for _, c := range crops {
		offers = append(offers, request.Offer{Item: c.name, Most: max(1, most-c.daysToRipe/4), Price: c.sellPrice})
	}
	return offers
}

// errandBook is what the smith asks for in exchange for the stone axe:
// the farm's crops, before there is an axe to cut wood with.
var errandBook = request.Book{
	Offers: []request.Offer{{Item: "carrot", Most: 4, Price: crops[0].sellPrice}, {Item: "lettuce", Most: 2, Price: crops[1].sellPrice}},
	Kinds:  2,
}

// The seconds a villager waits, after an order is filled, before asking
// for something else: at least orderPause, and up to as long again.
const orderPause = 25
