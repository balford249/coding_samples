package main

import (
	"encoding/json"
	database "file-evaluator/db"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
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

var db database.Store

func getConfig() (Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	configData, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Errorf("Error reading config file: %v", err)
	}

	var cfg Config
	err = json.Unmarshal(configData, &cfg)
	return cfg, err
}

// TODO - Get group from the eval name
func getKafkaConsumer() *kafka.Consumer {
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("Error parsing config data: %v", err)
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Broker,
		"group.id":          "file-evaluator",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	return consumer
}

func eval(filePath string) bool {
	_, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %v: %v", filePath, err)
	}
	return true
}

func main() {
	db = *database.InitDB()
	var cfg Config
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("Error parsing config data: %v", err)
	}

	consumer := getKafkaConsumer()
	defer consumer.Close()

	err = consumer.Subscribe(cfg.Topic, nil)
	if err != nil {
		log.Fatalf("Error subscribing to topic %s: %v", cfg.Topic, err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		var message Message
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		log.Printf("%d", message.ID)
		log.Printf(message.FilePath)
		res := eval(message.FilePath)
		if err != nil {
			log.Printf("Error processing data: %v", err)
		}
		db.InsertResult(message.ID, res)

		fmt.Printf("Inserted data for ID %d\n", message.ID)
	}
}
