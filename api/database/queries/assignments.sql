-- name: CreateAssignment :one
INSERT INTO assignments (ticket_id, user_id, assigned_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteAssignment :exec
DELETE FROM assignments
WHERE id = $1;

-- name: DeleteAssignmentByTicketIDAndUserID :exec
DELETE FROM assignments
WHERE ticket_id = $1 AND user_id = $2;

-- name: GetAssignmentsByTicketID :many
SELECT assignments.*, sqlc.embed(users)
FROM assignments
JOIN users ON assignments.user_id = users.id
WHERE ticket_id = $1;