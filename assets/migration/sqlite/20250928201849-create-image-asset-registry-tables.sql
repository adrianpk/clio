-- +migrate Up

CREATE TABLE images (
    id TEXT PRIMARY KEY,
    short_id TEXT NOT NULL UNIQUE,
    content_hash TEXT NOT NULL UNIQUE,
    mime TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    filesize_bytes INTEGER NOT NULL,
    etag TEXT NOT NULL,

    -- Accessibility fields
    title TEXT,
    alt_text TEXT,
    alt_lang TEXT,
    long_description TEXT,
    caption TEXT,
    decorative BOOLEAN NOT NULL DEFAULT FALSE,
    described_by_id TEXT,

    -- Audit fields
    created_by TEXT NOT NULL,
    updated_by TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    FOREIGN KEY (described_by_id) REFERENCES images(id) ON DELETE SET NULL
);

CREATE INDEX idx_images_short_id ON images (short_id);
CREATE INDEX idx_images_content_hash ON images (content_hash);

CREATE TABLE image_variants (
    id TEXT PRIMARY KEY,
    image_id TEXT NOT NULL,
    kind TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    filesize_bytes INTEGER NOT NULL,
    mime TEXT NOT NULL,
    blob_ref TEXT NOT NULL,

    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    UNIQUE (image_id, kind)
);

CREATE INDEX idx_image_variants_image_id ON image_variants (image_id);

ALTER TABLE layout
ADD COLUMN header_image_id TEXT;

-- Add foreign key constraint after column creation
-- This will be handled in a separate step or by the migration tool if it supports it.
-- For now, just outlining the intent.
-- ALTER TABLE layout
-- ADD CONSTRAINT fk_layout_header_image
-- FOREIGN KEY (header_image_id) REFERENCES images(id) ON DELETE RESTRICT;

-- +migrate Down

DROP TABLE image_variants;
DROP TABLE images;
ALTER TABLE layout
DROP COLUMN header_image_id;
