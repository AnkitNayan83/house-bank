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


-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    email_verified_at = COALESCE(sqlc.narg(email_verified_at), email_verified_at),
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password)
WHERE username = sqlc.narg(username)
RETURNING *; 