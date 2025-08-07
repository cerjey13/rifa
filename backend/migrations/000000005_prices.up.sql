CREATE TABLE IF NOT EXISTS prices (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid7(),
    bs_amount   NUMERIC(10,2) NOT NULL,
    usd_amount  NUMERIC(10,2) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default values only if table is empty
INSERT INTO prices (bs_amount, usd_amount)
SELECT 150, 1
WHERE NOT EXISTS (
    SELECT 1 FROM prices
);