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
	channel chan []*player.Player
}

func NewGameManager() *GameManager {
	return &GameManager{
		games:   make(map[string]*Game),
		gamesMu: sync.Mutex{},
		channel: make(chan []*player.Player),
	}
}

func (gm *GameManager) Run() {
	for waitingRoom := range gm.channel {
		gm.gamesMu.Lock()
		var gameID string
		for {
			gameID = utils.RandSeq(7)
			if _, ok := gm.games[gameID]; !ok {
				break
			}
		}
		playerChannels := make([]chan *map[string]string, 0)
		for _, p := range waitingRoom {
			playerChannels = append(playerChannels, p.GameToPlayerChannel())
		}
		g := NewGame(gameID, playerChannels)
		gm.games[gameID] = g
		gm.gamesMu.Unlock()
		for _, p := range waitingRoom {
			p.GameID = g.ID()
			p.NewGameChannel() <- g.channel
		}
		log.Printf("created game %s", gameID)
	}
}

func (gm *GameManager) Channel() chan []*player.Player {
	return gm.channel
}
