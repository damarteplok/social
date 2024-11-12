CREATE TABLE IF NOT EXISTS kantor_ngetes_id (
	id BIGSERIAL PRIMARY KEY,
	process_definition_key BIGINT NOT NULL,
	version INT NOT NULL,
	resource_name VARCHAR(256) NOT NULL,
	process_instance_key BIGINT,
	task_definition_id VARCHAR(256),
	task_state VARCHAR(20) NOT NULL DEFAULT 'CREATED', 
	created_by BIGINT NOT NULL,
	updated_by BIGINT,
	created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP(0) WITH TIME ZONE,
	CONSTRAINT task_state_check CHECK (task_state IN ('CREATED', 'COMPLETED', 'CANCELED', 'FAILED'))
)

DROP TABLE IF EXISTS kantor_ngetes_id;
	