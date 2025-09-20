-- +goose Up
ALTER TABLE boards ALTER COLUMN board TYPE JSONB USING
  REPLACE(REPLACE(board::TEXT, '{', '['), '}', ']')::JSONB;

-- +goose Down
ALTER TABLE boards ALTER COLUMN board TYPE INTEGER[3][3];
