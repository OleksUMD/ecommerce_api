ALTER TABLE users ADD COLUMN is_mfa_enabled BOOLEAN default false;
ALTER TABLE users ADD COLUMN email_confirmed BOOLEAN default false;
ALTER TABLE users ADD COLUMN totp_secret VARCHAR(255) NOT NULL;

CREATE TABLE user_mfa_backups (
	id SERIAL PRIMARY KEY,
	user_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	backup_code VARCHAR(255) NOT NULL,
	is_used BOOLEAN DEFAULT false,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_mfa_secret ON users(totp_secret);
CREATE INDEX IF NOT EXISTS idx_users_mfa_backup_code ON user_mfa_backups(backup_code);
