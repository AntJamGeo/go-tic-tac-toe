package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/utils"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(filepath.Join(utils.Root, "frontend"))))
	log.Fatalf("stopped listening: %s", http.ListenAndServe(":3000", nil))
}
