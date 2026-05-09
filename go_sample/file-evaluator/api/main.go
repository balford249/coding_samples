package main

import (
	database "file-evaluator/db"
	"fmt"
	"log"
	"net/http"
	utils "file-evaluator/utils"
	producer "file-evaluator/api/kafka"
	handlers "file-evaluator/api/handlers"
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
