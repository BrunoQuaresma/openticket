CREATE TABLE IF NOT EXISTS assignments (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER REFERENCES tickets (id) NOT NULL,
    user_id INTEGER REFERENCES users (id) NOT NULL,
    assigned_by INTEGER REFERENCES users (id) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, ticket_id)
);