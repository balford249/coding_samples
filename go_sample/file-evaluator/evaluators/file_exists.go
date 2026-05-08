package evaluators

import (
	"os"
)


// Checks file exists
type FileExistsEvaluator struct {
}

func (s FileExistsEvaluator) Name() string {
	return "FileExists"
}

func (s FileExistsEvaluator) Evaluate(path string) (EvaluationResult, error) {
	var exists = true
	_, err := os.Stat(path)
	if err != nil {
		exists = false
	}
	return EvaluationResult{exists, "FileExists"}, nil
}
