package processer

import (
	"fmt"
	"math"
	"strings"
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

func PadRight(s string, width int) string {
	return padRight(s, width)
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
