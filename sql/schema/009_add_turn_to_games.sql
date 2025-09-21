-- +goose Up
ALTER TABLE games
ADD COLUMN player_turn UUID REFERENCES players(id);


-- +goose Down
ALTER TABLE games
DROP COLUMN player_turn
