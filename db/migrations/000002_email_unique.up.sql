CREATE UNIQUE INDEX email_unique_idx ON users(email) WHERE archived_at IS NULL;
