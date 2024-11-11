CREATE TABLE IF NOT EXISTS pesen_ke_restorant (
	id BIGSERIAL PRIMARY KEY,
	process_definition_key BIGINT NOT NULL,
	version INT NOT NULL,
	resource_name VARCHAR(256) NOT NULL,
	process_instance_key BIGINT,
	created_by BIGINT NOT NULL,
	updated_by BIGINT,
	created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP(0) WITH TIME ZONE
)

DROP TABLE IF EXISTS pesen_ke_restorant;
	