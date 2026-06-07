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
	Type   string `json:"evalType"`
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

func (s *Store) CreateNewEvent() (int64, error) {
	var id int64
	err := s.DB.QueryRow(
		"INSERT INTO file_evaluation_event DEFAULT VALUES RETURNING id").Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("Error inserting into database: %v", err)
	}

	return id, nil
}

func (s *Store) GetEvent(id int64) ([]EvaluationResult, error) {
	rows, err := s.DB.Query(
		`SELECT t.type, COALESCE(r.status, 'pending') AS status
		FROM file_evaluation_type t
		LEFT JOIN file_evaluation_result r
		 ON r.type = t.type
		 AND r.eval_id = $1
		WHERE EXISTS (
		 SELECT 1 
    	FROM file_evaluation_event e
    	WHERE e.id = $1
		)
		ORDER BY t.type;`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []EvaluationResult{}

	for rows.Next() {
		var res EvaluationResult
		if err := rows.Scan(&res.Type, &res.Status); err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Store) InsertResult(processingId int64, result string, evalType string) {
	_, err := s.DB.Exec("INSERT INTO file_evaluation_result  VALUES ($1, $2, $3)", processingId, evalType, result)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
	}
}
