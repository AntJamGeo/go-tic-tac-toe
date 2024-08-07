package kafka

import kafka "github.com/segmentio/kafka-go"

var base64Mapping = createBase64Mapping()

func balancer(msg kafka.Message, partitions ...int) int {
	startCharVal := base64Mapping[rune(msg.Key[0])]
	return startCharVal / 16
}

func createBase64Mapping() map[rune]int {
	mapping := make(map[rune]int)

	// Map A-Z to 0-25
	for i := 'A'; i <= 'Z'; i++ {
		mapping[i] = int(i - 'A')
	}

	// Map a-z to 26-51
	for i := 'a'; i <= 'z'; i++ {
		mapping[i] = int(i - 'a' + 26)
	}

	// Map 0-9 to 52-61
	for i := '0'; i <= '9'; i++ {
		mapping[i] = int(i - '0' + 52)
	}

	// Map + to 62
	mapping['+'] = 62

	// Map / to 63
	mapping['/'] = 63

	return mapping
}
