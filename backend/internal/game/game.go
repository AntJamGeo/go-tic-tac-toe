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
	for symbol, p := range g.players {
		var yourTurn string
		if symbol == "x" {
			yourTurn = "true"
		} else {
			yourTurn = "false"
		}
		p.Receive(
			map[string]string{
				"rspType":      "game-Start",
				"gameID":       g.gameID,
				"gameState":    "---------",
				"opponentName": p.Opponent().Name(),
				"yourTurn":     yourTurn,
			},
		)
	}

	for req := range g.channel {
		log.Printf("gameID: %s - game got a request: %v", g.gameID, req)
		switch req["reqType"] {
		case "game-Move":
			symbol := req["player"]
			p := g.players[symbol]
			if g.invalidPlayer(symbol) {
				log.Printf("gameID: %s - player \"%s\" attempted to make a move when it is not their turn", g.gameID, p.Name())
				continue
			}
			op := p.Opponent()
			cell, err := strconv.Atoi(req["cell"])
			if err != nil {
				log.Printf("gameID: %s - received non-numerical value \"%s\" for cell number from player \"%s\"", g.gameID, req["cell"], p.Name())
				continue
			}
			if g.invalidMove(cell) {
				log.Printf("gameID: %s - received invalid move from player \"%s\": cell %d is not a valid cell", g.gameID, p.Name(), cell)
				continue
			}
			g.turns++
			g.gameState = g.gameState[:cell] + symbol + g.gameState[cell+1:]
			if cells := g.gameWon(); cells != "" {
				p.Receive(
					map[string]string{
						"rspType":   "game-Won",
						"gameState": g.gameState,
						"winner":    "true",
						"cells":     cells,
					},
				)
				op.Receive(
					map[string]string{
						"rspType":   "game-Won",
						"gameState": g.gameState,
						"winner":    "false",
						"cells":     cells,
					},
				)
				g.deregister()
			} else if g.turns == 9 {
				p.Receive(
					map[string]string{
						"rspType":   "game-Drawn",
						"gameState": g.gameState,
					},
				)
				op.Receive(
					map[string]string{
						"rspType":   "game-Drawn",
						"gameState": g.gameState,
					},
				)
				g.deregister()
			} else {
				p.Receive(
					map[string]string{
						"rspType":   "game-Update",
						"gameState": g.gameState,
						"yourTurn":  "false",
					},
				)
				op.Receive(
					map[string]string{
						"rspType":   "game-Update",
						"gameState": g.gameState,
						"yourTurn":  "true",
					},
				)
			}
		case "game-Forfeit":
		}
	}
}

func (g *Game) gameWon() (cells string) {
	gs := g.gameState
	for _, combo := range winningCombinations {
		if gs[combo[0]] != '-' && gs[combo[0]] == gs[combo[1]] && gs[combo[1]] == gs[combo[2]] {
			return fmt.Sprintf("%d%d%d", combo[0], combo[1], combo[2])
		}
	}
	return ""
}

func (g *Game) invalidPlayer(symbol string) bool {
	var expectedSymbol string
	if g.turns%2 == 1 {
		expectedSymbol = "o"
	} else {
		expectedSymbol = "x"
	}
	return expectedSymbol != symbol
}

func (g *Game) invalidMove(cell int) bool {
	return cell < 0 || cell > 8 || g.gameState[cell] != '-'
}

func (g *Game) deregister() {

}
