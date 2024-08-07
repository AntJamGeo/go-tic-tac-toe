package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/game"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/kafka"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/matchmaker"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/playermanager"
	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
	"golang.org/x/net/websocket"
)

func main() {
	log.Print("CONNECTING\n")
	kafka.Init()
	log.Print("connected\n")
	gm := game.NewGameManager()
	mm := matchmaker.NewMatchmaker(gm.Ch())
	pm := playermanager.NewPlayerManager(mm.Ch())
	go mm.Run()
	go gm.Run()
	http.Handle("/", http.FileServer(http.Dir(filepath.Join(utils.Root, "frontend"))))
	http.Handle("/play", websocket.Handler(pm.HandlePlayer))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
