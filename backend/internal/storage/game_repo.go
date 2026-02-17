package storage

import (
	"context"
	"fmt"

	"codenames/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepo struct {
	pool *pgxpool.Pool
}

func NewGameRepo(pool *pgxpool.Pool) *GameRepo {
	return &GameRepo{pool: pool}
}

func (r *GameRepo) Create(ctx context.Context, roomID string, firstTeam model.Team) (model.Game, error) {
	var g model.Game
	err := r.pool.QueryRow(ctx, `
		INSERT INTO games (room_id, phase, current_team)
		VALUES ($1, 'playing', $2)
		RETURNING id, room_id, phase, current_team, current_clue, current_number, guesses_left, winner
	`, roomID, firstTeam).Scan(&g.ID, &g.RoomID, &g.Phase, &g.CurrentTeam, &g.CurrentClue, &g.CurrentNumber, &g.GuessesLeft, &g.Winner)
	if err != nil {
		return model.Game{}, fmt.Errorf("create game: %w", err)
	}
	return g, nil
}

func (r *GameRepo) GetActiveByRoomID(ctx context.Context, roomID string) (model.Game, error) {
	var g model.Game
	err := r.pool.QueryRow(ctx, `
		SELECT id, room_id, phase, current_team, current_clue, current_number, guesses_left, winner
		FROM games WHERE room_id = $1 AND phase != 'lobby'
		ORDER BY created_at DESC LIMIT 1
	`, roomID).Scan(&g.ID, &g.RoomID, &g.Phase, &g.CurrentTeam, &g.CurrentClue, &g.CurrentNumber, &g.GuessesLeft, &g.Winner)
	if err != nil {
		return model.Game{}, fmt.Errorf("get active game: %w", err)
	}
	return g, nil
}

func (r *GameRepo) Update(ctx context.Context, g model.Game) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE games SET phase=$2, current_team=$3, current_clue=$4, current_number=$5, guesses_left=$6, winner=$7
		WHERE id=$1
	`, g.ID, g.Phase, g.CurrentTeam, g.CurrentClue, g.CurrentNumber, g.GuessesLeft, g.Winner)
	return err
}

func (r *GameRepo) CreateCards(ctx context.Context, cards []model.Card) error {
	for _, c := range cards {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO cards (game_id, word, card_type, position)
			VALUES ($1, $2, $3, $4)
		`, c.GameID, c.Word, c.CardType, c.Position)
		if err != nil {
			return fmt.Errorf("create card: %w", err)
		}
	}
	return nil
}

func (r *GameRepo) GetCardsByGameID(ctx context.Context, gameID string) ([]model.Card, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, game_id, word, card_type, position, revealed, revealed_by
		FROM cards WHERE game_id = $1 ORDER BY position
	`, gameID)
	if err != nil {
		return nil, fmt.Errorf("get cards: %w", err)
	}
	defer rows.Close()

	var cards []model.Card
	for rows.Next() {
		var c model.Card
		if err := rows.Scan(&c.ID, &c.GameID, &c.Word, &c.CardType, &c.Position, &c.Revealed, &c.RevealedBy); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *GameRepo) RevealCard(ctx context.Context, cardID string, revealedBy model.Team) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE cards SET revealed = true, revealed_by = $2 WHERE id = $1
	`, cardID, revealedBy)
	return err
}

func (r *GameRepo) SetFinished(ctx context.Context, gameID string, winner model.Team) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE games SET phase = 'finished', winner = $2 WHERE id = $1
	`, gameID, winner)
	return err
}
