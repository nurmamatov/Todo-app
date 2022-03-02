CREATE UNIQUE INDEX IF NOT EXISTS unique_username ON assignees (username) WHERE deleted_at IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS unique_email ON assignees (email) WHERE deleted_at IS NOT NULL;