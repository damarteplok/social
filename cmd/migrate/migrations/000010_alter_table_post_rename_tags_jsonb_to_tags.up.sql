ALTER TABLE posts DROP COLUMN tags;
ALTER TABLE posts RENAME COLUMN tags_jsonb TO tags;