-- name: CreateEntry :one
INSERT INTO entries (
    account_id, amount
)
VALUES (
    $1, $2
)
RETURNING *;
