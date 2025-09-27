-- +goose Up
ALTER TABLE players
ADD COLUMN google_id TEXT UNIQUE,
ADD COLUMN email TEXT UNIQUE,
ALTER COLUMN hashed_password DROP NOT NULL,
ALTER COLUMN username TYPE VARCHAR(50);

-- +goose Down
ALTER TABLE players
DROP COLUMN google_id,
DROP COLUMN email,
ALTER COLUMN hashed_password SET NOT NULL;
