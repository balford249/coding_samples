package order

// attributes with lower case are private 
type Order struct {
	ID     string
	Symbol string
	Price  float32
	Size   int
}

// * is a Pointer so will update the original struct, non * makes a copy to the reference
func (o *Order) UpdateSize(newSize int) {
	o.Size = newSize
}

func (o *Order) ReduceSize(executionAmount int) {
	o.Size -= executionAmount
}
