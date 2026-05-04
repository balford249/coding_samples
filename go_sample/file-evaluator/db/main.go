package db

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
)

type Store struct {
	DB *sql.DB
}

type EvaluationResult struct {
	Status string `json:"status"`
	Result *bool  `json:"result,omitempty"`
}

func InitDB() *Store {
	var db *sql.DB

	var err error
	db, err = sql.Open("postgres", "user=appuser dbname=eval host=postgres password=password sslmode=disable")
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

func (s *Store) CreateNewEvent() int64 {
	var id int64
	err := s.DB.QueryRow(
		"INSERT INTO file_evaluation (status) values ('working') RETURNING id",
	).Scan(&id)

	if err != nil {
		log.Printf("Error inserting into database: %v", err)
		return 0
	}

	return id
}

func (s *Store) GetEvent(id int64) (*EvaluationResult, error) {
	var res EvaluationResult

	err := s.DB.QueryRow(
		`SELECT status, result 
		 FROM file_evaluation 
		 WHERE id = $1`, id,
	).Scan(&res.Status, &res.Result)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found
		}
		return nil, err
	}

	return &res, nil
}

func (s *Store) InsertResult(processingId int64, result bool) {
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling data into JSON: %v", err)
	}

	_, err = s.DB.Exec("UPDATE pitch_processing_event SET result = $1, WHERE id = $2", jsonData, processingId)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
	}
}
