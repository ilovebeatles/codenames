CREATE TABLE clue_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    team TEXT NOT NULL,
    clue TEXT NOT NULL,
    number INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_clue_history_game_id_created_at ON clue_history(game_id, created_at);
