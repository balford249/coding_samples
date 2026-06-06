package main

import (
	handlers "file-evaluator/api/handlers"
	producer "file-evaluator/api/kafka"
	database "file-evaluator/db"
	utils "file-evaluator/utils"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Broker  string `json:"broker"`
	Topic   string `json:"topic"`
	APIPort int    `json:"apiPort"`
}

func main() {
	var config Config
	utils.LoadConfig("config.json", &config)
	kafkaProducer := producer.InitKafkaProducer(config.Broker, config.Topic)
	defer kafkaProducer.Close()

	db := database.InitDB()
	defer db.DB.Close()

	h := handlers.HttpHandler{DB: db, KafkaProducer: kafkaProducer}

	http.HandleFunc("/evaluate", h.FileEvalHandler)

	url := fmt.Sprintf("0.0.0.0:%d", config.APIPort)
	log.Printf("Starting server on port %s...", url)
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
