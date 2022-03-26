create table assignees();

ALTER TABLE assignees
ADD COLUMN refresh_token TEXT,
ADD COLUMN accses_token TEXT;

CREATE INDEX username ON assignees (username);
CREATE INDEX email ON assignees (email);

CREATE UNIQUE INDEX unique_username ON assignees (username) WHERE deleted_at IS NOT NULL;
CREATE UNIQUE INDEX unique_email ON assignees (email) WHERE deleted_at IS NOT NULL;