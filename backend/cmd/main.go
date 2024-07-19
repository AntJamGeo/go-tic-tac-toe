package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/game"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/matchmaker"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/player"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
	"golang.org/x/net/websocket"
)

func main() {
	pm := player.NewPlayerManager()
	mm := matchmaker.NewMatchmaker()
	gm := game.NewGameManager()
	mmToGM := make(chan []*player.Player)
	gmToMM := make(chan *game.Game)
	go pm.Run(mm.Channel())
	go mm.Run(mmToGM, gmToMM)
	go gm.Run(mmToGM, gmToMM)
	http.Handle("/", http.FileServer(http.Dir(filepath.Join(utils.Root, "frontend"))))
	http.Handle("/play", websocket.Handler(pm.HandlePlayer))
	log.Fatalf("stopped listening: %s", http.ListenAndServe(":3000", nil))
}
