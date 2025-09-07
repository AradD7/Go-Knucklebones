-- name: CreateBoard :one
INSERT INTO boards (id, created_at, updated_at, board, player_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    ARRAY[]::INTEGER[][],
    $1
)
RETURNING *;
--

-- name: LinkGame :exec
UPDATE boards
SET game_id = $2, updated_at = NOW()
WHERE id = $1;
--

-- name: GetGamesWithPlayerId :many
SELECT game_id
FROM boards
WHERE player_id = $1;
--

-- name: GetPlayerUsernameByBoardId :one
SELECT players.username
FROM boards
LEFT JOIN players ON players.id = boards.player_id
WHERE boards.id = $1;
--

-- name: GetBoardByPlayerIdAndGameId :one
SELECT * FROM boards
WHERE game_id = $1 AND player_id = $2;
--

-- name: GetBoardById :one
SELECT * FROM boards
WHERE id = $1;
--

-- name: UpdateBoard :exec
UPDATE boards
SET board = $2, score = $3, updated_at = NOW()
WHERE id = $1;
--

