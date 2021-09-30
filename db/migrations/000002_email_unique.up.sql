CREATE UNIQUE INDEX uuid_unique_idx ON users(email) WHERE archived_at IS NULL;
