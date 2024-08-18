package kafka

const (
	kafkaAddr             = "kafka:9092"
	gameTopic             = "games"
	numPartitions         = 4
	replicationFactor     = 1
	numTopicCreateRetries = 5

	// Message Types
	register   = "register"
	deregister = "deregister"
	update     = "update"
)
