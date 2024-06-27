-- name: CreateAssignment :one
INSERT INTO assignments (ticket_id, user_id, assigned_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteAssignment :exec
DELETE FROM assignments
WHERE id = $1;