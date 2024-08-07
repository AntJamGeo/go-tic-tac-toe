package kafka

const kafkaAddr = "kafka:9092"
const gameTopic = "games"
const numPartitions = 4
const replicationFactor = 1
const numTopicCreateRetries = 5

// Message Types
const register = "register"
const deregister = "deregister"
const update = "update"
