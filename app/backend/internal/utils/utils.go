package utils

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"

	"path/filepath"
	"runtime"

	"github.com/AntJamGeo/go-tic-tac-toe/backend/internal/message"
	"golang.org/x/net/websocket"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "..", "..", "..")
)
var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq generates a random base-62 sequence of length n
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

// Read takes a websocket connection and reads JSON data sent from it into the provided buffer.
// It verifies that the data sent is OK, and if it is, unmarshals it into req.
func Read(ws *websocket.Conn, buf []byte, req *map[string]string) (ok bool, err error) {
	n, err := ws.Read(buf)
	if err != nil {
		if err == io.EOF {
			log.Printf("client disconnected unexpectedly")
			(*req)["reqType"] = message.PlayerDisconnect
			return true, nil
		}
		log.Printf("bad read: error receiving request: %v", err)
		return false, err
	}
	json.Unmarshal(buf[:n], req)
	_, ok = (*req)["reqType"]
	if !ok {
		log.Printf("bad read: found no reqType in request: %v", *req)
		return false, nil
	}
	return true, nil
}

// Write takes JSON data and sends it through the websocket back to the client.
func Write(ws *websocket.Conn, rsp map[string]string) (ok bool) {
	rspBytes, _ := json.Marshal(rsp)
	if _, err := ws.Write(rspBytes); err != nil {
		log.Printf("bad write: failed to send response %v. got error: %v", rsp, err)
		return false
	}
	return true
}
