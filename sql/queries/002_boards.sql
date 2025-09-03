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
SET game_id = $2
WHERE id = $1;
