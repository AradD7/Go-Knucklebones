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

-- name: GetGamesWithPlayerId :many
SELECT
    g.id AS game_id,
    g.created_at AS date,
    CASE
        WHEN b1.player_id = $1 THEN p2.display_name
        ELSE p1.display_name
    END::TEXT AS opponent_name,
    g.winner AS winner_id
FROM games g
JOIN boards b1 ON g.board1 = b1.id
JOIN boards b2 ON g.board2 = b2.id
JOIN players p1 ON b1.player_id = p1.id
JOIN players p2 ON b2.player_id = p2.id
WHERE b1.player_id = $1 OR b2.player_id = $1;
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
SET board2 = $2, updated_at = NOW()
WHERE id = $1;
--

-- name: SetPlayerTurn :exec
UPDATE games
SET player_turn = $2, updated_at = NOW()
WHERE id = $1;
--

-- name: DeleteEmptyBoardsForPlayer :exec
DELETE FROM games
USING boards
WHERE games.board1 = boards.id
  AND boards.player_id = $1
  AND games.board2 IS NULL;
--
