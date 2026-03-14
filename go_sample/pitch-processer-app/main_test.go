package main

import (
	"fmt"
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

func Add(orderId string, symbol string, price int, size int) string {
	return "A" +
		padRight(orderId, 10) +
		padRight(symbol, 4) +
		padInt(price, 8) +
		padInt(size, 6)
}

func Modify(orderId string, price int, size int) string {
	return "M" +
		padRight(orderId, 10) +
		padRight("", 4) +
		padInt(price, 8) +
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

func Trade(symbol string, price int, size int) string {
	return "T" +
		padRight("", 10) + // orderId unused
		padRight(symbol, 4) +
		padInt(price, 8) +
		padInt(size, 6)
}

func Lines(events ...string) string {
	return strings.Join(events, "\n")
}

func TestProcessPitchFile_SimpleScenario(t *testing.T) {

	pitchData := Lines(
		Add("1", "AAPL", 10000, 100),
		Add("2", "MSFT", 20000, 80),
		Execute("1", 50),
		Modify("2", 21000, 100),
		Execute("2", 40),
		Cancel("2"),
	)

	file, err := os.CreateTemp("", "pitch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(pitchData)
	file.Close()

	result := processPitchFile(file.Name(), "testdata/parser.json")

	if len(result) == 0 {
		t.Fatal("expected traded symbols")
	}

	if result[0].Symbol != "AAPL" {
		t.Errorf("expected AAPL got %s", result[0].Symbol)
	}
}
