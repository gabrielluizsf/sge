package main

import "github.com/gabrielluizsf/sge/request"

// orderBook is what villagers ask for: wood, carrots and tomatoes, one or
// two of them, paying half as much again as the stall.
var orderBook = request.Book{
	Offers: []request.Offer{{Item: wood, Most: 6, Price: woodPrice}, {Item: "carrot", Most: 4, Price: 12}, {Item: "tomato", Most: 3, Price: 25}},
	Kinds:  2,
	Bonus:  1.5,
}

// errandBook is what the smith asks for in exchange for the stone axe:
// the farm's crops, before there is an axe to cut wood with.
var errandBook = request.Book{
	Offers: []request.Offer{{Item: "carrot", Most: 4, Price: 12}, {Item: "tomato", Most: 2, Price: 25}},
	Kinds:  2,
}

// The seconds a villager waits, after an order is filled, before asking
// for something else: at least orderPause, and up to as long again.
const orderPause = 25
