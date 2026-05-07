package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	
)


type KafkaProducer struct {
	producer *kafka.Producer
	topic string
}

func InitKafkaProducer(broker string, topic string ) *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
	})
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	return &KafkaProducer{producer: producer, topic: topic}
}

func (p *KafkaProducer) ProduceEvent(event []byte) error {
	defer p.producer.Close()
	deliveryChan := make(chan kafka.Event, 1)

	p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          event,
	}, deliveryChan)

	e := <-deliveryChan
	msg := e.(*kafka.Message)
	close(deliveryChan)

	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}
	return nil
}