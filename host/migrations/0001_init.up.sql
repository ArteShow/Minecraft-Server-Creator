CREATE TABLE IF NOT EXISTS servers (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL,
    container_id TEXT,
    port INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
