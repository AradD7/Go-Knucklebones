-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, player_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() AT TIME ZONE 'UTC' + INTERVAL '7 days'
)
RETURNING *;
--

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;
--

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW() AT TIME ZONE 'UTC', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE token = $1;
--

-- name: GetRefreshTokenFromPlayerId :one
SELECT * FROM refresh_tokens
WHERE player_id = $1 AND (revoked_at != NULL OR expires_at >= NOW() AT TIME ZONE 'UTC');
--

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;
--
