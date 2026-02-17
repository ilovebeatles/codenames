package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"codenames/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepo struct {
	pool *pgxpool.Pool
}

func NewRoomRepo(pool *pgxpool.Pool) *RoomRepo {
	return &RoomRepo{pool: pool}
}

func (r *RoomRepo) Create(ctx context.Context) (model.Room, error) {
	id, err := generateRoomID()
	if err != nil {
		return model.Room{}, err
	}
	var room model.Room
	err = r.pool.QueryRow(ctx,
		`INSERT INTO rooms (id) VALUES ($1) RETURNING id, created_at`, id,
	).Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		return model.Room{}, fmt.Errorf("create room: %w", err)
	}
	return room, nil
}

func (r *RoomRepo) GetByID(ctx context.Context, id string) (model.Room, error) {
	var room model.Room
	err := r.pool.QueryRow(ctx,
		`SELECT id, created_at FROM rooms WHERE id = $1`, id,
	).Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		return model.Room{}, fmt.Errorf("get room: %w", err)
	}
	return room, nil
}

func generateRoomID() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
