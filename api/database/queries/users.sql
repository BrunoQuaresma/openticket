-- name: CreateUser :one
INSERT INTO users (name, username, email, hash, profile_picture_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, username, email, profile_picture_url, created_at, updated_at;