-- name: GetLabelByName :one
SELECT * FROM labels WHERE name = $1;

-- name: CreateLabelIfNotExists :one
INSERT INTO labels (name, created_by)
VALUES ($1, $2)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: AssignLabelToTicketIfNotAssigned :exec
INSERT INTO ticket_labels (ticket_id, label_id)
SELECT $1, l.id
FROM labels l
WHERE l.name = $2
AND NOT EXISTS (
  SELECT 1
  FROM ticket_labels tl
  WHERE tl.ticket_id = $1
  AND tl.label_id = l.id
);
