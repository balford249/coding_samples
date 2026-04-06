package orderbook

import (
	"testing"
)

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

func TestAddOrder(t *testing.T) {
	ob := NewOrderBook()
	o := NewStandardTestOrder()

	err := ob.AddOrder(&o)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := ob.Book["1"]; !exists {
		t.Fatalf("expected order to be added")
	}
}

func TestAddOrder_Duplicate(t *testing.T) {
	ob := NewOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(&o)
	err := ob.AddOrder(&o)

	if err == nil {
		t.Fatalf("expected duplicate order error")
	}
}

func TestModifyOrderSuccess(t *testing.T) {

	ob := NewOrderBook()

	order := Order{
		ID:     "1",
		Symbol: "AAPL",
		Price:  100,
		Size:   10,
	}

	ob.AddOrder(&order)

	err := ob.ModifyOrder("1", 20, 150)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	o, exists := ob.Book["1"]
	if !exists {
		t.Fatalf("order should exist in book")
	}

	if o.Size != 20 {
		t.Errorf("expected size 20, got %d", o.Size)
	}

	if o.Price != 150 {
		t.Errorf("expected price 150, got %f", o.Price)
	}
}

func TestModifyOrderOrderNotFound(t *testing.T) {

	ob := NewOrderBook()

	err := ob.ModifyOrder("999", 10, 100)

	if err == nil {
		t.Fatalf("expected error for missing order")
	}
}

func TestModifyOrder_InvalidInputs(t *testing.T) {

	tests := []struct {
		size  int
		price float32
	}{
		{0, 100},
		{10, 0},
	}

	for _, tt := range tests {
		ob := NewOrderBook()
		o := NewStandardTestOrder()

		ob.AddOrder(&o)

		err := ob.ModifyOrder("1", tt.size, tt.price)

		if err == nil {
			t.Fatalf("expected error")
		}
	}
}

func TestRemoveOrder(t *testing.T) {
	ob := NewOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(&o)
	err := ob.RemoveOrder("1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := ob.Book["1"]; exists {
		t.Fatalf("expected order to be removed")
	}
}

func TestRemoveOrder_NotFound(t *testing.T) {
	ob := NewOrderBook()

	err := ob.RemoveOrder("missing")

	if err == nil {
		t.Fatalf("expected order not found error")
	}
}

func TestExecuteOrder_PartialFill(t *testing.T) {
	ob := NewOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(&o)
	err := ob.ExecuteOrder("1", 4)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ob.Book["1"].Size != 6 {
		t.Fatalf("expected size 6 got %d", ob.Book["1"].Size)
	}

	if ob.VolumeTradedBySymbol["AAPL"] != 400 {
		t.Fatalf("expected traded value 4 got %f", ob.VolumeTradedBySymbol["AAPL"])
	}
}

func TestExecuteOrder_FullFill(t *testing.T) {
	ob := NewOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(&o)
	err := ob.ExecuteOrder("1", 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := ob.Book["1"]; exists {
		t.Fatalf("expected order removed when fully filled")
	}

	if ob.VolumeTradedBySymbol["AAPL"] != 1000 {
		t.Fatalf("expected traded quantity 10 got %f", ob.VolumeTradedBySymbol["AAPL"])
	}
}

func TestExecuteOrder_NotFound(t *testing.T) {
	ob := NewOrderBook()

	err := ob.ExecuteOrder("missing", 5)

	if err == nil {
		t.Fatalf("expected order not found error")
	}
}

func TestExecuteOrder_InvalidInputs(t *testing.T) {

	tests := []int{10000, 0}

	for _, tt := range tests {
		ob := NewOrderBook()
		o := NewStandardTestOrder()

		ob.AddOrder(&o)

		err := ob.ExecuteOrder("1", tt)

		if err == nil {
			t.Fatalf("expected error")
		}
	}
}

func TestHandleTrade(t *testing.T) {
	ob := NewOrderBook()

	err := ob.HandleTrade("AAPL", 20, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ob.VolumeTradedBySymbol["AAPL"] != 200 {
		t.Fatalf("expected traded quantity 20 got %f", ob.VolumeTradedBySymbol["AAPL"])
	}
}

func TestHandleTrade_InvalidSize(t *testing.T) {
	ob := NewOrderBook()

	err := ob.HandleTrade("AAPL", 0, 0)

	if err == nil {
		t.Fatalf("expected error for invalid trade size")
	}
}

func TestSymbolVolumes(t *testing.T) {
	ob := NewOrderBook()

	ob.VolumeTradedBySymbol["AAPL"] = 100
	ob.VolumeTradedBySymbol["MSFT"] = 50
	ob.VolumeTradedBySymbol["GOOG"] = 200

	top := ob.SymbolVolumes()

	if len(top) != 3 {
		t.Fatalf("expected 3 results got %d", len(top))
	}

	if top[0].Symbol != "GOOG" {
		t.Fatalf("expected GOOG first got %s", top[0].Symbol)
	}

	if top[1].Symbol != "AAPL" {
		t.Fatalf("expected AAPL second got %s", top[1].Symbol)
	}
}
