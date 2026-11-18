DROP INDEX IF EXISTS idx_payment_events_id;
DROP INDEX IF EXISTS  idx_payment_events_processed_at;
DROP INDEX IF EXISTS  idx_orders_payment_intent_id;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_intent_id;
DROP TABLE IF EXISTS payment_events;
