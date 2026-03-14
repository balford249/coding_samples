package pitchparser

import (
	"strconv"
	"os"
	"encoding/json"
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

type PitchFileParser struct {
	EventChars      EventChars              `json:"eventChars"`
	EventTypeOffset FieldOffset             `json:"eventTypeOffset"`
	AddOrderOffsets AddOrderEventOffsets    `json:"addOrderOffsets"`
	ModifyOrderOffsets ModifyOrderEventOffsets `json:"modifyOrderOffsets"`
	ExecuteOrderOffsets ExecuteOrderEventOffsets `json:"executeOrderOffsets"`
	CancelOrderOffsets CancelOrderEventOffsets `json:"cancelOrderOffsets"`
	TradeOffsets TradeEventOffsets          `json:"tradeOffsets"`
}

func LoadParserConfig(path string) (PitchFileParser, error) {

	var parser PitchFileParser

	data, err := os.ReadFile(path)
	if err != nil {
		return parser, err
	}

	err = json.Unmarshal(data, &parser)
	if err != nil {
		return parser, err
	}

	return parser, nil
}

func NewPitchParser(filePath string) PitchFileParser {

	parser, err := LoadParserConfig(filePath)
	if err != nil {
		panic(err)
	}

	return parser
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

type ExecutionDetails struct {
	OrderId string
	Size    int
}

func extractField(line string, f FieldOffset) string {
	return line[f.Start:f.End]
}

func (pp PitchFileParser) GetEvent(line string) EventType {
	eventVal := extractField(line, pp.EventTypeOffset)
	switch eventVal {
	case pp.EventChars.AddOrder:
		return AddOrder
	case pp.EventChars.ModifyOrder:
		return ModifyOrder
	case pp.EventChars.CancelOrder:
		return CancelOrder
	case pp.EventChars.ExecuteOrder:
		return ExecuteOrder
	case pp.EventChars.Trade:
		return Trade
	}
	return UnknownEvent
}

func (pp PitchFileParser) GetAddOrderDetails(line string) AddOrderDetails {
	OrderId := extractField(line, pp.AddOrderOffsets.OrderId)
	Symbol := extractField(line, pp.AddOrderOffsets.Symbol)
	Price64, _ := strconv.ParseFloat(extractField(line, pp.AddOrderOffsets.Price), 64)
	Price := float32(Price64)
	Size, _ := strconv.Atoi(extractField(line, pp.AddOrderOffsets.Size))
	return AddOrderDetails{OrderId: OrderId, Symbol: Symbol, Price: Price, Size: Size}
}

func (pp PitchFileParser) GetModifyOrderDetails(line string) ModifyOrderDetails {
	OrderId := extractField(line, pp.ModifyOrderOffsets.OrderId)
	Price64, _ := strconv.ParseFloat(extractField(line, pp.ModifyOrderOffsets.Price), 64)
	Price := float32(Price64)
	Size, _ := strconv.Atoi(extractField(line, pp.ModifyOrderOffsets.Size))
	return ModifyOrderDetails{OrderId: OrderId, Price: Price, Size: Size}
}

func (pp PitchFileParser) GetOrderExecutedDetails(line string) ExecutionDetails {
	OrderId := extractField(line, pp.ExecuteOrderOffsets.OrderId)
	Size, _ := strconv.Atoi(extractField(line, pp.ExecuteOrderOffsets.Size))
	return ExecutionDetails{OrderId: OrderId, Size: Size}
}


func (pp PitchFileParser) GetCancelOrderDetails(line string) CancelOrderDetails {
	OrderId := extractField(line, pp.CancelOrderOffsets.OrderId)
	return CancelOrderDetails{OrderId: OrderId}
}

func (pp PitchFileParser) GetTradeDetails(line string) TradeDetails {
	Symbol := extractField(line, pp.TradeOffsets.Symbol)
	Price64, _ := strconv.ParseFloat(extractField(line, pp.TradeOffsets.Price), 64)
	Price := float32(Price64)
	Size, _ := strconv.Atoi(extractField(line, pp.TradeOffsets.Size))
	return TradeDetails{Symbol: Symbol, Price: Price, Size: Size}
}