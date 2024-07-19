package game

type Game struct {
	gameID         string
	playerChannels []chan *map[string]string
	channel        chan *map[string]string
}

func NewGame(gameID string, playerChannels []chan *map[string]string) *Game {
	return &Game{
		gameID:         gameID,
		playerChannels: playerChannels,
		channel:        make(chan *map[string]string),
	}
}

func (g *Game) ID() string {
	return g.gameID
}

func (g *Game) Channel() chan *map[string]string {
	return g.channel
}
