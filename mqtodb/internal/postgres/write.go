package postgres

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/AntJamGeo/go-tic-tac-toe/mqtodb/internal/message"
	_ "github.com/lib/pq"
)

func Write(ch chan message.Message) {
	connStr := "user=postgres password=postgres dbname=mydb sslmode=disable host=db port=5432"
	var db *sql.DB
	var err error
	for i := 0; i < numDbConnectRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			break
		}
		log.Printf("failed connection attempt %d/%d to db", i+1, numDbConnectRetries)
		time.Sleep(time.Second * 10)
	}
	defer db.Close()

	var msgData map[string]string
	for msg := range ch {
		if err := json.Unmarshal(msg, &msgData); err != nil {
			log.Printf("unable to unmarshall: %v", err)
			continue
		}
		gameID := msgData["gameID"]
		switch msgData["type"] {
		case register:
			log.Printf("received register: %s", msgData)
			query := `
				INSERT INTO games (game_id, player_id1, player_id2, game_state, live)
				VALUES ($1, $2, $3, $4, $5)
			`
			playerID1 := msgData["playerID1"]
			playerID2 := msgData["playerID2"]
			_, err := db.Exec(query, gameID, playerID1, playerID2, "---------", true)
			if err != nil {
				log.Printf("unable to write data to db: %v", err)
			}
		case deregister:
			log.Printf("received deregister: %s", msgData)
			query := `
				UPDATE games
				SET live = $1
				WHERE game_id = $2
			`
			_, err := db.Exec(query, false, gameID)
			if err != nil {
				log.Printf("unable to write data to db: %v", err)
			}
		case update:
			log.Printf("received update: %s", msgData)
			query := `
				UPDATE games
				SET game_state = $1
				WHERE game_id = $2
			`
			gameState := msgData["gameState"]
			_, err := db.Exec(query, gameState, gameID)
			if err != nil {
				log.Printf("unable to write data to db: %v", err)
			}
		}
	}
}
