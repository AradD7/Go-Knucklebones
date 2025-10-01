-- +goose Up
ALTER TABLE players
ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE players
DROP COLUMN email_verified;
