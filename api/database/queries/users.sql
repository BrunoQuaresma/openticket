-- name: HasFirstUser :one
SELECT EXISTS (SELECT 1 FROM users LIMIT 1);

-- name: CreateUser :one
INSERT INTO users (name, username, email, password_hash, profile_picture_url, role)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserByID :one
UPDATE users
SET name = $2, username = $3, email = $4, profile_picture_url = $5, role = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CountAdmins :one
SELECT COUNT(*) FROM users WHERE role = 'admin';