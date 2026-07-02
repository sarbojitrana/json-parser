CREATE TABLE services(
    service_group_id BIGINT PRIMARY KEY,
    service_name TEXT NOT NULL
);
CREATE OR REPLACE FUNCTION prevent_services_update()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'services is read-only. Updates are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER services_no_update
BEFORE UPDATE
ON services
FOR EACH ROW
EXECUTE FUNCTION prevent_services_update();

CREATE OR REPLACE FUNCTION prevent_services_delete()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'services is read-only. Deletes are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER services_no_delete
BEFORE DELETE
ON services
FOR EACH ROW
EXECUTE FUNCTION prevent_services_delete();

CREATE TABLE service_mappings (
    id BIGSERIAL PRIMARY KEY,

    service_group_id BIGINT NOT NULL
        REFERENCES services(service_group_id),

    section_name TEXT,
    section_id BIGINT,

    field_id TEXT NOT NULL,
    field_name TEXT NOT NULL,

    input_type TEXT,
    field_set_id BIGINT,

    UNIQUE(service_group_id, field_id)
);
CREATE OR REPLACE FUNCTION prevent_service_mappings_update()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'service_mappings is read-only. Updates are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER service_mappings_no_update
BEFORE UPDATE
ON service_mappings
FOR EACH ROW
EXECUTE FUNCTION prevent_service_mappings_update();

CREATE OR REPLACE FUNCTION prevent_service_mappings_delete()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'service_mappings is read-only. Deletes are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER service_mappings_no_delete
BEFORE DELETE
ON service_mappings
FOR EACH ROW
EXECUTE FUNCTION prevent_service_mappings_delete();

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

CREATE TABLE applications (
    id BIGSERIAL PRIMARY KEY,

    appl_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,

    root_type TEXT NOT NULL,

    app_ref_no TEXT,
    service_name TEXT,

    submission_location TEXT,
    submitted_by TEXT,

    submission_date TIMESTAMPTZ,

    status TEXT,
    action_no INT,

    applicant_name TEXT,

    district TEXT,
    district_lgd_code TEXT,

    sub_division TEXT,
    sub_division_lgd_code TEXT,

    block TEXT,
    block_lgd_code TEXT,

    pincode TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE OR REPLACE FUNCTION prevent_applications_update()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'applications is read-only. Updates are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER applications_no_update
BEFORE UPDATE
ON applications
FOR EACH ROW
EXECUTE FUNCTION prevent_applications_update();

CREATE OR REPLACE FUNCTION prevent_applications_delete()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'applications is read-only. Deletes are not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER applications_no_delete
BEFORE DELETE
ON applications
FOR EACH ROW
EXECUTE FUNCTION prevent_applications_delete();

CREATE INDEX idx_applications_submission_date
ON applications(submission_date DESC);

CREATE INDEX idx_applications_app_id
ON applications(appl_id);

CREATE INDEX idx_applications_service_id
ON applications(service_id);

CREATE TABLE logs(
    id BIGSERIAL PRIMARY KEY,

    level TEXT NOT NULL,

    source TEXT NOT NULL,

    message TEXT NOT NULL,

    metadata JSONB,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_workflow_app_id
ON workflow_events(appl_id);

CREATE INDEX idx_workflow_service_id
ON workflow_events(service_id);

CREATE INDEX idx_service_mapping_service_group_id
ON service_mappings(service_group_id);

CREATE INDEX idx_service_mapping_field_id
ON service_mappings(field_id);
