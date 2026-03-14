package pitchprocesser


type EventChars struct {
	AddOrder string
	ModifyOrder string
	ExecuteOrder string	
	Trade string
}

type FieldOffset struct {
	Start int
	End   int
}

type AddOrderEventOffsets struct {
	orderId FieldOffset
	symbol FieldOffset
	price FieldOffset
	size FieldOffset
}

type ModifyOrderEventOffsets struct {
	orderId FieldOffset
	price FieldOffset
	size FieldOffset
}

type ExecuteOrderEventOffsets struct {
	orderId FieldOffset
	size FieldOffset
}

type TradeEventOffsets struct {
	symbol FieldOffset
	size FieldOffset
}


type PitchFileParser struct {
	EventChars EventChars
	EventTypeOffset FieldOffset

	AddOrderOffsets     AddOrderEventOffsets
	ModifyOrderOffsets  ModifyOrderEventOffsets
	ExecuteOrderOffsets ExecuteOrderEventOffsets
	TradeOffsets        TradeEventOffsets
}

