package consumer


import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"context"
	"time"
	"encoding/json"
)


type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

type FileEvalEvent struct {
	ID       int64  `json:"id"`
	FilePath string `json:"path"`
	EvalType string `json:"type"`
}

func InitKafkaConsumer(broker string) *KafkaConsumer {
	// TODO - create better config 
	consumer , err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "file-evaluator",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	return &KafkaConsumer{Consumer: consumer}
}

func MapFileEvalEvent(msg *kafka.Message) (any, error) {
	var event FileEvalEvent

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}


func (c *KafkaConsumer) ConsumeMapped(
	ctx context.Context,
	topics []string,
	mapper func(*kafka.Message) (any, error),
) (<-chan any, error) {

	defer c.Consumer.Close()
	out := make(chan any)

	err := c.Consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("failed to subscribe: %w", err)
	}

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msg, err := c.Consumer.ReadMessage(time.Second)
			if err != nil {
				continue
			}

			mapped, err := mapper(msg)
			if err != nil {
				// Deliberately fail fast 
				log.Fatalf("Failed to map Kafka message to struct: %w", err)
			}

			select {
			case out <- mapped:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}