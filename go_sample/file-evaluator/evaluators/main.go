package main

import (
	consumer "file-evaluator/consumer"
	db "file-evaluator/db"
	types "file-evaluator/evaluators/types"
	utils "file-evaluator/utils"
	"flag"
)

var Registry = map[string]consumer.ConsumerType{
	"FileExists": types.FileExistsConsumerType{DB: *db.InitDB()},
	"IsTxtFile":  types.IsTxtFileConsumerType{DB: *db.InitDB()},
}

type EvaluatorConfig struct {
	Broker   string `json:"broker"`
	Topic    string `json:"topic"`
	GroupID  string `json:"group"`
	EvalType string `json:"type"`
}

func getConfigFile() string {
	configFilePath := flag.String("config", "", "Filepath for the JSON config")

	flag.Parse()

	if *configFilePath == "" {
		panic("--config is required")
	}

	return *configFilePath
}

func main() {
	var config EvaluatorConfig
	utils.LoadConfig(getConfigFile(), &config)
	kafkaRunner := consumer.KafkaRunner{
		ConsumerType: Registry[config.EvalType],
		Config:       consumer.ConsumerConfig{Broker: config.Broker, Topic: config.Topic, GroupID: config.GroupID}}
	kafkaRunner.Run()
}

