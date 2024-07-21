package player

type Player struct {
	readFromGameChannel chan map[string]string
	writeToGameChannel  chan map[string]string
	newGameChannel      chan chan map[string]string
	name                string
	opponent            *Player
	symbol              string
}

func NewPlayer() *Player {
	return &Player{
		readFromGameChannel: make(chan map[string]string),
		newGameChannel:      make(chan chan map[string]string),
	}
}

func (p *Player) ReadFromGameChannel() chan map[string]string {
	return p.readFromGameChannel
}

func (p *Player) WriteToGameChannel() chan map[string]string {
	return p.writeToGameChannel
}

func (p *Player) NewGameChannel() chan chan map[string]string {
	return p.newGameChannel
}

func (p *Player) AwaitGame() {
	p.writeToGameChannel = <-p.newGameChannel
}

func (p *Player) ConnectToGame(gameChannel chan map[string]string) {
	p.newGameChannel <- gameChannel
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) SetName(name string) {
	p.name = name
}

func (p *Player) Opponent() *Player {
	return p.opponent
}

func (p *Player) SetOpponent(opponent *Player) {
	p.opponent = opponent
}

func (p *Player) Symbol() string {
	return p.symbol
}

func (p *Player) SetSymbol(symbol string) {
	p.symbol = symbol
}
