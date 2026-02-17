package storage

import (
	"context"
	"fmt"

	"codenames/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerRepo struct {
	pool *pgxpool.Pool
}

func NewPlayerRepo(pool *pgxpool.Pool) *PlayerRepo {
	return &PlayerRepo{pool: pool}
}

func (r *PlayerRepo) Upsert(ctx context.Context, roomID, sessionID, name string) (model.Player, error) {
	var p model.Player
	err := r.pool.QueryRow(ctx, `
		INSERT INTO players (room_id, session_id, name, is_online)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (room_id, session_id) DO UPDATE SET name = EXCLUDED.name, is_online = true
		RETURNING id, room_id, session_id, name, team, role, is_online
	`, roomID, sessionID, name).Scan(&p.ID, &p.RoomID, &p.SessionID, &p.Name, &p.Team, &p.Role, &p.IsOnline)
	if err != nil {
		return model.Player{}, fmt.Errorf("upsert player: %w", err)
	}
	return p, nil
}

func (r *PlayerRepo) GetByRoomID(ctx context.Context, roomID string) ([]model.Player, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, room_id, session_id, name, team, role, is_online
		FROM players WHERE room_id = $1
		ORDER BY name
	`, roomID)
	if err != nil {
		return nil, fmt.Errorf("get players: %w", err)
	}
	defer rows.Close()

	var players []model.Player
	for rows.Next() {
		var p model.Player
		if err := rows.Scan(&p.ID, &p.RoomID, &p.SessionID, &p.Name, &p.Team, &p.Role, &p.IsOnline); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (r *PlayerRepo) GetBySessionAndRoom(ctx context.Context, sessionID, roomID string) (model.Player, error) {
	var p model.Player
	err := r.pool.QueryRow(ctx, `
		SELECT id, room_id, session_id, name, team, role, is_online
		FROM players WHERE session_id = $1 AND room_id = $2
	`, sessionID, roomID).Scan(&p.ID, &p.RoomID, &p.SessionID, &p.Name, &p.Team, &p.Role, &p.IsOnline)
	if err != nil {
		return model.Player{}, fmt.Errorf("get player by session: %w", err)
	}
	return p, nil
}

func (r *PlayerRepo) SetTeamRole(ctx context.Context, playerID string, team model.Team, role model.Role) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE players SET team = $2, role = $3 WHERE id = $1
	`, playerID, team, role)
	return err
}

func (r *PlayerRepo) SetOnline(ctx context.Context, playerID string, online bool) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE players SET is_online = $2 WHERE id = $1
	`, playerID, online)
	return err
}

func (r *PlayerRepo) ResetTeamsAndRoles(ctx context.Context, roomID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE players SET team = '', role = '' WHERE room_id = $1
	`, roomID)
	return err
}
