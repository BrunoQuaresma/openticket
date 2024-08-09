-- name: AssignLabelToTicket :exec
WITH label AS (
    SELECT id
    FROM labels
    WHERE name = @label_name
)
INSERT INTO ticket_labels (ticket_id, label_id)
SELECT @ticket_id, id
FROM label
RETURNING *;

-- name: UnassignLabelFromTicket :exec
DELETE FROM ticket_labels
WHERE ticket_id = @ticket_id
AND label_id = (
  SELECT id
  FROM labels
  WHERE name = @label_name
);

-- name: GetLabels :many
SELECT * FROM labels;

-- name: CreateLabel :one
INSERT INTO labels (name, created_by)
VALUES ($1, $2)
RETURNING *;