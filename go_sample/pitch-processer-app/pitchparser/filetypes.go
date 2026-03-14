package pitchprocesser


var FileTypeAParser = PitchFileParser{
	EventTypeOffset: FieldOffset{Start: 0, End: 0},

	AddOrderOffsets: AddOrderEventOffsets{
		orderId: FieldOffset{Start: 1, End: 11},
		symbol:  FieldOffset{Start: 11, End: 15},
		price:   FieldOffset{Start: 15, End: 23},
		size:    FieldOffset{Start: 23, End: 29},
	},

	ModifyOrderOffsets: ModifyOrderEventOffsets{
		orderId: FieldOffset{1, 11},
		price:   FieldOffset{15, 23},
		size:    FieldOffset{23, 29},
	},

	ExecuteOrderOffsets: ExecuteOrderEventOffsets{
		orderId: FieldOffset{1, 11},
		size:    FieldOffset{23, 29},
	},

	TradeOffsets: TradeEventOffsets{
		symbol: FieldOffset{11, 15},
		size:   FieldOffset{23, 29},
	},
}