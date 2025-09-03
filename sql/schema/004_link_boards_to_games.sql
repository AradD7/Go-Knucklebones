-- +goose Up
ALTER TABLE boards
ADD COLUMN game_id UUID REFERENCES games(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE boards
DROP COLUMN game_id;
