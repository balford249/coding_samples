package pitchprocesser

import (
	"strconv"
)

type EventType int

const (
	AddOrder EventType = iota
	ModifyOrder
	ExecuteOrder
	Trade
	UnknownEvent
)

type PitchProcesser struct {
	fileParser PitchFileParser
}

type AddOrderDetails struct {
	OrderId string
	Symbol  string
	Price   float32
	Size    int
}

func extractField(line string, f FieldOffset) string {
	return line[f.Start:f.End]
}

func (pp PitchProcesser) getEvent(line string) EventType {
	eventVal := extractField(line, pp.fileParser.EventTypeOffset)
	switch eventVal {
	case pp.fileParser.EventChars.AddOrder:
		return AddOrder
	case pp.fileParser.EventChars.ModifyOrder:
		return ModifyOrder
	case pp.fileParser.EventChars.ExecuteOrder:
		return ExecuteOrder
	case pp.fileParser.EventChars.Trade:
		return Trade
	}
	return UnknownEvent
}

func (pp PitchProcesser) getAddOrderDetails(line string) AddOrderDetails {
	orderid := extractField(line, pp.fileParser.AddOrderOffsets.orderId)
	symbol := extractField(line, pp.fileParser.AddOrderOffsets.symbol)
	price64, _ := strconv.ParseFloat(extractField(line, pp.fileParser.AddOrderOffsets.price), 64)
	price := float32(price64)
	size, _ := strconv.Atoi(extractField(line, pp.fileParser.AddOrderOffsets.size))
	return AddOrderDetails{OrderId: orderid, Symbol: symbol, Price: price, Size: size}
}
