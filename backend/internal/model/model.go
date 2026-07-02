package model

import "time"

type Service struct {
	ServiceGroupID int64  `db:"service_group_id"`
	ServiceName    string `db:"service_name"`
}

type ServiceMapping struct {
	ID int64 `db:"id"`

	ServiceGroupID int64 `db:"service_group_id"`

	SectionName string `db:"section_name"`
	SectionID   int64  `db:"section_id"`

	FieldID   string `db:"field_id"`
	FieldName string `db:"field_name"`

	InputType string `db:"input_type"`

	FieldSetID *int64 `db:"field_set_id"`
}

type WorkflowEvent struct {
	ID int64 `db:"id" json:"id"`

	ApplID int64 `db:"appl_id" json:"appl_id"`

	ServiceID int64 `db:"service_id" json:"service_id"`

	RootType string `db:"root_type" json:"root_type"`

	TaskName string `db:"task_name" json:"task_name"`

	ActionNo int `db:"action_no" json:"action_no"`

	TaskType int `db:"task_type" json:"task_type"`

	ReceivedTime *time.Time `db:"received_time" json:"received_time"`

	ExecutedTime *time.Time `db:"executed_time" json:"executed_time"`

	RawPayload []byte `db:"raw_payload" json:"raw_payload"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Application struct {
	RootType           string     `json:"root_type" db:"root_type"`
	ApplID             int64      `json:"appl_id" db:"appl_id"`
	AppRefNo           *string    `json:"app_ref_no" db:"app_ref_no"`
	ServiceID          int64      `json:"service_id" db:"service_id"`
	ServiceName        *string    `json:"service_name" db:"service_name"`
	SubmissionLocation *string    `json:"submission_location" db:"submission_location"`
	SubmittedBy        *string    `json:"submitted_by" db:"submitted_by"`
	SubmissionDate     *time.Time `json:"submission_date" db:"submission_date"`

	Status   *string `json:"status" db:"status"`
	ActionNo *int    `json:"action_no" db:"action_no"`

	ApplicantName *string `json:"applicant_name" db:"applicant_name"`

	District        *string `json:"district" db:"district"`
	DistrictLGDCode *string `json:"district_lgd_code" db:"district_lgd_code"`

	SubDivision        *string `json:"sub_division" db:"sub_division"`
	SubDivisionLGDCode *string `json:"sub_division_lgd_code" db:"sub_division_lgd_code"`

	Block        *string `json:"block" db:"block"`
	BlockLGDCode *string `json:"block_lgd_code" db:"block_lgd_code"`

	Pincode *string `json:"pincode" db:"pincode"`
}

type Log struct {
	ID int64 `db:"id" json:"id"`

	Level string `db:"level" json:"level"`

	Source string `db:"source" json:"source"`

	Message string `db:"message" json:"message"`

	Metadata []byte `db:"metadata" json:"metadata"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type AttributeIDs struct {
	District string `db:"district"`

	Block string `db:"block"`

	SubDivision string `db:"sub_division"`

	Pincode string `db:"pincode"`

	FirstName string `db:"first_name"`

	MiddleName string `db:"middle_name"`

	LastName string `db:"last_name"`

	Salutation string `db:"salutation"`
}

type PaginatedResponse[T interface{}] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
