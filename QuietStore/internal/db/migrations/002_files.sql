CREATE TABLE IF NOT EXISTS files (
  id             TEXT PRIMARY KEY,
  owner_user_id  TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  object_key     TEXT NOT NULL,
  original_name  TEXT NOT NULL,
  size_bytes     BIGINT NOT NULL,
  content_type   TEXT,
  sha256         TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at     TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_files_owner_created
  ON files(owner_user_id, created_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS idx_files_object_key
  ON files(object_key);
