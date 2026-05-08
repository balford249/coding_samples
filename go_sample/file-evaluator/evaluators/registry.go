package evaluators


var Registry = map[string]FileEvaluator{
	"FileExists": FileExistsEvaluator{},
}