-- name: CreateVerificationToken :one
INSERT INTO verification_tokens (token_hash, player_id, expires_at, created_at)
VALUES (
    $1,
    $2,
    NOW() AT TIME ZONE 'UTC' + INTERVAL '2 hours',
    NOW() AT TIME ZONE 'UTC'
)
RETURNING *;
--

-- name: DeleteVerificationToken :exec
DELETE FROM verification_tokens
WHERE token_hash = $1;
--

-- name: GetVerificationToken :one
SELECT * FROM verification_tokens
WHERE token_hash = $1;
--

-- name: GetVerificationTokenByPlayerId :one
SELECT * FROM verification_tokens
WHERE player_id = $1;
--
