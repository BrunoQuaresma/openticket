-- name: GetLabelByName :one
SELECT * FROM labels WHERE name = $1;

-- name: CreateLabel :one
INSERT INTO labels (name, created_by)
VALUES ($1, $2)
RETURNING *;

-- name: AssignLabelToTicket :exec
INSERT INTO ticket_labels (ticket_id, label_id)
VALUES ($1, $2);