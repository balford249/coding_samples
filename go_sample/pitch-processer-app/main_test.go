package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

func padRight(s string, width int) string {
	if len(s) > width {
		return s[:width]
	}
	return fmt.Sprintf("%-*s", width, s)
}

func padInt(n int, width int) string {
	return fmt.Sprintf("%0*d", width, n)
}

func padFloat32(f float32, width int) string {
	// 10.50 --> 00001050
	n := int(math.Round(float64(f) * 100))
	return fmt.Sprintf("%0*d", width, n)
}

func Add(orderId string, symbol string, price float32, size int) string {
	return "A" +
		padRight(orderId, 10) +
		padRight(symbol, 4) +
		padFloat32(price, 8) +
		padInt(size, 6)
}

func Modify(orderId string, price float32, size int) string {
	return "M" +
		padRight(orderId, 10) +
		padRight("", 4) +
		padFloat32(price, 8) +
		padInt(size, 6)
}

func Execute(orderId string, size int) string {
	return "E" +
		padRight(orderId, 10) +
		padRight("", 4) +
		padInt(0, 8) +
		padInt(size, 6)
}

func Cancel(orderId string) string {
	return "C" + padRight(orderId, 10)
}

func Trade(symbol string, price float32, size int) string {
	return "T" +
		padRight("", 10) + // orderId unused
		padRight(symbol, 4) +
		padFloat32(price, 8) +
		padInt(size, 6)
}

func Lines(events ...string) string {
	return strings.Join(events, "\n")
}

func TestProcessPitchFile_SimpleScenario(t *testing.T) {

	pitchData := Lines(
		Add("1", "AAPL", 5.5, 10),
		Add("2", "MSFT", 10.1, 8),
		Execute("1", 5),
		Modify("2", 10.5, 8),
		Execute("2", 4),
		Cancel("2"),
		Trade("VOD", 10.50, 10),
	)

	file, err := os.CreateTemp("", "pitch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(pitchData)
	file.Close()

	result := processPitchFile(file.Name(), "pitchparser/testdata/pitchFileTypeA.json")

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
