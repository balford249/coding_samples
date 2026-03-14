package pitchprocesser

import (
	"encoding/json"
	"os"
)

type FieldOffset struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type EventChars struct {
	AddOrder     string `json:"addOrder"`
	ModifyOrder  string `json:"modifyOrder"`
	ExecuteOrder string `json:"executeOrder"`
	CancelOrder  string `json:"cancelOrder"`
	Trade        string `json:"trade"`
}

type AddOrderEventOffsets struct {
	OrderId FieldOffset `json:"orderId"`
	Symbol  FieldOffset `json:"symbol"`
	Price   FieldOffset `json:"price"`
	Size    FieldOffset `json:"size"`
}

type ModifyOrderEventOffsets struct {
	OrderId FieldOffset `json:"orderId"`
	Price   FieldOffset `json:"price"`
	Size    FieldOffset `json:"size"`
}

type CancelOrderEventOffsets struct {
	OrderId FieldOffset `json:"orderId"`
}

type ExecuteOrderEventOffsets struct {
	OrderId FieldOffset `json:"orderId"`
	Size    FieldOffset `json:"size"`
}

type TradeEventOffsets struct {
	Symbol FieldOffset `json:"symbol"`
	Price  FieldOffset `json:"price"`
	Size   FieldOffset `json:"size"`
}

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