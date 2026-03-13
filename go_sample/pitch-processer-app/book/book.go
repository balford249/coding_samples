package book

import (
	"sort"
	"pitch-processer-app/order"
)

type OrderBook struct {
	Book     map[string]*order.Order
	QuantatiyTradedBySymbol map[string]int
}

type SymbolVolume struct {
    	Symbol   string
    	QuantatiyTraded int
}

func (ob *OrderBook) AddOrder(order order.Order) {
	// Check it already doesn't exist 
	ob.Book[order.ID] = &order
}

func (ob *OrderBook) RemoveOrder(order order.Order) {
	// Check it already hasnt been deleted
	delete(ob.Book, order.ID)
}

func (ob *OrderBook) ExecuteOrder(orderId string, size int) {
	o := ob.Book[orderId]
	o.ReduceSize(size)
	ob.QuantatiyTradedBySymbol[o.Symbol] += size
}

func (ob OrderBook) TopTradedSymbols(noOfSymbols int) []SymbolVolume{
	var items []SymbolVolume
	
	for k, v := range ob.QuantatiyTradedBySymbol {
    	items = append(items, SymbolVolume{k, v})
	}
	
	sort.Slice(items, func(i, j int) bool {
    return items[i].QuantatiyTraded > items[j].QuantatiyTraded
	})

	var topSymbols []SymbolVolume
	for i := 0; i < noOfSymbols && i < len(items); i++ {
		topSymbols = append(topSymbols, items[i])
	}

	return topSymbols
}

