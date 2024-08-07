package kafka

import (
	"context"
	"encoding/json"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func write(msg *map[string]string, key []byte) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	w := &kafka.Writer{
		Addr:         kafka.TCP(kafkaAddr),
		Topic:        gameTopic,
		Balancer:     kafka.BalancerFunc(balancer),
		BatchTimeout: time.Millisecond * 10,
	}
	err = w.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   key,
			Value: jsonData,
		},
	)
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}
