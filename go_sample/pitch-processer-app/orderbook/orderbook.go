package orderbook

import (
	"errors"
	"sort"
)

type Order struct {
	ID     string
	Symbol string
	Price  float32
	Size   int
}

type OrderBook struct {
	Book                    map[string] *Order
	VolumeTradedBySymbol map[string]float32
}

type SymbolVolume struct {
	Symbol          string
	VolumeTraded float32
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

func (ob  *OrderBook) ModifyOrder(orderId string, newSize int, newPrice float32) error {
	if newPrice <= 0 {
		return errors.New("price update must be greater than zero")
	}
	if newSize <= 0 {
		return errors.New("size update must be greater than zero")
	}
	o, err := ob.getOrder(orderId)
	if err != nil {
		return err
	}
	o.Price = newPrice
	o.Size = newSize
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

	if o.Size - size < 0 {
		return errors.New("Execution amount is greater than order volume")
	}

	o.Size -= size
	ob.VolumeTradedBySymbol[o.Symbol] += float32(size) * o.Price

	if o.Size == 0 {
		delete(ob.Book, orderId)
	}

	return nil
}

func (ob *OrderBook) HandleTrade(symbol string, size int, price float32) error {
	if size <= 0 {
		return errors.New("trade size must be greater than zero")
	}
	ob.VolumeTradedBySymbol[symbol] += float32(size) * price
	return nil
}

func (ob OrderBook) TopTradedSymbols() []SymbolVolume {
	var items []SymbolVolume

	for k, v := range ob.VolumeTradedBySymbol {
		items = append(items, SymbolVolume{k, v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].VolumeTraded > items[j].VolumeTraded
	})

	return items
}
