package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

type Payload struct {
	PitchFile  string `json:"pitchfile"`
	ConfigFile string `json:"configfile"`
}

type Config struct {
	Broker  string `json:"broker"`
	Topic   string `json:"topic"`
	APIPort string `json:"apiPort"`
}

var kafkaProducer *kafka.Producer
var config Config

func loadConfig() error {
	data, err := ioutil.ReadFile("config.json")
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
func pitchProcessHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a unique request ID, TODO - check if in Postgres table
	requestID := uuid.New().String() // Generate a unique UUID for this request

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to convert payload to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	kafkaMessageWithID := struct {
		ID      string  `json:"id"`
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
		ID string `json:"id"`
	}{
		ID: requestID,
	}

	json.NewEncoder(w).Encode(response)

	close(deliveryChan)
}

func main() {
	var err error
	kafkaProducer, err = initKafkaProducer()
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	http.HandleFunc("/pitchprocesser/", pitchProcessHandler)

	// Start the server
	port := ":8080"
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
