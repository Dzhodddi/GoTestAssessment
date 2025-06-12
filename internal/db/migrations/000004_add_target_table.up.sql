CREATE TABLE IF NOT EXISTS targets (
    id bigserial PRIMARY KEY,
    mission_id INTEGER NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    notes TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE
);