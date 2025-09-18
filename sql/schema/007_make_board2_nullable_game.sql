-- +goose Up
ALTER TABLE games
ALTER COLUMN board2 DROP NOT NULL;

-- +goose Down
ALTER TABLE games
ALTER COLUMN board2 SET NOT NULL;
