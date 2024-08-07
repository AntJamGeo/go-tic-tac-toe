package kafka

import (
	"log"
	"net"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func createGameTopic() error {
	log.Printf("kafkaAddr: %s", kafkaAddr)
	var conn *kafka.Conn
	var err error
	for i := 0; i < numTopicCreateRetries; i++ {
		conn, err = kafka.Dial("tcp", kafkaAddr)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 10)
		log.Printf("failed connection attempt %d/%d to kafkaAddr: %s", i, numTopicCreateRetries, kafkaAddr)
	}
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             gameTopic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}
