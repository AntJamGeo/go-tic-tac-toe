package utils

import (
	"encoding/json"
	"errors"
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

func validateMsg(msg *map[string]string) (msgType string, err error) {
	msgType, ok := (*msg)["msgType"]
	if !ok {
		return "", errors.New("msgType not found in msg")
	}
	return msgType, nil
}

func ReadMsg(ws *websocket.Conn, buf []byte, msg *map[string]string, expectedMsgType string) (ok bool) {
	n, err := ws.Read(buf)
	if err != nil {
		log.Printf("bad read: error receiving %s message: %s", expectedMsgType, err)
		return false
	}

	json.Unmarshal(buf[:n], msg)
	msgType, err := validateMsg(msg)
	if err != nil || msgType != expectedMsgType {
		log.Printf("bad read: expected message with msgType=\"%s\", got %s", expectedMsgType, msg)
		return false
	}
	return true
}

func GetMsgData(msg *map[string]string, msgType string, key string) (val string, ok bool) {
	val, ok = (*msg)[key]
	if !ok {
		log.Printf("bad get: no %s found in %s message", key, msgType)
	}
	return val, ok
}

func WriteMsg(ws *websocket.Conn, msg *map[string]string, msgType string, playerName string, playerID string) (ok bool) {
	msgBytes, _ := json.Marshal(msg)
	_, err := ws.Write(msgBytes)
	if err != nil {
		log.Printf("bad write: failed to send %s to %s:%s : %s", msgType, playerName, playerID, err)
		return false
	}
	return true
}
