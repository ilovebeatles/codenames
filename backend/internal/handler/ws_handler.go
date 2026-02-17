package handler

import (
	"log"
	"net/http"

	"codenames/internal/hub"
	"codenames/internal/storage"

	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

type WSHandler struct {
	hub        *hub.Hub
	playerRepo *storage.PlayerRepo
}

func NewWSHandler(h *hub.Hub, playerRepo *storage.PlayerRepo) *WSHandler {
	return &WSHandler{hub: h, playerRepo: playerRepo}
}

func (h *WSHandler) Handle(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "session_id query param required", http.StatusBadRequest)
		return
	}

	player, err := h.playerRepo.GetBySessionAndRoom(r.Context(), sessionID, roomID)
	if err != nil {
		http.Error(w, "player not found in this room", http.StatusNotFound)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("ws accept: %v", err)
		return
	}

	client := hub.NewClient(conn, h.hub, roomID, sessionID, player.ID)
	h.hub.Register(client)

	ctx := r.Context()
	go client.WritePump(ctx)
	client.ReadPump(ctx)
}
