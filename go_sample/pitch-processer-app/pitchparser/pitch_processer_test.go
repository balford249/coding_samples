package pitchprocesser

import "testing"

func newTestProcessor() PitchProcesser {

	parser, err := LoadParserConfig("testdata/pitchFileTypeA.json")
	if err != nil {
		panic(err)
	}

	return PitchProcesser{
		fileParser: parser,
	}
}

func TestGetEvent_AddOrder(t *testing.T) {
	pp := newTestProcessor()

	line := "A1234567890AAPL00012345000100"

	event := pp.getEvent(line)

	if event != AddOrder {
		t.Fatalf("expected AddOrder, got %v", event)
	}
}

func TestGetAddOrderDetails(t *testing.T) {
	pp := newTestProcessor()

	line := "A1234567890AAPL00012345000100"

	details := pp.getAddOrderDetails(line)

	if details.OrderId != "1234567890" {
		t.Errorf("expected orderId 1234567890, got %s", details.OrderId)
	}

	if details.Symbol != "AAPL" {
		t.Errorf("expected symbol AAPL, got %s", details.Symbol)
	}

	if details.Price != 12345 {
		t.Errorf("expected price 12345, got %f", details.Price)
	}

	if details.Size != 100 {
		t.Errorf("expected size 100, got %d", details.Size)
	}
}

func TestGetModifyOrderDetails(t *testing.T) {
	pp := newTestProcessor()

	line := "M1234567890AAPL00020000000200"

	details := pp.getModifyOrderDetails(line)

	if details.OrderId != "1234567890" {
		t.Errorf("expected orderId 1234567890, got %s", details.OrderId)
	}

	if details.Price != 20000 {
		t.Errorf("expected price 20000, got %f", details.Price)
	}

	if details.Size != 200 {
		t.Errorf("expected size 200, got %d", details.Size)
	}
}

func TestGetCancelOrderDetails(t *testing.T) {
	pp := newTestProcessor()

	line := "C1234567890AAPL00000000000000"

	details := pp.getCancelOrderDetails(line)

	if details.OrderId != "1234567890" {
		t.Errorf("expected orderId 1234567890, got %s", details.OrderId)
	}
}

func TestGetTradeDetails(t *testing.T) {
	pp := newTestProcessor()

	line := "T1234567890AAPL00010000000050"

	details := pp.getTradeDetails(line)

	if details.Symbol != "AAPL" {
		t.Errorf("expected symbol AAPL, got %s", details.Symbol)
	}

	if details.Price != 10 {
		t.Errorf("expected price 10000, got %f", details.Price)
	}

	if details.Size != 50 {
		t.Errorf("expected size 50, got %d", details.Size)
	}
}
