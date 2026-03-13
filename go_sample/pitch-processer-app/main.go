package main

import (
	"fmt"
	"pitch-processer-app/order"
	"pitch-processer-app/book"
)


func main() {
	ob := book.OrderBook{
		Book:                    make(map[string]*order.Order),
		QuantatiyTradedBySymbol: make(map[string]int),
	}
	o := order.Order{ID: "abc", Price: 1.0, Size: 5, Symbol: "VOD"}
	ob.AddOrder(o)
	ob.ExecuteOrder(o.ID, 1)

	o1 := order.Order{ID: "abcd", Price: 1.0, Size: 5, Symbol: "VOD1"}
	ob.AddOrder(o1)
	ob.ExecuteOrder(o1.ID, 2)
	fmt.Println(ob.TopTradedSymbols(3))


}
