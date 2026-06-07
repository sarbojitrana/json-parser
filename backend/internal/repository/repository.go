package repository

import (
	"context"

	"parser/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateService(
	ctx context.Context,
	service model.Service,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO services(
			service_id,
			service_name
		)
		VALUES($1, $2)
		ON CONFLICT(service_id)
		DO UPDATE SET
			service_name = EXCLUDED.service_name
		`,
		service.ServiceID,
		service.ServiceName,
	)

	return err
}

func (r *Repository) CreateMapping(
	ctx context.Context,
	mapping model.ServiceMapping,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO service_mappings(
			service_id,
			section_name,
			section_id,
			field_id,
			field_name,
			input_type,
			field_set_id
		)
		VALUES($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT(service_id, field_id)
		DO UPDATE SET
			section_name = EXCLUDED.section_name,
			section_id = EXCLUDED.section_id,
			field_name = EXCLUDED.field_name,
			input_type = EXCLUDED.input_type,
			field_set_id = EXCLUDED.field_set_id
		`,
		mapping.ServiceID,
		mapping.SectionName,
		mapping.SectionID,
		mapping.FieldID,
		mapping.FieldName,
		mapping.InputType,
		mapping.FieldSetID,
	)

	return err
}

func (r *Repository) CreateWorkflowEvent(
	ctx context.Context,
	event model.WorkflowEvent,
) (int64, error) {

	var id int64

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO workflow_events(
			appl_id,
			service_id,
			root_type,
			task_name,
			action_no,
			task_type,
			received_time,
			executed_time,
			raw_payload
		)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
		`,
		event.ApplID,
		event.ServiceID,
		event.RootType,
		event.TaskName,
		event.ActionNo,
		event.TaskType,
		event.ReceivedTime,
		event.ExecutedTime,
		event.RawPayload,
	).Scan(&id)

	return id, err
}

func (r *Repository) GetWorkflowEvent(
	ctx context.Context,
	applID int64,
	serviceID int64,
	rootType string,
) (*model.WorkflowEvent, error) {

	var event model.WorkflowEvent

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			appl_id,
			service_id,
			root_type,
			task_name,
			action_no,
			task_type,
			received_time,
			executed_time,
			raw_payload,
			created_at
		FROM workflow_events
		WHERE appl_id = $1
		AND service_id = $2
		AND root_type = $3
		ORDER BY id DESC
		LIMIT 1
		`,
		applID,
		serviceID,
		rootType,
	).Scan(
		&event.ID,
		&event.ApplID,
		&event.ServiceID,
		&event.RootType,
		&event.TaskName,
		&event.ActionNo,
		&event.TaskType,
		&event.ReceivedTime,
		&event.ExecutedTime,
		&event.RawPayload,
		&event.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &event, nil
}



func (r *Repository) GetMappingsByServiceID(
	ctx context.Context,
	serviceID int64,
) (map[string]string, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			field_id,
			field_name
		FROM service_mappings
		WHERE service_id = $1
		`,
		serviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mappings := make(map[string]string)

	for rows.Next() {
		var fieldID string
		var fieldName string

		if err := rows.Scan(
			&fieldID,
			&fieldName,
		); err != nil {
			return nil, err
		}

		mappings[fieldID] = fieldName
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return mappings, nil
}

func (r *Repository) DeleteApplication(
	ctx context.Context,
	applID int64,
	serviceID int64,
	rootType string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		DELETE
		FROM workflow_events
		WHERE appl_id = $1
		AND service_id = $2
		AND root_type = $3
		`,
		applID,
		serviceID,
		rootType,
	)

	return err
}