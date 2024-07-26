package player

type Player struct {
	ch        chan map[string]string
	gameCh    chan map[string]string
	newGameCh chan chan map[string]string
	name      string
	opponent  *Player
	symbol    string
}

func NewPlayer() *Player {
	return &Player{
		ch:        make(chan map[string]string),
		newGameCh: make(chan chan map[string]string),
	}
}

func (p *Player) Ch() chan map[string]string {
	return p.ch
}

func (p *Player) Receive(rsp map[string]string) {
	p.ch <- rsp
}

func (p *Player) SendToGame(req map[string]string) {
	p.gameCh <- req
}

func (p *Player) AwaitGame() {
	p.gameCh = <-p.newGameCh
}

func (p *Player) ConnectToGame(gameCh chan map[string]string) {
	p.newGameCh <- gameCh
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
