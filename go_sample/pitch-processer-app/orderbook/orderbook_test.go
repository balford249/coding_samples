package orderbook

import (
	"testing"
)

func newTestOrderBook() OrderBook {
	return OrderBook{
		Book:                    make(map[string]*Order),
		QuantatiyTradedBySymbol: make(map[string]int),
	}
}


func TestAddOrder(t *testing.T) {
	ob := newTestOrderBook()
	o := NewStandardTestOrder()

	err := ob.AddOrder(o)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := ob.Book["1"]; !exists {
		t.Fatalf("expected order to be added")
	}
}

func TestAddOrder_Duplicate(t *testing.T) {
	ob := newTestOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(o)
	err := ob.AddOrder(o)

	if err == nil {
		t.Fatalf("expected duplicate order error")
	}
}

func TestRemoveOrder(t *testing.T) {
	ob := newTestOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(o)
	err := ob.RemoveOrder("1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := ob.Book["1"]; exists {
		t.Fatalf("expected order to be removed")
	}
}

func TestRemoveOrder_NotFound(t *testing.T) {
	ob := newTestOrderBook()

	err := ob.RemoveOrder("missing")

	if err == nil {
		t.Fatalf("expected order not found error")
	}
}

func TestExecuteOrder_PartialFill(t *testing.T) {
	ob := newTestOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(o)
	err := ob.ExecuteOrder("1", 4)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ob.Book["1"].Size != 6 {
		t.Fatalf("expected size 6 got %d", ob.Book["1"].Size)
	}

	if ob.QuantatiyTradedBySymbol["AAPL"] != 4 {
		t.Fatalf("expected traded quantity 4 got %d", ob.QuantatiyTradedBySymbol["AAPL"])
	}
}

func TestExecuteOrder_FullFill(t *testing.T) {
	ob := newTestOrderBook()
	o := NewStandardTestOrder()

	ob.AddOrder(o)
	err := ob.ExecuteOrder("1", 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := ob.Book["1"]; exists {
		t.Fatalf("expected order removed when fully filled")
	}

	if ob.QuantatiyTradedBySymbol["AAPL"] != 10 {
		t.Fatalf("expected traded quantity 10 got %d", ob.QuantatiyTradedBySymbol["AAPL"])
	}
}

func TestExecuteOrder_NotFound(t *testing.T) {
	ob := newTestOrderBook()

	err := ob.ExecuteOrder("missing", 5)

	if err == nil {
		t.Fatalf("expected order not found error")
	}
}

func TestHandleTrade(t *testing.T) {
	ob := newTestOrderBook()

	err := ob.HandleTrade("AAPL", 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ob.QuantatiyTradedBySymbol["AAPL"] != 20 {
		t.Fatalf("expected traded quantity 20 got %d", ob.QuantatiyTradedBySymbol["AAPL"])
	}
}

func TestHandleTrade_InvalidSize(t *testing.T) {
	ob := newTestOrderBook()

	err := ob.HandleTrade("AAPL", 0)

	if err == nil {
		t.Fatalf("expected error for invalid trade size")
	}
}

func TestTopTradedSymbols(t *testing.T) {
	ob := newTestOrderBook()

	ob.QuantatiyTradedBySymbol["AAPL"] = 100
	ob.QuantatiyTradedBySymbol["MSFT"] = 50
	ob.QuantatiyTradedBySymbol["GOOG"] = 200

	top := ob.TopTradedSymbols(2)

	if len(top) != 2 {
		t.Fatalf("expected 2 results got %d", len(top))
	}

	if top[0].Symbol != "GOOG" {
		t.Fatalf("expected GOOG first got %s", top[0].Symbol)
	}

	if top[1].Symbol != "AAPL" {
		t.Fatalf("expected AAPL second got %s", top[1].Symbol)
	}
}