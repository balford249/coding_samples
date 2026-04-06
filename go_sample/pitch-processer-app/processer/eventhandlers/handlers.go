package processer

import (
	orderbook "pitch-processer-app/processer/orderbook"
	pitchparser "pitch-processer-app/processer/pitchparser"
)

func HandleAddOrder(ob *orderbook.OrderBook, parser pitchparser.PitchFileParser, line string) error {

	details, err := parser.GetAddOrderDetails(line)
	if err != nil {
		return err
	}

	return ob.AddOrder(&orderbook.Order{
		ID:     details.OrderID,
		Symbol: details.Symbol,
		Price:  details.Price,
		Size:   details.Size,
	})
}

func HandleCancelOrder(ob *orderbook.OrderBook, parser pitchparser.PitchFileParser, line string) error {

	details, err := parser.GetCancelOrderDetails(line)
	if err != nil {
		return err
	}

	return ob.RemoveOrder(details.OrderID)
}

func HandleExecuteOrder(ob *orderbook.OrderBook, parser pitchparser.PitchFileParser, line string) error {

	details, err := parser.GetExecutionDetails(line)
	if err != nil {
		return err
	}

	return ob.ExecuteOrder(details.OrderID, details.Size)
}

func HandleModifyOrder(ob *orderbook.OrderBook, parser pitchparser.PitchFileParser, line string) error {

	details, err := parser.GetModifyOrderDetails(line)
	if err != nil {
		return err
	}

	return ob.ModifyOrder(details.OrderID, details.Size, details.Price)
}

func HandleTrade(ob *orderbook.OrderBook, parser pitchparser.PitchFileParser, line string) error {

	details, err := parser.GetTradeDetails(line)
	if err != nil {
		return err
	}

	return ob.HandleTrade(details.Symbol, details.Size, details.Price)
}