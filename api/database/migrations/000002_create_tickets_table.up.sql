CREATE TABLE IF NOT EXISTS labels (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  created_by INTEGER REFERENCES users (id) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tickets (
  id SERIAL PRIMARY KEY,
  title VARCHAR(70) NOT NULL,
  description TEXT NOT NULL,
  opened_by INTEGER REFERENCES users (id) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ticket_labels (
  ticket_id INTEGER REFERENCES tickets (id) ON DELETE CASCADE,
  label_id INTEGER REFERENCES labels (id) ON DELETE CASCADE,
  PRIMARY KEY (ticket_id, label_id)
);