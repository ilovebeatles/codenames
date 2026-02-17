package model

type Team string

const (
	TeamRed  Team = "red"
	TeamBlue Team = "blue"
)

func (t Team) Opposite() Team {
	if t == TeamRed {
		return TeamBlue
	}
	return TeamRed
}

type Role string

const (
	RoleSpymaster  Role = "spymaster"
	RoleOperative  Role = "operative"
)

type Phase string

const (
	PhaseLobby    Phase = "lobby"
	PhasePlaying  Phase = "playing"
	PhaseFinished Phase = "finished"
)

type CardType string

const (
	CardTypeRed      CardType = "red"
	CardTypeBlue     CardType = "blue"
	CardTypeNeutral  CardType = "neutral"
	CardTypeAssassin CardType = "assassin"
)
