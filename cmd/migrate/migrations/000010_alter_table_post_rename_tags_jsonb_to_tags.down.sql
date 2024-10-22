ALTER TABLE posts ADD COLUMN tags jsonb;
ALTER TABLE posts RENAME COLUMN tags TO tags_jsonb;