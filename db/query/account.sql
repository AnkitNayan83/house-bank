-- name: CreateAccount :one
INSERT INTO accounts (owner,balance,currency)
VALUES ($1,$2,$3)
RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByIdForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccounts :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetUsersAccounts :many
SELECT a.*
FROM accounts a
JOIN users u ON u.username = a.owner
WHERE u.username = $1;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;