-- +migrate Up
CREATE TABLE content (
    id TEXT PRIMARY KEY,
    short_id TEXT NOT NULL DEFAULT '',
    user_id TEXT NOT NULL,
    section_id TEXT NOT NULL,
    heading TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT '',
    created_by TEXT NOT NULL DEFAULT '',
    updated_by TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE content;
