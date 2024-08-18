package player

// Player holds information about a player.
type Player struct {
	// gameReceiverCh is the player's receiving channel, listening out for
	// updates from the game
	gameReceiverCh chan map[string]string

	// gameSenderCh is the channel used by the player to send updates
	// to the game
	gameSenderCh chan map[string]string

	// newGameCh is the channel through which a player can receive
	// a gameCh. When a match has been found, the new game's receiving
	// channel is sent through here so that the player knows where to
	// send its updates.
	newGameCh chan chan map[string]string

	// name is the player's name set by the client.
	name string

	// opponent is the player's opponent in the game.
	opponent *Player

	// symbol is the symbol "x" or "o" that the player is playing as in
	// the game.
	symbol string
}

// NewPlayer returns a new Player.
//
// Initially, only the player's receiving channels are initialised.
// A name is set for the player upon receiving a username via the
// PlayerManager, and the other information needs to be filled in after
// a game has been found for the player.
func NewPlayer() *Player {
	return &Player{
		gameReceiverCh: make(chan map[string]string),
		newGameCh:      make(chan chan map[string]string),
	}
}

// GameReceiverCh returns the player's receiving channel, listening out for
// updates from the game.
func (p *Player) GameReceiverCh() chan map[string]string {
	return p.gameReceiverCh
}

// ReceiveFromGame takes updates from the game so that they can then be
// processed.
func (p *Player) ReceiveFromGame(rsp map[string]string) {
	p.gameReceiverCh <- rsp
}

// SendToGame sends updates from the client to the game.
func (p *Player) SendToGame(req map[string]string) {
	p.gameSenderCh <- req
}

// AwaitGame waits for a new game to be made for the player.
func (p *Player) AwaitGame() {
	p.gameSenderCh = <-p.newGameCh
}

// ConnectToGame connects the player to the game. It takes a newly-created
// game's receiving channel so that the updates from the client can be sent
// to it.
func (p *Player) ConnectToGame(gameSenderCh chan map[string]string) {
	p.newGameCh <- gameSenderCh
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
