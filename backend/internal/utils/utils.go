package utils

import (
	"encoding/json"
	"log"
	"math/rand"

	"path/filepath"
	"runtime"

	"golang.org/x/net/websocket"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "..", "..", "..")
)
var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func Read(ws *websocket.Conn, buf []byte, req *map[string]string) (ok bool) {
	n, err := ws.Read(buf)
	if err != nil {
		log.Printf("bad read: error receiving request: %s", err)
		return false
	}
	json.Unmarshal(buf[:n], req)
	_, ok = (*req)["reqType"]
	if !ok {
		log.Printf("bad read: found no reqType in request: %v", *req)
		return false
	}
	return true
}

func Write(ws *websocket.Conn, rsp map[string]string) (ok bool) {
	rspBytes, _ := json.Marshal(rsp)
	if _, err := ws.Write(rspBytes); err != nil {
		log.Printf("bad write: failed to send response %v. got error: %s", rsp, err)
		return false
	}
	return true
}
