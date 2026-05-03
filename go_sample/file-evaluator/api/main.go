package main

import (
	"encoding/json"
	database "file-evaluator/db"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Payload struct {
	FilePath string `json:"path"`
}

type Config struct {
	Broker  string `json:"broker"`
	Topic   string `json:"topic"`
	APIPort int    `json:"apiPort"`
}

var kafkaProducer *kafka.Producer
var config Config
var db database.Store

func loadConfig() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}

func initKafkaProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Broker,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	return producer, nil
}

// POST handler for /pitchprocesser/
func fileEvalHandler(w http.ResponseWriter, r *http.Request) {

	var payload Payload

	// Step 2: Decode the request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON body: %v", err), http.StatusBadRequest)
		return
	}

	// Step 3: Additional field validation (optional, depends on your use case)
	if payload.FilePath == "" {
		http.Error(w, "Missing required field: path", http.StatusBadRequest)
		return
	}

	// Step 4: Marshal the payload back to check for any potential errors
	_, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to convert payload to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate a unique request ID, TODO - check if in Postgres table
	requestID := db.CreateNewEvent()

	kafkaMessageWithID := struct {
		ID      int64   `json:"id"`
		Payload Payload `json:"payload"`
	}{
		ID:      requestID,
		Payload: payload,
	}

	kafkaMessageWithIDBytes, err := json.Marshal(kafkaMessageWithID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to include ID in Kafka message: %v", err), http.StatusInternalServerError)
		return
	}

	kafkaTopic := config.Topic

	// GO uses channels to communicate between goroutines.
	// This creates a buffered channel that holds one event of type kafka.Event. The Kafka producer sends a delivery
	// report to this channel after it produces a message.
	deliveryChan := make(chan kafka.Event, 1)
	err = kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          kafkaMessageWithIDBytes,
	}, deliveryChan)

	// gets the delivery back from the kafka producer and asserts thats it can be constructed into a Kafka.Message
	e := <-deliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		http.Error(w, fmt.Sprintf("Failed to send message to Kafka: %v", msg.TopicPartition.Error), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := struct {
		ID int64 `json:"id"`
	}{
		ID: requestID,
	}

	json.NewEncoder(w).Encode(response)

	close(deliveryChan)
}

func main() {
	var err error
	loadConfig()
	db = *database.InitDB("user=appuser dbname=eval password=password")
	kafkaProducer, err = initKafkaProducer()
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	http.HandleFunc("/evaluate", fileEvalHandler)

	url := fmt.Sprintf("0.0.0.0:%d", config.APIPort)
	log.Printf("Starting server on port %s...", url)
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
