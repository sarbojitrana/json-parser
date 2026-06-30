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
			service_group_id,
			service_name
		)
		VALUES($1, $2)
		ON CONFLICT(service_group_id)
		DO UPDATE SET
			service_name = EXCLUDED.service_name
		`,
		service.ServiceGroupID,
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
			service_group_id,
			section_name,
			section_id,
			field_id,
			field_name,
			input_type,
			field_set_id
		)
		VALUES($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT(service_group_id, field_id)
		DO UPDATE SET
			section_name = EXCLUDED.section_name,
			section_id = EXCLUDED.section_id,
			field_name = EXCLUDED.field_name,
			input_type = EXCLUDED.input_type,
			field_set_id = EXCLUDED.field_set_id
		`,
		mapping.ServiceGroupID,
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

func(r *Repository) ServiceGroupExists(
	ctx context.Context,
	serviceGroupID int64,
)(bool,error){
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM services
			WHERE service_group_id = $1
		)	
		`,
		serviceGroupID,
	).Scan(&exists)

	return exists,err
}

func (r *Repository) CreateApplicationInitiated(
	ctx context.Context,
	app model.ApplicationInitiated,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO application_initiated(
			appl_id,
			service_id,
			service_name,
			appl_ref_no,
			submission_date,
			submission_location,
			applied_by,
			payment_mode,
			amount
		)
		VALUES(
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9
		)
		ON CONFLICT(
			appl_id,
			service_id
		)
		DO UPDATE SET
			service_name =
				EXCLUDED.service_name,

			appl_ref_no =
				EXCLUDED.appl_ref_no,

			submission_date =
				EXCLUDED.submission_date,

			submission_location =
				EXCLUDED.submission_location,

			applied_by =
				EXCLUDED.applied_by,

			payment_mode =
				EXCLUDED.payment_mode,

			amount =
				EXCLUDED.amount

		`,
		app.ApplID,
		app.ServiceID,
		app.ServiceName,
		app.ApplRefNo,
		app.SubmissionDate,
		app.SubmissionLocation,
		app.AppliedBy,
		app.PaymentMode,
		app.Amount,
	)

	return err
}

func (r *Repository) CreateApplicationExecution(
	ctx context.Context,
	app model.ApplicationExecution,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO application_execution(
			appl_id,
			service_id,
			task_name,
			action_no,
			action_taken,
			task_type,
			user_name,
			designation,
			location_name,
			received_time,
			executed_time,
			remarks
		)
		VALUES(
			$1,$2,$3,$4,$5,$6,
			$7,$8,$9,$10,$11,
			$12
		)
		ON CONFLICT(
			appl_id,
			service_id,
			action_no
		)
		DO UPDATE SET
			task_name =
				EXCLUDED.task_name,

			action_taken =
				EXCLUDED.action_taken,

			task_type =
				EXCLUDED.task_type,

			user_name =
				EXCLUDED.user_name,

			designation =
				EXCLUDED.designation,

			location_name =
				EXCLUDED.location_name,

			received_time =
				EXCLUDED.received_time,

			executed_time =
				EXCLUDED.executed_time,

			remarks =
				EXCLUDED.remarks
		`,
		app.ApplID,
		app.ServiceID,
		app.TaskName,
		app.ActionNo,
		app.ActionTaken,
		app.TaskType,
		app.UserName,
		app.Designation,
		app.LocationName,
		app.ReceivedTime,
		app.ExecutedTime,
		app.Remarks,
	)

	return err
}

func (r *Repository) GetWorkflowEvents(
	ctx context.Context,
	applID int64,
	serviceID int64,
	rootType string,
) ([]model.WorkflowEvent, error) {

	rows, err := r.db.Query(
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
		ORDER BY action_no ASC, id ASC
		`,
		applID,
		serviceID,
		rootType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.WorkflowEvent

	for rows.Next() {

		var event model.WorkflowEvent

		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}

		events = append(
			events,
			event,
		)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if len(events) == 0 {
		return nil, nil
	}

	return events, nil
}


func (r *Repository) GetApplicationInitiated(
	ctx context.Context,
	applID int64,
	serviceID int64,
) (*model.ApplicationInitiated, error) {

	var app model.ApplicationInitiated

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			appl_id,
			service_id,
			service_name,
			appl_ref_no,
			submission_date,
			submission_location,
			applied_by,
			payment_mode,
			amount,
			created_at
		FROM application_initiated
		WHERE appl_id = $1
		AND service_id = $2
		`,
		applID,
		serviceID,
	).Scan(
		&app.ID,
		&app.ApplID,
		&app.ServiceID,
		&app.ServiceName,
		&app.ApplRefNo,
		&app.SubmissionDate,
		&app.SubmissionLocation,
		&app.AppliedBy,
		&app.PaymentMode,
		&app.Amount,
		&app.CreatedAt,
	)

	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &app, nil
}

func (r *Repository) GetApplicationExecution(
	ctx context.Context,
	applID int64,
	serviceID int64,
	actionNo int,
) (*model.ApplicationExecution, error) {

	var execution model.ApplicationExecution

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			appl_id,
			service_id,
			task_name,
			action_no,
			action_taken,
			task_type,
			user_name,
			designation,
			location_name,
			received_time,
			executed_time,
			remarks,
			created_at
		FROM application_execution
		WHERE appl_id = $1
		AND service_id = $2
		AND action_no = $3
		`,
		applID,
		serviceID,
		actionNo,
	).Scan(
		&execution.ID,
		&execution.ApplID,
		&execution.ServiceID,
		&execution.TaskName,
		&execution.ActionNo,
		&execution.ActionTaken,
		&execution.TaskType,
		&execution.UserName,
		&execution.Designation,
		&execution.LocationName,
		&execution.ReceivedTime,
		&execution.ExecutedTime,
		&execution.Remarks,
		&execution.CreatedAt,
	)

	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &execution, nil
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

//func (r *Repository) DeleteApplication(
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
	)//

	return err//
}

func (r *Repository) CreateLog(
	ctx context.Context,
	log model.Log,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO logs(
			level,
			source,
			message,
			metadata
		)
		VALUES($1,$2,$3,$4)
		`,
		log.Level,
		log.Source,
		log.Message,
		log.Metadata,
	)

	return err
}

func (r *Repository) GetLogs(
	ctx context.Context,
) ([]model.Log, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			level,
			source,
			message,
			metadata,
			created_at
		FROM logs
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.Log

	for rows.Next() {

		var log model.Log

		if err := rows.Scan(
			&log.ID,
			&log.Level,
			&log.Source,
			&log.Message,
			&log.Metadata,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}

		logs = append(
			logs,
			log,
		)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return logs, nil
}