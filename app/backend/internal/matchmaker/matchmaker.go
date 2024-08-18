package matchmaker

import (
	"log"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
)

const maxRoomSize = 2

// Matchmaker collects players in a waiting room and once filled,
// will send the room to a GameManager to create a game.
type Matchmaker struct {
	// waitingRoom is where players currently waiting for a match
	// are stored
	waitingRoom []*player.Player

	// ch is the Matchmaker's receiving channel. It will take players
	// from the channel and collect them in the waiting room.
	ch chan *player.Player

	// gmCh is the GameManager channel. Once a waiting room is filled,
	// it is sent through the GameManager channel in order to admit the
	// players to a game.
	gmCh chan []*player.Player
}

// NewMatchmaker returns a new Matchmaker.
//
// It requires the channel from the GameManager service so that it can
// send a filled waiting room over through the channel for a game to be
// made.
// NewMatchmaker also creates the matchmaker's receiving channel which is
// to be passed to a PlayerManager so that players can be sent to the
// Matchmaker.
func NewMatchmaker(gmCh chan []*player.Player) *Matchmaker {
	return &Matchmaker{
		waitingRoom: make([]*player.Player, 0),
		ch:          make(chan *player.Player),
		gmCh:        gmCh,
	}
}

// Run is the main loop of the Matchmaking service. It listens for new
// players and adds them to the waiting room until full, before sending
// them to a game.
func (mm *Matchmaker) Run() {
	for p := range mm.ch {
		// Add player p to the waiting room
		mm.waitingRoom = append(mm.waitingRoom, p)

		// If mm.waitingRoom is full, send it off to be made into a game
		if len(mm.waitingRoom) == maxRoomSize {
			filledRoom := mm.waitingRoom[:maxRoomSize]
			mm.waitingRoom = mm.waitingRoom[maxRoomSize:]
			log.Printf("room filled, creating game")
			mm.gmCh <- filledRoom
		}
	}
}

// Ch returns ch: the Matchmaker's receiving channel. It will take players
// from the channel and collect them in the waiting room.
func (mm *Matchmaker) Ch() chan *player.Player {
	return mm.ch
}
