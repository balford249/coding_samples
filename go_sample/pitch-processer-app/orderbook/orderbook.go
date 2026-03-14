package orderbook

import (
	"errors"
	"sort"
)

type OrderBook struct {
	Book                    map[string] *Order
	QuantatiyTradedBySymbol map[string]int
}

type SymbolVolume struct {
	Symbol          string
	QuantatiyTraded int
}

func (ob *OrderBook) getOrder(orderID string) (*Order, error) {
	o, exists := ob.Book[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (ob *OrderBook) AddOrder(order Order) error {

	if _, exists := ob.Book[order.ID]; exists {
		return errors.New("order already exists")
	}

	ob.Book[order.ID] = &order
	return nil
}

func (ob *OrderBook) RemoveOrder(orderId string) error {

	_, err := ob.getOrder(orderId)
	if err != nil {
		return err
	}

	delete(ob.Book, orderId)
	return nil
}

func (ob *OrderBook) ExecuteOrder(orderId string, size int) error {

	o, err := ob.getOrder(orderId)
	if err != nil {
		return err
	}

	o.Size -= size
	ob.QuantatiyTradedBySymbol[o.Symbol] += size

	if o.Size == 0 {
		delete(ob.Book, orderId)
	}

	return nil
}

func (ob *OrderBook) HandleTrade(symbol string, size int) error {
	if size <= 0 {
		return errors.New("trade size must be greater than zero")
	}
	ob.QuantatiyTradedBySymbol[symbol] += size
	return nil
}

func (ob OrderBook) TopTradedSymbols(noOfSymbols int) []SymbolVolume {
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
