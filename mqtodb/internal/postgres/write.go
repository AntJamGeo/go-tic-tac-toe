package postgres

import (
	"encoding/json"
	"log"

	"github.com/AntJamGeo/go-tic-tac-toe/mqtodb/internal/message"
)

func Write(ch chan message.Message) {
	var msgData map[string]string
	for msg := range ch {
		if err := json.Unmarshal(msg, &msgData); err != nil {
			log.Printf("unable to unmarshall: %v", err)
			continue
		}
		switch msgData["type"] {
		case register:
			log.Printf("received register: %s", msgData)
		case deregister:
			log.Printf("received deregister: %s", msgData)
		case update:
			log.Printf("received update: %s", msgData)
		}
	}
}
