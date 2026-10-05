DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_deleted_at;

DROP TABLE IF EXISTS refresh_tokens;
