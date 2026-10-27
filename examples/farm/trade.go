package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/vector"
)

// buy buys a packet of crop c's seeds, if the money is there, and tells
// so above at.
func (f *farm) buy(c int, at vector.Vector3) bool {
	price := crops[c].seedPrice
	if f.money < price {
		f.float(at, "SEM DINHEIRO", 'r')
		f.sound("deny", nil)
		return false
	}
	f.money -= price
	f.seeds[c]++
	f.float(at, fmt.Sprint("-$", price), 'r')
	f.sound("coin", nil)
	return true
}

// sell sells the whole harvest, and tells what it made above at.
func (f *farm) sell(at vector.Vector3) bool {
	earned := 0
	for c, n := range f.harvest {
		earned += n * crops[c].sellPrice
		f.sold += n
		f.harvest[c] = 0
	}
	if earned == 0 {
		f.float(at, "NADA A VENDER", 'r')
		f.sound("deny", nil)
		return false
	}
	f.money += earned
	f.float(at, fmt.Sprint("+$", earned), 'y')
	f.sound("sell", nil)
	return true
}
