package game

import "github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"

type Game struct {
	gameID  string
	players []*player.Player
	channel chan *map[string]string
}

func NewGame(gameID string, players []*player.Player) *Game {
	return &Game{
		gameID:  gameID,
		players: players,
		channel: make(chan *map[string]string),
	}
}

func (g *Game) ID() string {
	return g.gameID
}

func (g *Game) Channel() chan *map[string]string {
	return g.channel
}
