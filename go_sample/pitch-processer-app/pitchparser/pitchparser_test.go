package pitchparser

import (
	"testing"
	"pitch-processer-app/testutils"
)

func newTestParser(t *testing.T) PitchFileParser {
	pp, err := NewPitchParser("testdata/pitchFileFooExchange.json")
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	return pp
}

func TestGetEvent(t *testing.T) {
	pp := newTestParser(t)

	tests := []struct {
		name     string
		line     string
		expected EventType
	}{
		{"Add", testutils.Add("1234567890", "AAPL", 123.45, 100), AddOrder},
		{"Modify", testutils.Modify("1234567890", 200.00, 200), ModifyOrder},
		{"Cancel", testutils.Cancel("1234567890"), CancelOrder},
		{"Execute", testutils.Execute("1234567890", 50), ExecuteOrder},
		{"Trade", testutils.Trade("AAPL", 100.00, 50), Trade},
		{"Unknown", "X" + testutils.PadRight("1234567890", 10), UnknownEvent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			event, err := pp.GetEvent(tt.line)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if event != tt.expected {
				t.Fatalf("expected %v got %v", tt.expected, event)
			}
		})
	}
}

func TestAddOrderDetails(t *testing.T) {
	pp := newTestParser(t)

	tests := []struct {
		name      string
		line      string
		expectErr bool
		orderID   string
		symbol    string
		price     float32
		size      int
	}{
		{
			"Valid",
			testutils.Add("1234567890", "AAPL", 123.45, 100),
			false,
			"1234567890",
			"AAPL",
			123.45,
			100,
		},
		{
			"ShortLine",
			"A123",
			true,
			"",
			"",
			0,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			details, err := pp.GetAddOrderDetails(tt.line)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if details.OrderID != tt.orderID {
				t.Errorf("expected orderID %s got %s", tt.orderID, details.OrderID)
			}

			if details.Symbol != tt.symbol {
				t.Errorf("expected symbol %s got %s", tt.symbol, details.Symbol)
			}

			if details.Price != tt.price {
				t.Errorf("expected price %f got %f", tt.price, details.Price)
			}

			if details.Size != tt.size {
				t.Errorf("expected size %d got %d", tt.size, details.Size)
			}
		})
	}
}

func TestModifyOrderDetails(t *testing.T) {
	pp := newTestParser(t)

	tests := []struct {
		name    string
		line    string
		orderID string
		price   float32
		size    int
	}{
		{
			"Valid",
			testutils.Modify("1234567890", 200.00, 200),
			"1234567890",
			200.00,
			200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			details, err := pp.GetModifyOrderDetails(tt.line)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if details.OrderID != tt.orderID {
				t.Errorf("expected orderID %s got %s", tt.orderID, details.OrderID)
			}

			if details.Price != tt.price {
				t.Errorf("expected price %f got %f", tt.price, details.Price)
			}

			if details.Size != tt.size {
				t.Errorf("expected size %d got %d", tt.size, details.Size)
			}
		})
	}
}

func TestExecutionDetails(t *testing.T) {
	pp := newTestParser(t)

	line := testutils.Execute("1234567890", 50)

	details, err := pp.GetExecutionDetails(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if details.OrderID != "1234567890" {
		t.Errorf("unexpected orderID")
	}

	if details.Size != 50 {
		t.Errorf("unexpected size")
	}
}

func TestTradeDetails(t *testing.T) {
	pp := newTestParser(t)

	line := testutils.Trade("AAPL", 100.00, 50)

	details, err := pp.GetTradeDetails(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if details.Symbol != "AAPL" {
		t.Errorf("unexpected symbol")
	}

	if details.Price != 100.00 {
		t.Errorf("unexpected price")
	}

	if details.Size != 50 {
		t.Errorf("unexpected size")
	}
}