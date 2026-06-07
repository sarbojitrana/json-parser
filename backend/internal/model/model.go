package model

import "time"

type Service struct {
	ServiceID   int64  `db:"service_id"`
	ServiceName string `db:"service_name"`
}

type ServiceMapping struct {
	ID int64 `db:"id"`

	ServiceID int64 `db:"service_id"`

	SectionName string `db:"section_name"`
	SectionID   int64  `db:"section_id"`

	FieldID   string `db:"field_id"`
	FieldName string `db:"field_name"`

	InputType string `db:"input_type"`

	FieldSetID *int64 `db:"field_set_id"`
}

type WorkflowEvent struct {
	ID int64 `db:"id"`

	ApplID int64 `db:"appl_id"`

	ServiceID int64 `db:"service_id"`

	RootType string `db:"root_type"`

	TaskName string `db:"task_name"`

	ActionNo int `db:"action_no"`
	TaskType int `db:"task_type"`

	ReceivedTime *time.Time `db:"received_time"`
	ExecutedTime *time.Time `db:"executed_time"`

	RawPayload []byte `db:"raw_payload"`

	CreatedAt time.Time `db:"created_at"`
}