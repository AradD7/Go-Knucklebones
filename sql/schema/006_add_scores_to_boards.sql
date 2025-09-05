-- +goose Up
ALTER TABLE boards
ADD COLUMN score INTEGER DEFAULT 0;

-- +goose Down
ALTER TABLE boards
DROP COLUMN score;
