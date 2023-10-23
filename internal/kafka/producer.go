package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

type Producer struct {
	confluentProducer *kafka.Producer
	topic             string
}

func NewProducer(broker, topic string) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Error creating confluent kafka producer: %v", err)
	}
	return &Producer{confluentProducer: producer, topic: topic}, nil
}

func (p *Producer) Produce(_ context.Context, msg []byte) error {
	deliveryChan := make(chan kafka.Event)

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}

	// todo: add retry logic for retryable errors
	if err := p.confluentProducer.Produce(message, deliveryChan); err != nil {
		return fmt.Errorf("error producing message: %w", err)
	}

	// handle message delivery async
	// todo: make configurable
	timeout := time.After(10 * time.Second)
	go p.handleMessageDelivery(timeout, deliveryChan)
	return nil
}

func (p *Producer) handleMessageDelivery(timeout <-chan time.Time, deliveryChan chan kafka.Event) {
	select {
	case <-timeout:
		log.Println("timeout for waiting delivery expired")
		return
	case ev := <-deliveryChan:
		switch event := ev.(type) {
		case *kafka.Message:
			if event.TopicPartition.Error != nil {
				log.Printf("failed to deliver message: %v", event.TopicPartition)
				return
			}
			log.Printf("message successfully delivered to %v", event.TopicPartition)
		case kafka.Error:
			log.Printf("failed to deliver message: %v", event)
		default:
			log.Printf("ignored event: %v", event)
		}
	}
}

// todo: implement graceful shutdown (producer close)
