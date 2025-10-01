-- +goose Up
CREATE TABLE verification_tokens (
    token_hash VARCHAR(255) PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE verification_tokens;
