package pitchprocesser

import (
	"strconv"
)

type EventType int

const (
	AddOrder EventType = iota
	ModifyOrder
	ExecuteOrder
	CancelOrder
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

type ModifyOrderDetails struct {
	OrderId string
	Price   float32
	Size    int
}

type CancelOrderDetails struct {
	OrderId string
}

type TradeDetails struct {
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
	case pp.fileParser.EventChars.CancelOrder:
		return CancelOrder
	case pp.fileParser.EventChars.ExecuteOrder:
		return ExecuteOrder
	case pp.fileParser.EventChars.Trade:
		return Trade
	}
	return UnknownEvent
}

func (pp PitchProcesser) getAddOrderDetails(line string) AddOrderDetails {
	OrderId := extractField(line, pp.fileParser.AddOrderOffsets.OrderId)
	Symbol := extractField(line, pp.fileParser.AddOrderOffsets.Symbol)
	Price64, _ := strconv.ParseFloat(extractField(line, pp.fileParser.AddOrderOffsets.Price), 64)
	Price := float32(Price64)
	Size, _ := strconv.Atoi(extractField(line, pp.fileParser.AddOrderOffsets.Size))
	return AddOrderDetails{OrderId: OrderId, Symbol: Symbol, Price: Price, Size: Size}
}

func (pp PitchProcesser) getModifyOrderDetails(line string) ModifyOrderDetails {
	OrderId := extractField(line, pp.fileParser.ModifyOrderOffsets.OrderId)
	Price64, _ := strconv.ParseFloat(extractField(line, pp.fileParser.ModifyOrderOffsets.Price), 64)
	Price := float32(Price64)
	Size, _ := strconv.Atoi(extractField(line, pp.fileParser.ModifyOrderOffsets.Size))
	return ModifyOrderDetails{OrderId: OrderId, Price: Price, Size: Size}
}

func (pp PitchProcesser) getCancelOrderDetails(line string) CancelOrderDetails {
	OrderId := extractField(line, pp.fileParser.CancelOrderOffsets.OrderId)
	return CancelOrderDetails{OrderId: OrderId}
}

func (pp PitchProcesser) getTradeDetails(line string) TradeDetails {
	Symbol := extractField(line, pp.fileParser.TradeOffsets.Symbol)
	Price64, _ := strconv.ParseFloat(extractField(line, pp.fileParser.TradeOffsets.Price), 64)
	Price := float32(Price64)
	Size, _ := strconv.Atoi(extractField(line, pp.fileParser.TradeOffsets.Size))
	return TradeDetails{Symbol: Symbol, Price: Price, Size: Size}
}