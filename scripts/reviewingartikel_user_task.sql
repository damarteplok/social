CREATE TABLE IF NOT EXISTS reviewingartikel (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(256) NOT NULL,
	form_id VARCHAR(256),
	properties JSONB,
	created_by BIGINT NOT NULL,
	updated_by BIGINT,
	created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP(0) WITH TIME ZONE
);

DROP TABLE IF EXISTS reviewingartikel;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_reviewingartikel_properties ON reviewingartikel USING gin (properties);

	