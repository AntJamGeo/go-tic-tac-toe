package player

import (
	"log"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
	"golang.org/x/net/websocket"
)

type PlayerManager struct {
	matchmaker chan *Player
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		matchmaker: make(chan *Player),
	}
}

func (pm *PlayerManager) Run(mmChannel chan *Player) {
	for p := range pm.matchmaker {
		mmChannel <- p
	}
}

func (pm *PlayerManager) HandlePlayer(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	var msg map[string]string
	var msgType string
	playerID := utils.RandSeq(8)

	// Receive player username
	msgType = "connecting"
	if ok := utils.ReadMsg(ws, buf, &msg, msgType); !ok {
		return
	}
	playerName, ok := utils.GetMsgData(&msg, msgType, "playerName")
	if !ok {
		return
	}

	// Create Player
	p := NewPlayer(playerID, playerName)
	log.Printf("created player %s:%s", p.name, p.id)

	// Send for matchmaking
	pm.matchmaker <- p

	// Receive a new game channel and set it so that the player knows where to send messages
	p.playerToGameChannel = <-p.newGameChannel
	log.Printf("%s added to game %s", p.name, p.GameID)

	// Game Play
}
