package pitchparser

import (
	"encoding/json"
	"fmt"
	"os"
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

type PitchFileParser struct {
	EventChars          EventChars               `json:"eventChars"`
	EventTypeOffset     FieldOffset              `json:"eventTypeOffset"`
	AddOrderOffsets     AddOrderEventOffsets     `json:"addOrderOffsets"`
	ModifyOrderOffsets  ModifyOrderEventOffsets  `json:"modifyOrderOffsets"`
	ExecuteOrderOffsets ExecuteOrderEventOffsets `json:"executeOrderOffsets"`
	CancelOrderOffsets  CancelOrderEventOffsets  `json:"cancelOrderOffsets"`
	TradeOffsets        TradeEventOffsets        `json:"tradeOffsets"`
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

func NewPitchParser(filePath string) (PitchFileParser, error) {

	parser, err := LoadParserConfig(filePath)
	if err != nil {
		return parser, err
	}

	return parser, nil
}

type AddOrderDetails struct {
	OrderID string
	Symbol  string
	Price   float32
	Size    int
}

type ModifyOrderDetails struct {
	OrderID string
	Price   float32
	Size    int
}

type CancelOrderDetails struct {
	OrderID string
}

type TradeDetails struct {
	Symbol string
	Price  float32
	Size   int
}

type ExecutionDetails struct {
	OrderID string
	Size    int
}

func extractField(line string, f FieldOffset) (string, error) {
	// I put the offset validation here in the helper function instead of validating the PitchParser offsets after creation.
	// You could hard code the offsets it would check but if you added a field offset to AddOrderDetails for example, that would not be included.
	//  The alternative was reflection which reduced readiblity.
	if f.Start > f.End {
		return "", fmt.Errorf("invalid offset: end %d is smaller than start %d", f.End, f.Start)
	}

	if f.End > len(line) {
		return "", fmt.Errorf("offset end %d exceeds line length %d", f.End, len(line))
	}

	return line[f.Start:f.End], nil
}

func (pp PitchFileParser) GetEvent(line string) (EventType, error) {

	eventVal, err := extractField(line, pp.EventTypeOffset)
	if err != nil {
		return UnknownEvent, err
	}

	switch eventVal {
	case pp.EventChars.AddOrder:
		return AddOrder, nil
	case pp.EventChars.ModifyOrder:
		return ModifyOrder, nil
	case pp.EventChars.CancelOrder:
		return CancelOrder, nil
	case pp.EventChars.ExecuteOrder:
		return ExecuteOrder, nil
	case pp.EventChars.Trade:
		return Trade, nil
	default:
		return UnknownEvent, nil
	}
}

func (pp PitchFileParser) GetAddOrderDetails(line string) (AddOrderDetails, error) {

	orderID, err := extractField(line, pp.AddOrderOffsets.OrderID)
	if err != nil {
		return AddOrderDetails{}, err
	}

	symbol, err := extractField(line, pp.AddOrderOffsets.Symbol)
	if err != nil {
		return AddOrderDetails{}, err
	}

	priceStr, err := extractField(line, pp.AddOrderOffsets.Price)
	if err != nil {
		return AddOrderDetails{}, err
	}

	price, err := convertStringToFloat32(priceStr)
	if err != nil {
		return AddOrderDetails{}, err
	}

	sizeStr, err := extractField(line, pp.AddOrderOffsets.Size)
	if err != nil {
		return AddOrderDetails{}, err
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return AddOrderDetails{}, err
	}

	return AddOrderDetails{
		OrderID: orderID,
		Symbol:  symbol,
		Price:   price,
		Size:    size,
	}, nil
}

func (pp PitchFileParser) GetModifyOrderDetails(line string) (ModifyOrderDetails, error) {

	orderID, err := extractField(line, pp.ModifyOrderOffsets.OrderID)
	if err != nil {
		return ModifyOrderDetails{}, err
	}

	priceStr, err := extractField(line, pp.ModifyOrderOffsets.Price)
	if err != nil {
		return ModifyOrderDetails{}, err
	}

	price, err := convertStringToFloat32(priceStr)
	if err != nil {
		return ModifyOrderDetails{}, err
	}

	sizeStr, err := extractField(line, pp.ModifyOrderOffsets.Size)
	if err != nil {
		return ModifyOrderDetails{}, err
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return ModifyOrderDetails{}, err
	}

	return ModifyOrderDetails{
		OrderID: orderID,
		Price:   price,
		Size:    size,
	}, nil
}

func (pp PitchFileParser) GetExecutionDetails(line string) (ExecutionDetails, error) {

	orderID, err := extractField(line, pp.ExecuteOrderOffsets.OrderID)
	if err != nil {
		return ExecutionDetails{}, err
	}

	sizeStr, err := extractField(line, pp.ExecuteOrderOffsets.Size)
	if err != nil {
		return ExecutionDetails{}, err
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return ExecutionDetails{}, err
	}

	return ExecutionDetails{
		OrderID: orderID,
		Size:    size,
	}, nil
}

func (pp PitchFileParser) GetCancelOrderDetails(line string) (CancelOrderDetails, error) {

	orderID, err := extractField(line, pp.CancelOrderOffsets.OrderID)
	if err != nil {
		return CancelOrderDetails{}, err
	}

	return CancelOrderDetails{
		OrderID: orderID,
	}, nil
}

func (pp PitchFileParser) GetTradeDetails(line string) (TradeDetails, error) {

	symbol, err := extractField(line, pp.TradeOffsets.Symbol)
	if err != nil {
		return TradeDetails{}, err
	}

	priceStr, err := extractField(line, pp.TradeOffsets.Price)
	if err != nil {
		return TradeDetails{}, err
	}

	price, err := convertStringToFloat32(priceStr)
	if err != nil {
		return TradeDetails{}, err
	}

	sizeStr, err := extractField(line, pp.TradeOffsets.Size)
	if err != nil {
		return TradeDetails{}, err
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return TradeDetails{}, err
	}

	return TradeDetails{
		Symbol: symbol,
		Price:  price,
		Size:   size,
	}, nil
}

func convertStringToFloat32(val string) (float32, error) {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return float32(i) / 100.0, nil
}
