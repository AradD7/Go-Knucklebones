-- +goose Up
CREATE TABLE boards (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    cells       INTEGER[3][3] DEFAULT ARRAY[]::INTEGER[][],
    player_id   UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE boards;
