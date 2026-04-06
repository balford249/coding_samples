package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	processor "pitch-processer-app/processer"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
)

// Config structure for reading Kafka configuration from JSON
type Config struct {
	Broker string `json:"broker"`
	Topic  string `json:"topic"`
}

type Message struct {
	ID        int    `json:"id"`
	Pitchfile string `json:"pitchfile"`
	Configfile string `json:"configfile"`
}

// NewDB creates and returns a new database connection
func NewDB() (*sql.DB, error) {
	// Replace this connection string with your actual DB credentials
	connStr := "user=youruser dbname=yourdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}
	return db, nil
}

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

func getKafkaConsumer() (*kafka.Consumer){
	cfg, err :=  getConfig()
	if err != nil {
		log.Fatalf("Error parsing config data: %v", err)
	}
	
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Broker,
		"group.id":          "pitch-processor-app", 
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	return consumer
}

func main() {
	var cfg Config
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("Error parsing config data: %v", err)
	}

	db, err := NewDB()
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

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

		data, err := processor.ProcessPitchFile(message.Pitchfile, message.Configfile)
		if err != nil {
			log.Printf("Error processing data: %v", err)
		}
				if err != nil {
			log.Printf("Error processing data: %v", err)
			continue
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling data into JSON: %v", err)
			continue
		}

		_, err = db.Exec("INSERT INTO data_table (id, json_data) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET json_data = excluded.json_data", message.ID, jsonData)
		if err != nil {
			log.Printf("Error inserting into database: %v", err)
			continue
		}

		fmt.Printf("Inserted data for ID %d\n", message.ID)
	}
}