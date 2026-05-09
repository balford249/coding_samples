package consumer

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerConfig struct {
	Broker  string
	Topic   string
	GroupID string
}

type Event any

type ConsumerType interface {
	EventStruct() Event
	EventMapper(msg *kafka.Message) (Event, error)
	EventHandler(Event) error
}

type KafkaRunner struct {
	ConsumerType ConsumerType
	Config       ConsumerConfig
}

func (k *KafkaRunner) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	consumer := InitKafkaConsumer(k.Config)

	ch, err := consumer.ConsumeMapped(
		ctx,
		[]string{k.Config.Topic},
		k.ConsumerType.EventMapper,
	)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for msg := range ch {
			if err := k.ConsumerType.EventHandler(msg); err != nil {
				log.Printf("handler error: %v", err)
			}
		}
	}()

	// wait for shutdown signal
	<-sigChan
	log.Println("shutdown signal received")

	// STEP 1: stop polling
	cancel()

	// STEP 2: wait for channel to drain
	wg.Wait()

	// STEP 3: close kafka consumer (IMPORTANT ORDER)
	consumer.Close()

	log.Println("consumer shutdown complete")
}

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

func InitKafkaConsumer(config ConsumerConfig) *KafkaConsumer {
	log.Println(config)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":           config.Broker,
		"group.id":                    config.GroupID,
		"auto.offset.reset":           "earliest",
		"session.timeout.ms":          6000,
		"metadata.request.timeout.ms": 5000,
		"request.timeout.ms":          5000,
	})
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	return &KafkaConsumer{Consumer: consumer}
}

func (c *KafkaConsumer) Close() {
	c.Consumer.Close()
}

func (k *KafkaConsumer) ConsumeMapped(
	ctx context.Context,
	topics []string,
	mapper func(*kafka.Message) (Event, error),
) (<-chan Event, error) {
	log.Println(topics)
	if err := k.Consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	out := make(chan Event)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msg, err := k.Consumer.ReadMessage(5000 * time.Millisecond)
			if err != nil {
				// safe: just log and continue
				if kafkaErr, ok := err.(kafka.Error); ok {
					log.Printf("kafka error: %v", kafkaErr)
				}
				continue
			}

			mapped, err := mapper(msg)
			if err != nil {
				log.Printf("mapping error: %v", err)
				continue
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
