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

// Game holds information about a game
type Game struct {
	// gameID is a unique identifier for a game
	gameID string

	// ch is the game's receiving channel, listening out for
	// interactions from the players
	ch chan map[string]string

	// players holds each player identified by their symbol "x" or "o"
	players map[string]*player.Player

	// gameState holds the current state of the game, a string of nine
	// characters, where "-" indicates a cell is empty while "x" or "o"
	// indicate that cell is occupied by "x" or "o"
	gameState string

	// turns is the number of turns played in the game
	turns int
}

// NewGame creates a new game from a filled waiting room.
//
// It creates a receiving channel to listen out for player interactions,
// and connects players to the game by passing the channel to each player.
func NewGame(gameID string, players []*player.Player) *Game {
	ch := make(chan map[string]string)
	playersMap := make(map[string]*player.Player)
	for i, p := range players {
		playersMap[symbols[i]] = p
		p.SetOpponent(players[1-i])
		p.SetSymbol(symbols[i])
		defer p.ConnectToGame(ch)
	}
	return &Game{
		gameID:    gameID,
		ch:        ch,
		players:   playersMap,
		gameState: "---------",
		turns:     0,
	}
}

// Run is the main loop of the Game. It listens for player interactions
// and updates the game state accordingly, before sending out the updated
// game state to each player.
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

	for req := range g.ch {
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
				// If the game has not been won after 9 turns, it is a draw
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

// gameWon checks if the gameState contains one of the winning combinations.
// It returns the winning cells as a concatenated string, or an empty string,
// if the game has not been won.
func (g *Game) gameWon() (cells string) {
	gs := g.gameState
	for _, combo := range winningCombinations {
		if gs[combo[0]] != '-' && gs[combo[0]] == gs[combo[1]] && gs[combo[1]] == gs[combo[2]] {
			return fmt.Sprintf("%d%d%d", combo[0], combo[1], combo[2])
		}
	}
	return ""
}

// invalidPlayer checks if a player is attempting to move when it is not their
// turn.
func (g *Game) invalidPlayer(symbol string) bool {
	var expectedSymbol string
	if g.turns%2 == 1 {
		expectedSymbol = "o"
	} else {
		expectedSymbol = "x"
	}
	return expectedSymbol != symbol
}

// invalidMove checks if a player is attempting to fill a cell that is already
// occupied or not within the range of the grid.
func (g *Game) invalidMove(cell int) bool {
	return cell < 0 || cell > 8 || g.gameState[cell] != '-'
}

// deregister informs the GameManager that the game is no longer active.
func (g *Game) deregister() {
	log.Printf("TODO - deregistering game:%s", g.gameID)
}
