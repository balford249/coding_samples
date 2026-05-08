package processer

import (
	"testing"

	orderbook "pitch-processer-app/orderbook"
	pitchparser "pitch-processer-app/pitchparser"
	testutils "pitch-processer-app/testutils"
)

func newTestEnv(t *testing.T) (*orderbook.OrderBook, pitchparser.PitchFileParser) {
	parser, err := pitchparser.NewPitchParser("../pitchparser/testdata/pitchFileFooExchange.json")
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	ob := orderbook.NewOrderBook()
	return ob, parser
}

func TestAddOrder(t *testing.T) {
	ob, parser := newTestEnv(t)

	line := testutils.Add("1234567890", "AAPL", 123.45, 100)
	if err := HandleAddOrder(ob, parser, line); err != nil {
		t.Fatalf("handleAddOrder failed: %v", err)
	}

	order := ob.Book["1234567890"]

	if order.Symbol != "AAPL" {
		t.Errorf("expected symbol AAPL, got %s", order.Symbol)
	}
	if order.Price != 123.45 {
		t.Errorf("expected price 123.45, got %f", order.Price)
	}
	if order.Size != 100 {
		t.Errorf("expected size 100, got %d", order.Size)
	}
}

func TestModifyOrder(t *testing.T) {
	ob, parser := newTestEnv(t)

	ob.AddOrder(&orderbook.Order{
		ID: "1234567890", Symbol: "AAPL", Price: 100, Size: 100,
	})

	line := testutils.Modify("1234567890", 200.00, 300)
	if err := HandleModifyOrder(ob, parser, line); err != nil {
		t.Fatalf("handleModifyOrder failed: %v", err)
	}

	order := ob.Book["1234567890"]

	if order.Price != 200.00 {
		t.Errorf("expected price 200, got %f", order.Price)
	}
	if order.Size != 300 {
		t.Errorf("expected size 300, got %d", order.Size)
	}
}

func TestCancelOrder(t *testing.T) {
	ob, parser := newTestEnv(t)

	ob.AddOrder(&orderbook.Order{
		ID: "1234567890", Symbol: "AAPL", Price: 100, Size: 100,
	})

	line := testutils.Cancel("1234567890")
	if err := HandleCancelOrder(ob, parser, line); err != nil {
		t.Fatalf("handleCancelOrder failed: %v", err)
	}

	//order := ob.Book["1234567890"]

}

func TestExecuteOrder(t *testing.T) {
	ob, parser := newTestEnv(t)

	ob.AddOrder(&orderbook.Order{
		ID: "1234567890", Symbol: "AAPL", Price: 10, Size: 100,
	})

	line := testutils.Execute("1234567890", 40)
	if err := HandleExecuteOrder(ob, parser, line); err != nil {
		t.Fatalf("handleExecuteOrder failed: %v", err)
	}

	order := ob.Book["1234567890"]
	if order.Size != 60 {
		t.Errorf("expected size 60, got %d", order.Size)
	}
}

func TestTrade(t *testing.T) {
	ob, parser := newTestEnv(t)

	line := testutils.Trade("AAPL", 100.00, 50)
	if err := HandleTrade(ob, parser, line); err != nil {
		t.Fatalf("handleTrade failed: %v", err)
	}

	vol := ob.SymbolVolumes()
	if len(vol) != 1 {
		t.Fatalf("expected one symbol volume")
	}
	if vol[0].Symbol != "AAPL" {
		t.Errorf("expected symbol AAPL, got %s", vol[0].Symbol)
	}
	if vol[0].VolumeTraded != 5000 {
		t.Errorf("expected volume 5000, got %f", vol[0].VolumeTraded)
	}
}
