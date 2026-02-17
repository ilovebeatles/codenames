package game

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"codenames/internal/model"
	"codenames/internal/storage"
)

type Engine struct {
	gameRepo   *storage.GameRepo
	playerRepo *storage.PlayerRepo
}

func NewEngine(gameRepo *storage.GameRepo, playerRepo *storage.PlayerRepo) *Engine {
	return &Engine{gameRepo: gameRepo, playerRepo: playerRepo}
}

// CanStartGame checks if the room has enough players to start.
func (e *Engine) CanStartGame(players []model.Player) error {
	redSpy, redOp, blueSpy, blueOp := 0, 0, 0, 0
	for _, p := range players {
		switch {
		case p.Team == model.TeamRed && p.Role == model.RoleSpymaster:
			redSpy++
		case p.Team == model.TeamRed && p.Role == model.RoleOperative:
			redOp++
		case p.Team == model.TeamBlue && p.Role == model.RoleSpymaster:
			blueSpy++
		case p.Team == model.TeamBlue && p.Role == model.RoleOperative:
			blueOp++
		}
	}
	if redSpy != 1 {
		return errors.New("red team needs exactly 1 spymaster")
	}
	if blueSpy != 1 {
		return errors.New("blue team needs exactly 1 spymaster")
	}
	if redOp < 1 {
		return errors.New("red team needs at least 1 operative")
	}
	if blueOp < 1 {
		return errors.New("blue team needs at least 1 operative")
	}
	return nil
}

// StartGame creates a new game with a random first team and generated board.
func (e *Engine) StartGame(ctx context.Context, roomID string) (model.Game, []model.Card, error) {
	firstTeam := model.TeamRed
	if rand.Intn(2) == 0 {
		firstTeam = model.TeamBlue
	}

	game, err := e.gameRepo.Create(ctx, roomID, firstTeam)
	if err != nil {
		return model.Game{}, nil, err
	}

	cards := GenerateBoard(game.ID, firstTeam)
	if err := e.gameRepo.CreateCards(ctx, cards); err != nil {
		return model.Game{}, nil, err
	}

	// re-read cards to get IDs
	cards, err = e.gameRepo.GetCardsByGameID(ctx, game.ID)
	if err != nil {
		return model.Game{}, nil, err
	}

	return game, cards, nil
}

// GiveClue sets the current clue and number for the active team.
func (e *Engine) GiveClue(ctx context.Context, game model.Game, clue string, number int) (model.Game, error) {
	if game.Phase != model.PhasePlaying {
		return game, errors.New("game is not in playing phase")
	}
	if game.CurrentClue != "" {
		return game, errors.New("already gave a clue this turn")
	}
	if clue == "" {
		return game, errors.New("clue cannot be empty")
	}
	if number < 0 {
		return game, errors.New("number must be >= 0")
	}

	game.CurrentClue = clue
	game.CurrentNumber = number
	if number == 0 {
		game.GuessesLeft = 25 // unlimited basically
	} else {
		game.GuessesLeft = number + 1
	}

	if err := e.gameRepo.Update(ctx, game); err != nil {
		return game, err
	}
	return game, nil
}

// GuessCard reveals a card and returns the updated game state.
// Returns (game, cards, error). The game may be finished after this.
func (e *Engine) GuessCard(ctx context.Context, game model.Game, cards []model.Card, cardID string, team model.Team) (model.Game, []model.Card, error) {
	if game.Phase != model.PhasePlaying {
		return game, cards, errors.New("game is not in playing phase")
	}
	if game.CurrentClue == "" {
		return game, cards, errors.New("no clue given yet")
	}
	if team != game.CurrentTeam {
		return game, cards, errors.New("not your team's turn")
	}
	if game.GuessesLeft <= 0 {
		return game, cards, errors.New("no guesses left")
	}

	// Find the card
	var card *model.Card
	cardIdx := -1
	for i := range cards {
		if cards[i].ID == cardID {
			card = &cards[i]
			cardIdx = i
			break
		}
	}
	if card == nil {
		return game, cards, errors.New("card not found")
	}
	if card.Revealed {
		return game, cards, errors.New("card already revealed")
	}

	// Reveal the card
	card.Revealed = true
	card.RevealedBy = team
	if err := e.gameRepo.RevealCard(ctx, cardID, team); err != nil {
		return game, cards, err
	}
	cards[cardIdx] = *card

	// Check assassin
	if card.CardType == model.CardTypeAssassin {
		game.Phase = model.PhaseFinished
		game.Winner = team.Opposite()
		if err := e.gameRepo.SetFinished(ctx, game.ID, game.Winner); err != nil {
			return game, cards, err
		}
		return game, cards, nil
	}

	// Check if a team has all their cards revealed
	if winner := checkAllRevealed(cards); winner != "" {
		game.Phase = model.PhaseFinished
		game.Winner = winner
		if err := e.gameRepo.SetFinished(ctx, game.ID, game.Winner); err != nil {
			return game, cards, err
		}
		return game, cards, nil
	}

	// Determine what happens next
	if model.CardType(team) == card.CardType {
		// Correct guess
		game.GuessesLeft--
		if game.GuessesLeft <= 0 {
			e.endTurn(&game)
		}
	} else {
		// Wrong guess (neutral or opponent's card) — end turn
		e.endTurn(&game)
	}

	if err := e.gameRepo.Update(ctx, game); err != nil {
		return game, cards, err
	}
	return game, cards, nil
}

// EndGuessing ends the current team's guessing phase.
func (e *Engine) EndGuessing(ctx context.Context, game model.Game) (model.Game, error) {
	if game.Phase != model.PhasePlaying {
		return game, errors.New("game is not in playing phase")
	}
	if game.CurrentClue == "" {
		return game, errors.New("no clue given yet")
	}
	e.endTurn(&game)
	if err := e.gameRepo.Update(ctx, game); err != nil {
		return game, fmt.Errorf("end guessing: %w", err)
	}
	return game, nil
}

func (e *Engine) endTurn(game *model.Game) {
	game.CurrentTeam = game.CurrentTeam.Opposite()
	game.CurrentClue = ""
	game.CurrentNumber = 0
	game.GuessesLeft = 0
}

func checkAllRevealed(cards []model.Card) model.Team {
	redTotal, redRevealed := 0, 0
	blueTotal, blueRevealed := 0, 0
	for _, c := range cards {
		switch c.CardType {
		case model.CardTypeRed:
			redTotal++
			if c.Revealed {
				redRevealed++
			}
		case model.CardTypeBlue:
			blueTotal++
			if c.Revealed {
				blueRevealed++
			}
		}
	}
	if redTotal > 0 && redRevealed == redTotal {
		return model.TeamRed
	}
	if blueTotal > 0 && blueRevealed == blueTotal {
		return model.TeamBlue
	}
	return ""
}
