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
	var rawCards []model.Card
	var redCardsLeft int
	var blueCardsLeft int
	var redClueHistory []model.ClueHistoryEntry
	var blueClueHistory []model.ClueHistoryEntry
	g, err := h.gameRepo.GetActiveByRoomID(r.Context(), id)
	if err == nil {
		game = &g
		rawCards, _ = h.gameRepo.GetCardsByGameID(r.Context(), game.ID)
		clueHistory, _ := h.gameRepo.GetClueHistoryByGameID(r.Context(), game.ID)
		redClueHistory, blueClueHistory = model.SplitClueHistoryByTeam(clueHistory)
	}

	var cards []model.CardView
	if game != nil {
		for _, c := range rawCards {
			cards = append(cards, model.CardToView(c, false))
			if !c.Revealed && c.CardType == model.CardTypeRed {
				redCardsLeft++
			}
			if !c.Revealed && c.CardType == model.CardTypeBlue {
				blueCardsLeft++
			}
		}
	}

	state := model.RoomState{
		Room:            room,
		Players:         players,
		Game:            game,
		Cards:           cards,
		RedCardsLeft:    redCardsLeft,
		BlueCardsLeft:   blueCardsLeft,
		RedClueHistory:  redClueHistory,
		BlueClueHistory: blueClueHistory,
	}
	writeJSON(w, http.StatusOK, state)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
