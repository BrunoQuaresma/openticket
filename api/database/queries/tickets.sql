-- name: CreateTicket :one
INSERT INTO tickets (title, created_by)
VALUES ($1, $2)
RETURNING *;

-- name: GetTicketByID :one
SELECT
  tickets.*,
  sqlc.embed(users),
  array_remove(array_agg(DISTINCT labels.name), NULL)::text[] AS labels,
  array_remove(array_agg(DISTINCT assignments.user_id), NULL)::integer[] AS assigned_to
FROM tickets
LEFT JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id
LEFT JOIN labels ON ticket_labels.label_id = labels.id
LEFT JOIN users ON tickets.created_by = users.id
LEFT JOIN assignments ON tickets.id = assignments.ticket_id
WHERE tickets.id = @id
GROUP BY tickets.id, users.id
LIMIT 1;

-- name: DeleteTicketByID :exec
DELETE FROM tickets
WHERE id = @id;

-- name: UpdateTicketByID :one
UPDATE tickets
SET title = $1
WHERE id = @id
RETURNING *;

-- name: UpdateTicketStatusByID :one
UPDATE tickets
SET status = $1
WHERE id = @id
RETURNING *;

-- name: GetTickets :many
SELECT
  tickets.*,
  sqlc.embed(users),
  array_remove(array_agg(DISTINCT labels.name), NULL)::text[] AS labels,
  array_remove(array_agg(DISTINCT assignments.user_id), NULL)::integer[] AS assigned_to
FROM tickets
LEFT JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id
LEFT JOIN labels ON ticket_labels.label_id = labels.id
LEFT JOIN users ON tickets.created_by = users.id
LEFT JOIN assignments ON tickets.id = assignments.ticket_id
WHERE
  CASE 
    WHEN @title::text != '' THEN
      tickets.title ILIKE concat('%', @title, '%')
    ELSE true
  END
  AND CASE
    WHEN @status::text != '' THEN
      tickets.status = @status::ticket_status
    ELSE true
  END
  AND CASE
    WHEN @createdBy::int != 0 THEN
      tickets.created_by = @createdBy
    ELSE true
  END
  AND CASE
    WHEN cardinality(@labels::text[]) > 0 THEN
      labels.name = ANY(@labels)
    ELSE true
  END
GROUP BY tickets.id, users.id;