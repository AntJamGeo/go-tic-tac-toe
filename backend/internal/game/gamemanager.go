package game

import (
	"log"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
)

// GameManager produces and stores new games.
type GameManager struct {
	// games is where the running games are stored. They are
	// identified by their unique ID.
	games map[string]*Game

	// gamesMu is required as the games map will have concurrent
	// reads and writes.
	gamesMu sync.Mutex

	// ch is the GameManager's receiving channel. It takes a list
	// of players and creates a game for them.
	ch chan []*player.Player
}

// NewGameManager returns a new GameManager.
//
// It creates the GameManager's receiving channel which is to be passed
// to a Matchmaker so that filled waiting rooms can be received for
// game creation.
//
// It also creates a new games map within which all running games are stored.
func NewGameManager() *GameManager {
	return &GameManager{
		games:   make(map[string]*Game),
		gamesMu: sync.Mutex{},
		ch:      make(chan []*player.Player),
	}
}

// Run is the main loop of the GameManager service. It listens for newly
// filled waiting rooms and creates a new game for the waiting players.
func (gm *GameManager) Run() {
	for waitingRoom := range gm.ch {
		gm.gamesMu.Lock()
		var gameID string
		for {
			gameID = utils.RandSeq(7)
			if _, ok := gm.games[gameID]; !ok {
				break
			}
		}
		g := NewGame(gameID, waitingRoom)
		gm.games[gameID] = g
		gm.gamesMu.Unlock()
		go g.Run()
		log.Printf("created game %s", gameID)
	}
}

// Ch returns ch: the GameManager's receiving channel. It takes a list
// of players and creates a game for them.
func (gm *GameManager) Ch() chan []*player.Player {
	return gm.ch
}
