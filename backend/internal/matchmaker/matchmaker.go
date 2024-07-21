package matchmaker

import (
	"log"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
)

const maxRoomSize = 2

type Matchmaker struct {
	waitingRoom   []*player.Player
	waitingRoomMu sync.Mutex
	channel       chan *player.Player
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		waitingRoom:   make([]*player.Player, 0),
		waitingRoomMu: sync.Mutex{},
		channel:       make(chan *player.Player),
	}
}

func (mm *Matchmaker) Run(gmChannel chan []*player.Player) {
	for p := range mm.channel {
		// Add player p to the waiting room
		mm.waitingRoomMu.Lock()
		mm.waitingRoom = append(mm.waitingRoom, p)

		// If mm.waitingRoom is full, send it off to be made into a game
		if len(mm.waitingRoom) == maxRoomSize {
			filledRoom := mm.waitingRoom[:maxRoomSize]
			mm.waitingRoom = mm.waitingRoom[maxRoomSize:]
			mm.waitingRoomMu.Unlock()
			log.Printf("room filled, creating game")
			gmChannel <- filledRoom
		} else {
			mm.waitingRoomMu.Unlock()
		}
	}
}

func (mm *Matchmaker) Channel() chan *player.Player {
	return mm.channel
}
