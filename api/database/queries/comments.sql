-- name: CreateComment :one
INSERT INTO comments (ticket_id, user_id, content, reply_to)
VALUES ($1, $2, $3, $4)
RETURNING *;