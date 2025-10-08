-- +goose Up
CREATE TABLE boards (
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    board       INTEGER[3][3] NOT NULL,
    player_id   UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE boards;
