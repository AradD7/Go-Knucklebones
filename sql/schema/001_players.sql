-- +goose Up
CREATE TABLE players(
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL,
    username        VARCHAR(8) NOT NULL UNIQUE,
    avatar          TEXT,
    hashed_password TEXT NOT NULL
);

-- +goose Down
DROP TABLE players;
