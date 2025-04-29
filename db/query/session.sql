-- name: GetSession :one
SELECT * FROM sessions
WHERE ID = $1 
AND is_blocked = false
AND expired_at > now();

-- name: UpdateSession :exec
UPDATE sessions
SET is_blocked = true
WHERE ID = $1;

-- name: CreateSession :one
INSERT INTO sessions (
    id, 
    username, 
    refresh_token, 
    user_agent, 
    client_id, 
    is_blocked, 
    expired_at
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;