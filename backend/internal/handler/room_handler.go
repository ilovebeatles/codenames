package handler

import (
	"encoding/json"
	"net/http"

	"codenames/internal/model"
	"codenames/internal/storage"

	"github.com/go-chi/chi/v5"
)

type RoomHandler struct {
	roomRepo   *storage.RoomRepo
	playerRepo *storage.PlayerRepo
	gameRepo   *storage.GameRepo
}

func NewRoomHandler(roomRepo *storage.RoomRepo, playerRepo *storage.PlayerRepo, gameRepo *storage.GameRepo) *RoomHandler {
	return &RoomHandler{roomRepo: roomRepo, playerRepo: playerRepo, gameRepo: gameRepo}
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	room, err := h.roomRepo.Create(r.Context())
	if err != nil {
		http.Error(w, "failed to create room", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, room)
}

func (h *RoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	room, err := h.roomRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	players, err := h.playerRepo.GetByRoomID(r.Context(), id)
	if err != nil {
		players = []model.Player{}
	}

	var game *model.Game
	g, err := h.gameRepo.GetActiveByRoomID(r.Context(), id)
	if err == nil {
		game = &g
	}

	var cards []model.CardView
	if game != nil {
		rawCards, _ := h.gameRepo.GetCardsByGameID(r.Context(), game.ID)
		for _, c := range rawCards {
			cards = append(cards, model.CardToView(c, false))
		}
	}

	state := model.RoomState{
		Room:    room,
		Players: players,
		Game:    game,
		Cards:   cards,
	}
	writeJSON(w, http.StatusOK, state)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
