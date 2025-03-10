-- name: CreateEntry :one
INSERT INTO entries (
    account_id, amount
)
VALUES (
    $1, $2
)
RETURNING *;

-- name: GetEntriesByAccountId :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetEntryById :one
SELECT * FROM entries
WHERE id = $1;

 