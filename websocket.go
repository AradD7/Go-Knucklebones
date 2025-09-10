package main

import (
	"net/http"
	"slices"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type PlayerMessage struct {
	Type 	string `json:"type"`
	Token 	string `json:"token"`
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
	for _, conn := range gs.connections[gameId.String()] {
		conn.WriteJSON(PlayerMessage{
			Type: "refresh",
		})
	}
	gs.rwMux.RUnlock()
}
