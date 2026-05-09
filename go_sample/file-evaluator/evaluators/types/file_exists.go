package evaluators

import (
	"encoding/json"
	consumer "file-evaluator/consumer"
	db "file-evaluator/db"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Checks file exists
type FileExistsConsumerType struct {
	DB db.Store
}

// Event type
type FileEvalEvent struct {
	ID       int64  `json:"id"`
	FilePath string `json:"path"`
	EvalType string `json:"type"`
}

func (f FileExistsConsumerType) EventStruct() consumer.Event {
	return FileEvalEvent{}
}

func (f FileExistsConsumerType) EventMapper(msg *kafka.Message) (consumer.Event, error) {
	var event FileEvalEvent

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (f FileExistsConsumerType) EventHandler(e consumer.Event) error {
	event := e.(FileEvalEvent)
	fileExists := true
	_, err := os.Stat(event.FilePath)
	if err != nil {
		fileExists = false
	}
	f.DB.InsertResult(event.ID, fileExists)
	return nil
}
