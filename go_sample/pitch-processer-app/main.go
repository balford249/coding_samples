package main

import (
	"bufio"
	"fmt"
	"os"
)

// Logging
func main() {
	// Load the JSON file and create a parser
	// Create a book struct
	// For each line in the file:
	// get the event
	// switch off the event to create an order object etc and then feed that into the book
	// Get the top n symbols, if n not provided output them all
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filepath>")
		return
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
