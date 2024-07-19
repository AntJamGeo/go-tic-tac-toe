package game

import (
	"log"
	"strconv"
)

type Game struct {
	gameID         string
	playerNames    map[string]string
	playerSymbols  map[string]string
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
			g.playerSymbols = map[string]string{playerID: "x", opponentID: "o"}
		case "gameUpdate":
			playerID := (*msg)["playerID"]
			opponentID := g.opponentMap[playerID]
			cell, err := strconv.Atoi((*msg)["cell"])
			if err != nil {
				log.Printf("received non-numerical value for cell number: %s", (*msg)["cell"])
			}
			g.gameState = g.gameState[:cell] + g.playerSymbols[playerID] + g.gameState[cell+1:]
			g.playerChannels[opponentID] <- &map[string]string{"msgType": "gameUpdate", "gameState": g.gameState, "yourTurn": "true"}
			g.playerChannels[playerID] <- &map[string]string{"msgType": "gameUpdate", "gameState": g.gameState, "yourTurn": "false"}
		}
	}
}
