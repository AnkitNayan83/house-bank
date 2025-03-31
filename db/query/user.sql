-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: UpdateUserEmailVerification :one
UPDATE users 
set email_verified_at = $2
WHERE username = $1
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ChangePassword :one
UPDATE users 
SET hashed_password = $2, password_changed_at = $3
WHERE username = $1
RETURNING *;