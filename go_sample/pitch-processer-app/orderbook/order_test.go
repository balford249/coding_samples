package orderbook

import "testing"


func TestModifyOrder_Success(t *testing.T) {
	order := NewStandardTestOrder()

	err := order.ModifyOrder(20, 150.0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.Size != 20 {
		t.Fatalf("expected size 20, got %d", order.Size)
	}

	if order.Price != 150.0 {
		t.Fatalf("expected price 150.0, got %f", order.Price)
	}
}

func TestModifyOrder_InvalidPrice(t *testing.T) {
	order :=NewStandardTestOrder()

	err := order.ModifyOrder(20, 0)

	if err == nil {
		t.Fatalf("expected error when price <= 0")
	}
}
