-- name: CreateNewGame :one
INSERT INTO games(id, created_at, updated_at, board1, board2)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: SetGameWinner :exec
UPDATE games
SET winner = $2, updated_at = NOW()
WHERE id = $1;
--

-- name: GetGameById :one
SELECT * FROM games
WHERE id = $1;
--

-- name: UpdateGame :exec
UPDATE games
SET updated_at = NOW()
WHERE id = $1;
--

-- name: JoinGame :exec
UPDATE games
SET board2 = $2
WHERE id = $1;
--
