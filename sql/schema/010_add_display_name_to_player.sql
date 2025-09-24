-- +goose Up
ALTER TABLE players
ADD COLUMN display_name TEXT;

-- +goose Down
ALTER TABLE players
DROP COLUMN display_name;
