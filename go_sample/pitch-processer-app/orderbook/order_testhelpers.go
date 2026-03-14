package orderbook



func NewTestOrder(id, symbol string, price float32, size int) Order {
	return Order{
		ID:     id,
		Symbol: symbol,
		Price:  price,
		Size:   size,
	}
}

func NewStandardTestOrder() Order {
	return NewTestOrder("1", "AAPL", 100, 10)
}