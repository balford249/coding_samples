package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

	err = pingWithRetry(db, 3)
	if err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	return &Store{db}
}

func pingWithRetry(db *sql.DB, maxAttempts int) error {
	var err error

	baseDelay := 1 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = db.Ping()
		if err == nil {
			return nil
		}

		log.Printf("Ping attempt %d failed: %v", attempt, err)

		if attempt < maxAttempts {
			// exponential backoff: 1s, 2s, 4s...
			backoff := baseDelay * time.Duration(1<<(attempt-1))
			log.Printf("Retrying in %v...", backoff)
			time.Sleep(backoff)
		}
	}

	return err
}

func (s *Store) CreateNewEvent(evalType string) (int64, error) {
	var id int64
	err := s.DB.QueryRow(
		"INSERT INTO file_evaluation_event RETURNING id",
		evalType,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("Error inserting into database: %v", err)
	}

	return id, nil
}

func (s *Store) GetEvent(id int64) (*EvaluationResult, error) {
	var res EvaluationResult
	err := s.DB.QueryRow(
		`SELECT type, upper(coalesce(status, 'pending'))
		 FROM file_evaluation_type
		 JOIN  file_evaluation_result using (type)
		 WHERE eval_id = $1`, id,
	).Scan(&res.Status, &res.Result)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found
		}
		return nil, err
	}

	return &res, nil
}

func (s *Store) InsertResult(processingId int64, result string, evalType string) {
	_, err := s.DB.Exec("INSERT INTO file_evaluation_result  VALUES ($1, $2, $3)", processingId, evalType, result)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
	}
}
