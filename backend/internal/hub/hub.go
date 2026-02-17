package hub

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"codenames/internal/game"
	"codenames/internal/model"
	"codenames/internal/storage"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client

	roomRepo   *storage.RoomRepo
	playerRepo *storage.PlayerRepo
	gameRepo   *storage.GameRepo
	engine     *game.Engine
}

func NewHub(roomRepo *storage.RoomRepo, playerRepo *storage.PlayerRepo, gameRepo *storage.GameRepo, engine *game.Engine) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		roomRepo:   roomRepo,
		playerRepo: playerRepo,
		gameRepo:   gameRepo,
		engine:     engine,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.mu.Unlock()

			ctx := context.Background()
			_ = h.playerRepo.SetOnline(ctx, client.playerID, true)
			h.broadcastRoomState(ctx, client.roomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.roomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
				}
				if len(clients) == 0 {
					delete(h.rooms, client.roomID)
				}
			}
			h.mu.Unlock()

			ctx := context.Background()
			_ = h.playerRepo.SetOnline(ctx, client.playerID, false)
			h.broadcastRoomState(ctx, client.roomID)
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) HandleMessage(ctx context.Context, client *Client, msg IncomingMessage) {
	switch msg.Type {
	case MsgJoinTeam:
		h.handleJoinTeam(ctx, client, msg)
	case MsgSetRole:
		h.handleSetRole(ctx, client, msg)
	case MsgStartGame:
		h.handleStartGame(ctx, client)
	case MsgGiveClue:
		h.handleGiveClue(ctx, client, msg)
	case MsgGuessCard:
		h.handleGuessCard(ctx, client, msg)
	case MsgEndGuessing:
		h.handleEndGuessing(ctx, client)
	case MsgNewGame:
		h.handleNewGame(ctx, client)
	default:
		client.SendError("unknown message type: " + msg.Type)
	}
}

func (h *Hub) handleJoinTeam(ctx context.Context, client *Client, msg IncomingMessage) {
	team := model.Team(msg.Team)
	if team != model.TeamRed && team != model.TeamBlue && team != "" {
		client.SendError("invalid team")
		return
	}
	role := model.Role(msg.Role)
	if err := h.playerRepo.SetTeamRole(ctx, client.playerID, team, role); err != nil {
		client.SendError("failed to join team")
		return
	}
	h.broadcastRoomState(ctx, client.roomID)
}

func (h *Hub) handleSetRole(ctx context.Context, client *Client, msg IncomingMessage) {
	player, err := h.playerRepo.GetBySessionAndRoom(ctx, client.sessionID, client.roomID)
	if err != nil {
		client.SendError("player not found")
		return
	}
	role := model.Role(msg.Role)
	if role != model.RoleSpymaster && role != model.RoleOperative {
		client.SendError("invalid role")
		return
	}
	if err := h.playerRepo.SetTeamRole(ctx, client.playerID, player.Team, role); err != nil {
		client.SendError("failed to set role")
		return
	}
	h.broadcastRoomState(ctx, client.roomID)
}

func (h *Hub) handleStartGame(ctx context.Context, client *Client) {
	players, err := h.playerRepo.GetByRoomID(ctx, client.roomID)
	if err != nil {
		client.SendError("failed to get players")
		return
	}
	if err := h.engine.CanStartGame(players); err != nil {
		client.SendError(err.Error())
		return
	}
	_, _, err = h.engine.StartGame(ctx, client.roomID)
	if err != nil {
		client.SendError("failed to start game")
		return
	}
	h.broadcastRoomState(ctx, client.roomID)
}

func (h *Hub) handleGiveClue(ctx context.Context, client *Client, msg IncomingMessage) {
	player, err := h.playerRepo.GetBySessionAndRoom(ctx, client.sessionID, client.roomID)
	if err != nil {
		client.SendError("player not found")
		return
	}
	if player.Role != model.RoleSpymaster {
		client.SendError("only spymasters can give clues")
		return
	}

	g, err := h.gameRepo.GetActiveByRoomID(ctx, client.roomID)
	if err != nil {
		client.SendError("no active game")
		return
	}
	if player.Team != g.CurrentTeam {
		client.SendError("not your team's turn")
		return
	}

	g, err = h.engine.GiveClue(ctx, g, msg.Clue, msg.Number)
	if err != nil {
		client.SendError(err.Error())
		return
	}
	h.broadcastRoomState(ctx, client.roomID)
}

func (h *Hub) handleGuessCard(ctx context.Context, client *Client, msg IncomingMessage) {
	player, err := h.playerRepo.GetBySessionAndRoom(ctx, client.sessionID, client.roomID)
	if err != nil {
		client.SendError("player not found")
		return
	}
	if player.Role != model.RoleOperative {
		client.SendError("only operatives can guess")
		return
	}

	g, err := h.gameRepo.GetActiveByRoomID(ctx, client.roomID)
	if err != nil {
		client.SendError("no active game")
		return
	}

	cards, err := h.gameRepo.GetCardsByGameID(ctx, g.ID)
	if err != nil {
		client.SendError("failed to get cards")
		return
	}

	g, cards, err = h.engine.GuessCard(ctx, g, cards, msg.CardID, player.Team)
	if err != nil {
		client.SendError(err.Error())
		return
	}
	h.broadcastRoomState(ctx, client.roomID)
}

func (h *Hub) handleEndGuessing(ctx context.Context, client *Client) {
	player, err := h.playerRepo.GetBySessionAndRoom(ctx, client.sessionID, client.roomID)
	if err != nil {
		client.SendError("player not found")
		return
	}

	g, err := h.gameRepo.GetActiveByRoomID(ctx, client.roomID)
	if err != nil {
		client.SendError("no active game")
		return
	}
	if player.Team != g.CurrentTeam {
		client.SendError("not your team's turn")
		return
	}

	_, err = h.engine.EndGuessing(ctx, g)
	if err != nil {
		client.SendError(err.Error())
		return
	}
	h.broadcastRoomState(ctx, client.roomID)
}

func (h *Hub) handleNewGame(ctx context.Context, client *Client) {
	// Reset the game state — go back to lobby
	if err := h.playerRepo.ResetTeamsAndRoles(ctx, client.roomID); err != nil {
		client.SendError("failed to reset")
		return
	}
	h.broadcastRoomState(ctx, client.roomID)
}

func (h *Hub) broadcastRoomState(ctx context.Context, roomID string) {
	room, err := h.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		log.Printf("broadcast: get room: %v", err)
		return
	}

	players, err := h.playerRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		log.Printf("broadcast: get players: %v", err)
		return
	}

	var g *model.Game
	var cards []model.Card
	activeGame, err := h.gameRepo.GetActiveByRoomID(ctx, roomID)
	if err == nil {
		g = &activeGame
		cards, _ = h.gameRepo.GetCardsByGameID(ctx, activeGame.ID)
	}

	h.mu.RLock()
	clients := h.rooms[roomID]
	h.mu.RUnlock()

	for client := range clients {
		// Determine if this player is a spymaster
		isSpymaster := false
		for _, p := range players {
			if p.SessionID == client.sessionID && p.Role == model.RoleSpymaster {
				isSpymaster = true
				break
			}
		}

		// Build card views
		var cardViews []model.CardView
		for _, c := range cards {
			showType := isSpymaster || g.Phase == model.PhaseFinished
			cardViews = append(cardViews, model.CardToView(c, showType))
		}

		state := &model.RoomState{
			Room:    room,
			Players: players,
			Game:    g,
			Cards:   cardViews,
		}

		data, err := json.Marshal(OutgoingMessage{Type: MsgRoomState, State: state})
		if err != nil {
			continue
		}
		select {
		case client.send <- data:
		default:
		}
	}
}
