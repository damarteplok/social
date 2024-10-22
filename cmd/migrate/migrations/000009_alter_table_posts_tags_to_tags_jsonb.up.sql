ALTER TABLE posts ADD COLUMN tags_jsonb jsonb;
UPDATE posts SET tags_jsonb = array_to_json(tags);