package handler

import (
	"encoding/json"
	"net/http"

	"codenames/internal/storage"
)

type PlayerHandler struct {
	playerRepo *storage.PlayerRepo
}

func NewPlayerHandler(playerRepo *storage.PlayerRepo) *PlayerHandler {
	return &PlayerHandler{playerRepo: playerRepo}
}

type createPlayerReq struct {
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
}

func (h *PlayerHandler) Create(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		http.Error(w, "X-Session-ID header required", http.StatusBadRequest)
		return
	}

	var req createPlayerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.RoomID == "" || req.Name == "" {
		http.Error(w, "room_id and name are required", http.StatusBadRequest)
		return
	}

	player, err := h.playerRepo.Upsert(r.Context(), req.RoomID, sessionID, req.Name)
	if err != nil {
		http.Error(w, "failed to create player", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, player)
}
