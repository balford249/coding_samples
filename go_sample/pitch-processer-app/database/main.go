package database

import (
	"database/sql"
	"encoding/json"
	"log"
	orderbook "pitch-processer-app/processer/orderbook"

	_ "github.com/lib/pq"
)

type Store struct {
	DB *sql.DB
}

func InitDB(connStr string) (*Store) {
	var db *sql.DB
	
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Optional but recommended
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	// Configure pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	return &Store{db}
}

func (s *Store) CreateNewProcessingEvent() int64 {
	var id int64
	err := s.DB.QueryRow(
		"INSERT INTO pitch_processing_event (status) RETURNING id",
	).Scan(&id)

	if err != nil {
		log.Printf("Error inserting into database: %v", err)
		return 0
	}

	return id
}

func (s *Store) InsertResult(processingId int64, result []orderbook.SymbolVolume) {
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling data into JSON: %v", err)
	}

	_, err = s.DB.Exec("UPDATE pitch_processing_event SET result = $1, WHERE id = $2", jsonData, processingId)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
	}
}
