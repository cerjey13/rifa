CREATE TABLE IF NOT EXISTS idempotency_keys (
    key TEXT PRIMARY KEY,
    body_hash TEXT NOT NULL,
    status_code INT NOT NULL,
    response_body BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);