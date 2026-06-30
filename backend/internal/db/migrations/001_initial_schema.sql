CREATE TABLE services(
    service_group_id BIGINT PRIMARY KEY,
    service_name TEXT NOT NULL
);

CREATE TABLE service_mappings (
    id BIGSERIAL PRIMARY KEY,

    service_group_id BIGINT NOT NULL
        REFERENCES services(service_group_id)
        ON DELETE CASCADE,

    section_name TEXT,
    section_id BIGINT,

    field_id TEXT NOT NULL,
    field_name TEXT NOT NULL,

    input_type TEXT,
    field_set_id BIGINT,

    UNIQUE(service_group_id, field_id)
);

CREATE TABLE workflow_events(
    id BIGSERIAL PRIMARY KEY,

    appl_id BIGINT NOT NULL,

    service_id BIGINT NOT NULL,

    root_type TEXT NOT NULL,

    task_name TEXT,
    action_no INT,
    task_type INT,

    received_time TIMESTAMP,
    executed_time TIMESTAMP,

    raw_payload JSONB NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
);
CREATE OR REPLACE FUNCTION prevent_workflow_update()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'workflow_events is read-only. Updates are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER workflow_events_no_update
BEFORE UPDATE
ON workflow_events
FOR EACH ROW
EXECUTE FUNCTION prevent_workflow_update();

CREATE OR REPLACE FUNCTION prevent_workflow_delete()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'workflow_events is read-only. Deletes are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER workflow_events_no_delete
BEFORE DELETE
ON workflow_events
FOR EACH ROW
EXECUTE FUNCTION prevent_workflow_delete();

CREATE TABLE application_initiated(
    id BIGSERIAL PRIMARY KEY,

    appl_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,

    service_name TEXT,

    appl_ref_no TEXT,

    submission_date TIMESTAMP,

    submission_location TEXT,

    applied_by TEXT,

    payment_mode TEXT,

    amount NUMERIC(12,2),

    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(appl_id, service_id)
);

CREATE TABLE application_execution(
    id BIGSERIAL PRIMARY KEY,

    appl_id BIGINT NOT NULL,

    service_id BIGINT NOT NULL,

    task_name TEXT,

    action_no INT NOT NULL,

    action_taken TEXT,

    task_type INT,

    user_name TEXT,

    designation TEXT,

    location_name TEXT,

    received_time TIMESTAMP,

    executed_time TIMESTAMP,

    remarks TEXT,

    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(
        appl_id,
        service_id,
        action_no
    )
);

CREATE TABLE logs(
    id BIGSERIAL PRIMARY KEY,

    level TEXT NOT NULL,

    source TEXT NOT NULL,

    message TEXT NOT NULL,

    metadata JSONB,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_workflow_appl_id
ON workflow_events(appl_id);

CREATE INDEX idx_workflow_service_id
ON workflow_events(service_id);

CREATE INDEX idx_service_mapping_service_group_id
ON service_mappings(service_group_id);

CREATE INDEX idx_service_mapping_field_id
ON service_mappings(field_id);

CREATE INDEX idx_app_initiated_appl_id
ON application_initiated(appl_id);

CREATE INDEX idx_app_initiated_service_id
ON application_initiated(service_id);

CREATE INDEX idx_app_execution_appl_id
ON application_execution(appl_id);

CREATE INDEX idx_app_execution_service_id
ON application_execution(service_id);

CREATE INDEX idx_app_execution_action_no
ON application_execution(action_no);

CREATE INDEX idx_app_execution_task_name
ON application_execution(task_name);