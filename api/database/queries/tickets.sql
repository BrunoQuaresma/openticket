-- name: CreateTicket :one
INSERT INTO tickets (title, description, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTickets :many
SELECT tickets.*, sqlc.embed(users), array_agg(labels.name)::text[] AS labels
FROM tickets 
JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id 
JOIN labels ON ticket_labels.label_id = labels.id 
JOIN users ON tickets.created_by = users.id
WHERE labels.name = ANY(@labels::text[]) OR @labels::text[] = '{}'
GROUP BY tickets.id, users.id;