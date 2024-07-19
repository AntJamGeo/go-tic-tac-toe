package game

type Game struct {
	gameID         string
	playerNames    map[string]string
	opponentMap    map[string]string
	playerChannels map[string]chan *map[string]string
	channel        chan *map[string]string
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
		playerChannels: playerChannels,
		opponentMap:    opponentMap,
		channel:        make(chan *map[string]string),
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
		case "getOpponentName":
			playerID := (*msg)["playerID"]
			rsp := map[string]string{"opponentName": g.playerNames[g.opponentMap[playerID]]}
			g.playerChannels[playerID] <- &rsp
		}
	}
}
