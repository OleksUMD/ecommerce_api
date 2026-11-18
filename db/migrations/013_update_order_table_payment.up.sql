ALTER TYPE order_status ADD VALUE 'paid';
ALTER TABLE orders ADD COLUMN payment_intent_id VARCHAR(255) UNIQUE;

CREATE TABLE payment_events (
    event_id VARCHAR(255) PRIMARY KEY,
    processed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_orders_payment_intent_id ON orders(payment_intent_id);
CREATE INDEX idx_payment_events_id ON payment_events(event_id);
CREATE INDEX idx_payment_events_processed_at ON payment_events(processed_at);
