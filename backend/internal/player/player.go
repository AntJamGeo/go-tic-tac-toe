package player

type Player struct {
	id                  string
	name                string
	newGameChannel      chan chan *map[string]string
	gameToPlayerChannel chan *map[string]string
	playerToGameChannel chan *map[string]string
	GameID              string
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		id:                  id,
		name:                name,
		newGameChannel:      make(chan chan *map[string]string),
		gameToPlayerChannel: make(chan *map[string]string),
		playerToGameChannel: nil,
		GameID:              "",
	}
}

func (p *Player) ID() string {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) NewGameChannel() chan chan *map[string]string {
	return p.newGameChannel
}

func (p *Player) GameToPlayerChannel() chan *map[string]string {
	return p.gameToPlayerChannel
}
