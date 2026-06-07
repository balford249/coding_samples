package evaluators

import (
	consumer "file-evaluator/consumer"
	db "file-evaluator/db"
	"os"
)

type FileExistsConsumerType struct {
	BaseFileEvaluator
}

func NewFileExistsConsumer(store *db.Store) FileExistsConsumerType {
	return FileExistsConsumerType{
		BaseFileEvaluator: BaseFileEvaluator{DB: store},
	}
}

func (f FileExistsConsumerType) EventHandler(e consumer.Event) error {
	event := e.(FileEvalEvent)
	fileExists := true
	_, err := os.Stat(event.FilePath)
	if err != nil {
		fileExists = false
	}
	return f.insertResult(event.ID, fileExists, "FileExists")
}
