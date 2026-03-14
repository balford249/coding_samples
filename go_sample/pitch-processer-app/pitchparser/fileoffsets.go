package pitchparser


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

