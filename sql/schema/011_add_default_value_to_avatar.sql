-- +goose Up
ALTER TABLE players
ALTER COLUMN avatar SET DEFAULT '008';

-- +goose Down
ALTER TABLE players
ALTER COLUMN avatar DROP DEFAULT;
