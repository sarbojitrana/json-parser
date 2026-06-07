package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"parser/internal/model"
)

type Parser struct{}

type Payload struct {
	ExecutionData []map[string]any `json:"execution_data"`
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(
	raw []byte,
	rootType string,
) ([]model.WorkflowEvent, error){
	var payload Payload

	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	events := make([]model.WorkflowEvent, 0, len(payload.ExecutionData))

	for _, item := range payload.ExecutionData {

		event, err := p.extractEvent(item, rootType)
		if err != nil {
			return nil, err
		}

		eventRaw, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("marshal event payload: %w", err)
		}

		event.RawPayload = eventRaw

		events = append(events, event)
	}

	return events, nil
}

func (p *Parser) extractEvent(data map[string]any, rootType string) (model.WorkflowEvent, error) {
	taskDetails, ok := data["task_details"].(map[string]any)
	if !ok {
		return model.WorkflowEvent{}, fmt.Errorf("task_details missing")
	}

	return model.WorkflowEvent{
		ApplID:       toInt64(taskDetails["appl_id"]),
		ServiceID:    toInt64(taskDetails["service_id"]),
		TaskName:     toString(taskDetails["task_name"]),
		ActionNo:     int(toInt64(taskDetails["action_no"])),
		TaskType:     int(toInt64(taskDetails["task_type"])),
		RootType:     rootType,
		ReceivedTime: parseTime(toString(taskDetails["received_time"])),
		ExecutedTime: parseTime(toString(taskDetails["executed_time"])),
	}, nil
}

func toString(v any) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

func toInt64(v any) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int:
		return int64(x)
	case int64:
		return x
	default:
		return 0
	}
}

func parseTime(value string) *time.Time {
	if value == "" {
		return nil
	}

	t, err := time.Parse(
		"02-01-2006 15:04:05",
		value,
	)
	if err != nil {
		return nil
	}

	return &t
}