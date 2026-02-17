package hub

import "codenames/internal/model"

// Client-to-server message types
const (
	MsgJoinTeam    = "join_team"
	MsgSetRole     = "set_role"
	MsgStartGame   = "start_game"
	MsgGiveClue    = "give_clue"
	MsgGuessCard   = "guess_card"
	MsgEndGuessing = "end_guessing"
	MsgNewGame     = "new_game"
)

// Server-to-client message types
const (
	MsgRoomState = "room_state"
	MsgError     = "error"
)

// IncomingMessage is a message from a client.
type IncomingMessage struct {
	Type   string `json:"type"`
	Team   string `json:"team,omitempty"`
	Role   string `json:"role,omitempty"`
	Clue   string `json:"clue,omitempty"`
	Number int    `json:"number,omitempty"`
	CardID string `json:"card_id,omitempty"`
}

// OutgoingMessage is a message to a client.
type OutgoingMessage struct {
	Type  string           `json:"type"`
	State *model.RoomState `json:"state,omitempty"`
	Error string           `json:"error,omitempty"`
}
