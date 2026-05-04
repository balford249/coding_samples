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

func handlePost(w http.ResponseWriter, r *http.Request) {
	var payload Payload

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON body: %v", err), http.StatusBadRequest)
		return
	}

	if payload.FilePath == "" {
		http.Error(w, "Missing required field: path", http.StatusBadRequest)
		return
	}

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
	deliveryChan := make(chan kafka.Event, 1)

	err = kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          kafkaMessageWithIDBytes,
	}, deliveryChan)

	e := <-deliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		http.Error(w, fmt.Sprintf("Failed to send message to Kafka: %v", msg.TopicPartition.Error), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{
		ID: requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	close(deliveryChan)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// Expect: /evaluate?id=123
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var id int64
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	// TODO: implement this in your db package
	result, err := db.GetEvent(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch event: %v", err), http.StatusInternalServerError)
		return
	}

	if result == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func fileEvalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		handlePost(w, r)

	case http.MethodGet:
		handleGet(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
