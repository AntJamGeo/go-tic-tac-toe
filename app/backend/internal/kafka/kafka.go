package kafka

import (
	"log"
)

func Init() {
	if err := createGameTopic(); err != nil {
		log.Fatalf("failed to create game topic: %v", err)
	}
}

func RegisterGame(gameID string, playerID1 string, playerID2 string) error {
	msg := map[string]string{"type": register, "gameID": gameID, "playerID1": playerID1, "playerID2": playerID2}
	return write(&msg, []byte(gameID))
}

func DeregisterGame(gameID string) error {
	msg := map[string]string{"type": deregister, "gameID": gameID}
	return write(&msg, []byte(gameID))
}

func UpdateGame(gameID string, gameState string) error {
	msg := map[string]string{"type": update, "gameID": gameID, "gameState": gameState}
	return write(&msg, []byte(gameID))
}
