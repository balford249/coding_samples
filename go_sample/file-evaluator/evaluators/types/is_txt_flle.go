package evaluators

import (
	"encoding/json"
	consumer "file-evaluator/consumer"
	db "file-evaluator/db"
	"path/filepath"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Checks file exists
type IsTxtFileConsumerType struct {
	DB  db.Store
}


func (f IsTxtFileConsumerType) EventStruct() (consumer.Event) {
	return FileEvalEvent{}
}

func (f IsTxtFileConsumerType) EventMapper(msg *kafka.Message) (consumer.Event, error) {
	var event FileEvalEvent

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (f IsTxtFileConsumerType) EventHandler(e consumer.Event) (error) {
	event := e.(FileEvalEvent)
	res := filepath.Ext(event.FilePath) == ".txt"

	var result string
	if res {
		result = "passed"
	} else {
		result = "failed"
	}

	f.DB.InsertResult(event.ID, result, "IsTxtFile")
	return nil
}

