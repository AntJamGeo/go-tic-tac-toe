package kafka

import (
	"context"
	"log"

	"github.com/AntJamGeo/go-tic-tac-toe/mqtodb/internal/message"
	kafka "github.com/segmentio/kafka-go"
)

func Read(ch chan message.Message) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddr},
		Topic:   gameTopic,
		GroupID: dbGroupID,
	})

	for {
		log.Printf("waiting for next msg...")
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from kafka: %v", err)
			continue
		}
		log.Printf("message received on partition: %d", msg.Partition)
		ch <- msg.Value
	}
}
