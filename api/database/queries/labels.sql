-- name: AssignLabelToTicket :exec
INSERT INTO ticket_labels (ticket_id, label_id)
VALUES (
  @ticket_id,
  (
    SELECT id
    FROM labels
    WHERE name = @label_name
  )
);

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

