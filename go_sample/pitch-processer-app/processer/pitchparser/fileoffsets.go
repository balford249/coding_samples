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
	OrderID FieldOffset `json:"OrderID"`
	Symbol  FieldOffset `json:"symbol"`
	Price   FieldOffset `json:"price"`
	Size    FieldOffset `json:"size"`
}

type ModifyOrderEventOffsets struct {
	OrderID FieldOffset `json:"OrderID"`
	Price   FieldOffset `json:"price"`
	Size    FieldOffset `json:"size"`
}

type CancelOrderEventOffsets struct {
	OrderID FieldOffset `json:"OrderID"`
}

type ExecuteOrderEventOffsets struct {
	OrderID FieldOffset `json:"OrderID"`
	Size    FieldOffset `json:"size"`
}

type TradeEventOffsets struct {
	Symbol FieldOffset `json:"symbol"`
	Price  FieldOffset `json:"price"`
	Size   FieldOffset `json:"size"`
}

