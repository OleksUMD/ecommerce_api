DROP INDEX IF EXISTS idx_users_mfa_secret;
DROP INDEX IF EXISTS  idx_users_mfa_backup_code;
ALTER TABLE users DROP COLUMN IF EXISTS is_mfa_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS email_confirmed;
ALTER TABLE users DROP COLUMN IF EXISTS totp_secret;
DROP TABLE IF EXISTS user_mfa_backups;
