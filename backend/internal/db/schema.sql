CREATE TABLE services(
    service_id BIGINT PRIMARY KEY,
    service_name TEXT NOT NULL
);

CREATE TABLE service_mappings (
    id BIGSERIAL PRIMARY KEY,

    service_id BIGINT NOT NULL
        REFERENCES services(service_id)
        ON DELETE CASCADE,

    section_name TEXT,
    section_id BIGINT,

    field_id TEXT NOT NULL,
    field_name TEXT NOT NULL,

    input_type TEXT,
    field_set_id BIGINT,

    UNIQUE(service_id, field_id)
);

CREATE TABLE workflow_events(
    id BIGSERIAL PRIMARY KEY,

    appl_id BIGINT NOT NULL,

    service_id BIGINT NOT NULL
        REFERENCES services(service_id),

    root_type TEXT NOT NULL,

    task_name TEXT,
    action_no INT,
    task_type INT,

    received_time TIMESTAMP,
    executed_time TIMESTAMP,

    raw_payload JSONB NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_workflow_appl_id
ON workflow_events(appl_id);

CREATE INDEX idx_workflow_service_id
ON workflow_events(service_id);

CREATE INDEX idx_service_mapping_service_id
ON service_mappings(service_id);

CREATE INDEX idx_service_mapping_field_id
ON service_mappings(field_id);