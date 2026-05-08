package evaluators


type FileEvaluator interface {
	Name() string
	Evaluate(path string) (EvaluationResult, error)
}

type EvaluationResult struct {
	Passed         bool
	EvaluationType string
}

