-- +goose Up
CREATE TABLE players(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    name        VARCHAR(8),
    avatar      TEXT
);

-- +goose Down
DROP TABLE players;
