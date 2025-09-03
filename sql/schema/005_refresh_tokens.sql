-- +goose Up
CREATE TABLE refresh_tokens(
    token       TEXT PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    player_id   UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    expires_at  TIMESTAMP NOT NULL,
    revoked_at  TIMESTAMP
);


-- +goose Down
DROP TABLE refresh_tokens;
