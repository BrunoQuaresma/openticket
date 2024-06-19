CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER REFERENCES tickets (id) NOT NULL,
    user_id INTEGER REFERENCES users (id) NOT NULL,
    reply_to INTEGER REFERENCES comments (id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);