package playermanager

import (
	"log"
	"maps"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/message"
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
	stopper := NewStopper()
	go stopper.Run()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		pm.listenToClient(ws, p, stopper)
		wg.Done()
	}()
	go func() {
		pm.listenToGame(ws, p, stopper)
		wg.Done()
	}()
	wg.Wait()
	log.Printf("successfully released player: %s", p.Name())
}

// listenToClient listens for requests from the client and forwards the
// message to the correct service.
func (pm *PlayerManager) listenToClient(ws *websocket.Conn, p *player.Player, stopper *Stopper) {
	stopperCh := stopper.Register()
	reqCh := make(chan map[string]string)
	connected := true

	go func() {
		defer func() {
			log.Printf("stopped reading")
		}()

		buf := make([]byte, 1024)
		var req map[string]string
		for {
			clear(req)
			ok, err := utils.Read(ws, buf, &req)
			if err != nil {
				stopper.Stop()
				return
			}
			if ok {
				reqCh <- maps.Clone(req)
			}
		}
	}()

listen:
	for {
		select {
		case <-stopperCh:
			log.Printf("stoppedClient")
			break listen
		case req := <-reqCh:
			log.Printf("reqType = %s", req["reqType"])
			switch req["reqType"] {
			case message.PlayerJoin:
				p.SetName(req["playerName"])
				pm.matchmakerCh <- p
			case message.PlayerCancel:
				log.Printf("TODO: handle PlayerCancel client message")
			case message.PlayerMove, message.PlayerForfeit:
				req["player"] = p.Symbol()
				p.SendToGame(req)
			case message.PlayerDisconnect:
				if !connected {
					break
				}
				connected = false
				req["player"] = p.Symbol()
				p.SendToGame(req)
				stopper.Stop()
			case message.PlayerLeave:
				if !connected {
					break
				}
				connected = false
				stopper.Stop()
			default:
				if !connected {
					break
				}
				connected = false
				stopper.Stop()
			}
		}
	}
	log.Printf("brokenClient")
}

// listenToGame listens for responses from the current game and forwards the
// response back to the client via the websocket connection.
func (pm *PlayerManager) listenToGame(ws *websocket.Conn, p *player.Player, stopper *Stopper) {
	stopperCh := stopper.Register()
	p.AwaitGame()

listen:
	for {
		select {
		case <-stopperCh:
			log.Printf("stoppedGame")
			break listen
		case rsp := <-p.GameReceiverCh():
			if ok := utils.Write(ws, rsp); !ok {
				log.Printf("TODO: handle non-ok writes called by pm.listenToGame")
			}
			if rsp["rspType"] == message.GameWon || rsp["rspType"] == message.GameDrawn {
				stopper.Stop()
			}
		}
	}
	log.Printf("brokenGame")
}

type Stopper struct {
	receiveCh chan struct{}
	sendChs   []chan struct{}
	sendChMu  sync.Mutex
	stopped   bool
}

func NewStopper() *Stopper {
	return &Stopper{
		receiveCh: make(chan struct{}),
		sendChs:   make([]chan struct{}, 0),
		sendChMu:  sync.Mutex{},
		stopped:   false,
	}
}

func (s *Stopper) Run() {
	log.Printf("stopper running...")
	<-s.receiveCh
	log.Printf("stop signal received")
	for _, ch := range s.sendChs {
		log.Printf("broadcasting")
		ch <- struct{}{}
	}
	log.Printf("broadcasted")
}

func (s *Stopper) Register() chan struct{} {
	s.sendChMu.Lock()
	defer s.sendChMu.Unlock()
	ch := make(chan struct{})
	s.sendChs = append(s.sendChs, ch)
	return ch
}

func (s *Stopper) Stop() {
	if s.stopped {
		return
	}
	s.stopped = true
	log.Printf("stop signal sending")
	s.receiveCh <- struct{}{}
	log.Printf("stop signal sent")
}
