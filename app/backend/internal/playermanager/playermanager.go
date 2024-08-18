package playermanager

import (
	"log"
	"maps"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
	"golang.org/x/net/websocket"
)

// PlayerManager handles communication between the client and server,
// reading incoming messages and sending them to the right service, and
// forwarding responses from the server back to the client.
type PlayerManager struct {
	// matchmakerCh is used to forward a player to the matchmaker service
	// to join a game
	matchmakerCh chan *player.Player
}

// NewPlayerManager returns a new PlayerManager.
//
// It requires the channel from a Matchmaker service to send players for
// matchmaking.
func NewPlayerManager(matchmakerCh chan *player.Player) *PlayerManager {
	return &PlayerManager{
		matchmakerCh: matchmakerCh,
	}
}

// HandlePlayer handles websocket connections between the client and the
// server. Requests will be read from the client and responses will be
// sent back through the given websocket.
func (pm *PlayerManager) HandlePlayer(ws *websocket.Conn) {
	p := player.NewPlayer()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go pm.listenToClient(ws, &wg, p)
	go pm.listenToGame(ws, &wg, p)
	wg.Wait()
}

// listenToClient listens for requests from the client and forwards the
// message to the correct service.
func (pm *PlayerManager) listenToClient(ws *websocket.Conn, wg *sync.WaitGroup, p *player.Player) {
	buf := make([]byte, 1024)
	var req map[string]string

listen:
	for {
		clear(req)
		if ok := utils.Read(ws, buf, &req); !ok {
			continue
		}
		switch req["reqType"] {
		case "game-Connect":
			p.SetName(req["playerName"])
			pm.matchmakerCh <- p
		case "game-Cancel":
			log.Printf("TODO: handle game-Cancel client message")
		case "game-Move", "game-Forfeit":
			req["player"] = p.Symbol()
			p.SendToGame(maps.Clone(req))
		case "chat-Message":
		case "chat-Report":
		case "disconnect":
			break listen
		default:
		}
	}
	wg.Done()
}

// listenToGame listens for responses from the current game and forwards the
// response back to the client via the websocket connection.
func (pm *PlayerManager) listenToGame(ws *websocket.Conn, wg *sync.WaitGroup, p *player.Player) {
	p.AwaitGame()
	for rsp := range p.Ch() {
		if ok := utils.Write(ws, rsp); !ok {
			log.Printf("TODO: handle non-ok writes called by pm.listenToGame")
		}
		if rsp["rspType"] == "game-Won" || rsp["rspType"] == "game-Drawn" {
			break
		}
	}
	wg.Done()
}
