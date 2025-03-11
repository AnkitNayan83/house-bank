-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetTransferById :one
SELECT * FROM transfers
WHERE id = $1;

-- name: GetAllTransferFromAAccount :many
SELECT * FROM transfers
WHERE from_account_id = $1
ORDER BY created_at desc
LIMIT $2
OFFSET $3;

-- name: GetAllTransfersBetweenTwoAccounts :many
SELECT * FROM transfers
WHERE from_account_id = $1 AND to_account_id = $2
ORDER BY created_at desc
LIMIT $2
OFFSET $3;

