-- +goose Up
CREATE TABLE games(
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    board1      UUID NOT NULL REFERENCES boards(id),
    board2      UUID NOT NULL REFERENCES boards(id),
    winner      UUID DEFAULT NULL REFERENCES players(id)
);

-- +goose Down
DROP TABLE games;
