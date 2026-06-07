package evaluators

import (
	consumer "file-evaluator/consumer"
	db "file-evaluator/db"
	"path/filepath"
)

type IsTxtFileConsumerType struct {
	BaseFileEvaluator
	DB *db.Store
}

func (f IsTxtFileConsumerType) EventHandler(e consumer.Event) error {
	event := e.(FileEvalEvent)
	res := filepath.Ext(event.FilePath) == ".txt"
	return f.insertResult(event.ID, res, "IsTxtFile")
}

func NewIsTxtFileConsumer(store *db.Store) IsTxtFileConsumerType {
	return IsTxtFileConsumerType{
		BaseFileEvaluator: BaseFileEvaluator{DB: store},
	}
}
