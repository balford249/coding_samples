package evaluators

import (
	"encoding/json"
	consumer "file-evaluator/consumer"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	db "file-evaluator/db"
)

type FileEvalEvent struct {
    ID       int64  `json:"id"`
    FilePath string `json:"path"`
    EvalType string `json:"type"`
}

type BaseFileEvaluator struct{
		DB *db.Store
}

func (BaseFileEvaluator) EventStruct() consumer.Event {
    return FileEvalEvent{}
}

func (BaseFileEvaluator) EventMapper(msg *kafka.Message) (consumer.Event, error) {
    var event FileEvalEvent
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        return nil, err
    }
    return event, nil
}

func (f BaseFileEvaluator) insertResult(processingId int64, passed bool, evalType string) error {
	var result string
	if passed {
		result = "passed"
	} else {
		result = "failed"
	}
    return f.DB.InsertResult(processingId, result, evalType)
}
	