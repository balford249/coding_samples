package orderbook

import "errors"

// attributes with lower case are private
type Order struct {
	ID     string
	Symbol string
	Price  float32
	Size   int
}

func (o *Order) ModifyOrder(newSize int, newPrice float32) error {
	if newPrice <= 0 {
		return errors.New("size update must be greater than zero")
	}
	o.Price = newPrice
	o.Size = newSize
	return nil
}
