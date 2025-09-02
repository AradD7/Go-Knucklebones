-- +goose Up
CREATE TABLE games(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    table_p1    UUID NOT NULL REFERENCES boards(id),
    table_p2    UUID NOT NULL REFERENCES boards(id)
);

-- +goose Down
DROP TABLE games;
