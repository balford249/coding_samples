package main

import (
	"os"
	"pitch-processer/testutils"
	"strings"
	"testing"
)

func TestProcessPitchFile_SimpleScenario(t *testing.T) {

	pitchData := testutils.Lines(
		// New order for AAPL with a price of 5.5 and size of 10
		testutils.Add("1", "AAPL", 5.5, 10),
		// New order for MSFT with a price of 10.1 and size of 8
		testutils.Add("2", "MSFT", 10.1, 8),
		// Execute 5 of AAPL order so should be 5 * 5.5
		testutils.Execute("1", 5),
		// Modify MSFT order with a new price of 10.5 
		testutils.Modify("2", 10.5, 8),
		// Execute MSFT order with the updated price of 10.5 and a size of 4
		testutils.Execute("2", 4),
		// Cancel MSFT order
		testutils.Cancel("2"),
		// Trade VOD for 10 at price 10.6. Trades can occur without orders as orders can be hidden. 
		testutils.Trade("VOD", 10.50, 10),
	)

	file, err := os.CreateTemp("", "pitch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(pitchData)
	file.Close()

	result, err := processPitchFile(file.Name(), "pitchparser/testdata/pitchFileFooExchange.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for index, item := range [][]interface{}{
		{"VOD", float32(105)},
		{"MSFT", float32(42)},
		{"AAPL", float32(27.5)},
	} {
		symbol := item[0].(string)
		volume := item[1].(float32)

		if strings.TrimSpace(result[index].Symbol) != symbol {
			t.Errorf("expected  %s got %s", symbol, result[index].Symbol)
		}

		if result[index].VolumeTraded != volume {
			t.Errorf("expected  %f got %f", volume, result[index].VolumeTraded)
		}
	}
}
