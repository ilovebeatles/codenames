CREATE TABLE rooms (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id TEXT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    phase TEXT NOT NULL DEFAULT 'lobby',
    current_team TEXT NOT NULL DEFAULT 'red',
    current_clue TEXT NOT NULL DEFAULT '',
    current_number INT NOT NULL DEFAULT 0,
    guesses_left INT NOT NULL DEFAULT 0,
    winner TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_games_room_id ON games(room_id);

CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id TEXT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    session_id TEXT NOT NULL,
    name TEXT NOT NULL,
    team TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT '',
    is_online BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(room_id, session_id)
);

CREATE INDEX idx_players_room_id ON players(room_id);
CREATE INDEX idx_players_session_id ON players(session_id);

CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    word TEXT NOT NULL,
    card_type TEXT NOT NULL,
    position INT NOT NULL,
    revealed BOOLEAN NOT NULL DEFAULT false,
    revealed_by TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_cards_game_id ON cards(game_id);
