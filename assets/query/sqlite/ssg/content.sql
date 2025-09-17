-- Res: Content
-- Table: content

-- Create
INSERT INTO content (
    id, short_id, user_id, section_id, heading, body, status, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :short_id, :user_id, :section_id, :heading, :body, :status, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT id, user_id, section_id, heading, body, status, short_id, created_by, updated_by, created_at, updated_at FROM content;

-- Get
SELECT id, user_id, section_id, heading, body, status, short_id, created_by, updated_by, created_at, updated_at FROM content WHERE id = :id;

-- Update
UPDATE content SET
    user_id = :user_id,
    section_id = :section_id,
    heading = :heading,
    body = :body,
    status = :status,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;

-- Delete
DELETE FROM content WHERE id = :id;

-- GetAllWithTags
SELECT
    c.id, c.user_id, c.section_id, c.heading, c.body, c.status, c.short_id, s.path AS section_path, s.name AS section_name,
    c.created_by, c.updated_by, c.created_at, c.updated_at,
    t.id AS tag_id, t.short_id AS tag_short_id, t.name AS tag_name, t.slug AS tag_slug,
    t.created_by AS tag_created_by, t.updated_by AS tag_updated_by, t.created_at AS tag_created_at, t.updated_at AS tag_updated_at
FROM
    content c
LEFT JOIN
    section s ON c.section_id = s.id
LEFT JOIN
    content_tag ct ON c.id = ct.content_id
LEFT JOIN
    tag t ON ct.tag_id = t.id
ORDER BY
    c.created_at DESC, t.name ASC;