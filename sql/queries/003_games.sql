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
SET winner = $2
WHERE id = $1;
