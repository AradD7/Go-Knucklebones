-- name: CreatePlayer :one
INSERT INTO players (id, created_at, updated_at, username, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;
--

-- name: GetPlayerByUsername :one
SELECT * FROM players
WHERE username = $1;
--

-- name: GetPlayerByPlayerId :one
SELECT * FROM players
WHERE id = $1;
--
