// Package request is what characters ask for: an Order of so many of
// some items, for a reward, drawn at random from a Book of what may be
// asked, and filled from a bag.
package request

import "github.com/gabrielluizsf/sge/inventory"

// Offer is an item that may be asked for: up to Most of it, each worth
// Price.
type Offer struct {
	Item  string
	Most  int
	Price int
}

// Book is what orders are drawn from: its offers, how many different
// items an order asks for at most, and how much more than their price an
// order pays, 1.5 for half as much again.
type Book struct {
	Offers []Offer
	Kinds  int
	Bonus  float64
}

// Order is a wish: the items wanted, and what fulfilling it pays.
type Order struct {
	Want   inventory.List
	Reward int
}
