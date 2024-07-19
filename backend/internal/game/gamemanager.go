package game

import (
	"log"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
)

type GameManager struct {
	games   map[string]*Game
	gamesMu sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		games:   make(map[string]*Game),
		gamesMu: sync.Mutex{},
	}
}

func (gm *GameManager) Run(mmToGM chan []*player.Player, gmToMM chan *Game) {
	for players := range mmToGM {
		gm.gamesMu.Lock()
		var gameID string
		for {
			gameID = utils.RandSeq(7)
			if _, ok := gm.games[gameID]; !ok {
				break
			}
		}
		g := NewGame(gameID, players)
		gm.games[gameID] = g
		gm.gamesMu.Unlock()
		log.Printf("created game %s", gameID)
		gmToMM <- g
	}
}
