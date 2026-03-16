package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"pitch-processer-app/eventhandlers"
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

	if *pitchFile == "" || *configFile == "" {
		fmt.Println("usage: app -pitchFile <file> -config <config>")
		os.Exit(1)
	}

	return Args{
		PitchFile:  *pitchFile,
		ConfigFile: *configFile,
	}
}

func failOnErr(err error, context string) {
	if err != nil {
		fmt.Printf("%s failed: %v\n", context, err)
		os.Exit(1)
	}
}

type eventHandler func(
	ob *orderbook.OrderBook,
	parser pitchparser.PitchFileParser,
	line string,
) error

var handlers = map[byte]eventHandler{
	byte(pitchparser.AddOrder):     eventhandlers.HandleAddOrder,
	byte(pitchparser.ModifyOrder):  eventhandlers.HandleModifyOrder,
	byte(pitchparser.CancelOrder):  eventhandlers.HandleCancelOrder,
	byte(pitchparser.ExecuteOrder): eventhandlers.HandleExecuteOrder,
	byte(pitchparser.Trade):        eventhandlers.HandleTrade,
}

func processPitchFile(pitchFilePath string, configPath string) ([]orderbook.SymbolVolume, error) {

	parser, err := pitchparser.NewPitchParser(configPath)
	if err != nil {
		return nil, err
	}

	ob := orderbook.NewOrderBook()

	file, err := os.Open(pitchFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()

		if len(line) == 0 {
			return nil, fmt.Errorf("unexpected empty line")
		}

		event, parserError := parser.GetEvent(line)
		if parserError != nil {
			return nil, fmt.Errorf("failed to get event for line: %s error: %w",line, err)
		}
		handler, exists := handlers[byte(event)]
		if !exists {
			return nil, fmt.Errorf("failed to get event handler for line: %s error: %w",line, err)
		}

		err := handler(ob, parser, line)
		if err != nil {
			return nil, fmt.Errorf("processing line failed, line: %s error: %w",line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ob.SymbolVolumes(), nil
}

func main() {
	args := parseArgs()

	results, err := processPitchFile(args.PitchFile, args.ConfigFile)
	fmt.Printf("process of pitch file failed: %v\n", err)
	os.Exit(1)

	for _, volume := range results {
		fmt.Printf("%s : %f\n", volume.Symbol, volume.VolumeTraded)
	}
}
