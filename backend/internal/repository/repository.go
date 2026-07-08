package repository

import (
	"context"
	"fmt"
	"time"

	"parser/internal/config"
	"parser/internal/model"
	"parser/internal/security"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func New(db *pgxpool.Pool, cfg *config.Config) *Repository {
	return &Repository{
		db:  db,
		cfg: cfg,
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
		DO NOTHING
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
		DO NOTHING
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

	if enc, err := security.EncryptPayload(string(event.RawPayload), r.cfg.Security.SecretKey); err == nil {
		event.RawPayload = []byte(enc)
	}

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

func (r *Repository) CreateApplication(
	ctx context.Context,
	app model.Application,
) error {

	app.AppRefNo = encryptPtr(app.AppRefNo, r.cfg.Security.SecretKey)
	app.ServiceName = encryptPtr(app.ServiceName, r.cfg.Security.SecretKey)
	app.SubmissionLocation = encryptPtr(app.SubmissionLocation, r.cfg.Security.SecretKey)
	app.SubmittedBy = encryptPtr(app.SubmittedBy, r.cfg.Security.SecretKey)
	app.Status = encryptPtr(app.Status, r.cfg.Security.SecretKey)
	app.ApplicantName = encryptPtr(app.ApplicantName, r.cfg.Security.SecretKey)
	app.District = encryptPtr(app.District, r.cfg.Security.SecretKey)
	app.DistrictLGDCode = encryptPtr(app.DistrictLGDCode, r.cfg.Security.SecretKey)
	app.SubDivision = encryptPtr(app.SubDivision, r.cfg.Security.SecretKey)
	app.SubDivisionLGDCode = encryptPtr(app.SubDivisionLGDCode, r.cfg.Security.SecretKey)
	app.Block = encryptPtr(app.Block, r.cfg.Security.SecretKey)
	app.BlockLGDCode = encryptPtr(app.BlockLGDCode, r.cfg.Security.SecretKey)
	app.Pincode = encryptPtr(app.Pincode, r.cfg.Security.SecretKey)

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO applications(
			root_type,
			appl_id,
			app_ref_no,
			service_id,
			service_name,
			submission_location,
			submitted_by,
			submission_date,
			status,
			action_no,
			applicant_name,
			district,
			district_lgd_code,
			sub_division,
			sub_division_lgd_code,
			block,
			block_lgd_code,
			pincode
		)
		VALUES(
			$1,$2,$3,$4,$5,$6,$7,$8,
			$9,$10,$11,$12,$13,$14,
			$15,$16,$17,$18
		)
		`,
		app.RootType,
		app.ApplID,
		app.AppRefNo,
		app.ServiceID,
		app.ServiceName,
		app.SubmissionLocation,
		app.SubmittedBy,
		app.SubmissionDate,
		app.Status,
		app.ActionNo,
		app.ApplicantName,
		app.District,
		app.DistrictLGDCode,
		app.SubDivision,
		app.SubDivisionLGDCode,
		app.Block,
		app.BlockLGDCode,
		app.Pincode,
	)

	return err
}

func (r *Repository) GetApplications(
	ctx context.Context,
	from time.Time,
	to time.Time,
	page int,
	limit int,
) (*model.PaginatedResponse[model.Application], error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM applications
		WHERE submission_date >= $1
		AND submission_date <= $2
		`,
		from,
		to,
	).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			root_type,
			appl_id,
			app_ref_no,
			service_id,
			service_name,
			submission_location,
			submitted_by,
			submission_date,
			status,
			action_no,
			applicant_name,
			district,
			district_lgd_code,
			sub_division,
			sub_division_lgd_code,
			block,
			block_lgd_code,
			pincode
		FROM applications
		WHERE submission_date >= $1
		AND submission_date <= $2
		ORDER BY
    		submission_date DESC,
    		appl_id DESC,
    		action_no ASC
		LIMIT $3
		OFFSET $4
		`,
		from,
		to,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []model.Application

	for rows.Next() {

		var app model.Application

		if err := rows.Scan(
			&app.RootType,
			&app.ApplID,
			&app.AppRefNo,
			&app.ServiceID,
			&app.ServiceName,
			&app.SubmissionLocation,
			&app.SubmittedBy,
			&app.SubmissionDate,
			&app.Status,
			&app.ActionNo,
			&app.ApplicantName,
			&app.District,
			&app.DistrictLGDCode,
			&app.SubDivision,
			&app.SubDivisionLGDCode,
			&app.Block,
			&app.BlockLGDCode,
			&app.Pincode,
		); err != nil {
			return nil, err
		}

		app.AppRefNo = decryptPtr(app.AppRefNo, r.cfg.Security.SecretKey)
		app.ServiceName = decryptPtr(app.ServiceName, r.cfg.Security.SecretKey)
		app.SubmissionLocation = decryptPtr(app.SubmissionLocation, r.cfg.Security.SecretKey)
		app.SubmittedBy = decryptPtr(app.SubmittedBy, r.cfg.Security.SecretKey)
		app.Status = decryptPtr(app.Status, r.cfg.Security.SecretKey)
		app.ApplicantName = decryptPtr(app.ApplicantName, r.cfg.Security.SecretKey)
		app.District = decryptPtr(app.District, r.cfg.Security.SecretKey)
		app.DistrictLGDCode = decryptPtr(app.DistrictLGDCode, r.cfg.Security.SecretKey)
		app.SubDivision = decryptPtr(app.SubDivision, r.cfg.Security.SecretKey)
		app.SubDivisionLGDCode = decryptPtr(app.SubDivisionLGDCode, r.cfg.Security.SecretKey)
		app.Block = decryptPtr(app.Block, r.cfg.Security.SecretKey)
		app.BlockLGDCode = decryptPtr(app.BlockLGDCode, r.cfg.Security.SecretKey)
		app.Pincode = decryptPtr(app.Pincode, r.cfg.Security.SecretKey)

		applications = append(
			applications,
			app,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &model.PaginatedResponse[model.Application]{
		Data:       applications,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *Repository) GetAttributeIDs(
	ctx context.Context,
	serviceGroupID int64,
) (model.AttributeIDs, error) {

	stmt := `
		SELECT field_name, field_id
		FROM service_mappings
		WHERE service_group_id = @service_group_id
		AND (
		    (
		        section_name ILIKE 'Applicant''s Address details'
		        AND (
		            (field_name = 'District' AND input_type = 'custLGDHierarchy')
		            OR
		            (field_name = 'Block' AND input_type = 'custLGDHierarchy')
		            OR
		            (field_name IN ('Sub-Division', 'Sub Division') AND input_type = 'custLGDHierarchy')
		            OR
		            (field_name IN ('Pincode & Post Office', 'Pincode And Postoffice', 'Post Office & Pincode') AND input_type = 'dropDown')
		        )
		    )
		    OR
		    (
		        section_name ILIKE 'Applicant''s Personal details'
		        AND (
		            field_name IN (
		                'Applicant''s Salutation',
		                'Applicant''s First Name',
		                'Applicant''s Middle Name',
		                'Applicant''s Last Name'
		            )
		        )
		    )
		)
	`

	rows, err := r.db.Query(
		ctx,
		stmt,
		pgx.NamedArgs{
			"service_group_id": serviceGroupID,
		},
	)
	if err != nil {
		return model.AttributeIDs{}, fmt.Errorf(
			"scan attribute ids: %w",
			err,
		)
	}
	defer rows.Close()

	var ids model.AttributeIDs

	for rows.Next() {

		var fieldName string
		var fieldID string

		if err := rows.Scan(
			&fieldName,
			&fieldID,
		); err != nil {
			return model.AttributeIDs{}, fmt.Errorf(
				"scan attribute ids: %w",
				err,
			)
		}

		switch fieldName {

		case "District":
			ids.District = fieldID

		case "Block":
			ids.Block = fieldID

		case "Sub Division":
			ids.SubDivision = fieldID

		case "Pincode & Post Office":
			ids.Pincode = fieldID

		case "Applicant's Salutation":
			ids.Salutation = fieldID

		case "Applicant's First Name":
			ids.FirstName = fieldID

		case "Applicant's Middle Name":
			ids.MiddleName = fieldID

		case "Applicant's Last Name":
			ids.LastName = fieldID
		}
	}

	if err := rows.Err(); err != nil {
		return model.AttributeIDs{}, fmt.Errorf(
			"iterate attribute ids: %w",
			err,
		)
	}

	return ids, nil
}

func (r *Repository) ServiceGroupExists(
	ctx context.Context,
	serviceGroupID int64,
) (bool, error) {
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

	return exists, err
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

		if dec, err := security.DecryptPayload(string(event.RawPayload), r.cfg.Security.SecretKey); err == nil {
			event.RawPayload = []byte(dec)
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

func encryptPtr(val *string, key string) *string {
	if val == nil {
		return nil
	}
	if enc, err := security.EncryptPayload(*val, key); err == nil {
		return &enc
	}
	return val
}

func decryptPtr(val *string, key string) *string {
	if val == nil {
		return nil
	}
	if dec, err := security.DecryptPayload(*val, key); err == nil {
		return &dec
	}
	return val
}
