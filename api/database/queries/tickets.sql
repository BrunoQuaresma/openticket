-- name: CreateTicket :one
INSERT INTO tickets (title, description, created_by)
VALUES ($1, $2, $3)
RETURNING *;