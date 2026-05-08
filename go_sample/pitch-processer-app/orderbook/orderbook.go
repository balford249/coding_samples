package processer

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

type SymbolVolume struct {
	Symbol       string
	VolumeTraded float32
}

type OrderBook struct {
	Book                 map[string]*Order
	VolumeTradedBySymbol map[string]float32
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Book:                 make(map[string]*Order),
		VolumeTradedBySymbol: make(map[string]float32),
	}
}

func (ob *OrderBook) getOrder(orderID string) (*Order, error) {
	o, exists := ob.Book[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (ob *OrderBook) AddOrder(order *Order) error {

	if order.ID == "" {
		return errors.New("order id required")
	}

	if order.Price <= 0 {
		return errors.New("price must be positive")
	}

	if order.Size <= 0 {
		return errors.New("size must be positive")
	}

	if _, exists := ob.Book[order.ID]; exists {
		return errors.New("order already exists")
	}

	ob.Book[order.ID] = order
	return nil
}

func (ob *OrderBook) ModifyOrder(orderID string, newSize int, newPrice float32) error {
	if newPrice <= 0 {
		return errors.New("price update must be greater than zero")
	}
	if newSize <= 0 {
		return errors.New("size update must be greater than zero")
	}
	o, err := ob.getOrder(orderID)
	if err != nil {
		return err
	}
	o.Price = newPrice
	o.Size = newSize
	return nil
}

func (ob *OrderBook) RemoveOrder(orderID string) error {

	_, err := ob.getOrder(orderID)
	if err != nil {
		return err
	}

	delete(ob.Book, orderID)
	return nil
}

func (ob *OrderBook) ExecuteOrder(orderID string, size int) error {

	o, err := ob.getOrder(orderID)
	if err != nil {
		return err
	}

	if size <= 0 {
		return errors.New("execution size must be greater than zero")
	}

	if size > o.Size {
		return errors.New("execution amount greater than order size")
	}

	o.Size -= size
	ob.VolumeTradedBySymbol[o.Symbol] += float32(size) * o.Price

	if o.Size == 0 {
		delete(ob.Book, orderID)
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

func (ob *OrderBook) SymbolVolumes() []SymbolVolume {
	var symbolVolumes []SymbolVolume

	for k, v := range ob.VolumeTradedBySymbol {
		symbolVolumes = append(symbolVolumes, SymbolVolume{
			Symbol:       k,
			VolumeTraded: v,
		})
	}

	sort.Slice(symbolVolumes, func(i, j int) bool {
		return symbolVolumes[i].VolumeTraded > symbolVolumes[j].VolumeTraded
	})

	return symbolVolumes
}
