package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"parser/internal/model"
)

type Parser struct{}

type ParseResult struct {
	WorkflowEvents []model.WorkflowEvent
	Application    []model.Application
}

type Payload struct {
	InitiatedData []map[string]any `json:"initiated_data"`
	ExecutionData []map[string]any `json:"execution_data"`
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(
	raw []byte,
	attributeIDs map[int64]model.AttributeIDs,
) (*ParseResult, error) {

	var payload Payload

	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf(
			"unmarshal payload: %w",
			err,
		)
	}

	result := &ParseResult{}

	// ---------------- Initiated ----------------

	for _, item := range payload.InitiatedData {

		event, err := p.extractInitiatedEvent(item)
		if err != nil {
			return nil, err
		}

		eventRaw, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf(
				"marshal payload: %w",
				err,
			)
		}

		event.RawPayload = eventRaw

		result.WorkflowEvents = append(
			result.WorkflowEvents,
			event,
		)

		application := model.Application{

			RootType: "initiated_data",

			ApplID: ToInt64(item["appl_id"]),

			ServiceID: ToInt64(item["service_id"]),

			ServiceName: stringPtr(
				toString(item["service_name"]),
			),

			AppRefNo: stringPtr(
				toString(item["appl_ref_no"]),
			),

			SubmissionLocation: stringPtr(
				toString(item["submission_location"]),
			),

			SubmittedBy: stringPtr(
				toString(item["applied_by"]),
			),

			SubmissionDate: parseTime(
				toString(item["submission_date"]),
			),

			Status: stringPtr("Applied"),

			ActionNo: intPtr(0),
		}

		attributeDetails, ok := item["attribute_details"].(map[string]any)
		if ok {
			serviceGroupID, _ := strconv.ParseInt(strconv.Itoa(int(application.ServiceID))[:4], 10, 64)
			if ids, ok := attributeIDs[serviceGroupID]; ok {
				populateApplication(
					&application,
					attributeDetails,
					ids,
				)
			}
		}

		result.Application = append(
			result.Application,
			application,
		)
	}

	// ---------------- Execution ----------------

	for _, item := range payload.ExecutionData {

		event, err := p.extractExecutionEvent(item)
		if err != nil {
			return nil, err
		}

		eventRaw, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf(
				"marshal payload: %w",
				err,
			)
		}

		event.RawPayload = eventRaw

		result.WorkflowEvents = append(
			result.WorkflowEvents,
			event,
		)

		taskDetails, ok := item["task_details"].(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"task_details missing",
			)
		}

		userDetail, _ := taskDetails["user_detail"].(map[string]any)

		application := model.Application{

			RootType: "execution_data",

			ApplID: ToInt64(
				taskDetails["appl_id"],
			),

			ServiceID: ToInt64(
				taskDetails["service_id"],
			),

			Status: stringPtr(
				toString(taskDetails["action_taken"]),
			),

			ActionNo: intPtr(
				int(ToInt64(taskDetails["action_no"])),
			),

			SubmittedBy: stringPtr(
				toString(taskDetails["user_name"]),
			),

			SubmissionLocation: stringPtr(
				toString(userDetail["location_name"]),
			),

			SubmissionDate: parseTime(
				toString(taskDetails["executed_time"]),
			),
		}

		result.Application = append(
			result.Application,
			application,
		)
	}

	return result, nil
}

func (p *Parser) extractInitiatedEvent(
	data map[string]any,
) (model.WorkflowEvent, error) {

	submissionDate := parseTime(
		toString(data["submission_date"]),
	)

	return model.WorkflowEvent{
		ApplID:    ToInt64(data["appl_id"]),
		ServiceID: ToInt64(data["service_id"]),

		RootType: "initiated_data",

		TaskName: "Application Submitted",

		ActionNo: 0,
		TaskType: 0,

		ReceivedTime: submissionDate,
		ExecutedTime: submissionDate,
	}, nil
}

func (p *Parser) extractExecutionEvent(
	data map[string]any,
) (model.WorkflowEvent, error) {

	taskDetails, ok := data["task_details"].(map[string]any)
	if !ok {
		return model.WorkflowEvent{},
			fmt.Errorf(
				"task_details missing",
			)
	}

	return model.WorkflowEvent{
		ApplID:    ToInt64(taskDetails["appl_id"]),
		ServiceID: ToInt64(taskDetails["service_id"]),
		TaskName:  toString(taskDetails["task_name"]),
		ActionNo:  int(ToInt64(taskDetails["action_no"])),
		TaskType:  int(ToInt64(taskDetails["task_type"])),
		RootType:  "execution_data",

		ReceivedTime: parseTime(
			toString(taskDetails["received_time"]),
		),

		ExecutedTime: parseTime(
			toString(taskDetails["executed_time"]),
		),
	}, nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int {
	return &i
}

func parseLGD(value string) (code, name string) {

	parts := strings.SplitN(
		value,
		"~",
		2,
	)

	if len(parts) != 2 {
		return "", value
	}

	return parts[0], parts[1]
}

func joinName(parts ...string) string {

	var result []string

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if part != "" {
			result = append(
				result,
				part,
			)
		}
	}

	return strings.Join(
		result,
		" ",
	)
}

func extractPincode(value string) string {

	fields := strings.Fields(value)

	for _, field := range fields {

		if len(field) != 6 {
			continue
		}

		if _, err := strconv.Atoi(field); err == nil {
			return field
		}
	}

	return ""
}

func populateApplication(
	app *model.Application,
	data map[string]any,
	ids model.AttributeIDs,
) {

	if value, ok := data[ids.District]; ok {

		code, name := parseLGD(
			toString(value),
		)

		app.District = stringPtr(name)
		app.DistrictLGDCode = stringPtr(code)
	}

	if value, ok := data[ids.SubDivision]; ok {

		code, name := parseLGD(
			toString(value),
		)

		app.SubDivision = stringPtr(name)
		app.SubDivisionLGDCode = stringPtr(code)
	}

	if value, ok := data[ids.Block]; ok {

		code, name := parseLGD(
			toString(value),
		)

		app.Block = stringPtr(name)
		app.BlockLGDCode = stringPtr(code)
	}

	if value, ok := data[ids.Pincode]; ok {

		_, office := parseLGD(
			toString(value),
		)

		pincode := extractPincode(
			office,
		)

		app.Pincode = stringPtr(
			pincode,
		)
	}

	salutation := ""

	if value, ok := data[ids.Salutation]; ok {
		_, salutation = parseLGD(toString(value))
	}

	applicantName := joinName(
		salutation,
		toString(data[ids.FirstName]),
		toString(data[ids.MiddleName]),
		toString(data[ids.LastName]),
	)

	app.ApplicantName = stringPtr(
		applicantName,
	)
}

func toString(v any) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

func ToInt64(v any) int64 {
	switch x := v.(type) {

	case float64:
		return int64(x)

	case int:
		return int64(x)

	case int64:
		return x

	case string:
		n, err := strconv.ParseInt(
			x,
			10,
			64,
		)
		if err != nil {
			return 0
		}

		return n

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

func toFloat64(v any) float64 {

	switch x := v.(type) {

	case float64:
		return x

	case string:

		f, err := strconv.ParseFloat(
			x,
			64,
		)

		if err != nil {
			return 0
		}

		return f

	default:
		return 0
	}
}
