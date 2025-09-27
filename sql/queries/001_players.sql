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

-- name: GetPlayerByRefreshToken :one
SELECT * FROM players
LEFT JOIN refresh_tokens ON players.id = refresh_tokens.player_id
WHERE refresh_tokens.token = $1;
--

-- name: UpdateProfile :exec
UPDATE players
SET display_name = $2, avatar = $3, updated_at = NOW()
WHERE id = $1;
--

-- name: CreatePlayerWithGoogle :one
INSERT INTO players (id, created_at, updated_at, username, google_id, email, display_name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;
--

-- name: GetPlayerByGoogleId :one
SELECT * FROM players
WHERE google_id = $1;
--
