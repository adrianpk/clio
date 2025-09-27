-- Res: ssg
-- Table: param
-- Create
INSERT INTO param (id, name, description, value, ref_key, created_by, updated_by, created_at, updated_at)
VALUES (:id, :name, :description, :value, :ref_key, :created_by, :updated_by, :created_at, :updated_at);

-- Res: ssg
-- Table: param
-- Get
SELECT id, name, description, value, ref_key, created_by, updated_by, created_at, updated_at
FROM param
WHERE id = ?;

-- Res: ssg
-- Table: param
-- GetByName
SELECT id, name, description, value, ref_key, created_by, updated_by, created_at, updated_at
FROM param
WHERE name = ?;

-- Res: ssg
-- Table: param
-- GetByRefKey
SELECT id, name, description, value, ref_key, created_by, updated_by, created_at, updated_at
FROM param
WHERE ref_key = ?;

-- Res: ssg
-- Table: param
-- List
SELECT id, name, description, value, ref_key, created_by, updated_by, created_at, updated_at
FROM param;

-- Res: ssg
-- Table: param
-- Update
UPDATE param
SET name = :name, description = :description, value = :value, ref_key = :ref_key, updated_by = :updated_by, updated_at = :updated_at
WHERE id = :id;

-- Res: ssg
-- Table: param
-- Delete
DELETE FROM param
WHERE id = ?;
