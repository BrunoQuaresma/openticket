-- name: CreateUser :one
INSERT INTO users (name, username, email, password_hash, profile_picture_url, role)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, username, email, profile_picture_url, created_at, updated_at, role;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;