package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"pitch-processer-app/orderbook"
	"pitch-processer-app/pitchparser"
)

type Args struct {
	PitchFile  string
	ConfigFile string
}

func parseArgs() Args {

	pitchFile := flag.String("pitchFile", "", "pitch file path")
	configFile := flag.String("config", "", "parser config json")

	flag.Parse()

	return Args{
		PitchFile:  *pitchFile,
		ConfigFile: *configFile,
	}
}

func processPitchFile(pitchFilePath string, pitchParserConfigFilePath string) []orderbook.SymbolVolume{
	parser := pitchparser.NewPitchParser(pitchParserConfigFilePath)
	ob := orderbook.OrderBook{
		Book:                    make(map[string]*orderbook.Order),
		QuantatiyTradedBySymbol: make(map[string]int),
	}

	file, err := os.Open(pitchFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		event := parser.GetEvent(line)
		switch event {
		case pitchparser.AddOrder:
			details := parser.GetAddOrderDetails(line)
			newOrder := orderbook.Order{
				ID:     details.OrderId,
				Symbol: details.Symbol,
				Price:  details.Price,
				Size:   details.Size,
			}
			ob.AddOrder(newOrder)
		case pitchparser.CancelOrder:
			details := parser.GetCancelOrderDetails(line)
			// Error handling 
			ob.RemoveOrder(details.OrderId)
		case pitchparser.ExecuteOrder:
			details := parser.GetOrderExecutedDetails(line)
			ob.ExecuteOrder(details.OrderId, details.Size)
		case pitchparser.ModifyOrder:
			details := parser.GetModifyOrderDetails(line)
			ob.ModifyOrder(details.OrderId, details.Size, details.Price)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return ob.TopTradedSymbols(5)

}

// Logging
// Naming
// Load the JSON file and create a parser
// Create a book struct
// For each line in the file:
// get the event
// switch off the event to create an order object etc and then feed that into the book
// Get the top n symbols, if n not provided output them all
func main() {
	args := parseArgs()
	processPitchFile(args.PitchFile, args.ConfigFile)
}
