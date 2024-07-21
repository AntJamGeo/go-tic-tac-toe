package playermanager

import (
	"log"
	"maps"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
	"golang.org/x/net/websocket"
)

type PlayerManager struct {
	matchmakerChannel chan *player.Player
}

func NewPlayerManager(matchmakerChannel chan *player.Player) *PlayerManager {
	return &PlayerManager{
		matchmakerChannel: matchmakerChannel,
	}
}

func (pm *PlayerManager) HandlePlayer(ws *websocket.Conn) {
	p := player.NewPlayer()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go pm.listenToClient(ws, &wg, p)
	go pm.listenToGame(ws, &wg, p)
	wg.Wait()
}

func (pm *PlayerManager) listenToClient(ws *websocket.Conn, wg *sync.WaitGroup, p *player.Player) {
	buf := make([]byte, 1024)
	var req map[string]string

listen:
	for {
		if ok := utils.Read(ws, buf, &req); !ok {
			return
		}
		log.Printf("got a request : %v", req)
		switch req["reqType"] {
		case "game-Connect":
			p.SetName(req["playerName"])
			log.Printf("sending %s to matchmaker", p.Name())
			pm.matchmakerChannel <- p
			log.Printf("sent %s to matchmaker", p.Name())
		case "game-Cancel":
			panic("TODO")
		case "game-Move", "game-Forfeit":
			req = maps.Clone(req)
			req["player"] = p.Symbol()
			p.WriteToGameChannel() <- req
		case "chat-Message":
		case "chat-Report":
		case "disconnect":
			break listen
		default:
		}
	}
	wg.Done()
	log.Printf("stopping listening to %s", p.Name())
}

func (pm *PlayerManager) listenToGame(ws *websocket.Conn, wg *sync.WaitGroup, p *player.Player) {
	log.Printf("awaiting game start")
	p.AwaitGame()
	log.Printf("game started")
	for rsp := range p.ReadFromGameChannel() {
		if ok := utils.Write(ws, rsp); !ok {
			panic("TODO")
		}
		if rsp["rspType"] == "game-Won" {
			break
		}
	}
	wg.Done()
}
