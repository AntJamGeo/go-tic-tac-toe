package main

import (
	"log"
	"sync"

	"github.com/AntJamGeo/go-tic-tac-toe/mqtodb/internal/kafka"
	"github.com/AntJamGeo/go-tic-tac-toe/mqtodb/internal/message"
	"github.com/AntJamGeo/go-tic-tac-toe/mqtodb/internal/postgres"
)

func main() {
	ch := make(chan message.Message)
	wg := sync.WaitGroup{}
	wg.Add(2)
	log.Printf("hi")
	go func() {
		kafka.Read(ch)
		wg.Done()
	}()
	go func() {
		postgres.Write(ch)
		wg.Done()
	}()
	log.Printf("started")
	wg.Wait()
}
