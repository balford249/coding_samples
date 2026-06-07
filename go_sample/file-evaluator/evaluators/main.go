package main

import (
	consumer "file-evaluator/consumer"
	db "file-evaluator/db"
	types "file-evaluator/evaluators/types"
	utils "file-evaluator/utils"
	"flag"
)

func initRegistry(store *db.Store) map[string]consumer.ConsumerType {
	return map[string]consumer.ConsumerType{
		"FileExists": types.NewFileExistsConsumer(store),
		"IsTxtFile":  types.NewIsTxtFileConsumer(store),
	}
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

	store := db.InitDB()
	registry := initRegistry(store)

	kafkaRunner := consumer.KafkaRunner{
		ConsumerType: registry[config.EvalType],
		Config:       consumer.ConsumerConfig{Broker: config.Broker, Topic: config.Topic, GroupID: config.GroupID}}
	kafkaRunner.Run()
}
