package game

import (
	"fmt"
	"log"
	"strconv"
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

type Game struct {
	gameID         string
	playerNames    map[string]string
	playerSymbols  map[string]byte
	playerChannels map[string]chan *map[string]string
	opponentMap    map[string]string
	channel        chan *map[string]string
	gameState      string
}

func NewGame(gameID string, playerNames map[string]string, playerChannels map[string]chan *map[string]string) *Game {
	opponentMap := make(map[string]string)
	for k := range playerChannels {
		for innerK := range playerChannels {
			if k != innerK {
				opponentMap[k] = innerK
			}
		}
	}
	return &Game{
		gameID:         gameID,
		playerNames:    playerNames,
		playerSymbols:  nil,
		playerChannels: playerChannels,
		opponentMap:    opponentMap,
		channel:        make(chan *map[string]string),
		gameState:      "---------",
	}
}

func (g *Game) ID() string {
	return g.gameID
}

func (g *Game) Channel() chan *map[string]string {
	return g.channel
}

func (g *Game) Run() {
	for msg := range g.channel {
		switch (*msg)["msgType"] {
		case "getGameData":
			playerID := (*msg)["playerID"]
			rsp := map[string]string{"opponentName": g.playerNames[g.opponentMap[playerID]], "gameID": g.gameID}
			g.playerChannels[playerID] <- &rsp
		case "ready":
			playerID := (*msg)["playerID"]
			opponentID := g.opponentMap[playerID]
			for {
				msg = <-g.channel
				if val := (*msg)["msgType"]; val == "ready" {
					if (*msg)["playerID"] == opponentID {
						break
					}
				}
				log.Printf("received unexpected message while waiting for ready confirmation: %s", (*msg))
			}
			log.Printf("both players ready, starting game %s", g.ID())
			g.playerChannels[playerID] <- &map[string]string{"msgType": "gameStart", "gameState": g.gameState, "yourTurn": "true"}
			g.playerChannels[opponentID] <- &map[string]string{"msgType": "gameStart", "gameState": g.gameState, "yourTurn": "false"}
			g.playerSymbols = map[string]byte{playerID: 'x', opponentID: 'o'}
		case "gameUpdate":
			playerID := (*msg)["playerID"]
			opponentID := g.opponentMap[playerID]
			cell, err := strconv.Atoi((*msg)["cell"])
			if err != nil {
				log.Printf("received non-numerical value for cell number: %s", (*msg)["cell"])
			}
			g.gameState = g.gameState[:cell] + string(g.playerSymbols[playerID]) + g.gameState[cell+1:]
			if winnerID, cells := g.gameWon(); winnerID != "" {
				g.playerChannels[winnerID] <- &map[string]string{"msgType": "gameWon", "gameState": g.gameState, "winner": "true", "cells": cells}
				g.playerChannels[g.opponentMap[winnerID]] <- &map[string]string{"msgType": "gameWon", "gameState": g.gameState, "winner": "false", "cells": cells}
				g.deregister()
			} else {
				g.playerChannels[opponentID] <- &map[string]string{"msgType": "gameUpdate", "gameState": g.gameState, "yourTurn": "true"}
				g.playerChannels[playerID] <- &map[string]string{"msgType": "gameUpdate", "gameState": g.gameState, "yourTurn": "false"}
			}
		}
	}
}

func (g *Game) gameWon() (winnerID string, cells string) {
	for _, combo := range winningCombinations {
		if g.gameState[combo[0]] != '-' && g.gameState[combo[0]] == g.gameState[combo[1]] && g.gameState[combo[1]] == g.gameState[combo[2]] {
			for playerID, symbol := range g.playerSymbols {
				if symbol == g.gameState[combo[0]] {
					winnerID = playerID
				} else {
					winnerID = g.opponentMap[playerID]
				}
				break
			}
			return winnerID, fmt.Sprintf("%d%d%d", combo[0], combo[1], combo[2])
		}
	}

	return "", ""
}

func (g *Game) deregister() {

}
