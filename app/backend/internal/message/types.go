package message

const (
	// Player Updates
	PlayerCancel     = "player-Cancel"
	PlayerDisconnect = "player-Disconnect"
	PlayerForfeit    = "player-Forfeit"
	PlayerJoin       = "player-Join"
	PlayerMove       = "player-Move"
	PlayerLeave      = "player-Leave"

	// Opponent Updates
	OpponentDisconnect = "opponent-Disconnect"

	// Game Updates
	GameDrawn  = "game-Drawn"
	GameStart  = "game-Start"
	GameUpdate = "game-Update"
	GameWon    = "game-Won"
)
