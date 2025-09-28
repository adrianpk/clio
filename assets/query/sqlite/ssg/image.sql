-- Res: ssg
-- Table: image
-- Create
INSERT INTO images (id, short_id, content_hash, mime, width, height, filesize_bytes, etag, title, alt_text, alt_lang, long_description, caption, decorative, described_by_id, created_by, updated_by, created_at, updated_at)
VALUES (:id, :short_id, :content_hash, :mime, :width, :height, :filesize_bytes, :etag, :title, :alt_text, :alt_lang, :long_description, :caption, :decorative, :described_by_id, :created_by, :updated_by, :created_at, :updated_at);

-- Res: ssg
-- Table: image
-- Get
SELECT id, short_id, content_hash, mime, width, height, filesize_bytes, etag, title, alt_text, alt_lang, long_description, caption, decorative, described_by_id, created_by, updated_by, created_at, updated_at
FROM images
WHERE id = ?;

-- Res: ssg
-- Table: image
-- Update
UPDATE images
SET content_hash = :content_hash, mime = :mime, width = :width, height = :height, filesize_bytes = :filesize_bytes, etag = :etag, title = :title, alt_text = :alt_text, alt_lang = :alt_lang, long_description = :long_description, caption = :caption, decorative = :decorative, described_by_id = :described_by_id, updated_by = :updated_by, updated_at = :updated_at
WHERE id = :id;

-- Res: ssg
-- Table: image
-- Delete
DELETE FROM images
WHERE id = ?;

-- Res: ssg
-- Table: image
-- List
SELECT id, short_id, content_hash, mime, width, height, filesize_bytes, etag, title, alt_text, alt_lang, long_description, caption, decorative, described_by_id, created_by, updated_by, created_at, updated_at
FROM images;
