package matchmaker

import (
	"log"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/game"
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

func (mm *Matchmaker) Run(mmToGM chan []*player.Player, gmToMM chan *game.Game) {
	for p := range mm.channel {
		// Add player p to the waiting room
		mm.waitingRoomMu.Lock()
		mm.waitingRoom = append(mm.waitingRoom, p)
		log.Printf("%s:%s has been admitted to the waiting room", p.Name(), p.ID())

		// If mm.waitingRoom is full, send it off to be made into a game
		if len(mm.waitingRoom) == maxRoomSize {
			filledRoom := mm.waitingRoom[:maxRoomSize]
			mm.waitingRoom = mm.waitingRoom[maxRoomSize:]
			mm.waitingRoomMu.Unlock()
			log.Printf("room filled, creating game")
			mmToGM <- filledRoom
			g := <-gmToMM
			for _, p := range filledRoom {
				go func() {
					p.NewGameChannel() <- g.Channel()
					p.GameID = g.ID()
				}()
			}
		} else {
			mm.waitingRoomMu.Unlock()
		}
	}
}

func (mm *Matchmaker) Channel() chan *player.Player {
	return mm.channel
}
