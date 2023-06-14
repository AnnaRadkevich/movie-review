ALTER TABLE movies ADD COLUMN version INTEGER NOT NULL DEFAULT 0;

---- create above / drop below ----

ALTER TABLE movies DROP COLUMN version;