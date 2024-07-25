package game

import (
	"fmt"
	"log"
	"strconv"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
)

var winningCombinations = [8][3]int{
	{0, 1, 2}, // Top row
	{3, 4, 5}, // Middle row
	{6, 7, 8}, // Bottom row
	{0, 3, 6}, // Left column
	{1, 4, 7}, // Middle column
	{2, 5, 8}, // Right column
	{0, 4, 8}, // Diagonal from top-left
	{2, 4, 6}, // Diagonal from top-right
}
var symbols = [2]string{"x", "o"}

type Game struct {
	gameID    string
	channel   chan map[string]string
	players   map[string]*player.Player
	gameState string
	turns     int
}

func NewGame(gameID string, players []*player.Player) *Game {
	channel := make(chan map[string]string)
	playersMap := make(map[string]*player.Player)
	for i, p := range players {
		playersMap[symbols[i]] = p
		p.SetOpponent(players[1-i])
		p.SetSymbol(symbols[i])
		defer p.ConnectToGame(channel)
	}
	return &Game{
		gameID:    gameID,
		channel:   channel,
		players:   playersMap,
		gameState: "---------",
		turns:     0,
	}
}

func (g *Game) Run() {
	var rsp map[string]string

	for symbol, p := range g.players {
		var yourTurn string
		if symbol == "x" {
			yourTurn = "true"
		} else {
			yourTurn = "false"
		}
		rsp = map[string]string{
			"rspType":      "game-Start",
			"gameID":       g.gameID,
			"gameState":    "---------",
			"opponentName": p.Opponent().Name(),
			"yourTurn":     yourTurn,
		}
		p.ReadFromGameChannel() <- rsp
	}

	for req := range g.channel {
		log.Printf("game got a request: %v", req)
		switch req["reqType"] {
		case "game-Move":
			g.turns++
			symbol := req["player"]
			p := g.players[symbol]
			op := p.Opponent()
			cell, err := strconv.Atoi(req["cell"])
			if err != nil {
				log.Printf("received non-numerical value for cell number: %s", req["cell"])
			}
			g.gameState = g.gameState[:cell] + symbol + g.gameState[cell+1:]
			if cells := g.gameWon(); cells != "" {
				p.ReadFromGameChannel() <- map[string]string{
					"rspType":   "game-Won",
					"gameState": g.gameState,
					"winner":    "true",
					"cells":     cells,
				}
				op.ReadFromGameChannel() <- map[string]string{
					"rspType":   "game-Won",
					"gameState": g.gameState,
					"winner":    "false",
					"cells":     cells,
				}
				g.deregister()
			} else if g.turns == 9 {
				p.ReadFromGameChannel() <- map[string]string{
					"rspType":   "game-Drawn",
					"gameState": g.gameState,
				}
				op.ReadFromGameChannel() <- map[string]string{
					"rspType":   "game-Drawn",
					"gameState": g.gameState,
				}
				g.deregister()
			} else {
				p.ReadFromGameChannel() <- map[string]string{
					"rspType":   "game-Update",
					"gameState": g.gameState,
					"yourTurn":  "false",
				}
				op.ReadFromGameChannel() <- map[string]string{
					"rspType":   "game-Update",
					"gameState": g.gameState,
					"yourTurn":  "true",
				}
			}
		case "game-Forfeit":
		}
	}
}

func (g *Game) gameWon() (cells string) {
	for _, combo := range winningCombinations {
		if g.gameState[combo[0]] != '-' && g.gameState[combo[0]] == g.gameState[combo[1]] && g.gameState[combo[1]] == g.gameState[combo[2]] {
			return fmt.Sprintf("%d%d%d", combo[0], combo[1], combo[2])
		}
	}
	return ""
}

func (g *Game) deregister() {

}
