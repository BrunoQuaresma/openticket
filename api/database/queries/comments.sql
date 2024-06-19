-- name: CreateComment :one
INSERT INTO comments (ticket_id, user_id, content, reply_to)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

-- name: GetCommentByID :one
SELECT * FROM comments WHERE id = $1 LIMIT 1;