package model

import "time"

type Room struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Player struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	SessionID string `json:"-"`
	Name      string `json:"name"`
	Team      Team   `json:"team"`
	Role      Role   `json:"role"`
	IsOnline  bool   `json:"is_online"`
}

type Game struct {
	ID            string `json:"id"`
	RoomID        string `json:"room_id"`
	Phase         Phase  `json:"phase"`
	CurrentTeam   Team   `json:"current_team"`
	CurrentClue   string `json:"current_clue"`
	CurrentNumber int    `json:"current_number"`
	GuessesLeft   int    `json:"guesses_left"`
	Winner        Team   `json:"winner"`
}

type Card struct {
	ID         string   `json:"id"`
	GameID     string   `json:"game_id"`
	Word       string   `json:"word"`
	CardType   CardType `json:"card_type"`
	Position   int      `json:"position"`
	Revealed   bool     `json:"revealed"`
	RevealedBy Team     `json:"revealed_by"`
}

// RoomState is the full state sent to clients via WebSocket.
type RoomState struct {
	Room    Room     `json:"room"`
	Players []Player `json:"players"`
	Game    *Game    `json:"game"`
	Cards   []CardView `json:"cards"`
}

// CardView is what the client sees — card_type may be hidden for operatives.
type CardView struct {
	ID         string   `json:"id"`
	Word       string   `json:"word"`
	CardType   CardType `json:"card_type"` // empty string for hidden
	Position   int      `json:"position"`
	Revealed   bool     `json:"revealed"`
	RevealedBy Team     `json:"revealed_by"`
}

func CardToView(c Card, showType bool) CardView {
	cv := CardView{
		ID:         c.ID,
		Word:       c.Word,
		Position:   c.Position,
		Revealed:   c.Revealed,
		RevealedBy: c.RevealedBy,
	}
	if showType || c.Revealed {
		cv.CardType = c.CardType
	}
	return cv
}
