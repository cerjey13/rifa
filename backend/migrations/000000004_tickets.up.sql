-- Create the lotteries table
CREATE TABLE IF NOT EXISTS lotteries (
    id UUID PRIMARY KEY DEFAULT uuid7(),
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Create the tickets table
CREATE TABLE IF NOT EXISTS tickets (
    id           SERIAL PRIMARY KEY,
    lottery_id   UUID NOT NULL REFERENCES lotteries(id) ON DELETE CASCADE,
    number       INT NOT NULL CHECK (number >= 0 AND number <= 9999),
    user_id      UUID REFERENCES users(id) ON DELETE SET NULL,
    purchase_id  UUID REFERENCES purchases(id) ON DELETE SET NULL,
    reserved_at  TIMESTAMP,
    status       VARCHAR(20) NOT NULL DEFAULT 'available' 
        CHECK (status IN ('available', 'sold')),
    UNIQUE(lottery_id, number)
);

-- Insert the initial lottery and capture its uuid7 id
DO $$
DECLARE
    initial_lottery_id UUID;
BEGIN
    -- Insert the initial lottery
    INSERT INTO lotteries (id, name)
    VALUES (uuid7(), 'Primera rifa')
    RETURNING id INTO initial_lottery_id;

    -- Seed 10,000 tickets for the initial lottery
    INSERT INTO tickets (lottery_id, number)
    SELECT initial_lottery_id, gs.num
    FROM generate_series(0, 9999) AS gs(num);
END
$$;