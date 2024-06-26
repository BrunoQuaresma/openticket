-- name: CreateTicket :one
INSERT INTO tickets (title, description, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTicketsByIDs :many
SELECT tickets.*, sqlc.embed(users), array_agg(labels.name)::text[] AS labels
FROM tickets
JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id
JOIN labels ON ticket_labels.label_id = labels.id
JOIN users ON tickets.created_by = users.id
WHERE tickets.id = ANY(@ids::int[])
GROUP BY tickets.id, users.id;

-- name: GetTicketByID :one
SELECT tickets.*, sqlc.embed(users), array_agg(labels.name)::text[] AS labels
FROM tickets
JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id
JOIN labels ON ticket_labels.label_id = labels.id
JOIN users ON tickets.created_by = users.id
WHERE tickets.id = @id
GROUP BY tickets.id, users.id;

-- name: DeleteTicketByID :exec
DELETE FROM tickets
WHERE id = @id;

-- name: UpdateTicketByID :one
UPDATE tickets
SET title = $1, description = $2
WHERE id = @id
RETURNING *;

-- name: UpdateTicketStatusByID :one
UPDATE tickets
SET status = $1
WHERE id = @id
RETURNING *;