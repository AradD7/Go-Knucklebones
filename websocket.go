package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from your React dev server
			origin := r.Header.Get("Origin")
			return origin == os.Getenv("FRONTEND_URL")		},
	}
)

type PlayerMessage struct {
	Type 		string `json:"type"`
	Token 		string `json:"token"`
	DisplayName string `json:"display_name"`
	Avatar 		string `json:"avatar"`
	Dice 		int    `json:"dice"`
}

func (cfg apiConfig) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("game_id")
	if gameId == ""{
		respondWithError(w, http.StatusBadRequest, "faild to get gameid from url", nil)
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var msg PlayerMessage
	if err := conn.ReadJSON(&msg); err != nil {
		return
	}

	_, err = auth.ValidateJWT(msg.Token, cfg.tokenSecret)
	if err != nil {
		return
	}

	cfg.gs.addConnection(gameId, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			cfg.gs.removeConnection(gameId, conn)
			return
		}
	}
}

func (gs *gameServer) addConnection(id string, conn *websocket.Conn) {

	gs.rwMux.Lock()
	gs.connections[id] = append(gs.connections[id], conn)
	gs.rwMux.Unlock()
}

func (gs *gameServer) removeConnection(id string, conn *websocket.Conn) {
	gs.rwMux.Lock()
	for i, connection := range gs.connections[id] {
		if connection == conn {
			gs.connections[id] = slices.Delete(gs.connections[id], i, i+1)
			break
		}
	}
	gs.rwMux.Unlock()
}

func (gs *gameServer) broadcastToGame(gameId uuid.UUID) {
	gs.rwMux.RLock()
	for i, conn := range gs.connections[gameId.String()] {
		err := conn.WriteJSON(PlayerMessage{
			Type: "refresh",
		})
		if err != nil {
			fmt.Printf("ERROR sending to connection %d: %v\n", i, err)
		}
	}
	gs.rwMux.RUnlock()
}

func (gs *gameServer) broadcastJoined(gameId uuid.UUID, displayName, avatar string) {
	gs.rwMux.RLock()
	for i, conn := range gs.connections[gameId.String()] {
		err := conn.WriteJSON(PlayerMessage{
			Type: 		 "joined",
			DisplayName: displayName,
			Avatar: 	 avatar,

		})
		if err != nil {
			fmt.Printf("ERROR sending to connection %d: %v\n", i, err)
		}
	}
	gs.rwMux.RUnlock()
}

func (gs *gameServer) broadcastRolled(gameId uuid.UUID, dice int) {
	gs.rwMux.RLock()
	for i, conn := range gs.connections[gameId.String()] {
		err := conn.WriteJSON(PlayerMessage{
			Type: "roll",
			Dice: dice,
		})
		if err != nil {
			fmt.Printf("ERROR sending to connection %d: %v\n", i, err)
		}
	}
	gs.rwMux.RUnlock()
}
