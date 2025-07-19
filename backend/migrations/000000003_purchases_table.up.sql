CREATE TABLE IF NOT EXISTS purchases (
    id UUID PRIMARY KEY DEFAULT uuid7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    monto_bs NUMERIC(12,2),
    monto_usd NUMERIC(12,2),
    payment_method TEXT NOT NULL,
    transaction_digits VARCHAR(6) NOT NULL,
    payment_screenshot BYTEA NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'cancelled')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);