package main

import (
	"context"
	Consumer "file-evaluator/consumer/generic_consumer"
	database "file-evaluator/db"
	utils "file-evaluator/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Config structure for reading Kafka configuration from JSON
type Config struct {
	Broker string `json:"broker"`
	Topic  string `json:"topic"`
}

type Message struct {
	ID       int64  `json:"id"`
	FilePath string `json:"path"`
}

func eval(filePath string) bool {
	_, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %v: %v", filePath, err)
	}
	return true
}

func main() {
	db := *database.InitDB()
	var config *Config
	utils.LoadConfig("config.json", config)


	consumer := Consumer.InitKafkaConsumer(config.Broker)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("shutdown signal received")
		cancel()
	}()

	ch, err := consumer.ConsumeMapped(ctx, []string{"customers-topic"}, Consumer.MapFileEvalEvent)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range ch {
		event := msg.(Message)
		res := eval(event.FilePath)
		db.InsertResult(event.ID, res)
		
	}

	fmt.Println("consumer stopped cleanly")

}
