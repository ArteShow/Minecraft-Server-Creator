package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/core"
	"github.com/gorilla/websocket"
)

var consoleUpgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func ServerConsoleWS(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimSpace(r.URL.Query().Get("server_id"))
	if serverID == "" {
		http.Error(w, "server_id is required", http.StatusBadRequest)
		return
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := parts[1]

	conn, err := consoleUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	conn.SetPongHandler(func(_ string) error {
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})

	// Reader goroutine keeps the connection alive by handling client pings/pongs.
	go func() {
		for {
			if _, _, readErr := conn.ReadMessage(); readErr != nil {
				_ = conn.Close()
				return
			}
		}
	}()

	lastPayload := ""
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		snapshot, snapErr := core.GetServerConsoleSnapshot(serverID, token, 250)
		if snapErr != nil {
			msg, _ := json.Marshal(map[string]any{"error": snapErr.Error()})
			if err = conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		} else {
			payloadBytes, _ := json.Marshal(snapshot)
			payload := string(payloadBytes)
			if payload != lastPayload {
				if err = conn.WriteMessage(websocket.TextMessage, payloadBytes); err != nil {
					return
				}
				lastPayload = payload
			}
		}

		if pingErr := conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(10*time.Second)); pingErr != nil {
			return
		}

		<-ticker.C
	}
}
